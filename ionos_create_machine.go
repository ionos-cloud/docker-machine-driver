package ionoscloud

import (
	"encoding/base64"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/docker/machine/libmachine/log"
	"github.com/ionos-cloud/docker-machine-driver/internal/pointer"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
)

func (d *Driver) getCubeTemplateUuid() (string, error) {
	templates, err := d.client().GetTemplates()
	if err != nil {
		return "", err
	}

	for _, template := range *templates.Items {
		if *template.Properties.Name == d.Template {
			return *template.Id, nil
		}
	}
	return "", err
}

func (d *Driver) addSSHUserToYaml() (string, error) {
	commonUser := map[interface{}]interface{}{
		"name":                d.SSHUser,
		"lock_passwd":         true,
		"sudo":                "ALL=(ALL) NOPASSWD:ALL",
		"create_groups":       false,
		"no_user_group":       true,
		"ssh_authorized_keys": []string{d.SSHKey},
	}

	return d.client().UpdateCloudInitFile(d.CloudInit, "users", []interface{}{commonUser}, false, "append")
}

func (d *Driver) CreateDataCenterIfNeeded() (err error) {
	var dc *sdkgo.Datacenter
	if d.DatacenterId == "" {
		d.DCExists = false
		var err error
		log.Debugf("Creating datacenter...")
		dc, err = d.client().CreateDatacenter(d.DatacenterName, d.Location)
		if err != nil {
			return fmt.Errorf("error creating datacenter: %w", err)
		}
	} else {
		d.DCExists = true
		log.Debugf("Getting existing datacenter..")
		dc, err = d.client().GetDatacenter(d.DatacenterId)
		if err != nil {
			return fmt.Errorf("error getting datacenter: %w", err)
		}
	}
	if dcId, ok := dc.GetIdOk(); ok && dcId != nil {
		d.DatacenterId = *dcId
		log.Debugf("Datacenter ID: %v", d.DatacenterId)
	}
	return nil
}

func (d *Driver) CreateLanIfNeeded() (err error) {
	if d.LanId == "" {
		lan, err := d.client().CreateLan(d.DatacenterId, d.LanName, !d.PrivateLan)
		if err != nil {
			err = fmt.Errorf("error creating LAN: %w", err)
			return err
			// TODO: <---
		}
		if lanId, ok := lan.GetIdOk(); ok && lanId != nil {
			d.LanId = *lanId
			log.Debugf("Lan ID: %v", d.LanId)
		}
	}
	return nil
}

func (d *Driver) GetFinalUserData() (userdata string, err error) {
	givenB64CloudInit, _ := base64.StdEncoding.DecodeString(d.CloudInitB64)
	if ud := getPropertyWithFallback(string(givenB64CloudInit), d.CloudInit, ""); ud != "" {
		// Provided B64 User Data has priority over UI provided User Data
		d.CloudInit = ud
	}

	if d.SSHUser != "root" || d.SSHInCloudInit {
		d.CloudInit, err = d.addSSHUserToYaml()
		if err != nil {
			return "", err
		}
	}
	d.CloudInit, err = d.client().UpdateCloudInitFile(
		d.CloudInit, "hostname", []interface{}{d.MachineName}, true, "skip",
	)
	if err != nil {
		return "", err
	}
	ud := base64.StdEncoding.EncodeToString([]byte(d.CloudInit))
	return ud, nil
}

