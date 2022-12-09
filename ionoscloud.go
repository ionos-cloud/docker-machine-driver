package ionoscloud

import (
	"context"
	"encoding/base64"
	b64 "encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/hashicorp/go-multierror"
	"github.com/ionos-cloud/docker-machine-driver/utils"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
)

const (
	flagEndpoint               = "ionoscloud-endpoint"
	flagUsername               = "ionoscloud-username"
	flagPassword               = "ionoscloud-password"
	flagToken                  = "ionoscloud-token"
	flagServerCores            = "ionoscloud-cores"
	flagServerRam              = "ionoscloud-ram"
	flagServerCpuFamily        = "ionoscloud-cpu-family"
	flagServerAvailabilityZone = "ionoscloud-server-availability-zone"
	flagDiskSize               = "ionoscloud-disk-size"
	flagDiskType               = "ionoscloud-disk-type"
	flagImage                  = "ionoscloud-image"
	flagImagePassword          = "ionoscloud-image-password"
	flagLocation               = "ionoscloud-location"
	flagDatacenterId           = "ionoscloud-datacenter-id"
	flagVolumeAvailabilityZone = "ionoscloud-volume-availability-zone"
	flagUserData               = "ionoscloud-user-data"
	flagSSHUser                = "ionoscloud-ssh-user"
	flagUserDataB64            = "ionoscloud-user-data-b64"
)

const (
	defaultRegion           = "us/las"
	defaultImageAlias       = "ubuntu:latest"
	defaultImagePassword    = "abcde12345" // Must contain both letters and numbers, at least 8 characters
	defaultCpuFamily        = "AMD_OPTERON"
	defaultAvailabilityZone = "AUTO"
	defaultDiskType         = "HDD"
	defaultSize             = 10
	driverName              = "ionoscloud"
)

const (
	rollingBackNotice = "WARNING: Error creating machine. Rolling back..."
	driverVersionDev  = "DEV"
)

// DriverVersion will be set at every new release
// For working locally with the Docker-Machine-Driver,
// it will be set to `DEV`.
var DriverVersion string

type Driver struct {
	*drivers.BaseDriver
	client func() utils.ClientService

	URL      string
	Username string
	Password string
	Token    string

	Ram                    int
	Cores                  int
	SSHKey                 string
	SSHUser                string
	DiskSize               int
	DiskType               string
	Image                  string
	ImagePassword          string
	Size                   int
	Location               string
	CpuFamily              string
	DCExists               bool
	UseAlias               bool
	VolumeAvailabilityZone string
	ServerAvailabilityZone string
	LanId                  string
	DatacenterId           string
	VolumeId               string
	NicId                  string
	ServerId               string
	IpBlockId              string
	UserData               string
	UserDataB64            string

	// Driver Version
	Version string
}

// NewDriver returns a new driver instance.
func NewDriver(hostName, storePath string) drivers.Driver {
	return NewDerivedDriver(hostName, storePath)
}

func NewDerivedDriver(hostName, storePath string) *Driver {
	var httpUserAgent string
	v := getDriverVersion(DriverVersion)
	driver := &Driver{
		Size:     defaultSize,
		Location: defaultRegion,
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
		Version: v,
	}
	if v != driverVersionDev {
		httpUserAgent = fmt.Sprintf("docker-machine-driver-ionoscloud/v%v", driver.Version)
	} else {
		httpUserAgent = fmt.Sprintf("docker-machine-driver-ionoscloud/%v", driver.Version)
	}
	driver.client = func() utils.ClientService {
		return utils.New(context.TODO(), driver.Username, driver.Password, driver.Token, driver.URL, httpUserAgent)
	}
	return driver
}