func (d *Driver) CreateIonosServer() (err error) {
	// Creating the server with the volume attached
	imageIdentifier, err := d.getImageIdOrAlias(d.Image)
	if err != nil {
		return fmt.Errorf("error getting image/alias %s: %w", d.Image, err)
	}

	ud, err := d.GetFinalUserData()
	if err != nil {
		return fmt.Errorf("error with user data: %w", err)
	}

	// Volume
	sshKeys := &[]string{}
	if !d.SSHInCloudInit {
		sshKeys = &[]string{d.SSHKey}
	}
	imagePassword := &d.ImagePassword
	if d.ImagePassword == "" {
		imagePassword = nil
	}
	floatDiskSize := float32(d.DiskSize)

	volumeProperties := sdkgo.VolumeProperties{
		Type:          &d.DiskType,
		Name:          &d.MachineName,
		ImagePassword: imagePassword,
		SshKeys:       sshKeys,
		UserData:      &ud,
	}

	if !d.UseAlias {
		log.Infof("Image Id: %v", imageIdentifier)
		volumeProperties.Image = &imageIdentifier
	} else {
		log.Infof("Image Alias: %v", imageIdentifier)
		volumeProperties.ImageAlias = &imageIdentifier
	}

	serverToCreate := sdkgo.Server{}
	if d.ServerType == "ENTERPRISE" {
		serverToCreate.Properties = &sdkgo.ServerProperties{
			Name:             &d.MachineName,
			Ram:              pointer.From(int32(d.Ram)),
			Cores:            pointer.From(int32(d.Cores)),
			CpuFamily:        &d.CpuFamily,
			AvailabilityZone: &d.ServerAvailabilityZone,
		}
		volumeProperties.Size = &floatDiskSize
		volumeProperties.AvailabilityZone = &d.VolumeAvailabilityZone
	} else {
		TemplateUuid, err := d.getCubeTemplateUuid()
		if err != nil {
			return fmt.Errorf("error getting CUBE Template UUID from Template %s: %w", d.Template, err)
		}
		serverToCreate.Properties = &sdkgo.ServerProperties{
			Name:         &d.MachineName,
			Type:         &d.ServerType,
			TemplateUuid: &TemplateUuid,
		}
		volumeProperties.Type = pointer.From("DAS")
	}

	attachedVolumes := sdkgo.NewAttachedVolumesWithDefaults()
	attachedVolumes.Items = &[]sdkgo.Volume{
		{
			Properties: &volumeProperties,
		},
	}
	serverToCreate.Entities = sdkgo.NewServerEntitiesWithDefaults()
	serverToCreate.Entities.SetVolumes(*attachedVolumes)

	var ipsForAttachedNic *[]string

	if len(d.NicIps) != 0 {
		ipsForAttachedNic = &d.NicIps // If IPs are provided use those
	} else if d.IsLanPrivate {
		ipsForAttachedNic = nil // Let CloudAPI generate an IP, which we can later use for the subnet
	} else {
		ipsForAttachedNic = d.ReservedIps // For public NICs we use the generated IPs
	}

	attachedNics := sdkgo.NewNicsWithDefaults()

	lanId, _ := strconv.Atoi(d.LanId)
	nicProperties := &sdkgo.NicProperties{
		Name: &d.MachineName,
		Lan:  sdkgo.PtrInt32(int32(lanId)),
		Ips:  ipsForAttachedNic,
		Dhcp: &d.NicDhcp,
	}

	attachedNics.Items = &[]sdkgo.Nic{
		{
			Properties: nicProperties,
		},
	}

	for _, additionalLanId := range d.AdditionalLansIds {
		additionalNic := sdkgo.Nic{
			Properties: &sdkgo.NicProperties{
				Name: sdkgo.PtrString(d.MachineName + " " + fmt.Sprint(additionalLanId)),
				Lan:  sdkgo.PtrInt32(int32(additionalLanId)),
				Ips:  nil,
				Dhcp: sdkgo.PtrBool(true),
			},
		}
		*attachedNics.Items = append(*attachedNics.Items, additionalNic)
	}

	serverToCreate.Entities.SetNics(*attachedNics)

	server, err := d.client().CreateServer(d.DatacenterId, serverToCreate)
	if err != nil {
		return err
	}
	if serverId, ok := server.GetIdOk(); ok && serverId != nil {
		d.ServerId = *serverId
		log.Debugf("Server ID: %v", d.ServerId)
	} else {
		return fmt.Errorf("error getting server: d.ServerId is empty")
	}
	return nil
}

func (d *Driver) CreateIonosNatAndSetIp() (err error) {
	nic, err := d.client().GetNic(d.DatacenterId, d.ServerId, d.NicId)
	if err != nil {
		return fmt.Errorf("error getting NIC: %w", err)
	}

	nicIps := &[]string{}
	if nicProp, ok := nic.GetPropertiesOk(); ok && nicProp != nil {
		nicIps = nicProp.GetIps()
	}

	// --- NAT ---
	if d.CreateNat {
		// TODO: Were CreateNat in a deeper scope, we wouldn't have the need of these variables (they are here to avoid function-wide side-effects)
		natPublicIps := d.ReservedIps
		if d.NatPublicIps != nil {
			natPublicIps = &d.NatPublicIps
		}
		natLansToGateways := &map[string][]string{"1": {"10.0.0.1"}} // User has to add this ip route to their cloud config if he doesn't set a custom gateway IP
		if d.NatLansToGateways != nil {
			natLansToGateways = &d.NatLansToGateways
		}
		sourceSubnet := net.ParseIP((*nicIps)[0]).Mask(net.CIDRMask(24, 32)).String() + "/24"
		nat, err := d.client().CreateNat(d.DatacenterId, d.NatName, *natPublicIps, d.NatFlowlogs, d.NatRules, *natLansToGateways, sourceSubnet, d.SkipDefaultNatRules)
		if err != nil {
			return err
		}
		log.Debugf("Nat ID: %s", *nat.Id)
		d.NatId = *nat.Id // NatId is used later to retrieve public IP, etc.
		d.IPAddress = (*natPublicIps)[0]
		log.Infof(d.IPAddress)
	} else if d.NatId != "" {
		nat, _ := d.client().GetNat(d.DatacenterId, d.NatId)

		foundNatLan := false
		lanIdInt, _ := strconv.Atoi(d.LanId)
		for _, natLan := range *nat.Properties.Lans {
			if *natLan.Id == int32(lanIdInt) {
				foundNatLan = true
				break
			}
		}
		if !foundNatLan {
			// connect lan to nat
			natLans := append(*nat.Properties.Lans, sdkgo.NatGatewayLanProperties{Id: pointer.From(int32(lanIdInt)), GatewayIps: nil})
			_, err = d.client().PatchNat(d.DatacenterId, d.NatId, *nat.Properties.Name, *nat.Properties.PublicIps, natLans)

			if err != nil {
				return err
			}
		}

		d.IPAddress = (*nat.Properties.PublicIps)[0]
		log.Infof(d.IPAddress)
	} else {
		if len(*nicIps) > 0 {
			d.IPAddress = (*nicIps)[0]
			log.Infof(d.IPAddress)
		}
	}
	return nil
}