// GetCreateFlags returns list of create flags driver accepts.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_ENDPOINT",
			Name:   flagEndpoint,
			Value:  sdkgo.DefaultIonosServerUrl,
			Usage:  "Ionos Cloud API Endpoint",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_USERNAME",
			Name:   flagUsername,
			Usage:  "Ionos Cloud Username",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_PASSWORD",
			Name:   flagPassword,
			Usage:  "Ionos Cloud Password",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_TOKEN",
			Name:   flagToken,
			Usage:  "Ionos Cloud Token",
		},
		mcnflag.IntFlag{
			EnvVar: "IONOSCLOUD_CORES",
			Name:   flagServerCores,
			Value:  4,
			Usage:  "Ionos Cloud Server Cores (2, 3, 4, 5, 6, etc.)",
		},
		mcnflag.IntFlag{
			EnvVar: "IONOSCLOUD_RAM",
			Name:   flagServerRam,
			Value:  2048,
			Usage:  "Ionos Cloud Server Ram in MB(1024, 2048, 3072, 4096, etc.)",
		},
		mcnflag.IntFlag{
			EnvVar: "IONOSCLOUD_DISK_SIZE",
			Name:   flagDiskSize,
			Value:  50,
			Usage:  "Ionos Cloud Volume Disk-Size in GB(10, 50, 100, 200, 400)",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_IMAGE",
			Name:   flagImage,
			Value:  defaultImageAlias,
			Usage:  "Ionos Cloud Image Id or Alias (ubuntu:latest, ubuntu:20.04)",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_IMAGE_PASSWORD",
			Name:   flagImagePassword,
			Value:  defaultImagePassword,
			Usage:  "Ionos Cloud Image Password to be able to access the server from DCD platform",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_LOCATION",
			Name:   flagLocation,
			Value:  defaultRegion,
			Usage:  "Ionos Cloud Location",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_DISK_TYPE",
			Name:   flagDiskType,
			Value:  defaultDiskType,
			Usage:  "Ionos Cloud Volume Disk-Type (HDD, SSD)",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_CPU_FAMILY",
			Name:   flagServerCpuFamily,
			Value:  defaultCpuFamily,
			Usage:  "Ionos Cloud Server CPU families (AMD_OPTERON, INTEL_XEON, INTEL_SKYLAKE)",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_DATACENTER_ID",
			Name:   flagDatacenterId,
			Usage:  "Ionos Cloud Virtual Data Center Id",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_VOLUME_ZONE",
			Name:   flagVolumeAvailabilityZone,
			Value:  defaultAvailabilityZone,
			Usage:  "Ionos Cloud Volume Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_SERVER_ZONE",
			Name:   flagServerAvailabilityZone,
			Value:  defaultAvailabilityZone,
			Usage:  "Ionos Cloud Server Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_USER_DATA",
			Name:   flagUserData,
			Usage:  "The cloud-init configuration for the volume as a multi-line string",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_USER_DATA_B64",
			Name:   flagUserDataB64,
			Usage:  "The cloud-init configuration for the volume as base64 encoded string",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_SSH_USER",
			Name:   flagSSHUser,
			Usage:  "The name of the user the driver will use for ssh",
		},
	}
}

// SetConfigFromFlags initializes driver values from the command line values.
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.URL = opts.String(flagEndpoint)
	d.Username = opts.String(flagUsername)
	d.Password = opts.String(flagPassword)
	d.Token = opts.String(flagToken)
	d.DiskSize = opts.Int(flagDiskSize)
	d.Image = opts.String(flagImage)
	d.ImagePassword = opts.String(flagImagePassword)
	d.Cores = opts.Int(flagServerCores)
	d.Ram = opts.Int(flagServerRam)
	d.Location = opts.String(flagLocation)
	d.DiskType = opts.String(flagDiskType)
	d.CpuFamily = opts.String(flagServerCpuFamily)
	d.DatacenterId = opts.String(flagDatacenterId)
	d.VolumeAvailabilityZone = opts.String(flagVolumeAvailabilityZone)
	d.ServerAvailabilityZone = opts.String(flagServerAvailabilityZone)
	d.UserData = opts.String(flagUserData)
	d.SSHUser = opts.String(flagSSHUser)
	d.UserDataB64 = opts.String(flagUserDataB64)

	d.SwarmMaster = opts.Bool("swarm-master")
	d.SwarmHost = opts.String("swarm-host")
	d.SwarmDiscovery = opts.String("swarm-discovery")
	d.SetSwarmConfigFromFlags(opts)

	if d.URL == "" {
		d.URL = sdkgo.DefaultIonosServerUrl
	}

	return nil
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return driverName
}

// PreCreateCheck validates if driver values are valid to create the machine.
func (d *Driver) PreCreateCheck() error {
	log.Infof("IONOS Cloud Driver Version: %s", d.Version)
	log.Infof("SDK-GO Version: %s", sdkgo.Version)
	if d.Token == "" {
		if d.Username == "" && d.Password == "" {
			return fmt.Errorf("please provide username($IONOSCLOUD_USERNAME) and password($IONOSCLOUD_PASSWORD) or token($IONOSCLOUD_TOKEN) to authenticate")
		}
		if d.Username == "" {
			return fmt.Errorf("please provide username as parameter --ionoscloud-username or as environment variable $IONOSCLOUD_USERNAME")
		}
		if d.Password == "" {
			return fmt.Errorf("please provide password as parameter --ionoscloud-password or as environment variable $IONOSCLOUD_PASSWORD")
		}
	}
	if d.DatacenterId != "" {
		d.DCExists = true
		dc, err := d.client().GetDatacenter(d.DatacenterId)
		if err != nil {
			return err
		}
		if dcProp, ok := dc.GetPropertiesOk(); ok && dcProp != nil {
			if name, ok := dcProp.GetNameOk(); ok && name != nil {
				log.Info("Creating machine under " + *name + " datacenter")
			}
			// If the datacenter already exists, update the driver location
			// from the default one to the datacenter's location
			if dcLocation, ok := dcProp.GetLocationOk(); ok && dcLocation != nil {
				d.Location = *dcLocation
			}
		}
	} else {
		d.DCExists = false
	}
	if imageId, err := d.getImageId(d.Image); err != nil && imageId == "" {
		return fmt.Errorf("error getting image/alias %s: %v", d.Image, err)
	}

	return nil
}