func (d *Driver) CreateIonosMachine() (err error) {
	log.Infof("Creating SSH key123...")
	if d.SSHKey == "" {
		d.SSHKey, err = d.createSSHKey()
		if err != nil {
			return fmt.Errorf("error creating SSH keys: %w", err)
		}
		log.Debugf("SSH Key generated in file: %v", d.publicSSHKeyPath())
	}
	err = d.CreateDataCenterIfNeeded()
	if err != nil {
		return err
	}
	err = d.CreateLanIfNeeded()
	if err != nil {
		return err
	}

	lan, err := d.client().GetLan(d.DatacenterId, d.LanId)
	if err != nil {
		return fmt.Errorf("error getting LAN: %w", err)
	}
	if lanProp, ok := lan.GetPropertiesOk(); ok && lanProp != nil {
		if public, ok := lanProp.GetPublicOk(); ok && public != nil {
			d.IsLanPrivate = !*public
		}
	}

	// Reserve IP if needed
	if !d.IsLanPrivate && !(len(d.NicIps) != 0) ||
		d.CreateNat && d.NatPublicIps == nil {
		ipBlock, err := d.client().CreateIpBlock(1, d.Location)
		if err != nil {
			return fmt.Errorf("error creating ipblock: %w", err)
		}
		if ipBlockId, ok := ipBlock.GetIdOk(); ok && ipBlockId != nil {
			d.IpBlockId = *ipBlockId
			log.Debugf("IpBlock ID: %v", d.IpBlockId)
		}
		d.ReservedIps, err = d.client().GetIpBlockIps(ipBlock)
		if err != nil {
			return err
		}
	}

	err = d.CreateIonosServer()
	if err != nil {
		return err
	}

	server, err := d.client().GetServer(d.DatacenterId, d.ServerId, 2)
	if err != nil {
		return fmt.Errorf("error getting server by id: %w", err)
	}
	d.VolumeId = *(*server.Entities.GetVolumes().Items)[0].GetId()
	log.Debugf("Volume ID: %v", d.VolumeId)

	nics := server.Entities.GetNics()
	for _, nic := range *nics.Items {
		if *nic.Properties.Name == d.MachineName {
			d.NicId = *nic.Id
			log.Debugf("Nic ID: %v", d.NicId)
		} else {
			d.AdditionalNicsIds = append(d.AdditionalNicsIds, *nic.Id)
		}
	}

	if d.WaitForIpChange {
		err := d.client().WaitForNicIpChange(d.DatacenterId, d.ServerId, d.NicId, d.WaitForIpChangeTimeout)
		if err != nil {
			return err
		}
	}

	err = d.CreateIonosNatAndSetIp()
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) getImageIdOrAlias(imageName string) (string, error) {
	d.UseAlias = false
	// First, look if the provided parameter matches an alias, if a match is found we return the image alias
	regionId, locationId := d.getRegionIdAndLocationId()
	location, err := d.client().GetLocationById(regionId, locationId)
	if err != nil {
		return "", err
	}
	if locationProp, ok := location.GetPropertiesOk(); ok && locationProp != nil {
		if imageAliases, ok := locationProp.GetImageAliasesOk(); ok && imageAliases != nil {
			for _, alias := range *imageAliases {
				if alias == imageName {
					d.UseAlias = true
					return imageName, nil
				}
			}
		}
	}
	// Second, check if the imageName provided is actually an imageId.
	// If an image is found, return the imageId
	imageFound, err := d.client().GetImageById(imageName)
	if err != nil {
		if !strings.Contains(err.Error(), "no image found") {
			return "", err
		}
	} else {
		if imageId, ok := imageFound.GetIdOk(); ok && imageId != nil {
			d.UseAlias = false
			return *imageId, nil
		}
	}
	// If no alias and id match, we do extended search, considering the image parameter
	// set by the user to be part of the image name and checking the location & image type.
	// If the extended search is successful, return the imageId.
	// Example: if the user sets: Ubuntu-20.04, the driver will know which image to use.
	images, err := d.client().GetImages()
	if err != nil {
		return "", err
	}

	if imagesItems, ok := images.GetItemsOk(); ok && imagesItems != nil {
		for _, image := range *imagesItems {
			imgName := ""
			if imgProp, ok := image.GetPropertiesOk(); ok && imgProp != nil {
				if name, ok := imgProp.GetNameOk(); ok && name != nil {
					if *name != "" {
						imgName = *name
					}
				}
			}
			diskType := d.DiskType
			if d.DiskType == "SSD" {
				diskType = defaultDiskType
			}
			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(imageName)) &&
				*image.Properties.ImageType == diskType && *image.Properties.Location == d.Location {
				d.UseAlias = false
				return *image.Id, nil
			}
		}
	}
	return "", nil
}

func (d *Driver) getRegionIdAndLocationId() (regionId, locationId string) {
	ids := strings.Split(d.Location, "/")
	// location has standard format: {regionId}/{locationId}
	if len(ids) != 2 {
		log.Errorf("error getting Region Id and Location Id from %s", d.Location)
		return "", ""
	}
	return ids[0], ids[1]
}