func (d *Driver) addSSHUserToYaml() (string, error) {
	var (
		sshUser     = d.SSHUser
		sshkey      = d.SSHKey
		yamlcontent = d.UserData
	)
	cf := make(map[interface{}]interface{})
	if err := yaml.Unmarshal([]byte(yamlcontent), &cf); err != nil {
		return "", err
	}

	commonUser := map[interface{}]interface{}{
		"name":        sshUser,
		"lock_passwd": true,
		"ssh_authorized_keys": []string{
			sshkey,
		},
	}

	switch "linux" {
	default:
		// implements https://github.com/canonical/cloud-init/blob/master/cloudinit/config/cc_users_groups.py#L28-L71
		// technically not in the spec, see this code for context
		// https://github.com/canonical/cloud-init/blob/master/cloudinit/distros/__init__.py#L394-L397
		commonUser["sudo"] = "ALL=(ALL) NOPASSWD:ALL"
		commonUser["create_groups"] = false
		commonUser["no_user_group"] = true

	// Administrator is the default ssh user on Windows Server 2019/2022
	// This implements cloudbase-init for Windows VMs as cloud-init doesn't support Windows
	// https://cloudbase-init.readthedocs.io/en/latest/
	// On Windows, primary_group and groups are concatenated.
	case "windows":
		commonUser["inactive"] = false
	}

	if val, ok := cf["users"]; ok {
		u := val.([]interface{})
		cf["users"] = append(u, commonUser)
	} else {
		users := make([]interface{}, 1)
		users[0] = commonUser
		cf["users"] = users
	}

	yaml, err := yaml.Marshal(cf)
	if err != nil {
		return "", err
	}
	return string(yaml), nil
}

func getPropertyWithFallback[T comparable](p1 T, p2 T, empty T) T {
	if p1 == empty {
		return p2
	}
	return p1
}

// Create creates the machine.
func (d *Driver) Create() error {
	var err error
	log.Infof("Creating SSH key...")
	if d.SSHKey == "" {
		d.SSHKey, err = d.createSSHKey()
		if err != nil {
			return fmt.Errorf("error creating SSH keys: %v", err)
		}
		log.Debugf("SSH Key generated in file: %v", d.publicSSHKeyPath())
	}

	rootSSHKey := d.SSHKey

	if d.SSHUser != "" {
		rootSSHKey = ""
		if d.UserData == "" {
			d.UserData = "#cloud-config\n" + d.UserData
		}
		newUserData, _ := d.addSSHUserToYaml()
		d.UserData = b64.StdEncoding.EncodeToString([]byte(d.UserData + newUserData))
	}

	result, err := d.getImageId(d.Image)
	if err != nil {
		return fmt.Errorf("error getting image/alias %s: %v", d.Image, err)
	}
	var alias string
	if d.UseAlias {
		alias = result
	}

	var dc *sdkgo.Datacenter
	if d.DatacenterId == "" {
		d.DCExists = false
		var err error
		log.Debugf("Creating datacenter...")
		dc, err = d.client().CreateDatacenter(d.MachineName, d.Location)
		if err != nil {
			return err
		}
	} else {
		d.DCExists = true
		log.Debugf("Getting existing datacenter..")
		dc, err = d.client().GetDatacenter(d.DatacenterId)
		if err != nil {
			return err
		}
	}
	if dcId, ok := dc.GetIdOk(); ok && dcId != nil {
		d.DatacenterId = *dcId
		log.Debugf("Datacenter ID: %v", d.DatacenterId)
	}

	ipBlock, err := d.client().CreateIpBlock(int32(1), d.Location)
	if err != nil {
		return err
	}
	if ipBlockId, ok := ipBlock.GetIdOk(); ok && ipBlockId != nil {
		d.IpBlockId = *ipBlockId
		log.Debugf("IpBlock ID: %v", d.IpBlockId)
	}

	lan, err := d.client().CreateLan(d.DatacenterId, d.MachineName, true)
	if err != nil {
		log.Warn(rollingBackNotice)
		if removeErr := d.Remove(); removeErr != nil {
			return fmt.Errorf("failed to create machine due to error: %v\n Removing created resources: %v", err, removeErr)
		}
		return err
	}
	if lanId, ok := lan.GetIdOk(); ok && lanId != nil {
		d.LanId = *lanId
		log.Debugf("Lan ID: %v", d.LanId)
	}

	server, err := d.client().CreateServer(d.DatacenterId, d.Location, d.MachineName, d.CpuFamily, d.ServerAvailabilityZone, int32(d.Ram), int32(d.Cores))
	if err != nil {
		log.Warn(rollingBackNotice)
		if removeErr := d.Remove(); removeErr != nil {
			return fmt.Errorf("failed to create machine due to error: %v\n Removing created resources: %v", err, removeErr)
		}
		return err
	}
	if serverId, ok := server.GetIdOk(); ok && serverId != nil {
		d.ServerId = *serverId
		log.Debugf("Server ID: %v", d.ServerId)
	}

	properties := utils.ClientVolumeProperties{
		DiskType:      d.DiskType,
		Name:          d.MachineName,
		ImagePassword: d.ImagePassword,
		Zone:          d.VolumeAvailabilityZone,
		SshKey:        rootSSHKey,
		DiskSize:      float32(d.DiskSize),
	}

	if ud := getPropertyWithFallback(base64.StdEncoding.EncodeToString([]byte(d.UserData)), d.UserDataB64, ""); ud != "" {
		log.Infof("Using user data: %s", ud)
		properties.UserData = ud
	}

	if !d.UseAlias {
		log.Infof("Image Id: %v", result)
		properties.ImageId = result
	} else {
		log.Infof("Image Alias: %v", alias)
		properties.ImageAlias = alias
	}
	volume, err := d.client().CreateAttachVolume(d.DatacenterId, d.ServerId, &properties)
	if err != nil {
		log.Warn(rollingBackNotice)
		if removeErr := d.Remove(); removeErr != nil {
			return fmt.Errorf("failed to create machine due to error: %v\n Removing created resources: %v", err, removeErr)
		}
		return err
	}
	if volumeId, ok := volume.GetIdOk(); ok && volumeId != nil {
		d.VolumeId = *volumeId
		log.Debugf("Volume ID: %v", d.VolumeId)
	}

	l, _ := strconv.Atoi(d.LanId)
	ips, err := d.client().GetIpBlockIps(ipBlock)
	if err != nil {
		return err
	}

	nic, err := d.client().CreateAttachNIC(d.DatacenterId, d.ServerId, d.MachineName, true, int32(l), ips)
	if err != nil {
		log.Warn(rollingBackNotice)
		if removeErr := d.Remove(); removeErr != nil {
			return fmt.Errorf("failed to create machine due to error: %v\n Removing created resources: %v", err, removeErr)
		}
		return err
	}
	if nicId, ok := nic.GetIdOk(); ok && nicId != nil {
		d.NicId = *nic.Id
		log.Debugf("Nic ID: %v", d.NicId)
	}

	if len(*ips) > 0 {
		ipBlockIps := *ips
		d.IPAddress = ipBlockIps[0]
		log.Info(d.IPAddress)
	}

	return nil
}

// Remove deletes the machine and resources associated to it.
func (d *Driver) Remove() error {
	var result *multierror.Error

	// NOTE:
	//   - if a resource is already gone or errors occur while deleting it, we
	//     continue removing other resources instead of failing

	log.Warn("NOTICE: Please check IONOS Cloud Console/CLI to ensure there are no leftover resources.")
	log.Info("Starting deleting resources...")

	log.Debugf("Datacenter Id: %v", d.DatacenterId)
	log.Debugf("Server Id: %v", d.ServerId)
	log.Debugf("Starting deleting Nic with Id: %v", d.NicId)
	err := d.client().RemoveNic(d.DatacenterId, d.ServerId, d.NicId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	log.Debugf("Starting deleting Volume with Id: %v", d.VolumeId)
	err = d.client().RemoveVolume(d.DatacenterId, d.VolumeId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	log.Debugf("Starting deleting Server with Id: %v", d.ServerId)
	err = d.client().RemoveServer(d.DatacenterId, d.ServerId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	log.Debugf("Starting deleting LAN with Id: %v", d.LanId)
	err = d.client().RemoveLan(d.DatacenterId, d.LanId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	// If the DataCenter existed before creating the machine, do not delete it at clean-up
	if !d.DCExists {
		log.Debugf("Starting deleting Datacenter with Id: %v", d.DatacenterId)
		err = d.client().RemoveDatacenter(d.DatacenterId)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	log.Debugf("Starting deleting IpBlock with Id: %v", d.IpBlockId)
	err = d.client().RemoveIpBlock(d.IpBlockId)
	if err != nil {
		result = multierror.Append(result, err)
	}

	if result != nil {
		return result.ErrorOrNil()
	}
	return nil
}

// Start issues a power on for the machine instance.
func (d *Driver) Start() error {
	serverState, err := d.GetState()
	if err != nil {
		return fmt.Errorf("error getting state: %v", err)
	}
	if serverState != state.Running {
		err = d.client().StartServer(d.DatacenterId, d.ServerId)
		if err != nil {
			return err
		}
	} else {
		log.Info("Host is already running or starting")
	}
	return nil
}

// Stop issues a power off for the machine instance.
func (d *Driver) Stop() error {
	vmState, err := d.GetState()
	if err != nil {
		return fmt.Errorf("error getting state: %v", err)
	}
	if vmState == state.Stopped {
		log.Infof("Host is already stopped")
		return nil
	}
	err = d.client().StopServer(d.DatacenterId, d.ServerId)
	if err != nil {
		return err
	}
	return nil
}

// Restart reboots the machine instance.
func (d *Driver) Restart() error {
	err := d.client().RestartServer(d.DatacenterId, d.ServerId)
	if err != nil {
		return err
	}
	return nil
}

// Kill stops the machine instance
func (d *Driver) Kill() error {
	err := d.client().StopServer(d.DatacenterId, d.ServerId)
	if err != nil {
		return err
	}
	return nil
}

// GetSSHHostname returns an IP address or hostname for the machine instance.
func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetSSHUsername() string {
	return d.SSHUser
}

// GetURL returns a socket address to connect to Docker engine of the machine instance.
func (d *Driver) GetURL() (string, error) {
	if err := drivers.MustBeRunning(d); err != nil {
		return "", err
	}
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

// GetIP returns public IP address or hostname of the machine instance.
func (d *Driver) GetIP() (string, error) {
	server, err := d.client().GetServer(d.DatacenterId, d.ServerId)
	if err != nil {
		return "", fmt.Errorf("error getting server by id: %v", err)
	}

	if serverEntities, ok := server.GetEntitiesOk(); ok && serverEntities != nil {
		if serverEntitiesNic, ok := serverEntities.GetNicsOk(); ok && serverEntitiesNic != nil {
			if serverEntitiesNicItems, ok := serverEntitiesNic.GetItemsOk(); ok && serverEntitiesNicItems != nil {
				entitiesNicItems := *serverEntitiesNicItems
				entityNic := entitiesNicItems[0]
				if nicProp, ok := entityNic.GetPropertiesOk(); ok && nicProp != nil {
					if nicIps, ok := nicProp.GetIpsOk(); ok && nicIps != nil {
						entityNicIps := *nicIps
						d.IPAddress = entityNicIps[0]
					}
				}
			}
		}
	}
	if d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return d.IPAddress, nil
}

// GetState returns the state of the machine role instance.
func (d *Driver) GetState() (state.State, error) {
	server, err := d.client().GetServer(d.DatacenterId, d.ServerId)
	if err != nil {
		return state.None, err
	}

	if metadata, ok := server.GetMetadataOk(); ok && metadata != nil {
		if metadataState, ok := metadata.GetStateOk(); ok && metadataState != nil {
			switch *metadataState {
			case "NOSTATE":
				return state.None, nil
			case "AVAILABLE":
				return state.Running, nil
			case "PAUSED":
				return state.Paused, nil
			case "BLOCKED":
				return state.Stopped, nil
			case "SHUTDOWN":
				return state.Stopped, nil
			case "SHUTOFF":
				return state.Stopped, nil
			case "CRASHED":
				return state.Error, nil
			case "INACTIVE":
				return state.Stopped, nil
			}
		}
	}
	return state.None, fmt.Errorf("error getting server information")
}

/*
	Private helper functions
*/

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}

func (d *Driver) createSSHKey() (string, error) {
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return "", err
	}
	publicKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return "", err
	}
	return string(publicKey), nil
}

func (d *Driver) isSwarmMaster() bool {
	return d.SwarmMaster
}

func (d *Driver) getImageId(imageName string) (string, error) {
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

func getDriverVersion(v string) string {
	if v == "" {
		return driverVersionDev
	} else {
		return v
	}
}
