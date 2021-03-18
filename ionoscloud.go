package ionoscloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/hashicorp/go-multierror"
	"github.com/ionos-cloud/docker-machine-driver/utils"
	sdkgo "github.com/ionos-cloud/sdk-go/v5"
)

const (
	flagEndpoint               = "ionoscloud-endpoint"
	flagUsername               = "ionoscloud-username"
	flagPassword               = "ionoscloud-password"
	flagServerCores            = "ionoscloud-cores"
	flagServerRam              = "ionoscloud-ram"
	flagServerCpuFamily        = "ionoscloud-cpu-family"
	flagServerAvailabilityZone = "ionoscloud-server-availability-zone"
	flagDiskSize               = "ionoscloud-disk-size"
	flagDiskType               = "ionoscloud-disk-type"
	flagImage                  = "ionoscloud-image"
	flagLocation               = "ionoscloud-location"
	flagDatacenterId           = "ionoscloud-datacenter-id"
	flagVolumeAvailabilityZone = "ionoscloud-volume-availability-zone"
)

const (
	defaultRegion           = "us/las"
	defaultApiEndpoint      = "https://api.ionos.com/cloudapi/v5"
	defaultImageAlias       = "ubuntu:latest"
	defaultCpuFamily        = "AMD_OPTERON"
	defaultAvailabilityZone = "AUTO"
	defaultDiskType         = "HDD"
	defaultSize             = 10
	driverName              = "ionoscloud"
)

type Driver struct {
	*drivers.BaseDriver
	client func() utils.ClientService

	URL      string
	Username string
	Password string

	Ram                    int
	Cores                  int
	SSHKey                 string
	DiskSize               int
	DiskType               string
	Image                  string
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
}

// NewDriver returns a new driver instance.
func NewDriver(hostName, storePath string) drivers.Driver {
	return NewDerivedDriver(hostName, storePath)
}

func NewDerivedDriver(hostName, storePath string) *Driver {
	driver := &Driver{
		Size:     defaultSize,
		Location: defaultRegion,
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
	driver.client = func() utils.ClientService {
		return utils.New(context.TODO(), driver.Username, driver.Password, driver.URL)
	}
	return driver
}

// GetCreateFlags returns list of create flags driver accepts.
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_ENDPOINT",
			Name:   flagEndpoint,
			Value:  defaultApiEndpoint,
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
			Usage:  "Ionos Cloud Server Ram (1024, 2048, 3072, 4096, etc.)",
		},
		mcnflag.IntFlag{
			EnvVar: "IONOSCLOUD_DISK_SIZE",
			Name:   flagDiskSize,
			Value:  50,
			Usage:  "Ionos Cloud Volume Disk-Size (10, 50, 100, 200, 400)",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_IMAGE",
			Name:   flagImage,
			Value:  defaultImageAlias,
			Usage:  "Ionos Cloud Image Alias",
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
	}
}

// SetConfigFromFlags initializes driver values from the command line values.
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.URL = opts.String(flagEndpoint)
	d.Username = opts.String(flagUsername)
	d.Password = opts.String(flagPassword)
	d.DiskSize = opts.Int(flagDiskSize)
	d.Image = opts.String(flagImage)
	d.Cores = opts.Int(flagServerCores)
	d.Ram = opts.Int(flagServerRam)
	d.Location = opts.String(flagLocation)
	d.DiskType = opts.String(flagDiskType)
	d.CpuFamily = opts.String(flagServerCpuFamily)
	d.DatacenterId = opts.String(flagDatacenterId)
	d.VolumeAvailabilityZone = opts.String(flagVolumeAvailabilityZone)
	d.ServerAvailabilityZone = opts.String(flagServerAvailabilityZone)

	d.SwarmMaster = opts.Bool("swarm-master")
	d.SwarmHost = opts.String("swarm-host")
	d.SwarmDiscovery = opts.String("swarm-discovery")
	d.SetSwarmConfigFromFlags(opts)

	if d.URL == "" {
		d.URL = defaultApiEndpoint
	}

	return nil
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return driverName
}

// PreCreateCheck validates if driver values are valid to create the machine.
func (d *Driver) PreCreateCheck() error {
	if d.Username == "" {
		return fmt.Errorf("please provide username as parameter --ionoscloud-username or as environment variable $IONOSCLOUD_USERNAME")
	}
	if d.Password == "" {
		return fmt.Errorf("please provide password as parameter --ionoscloud-password or as environment variable $IONOSCLOUD_PASSWORD")
	}
	if d.DatacenterId != "" {
		dc, err := d.client().GetDatacenter(d.DatacenterId)
		if err != nil {
			return err
		}
		if dcprop, ok := dc.GetPropertiesOk(); ok && dcprop != nil {
			if name, ok := dcprop.GetNameOk(); ok && name != nil {
				log.Info("Creating machine under " + *name + " datacenter")
			}
		}
	}
	if imageId, err := d.getImageId(d.Image); err != nil && imageId == "" {
		return fmt.Errorf("error getting image/alias %s: %v", d.Image, err)
	}

	return nil
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
	}

	result, err := d.getImageId(d.Image)
	if err != nil {
		return fmt.Errorf("error getting image/alias %s: %v", d.Image, err)
	}
	var alias string
	if d.UseAlias {
		alias = result
	}

	ipBlock, err := d.client().CreateIpBlock(int32(1), d.Location)
	if err != nil {
		return err
	}
	var dc *sdkgo.Datacenter
	if d.DatacenterId == "" {
		d.DCExists = false
		var err error
		dc, err = d.client().CreateDatacenter(d.MachineName, d.Location)
		if err != nil {
			return err
		}
	} else {
		d.DCExists = true
		dc, err = d.client().GetDatacenter(d.DatacenterId)
		if err != nil {
			return err
		}
	}
	if dcId, ok := dc.GetIdOk(); ok && dcId != nil {
		d.DatacenterId = *dcId
	}

	lan, err := d.client().CreateLan(d.DatacenterId, d.MachineName, true)
	if err != nil {
		return err
	}
	if lanId, ok := lan.GetIdOk(); ok && lanId != nil {
		d.LanId = *lanId
	}

	server, err := d.client().CreateServer(d.DatacenterId, d.Location, d.MachineName, d.CpuFamily, d.ServerAvailabilityZone, int32(d.Ram), int32(d.Cores))
	if err != nil {
		return err
	}
	if serverId, ok := server.GetIdOk(); ok && serverId != nil {
		d.ServerId = *serverId
	}

	volume, err := d.client().CreateAttachVolume(d.DatacenterId, d.ServerId, d.DiskType, d.MachineName, alias, d.VolumeAvailabilityZone, d.SSHKey, float32(d.DiskSize))
	if err != nil {
		return err
	}
	if volumeId, ok := volume.GetIdOk(); ok && volumeId != nil {
		d.VolumeId = *volumeId
	}

	l, _ := strconv.Atoi(d.LanId)
	ips, err := d.client().GetIpBlock(ipBlock)
	if err != nil {
		return err
	}

	nic, err := d.client().CreateAttachNIC(d.DatacenterId, d.ServerId, d.MachineName, true, int32(l), ips)
	if err != nil {
		return err
	}
	if nicId, ok := nic.GetIdOk(); ok && nicId != nil {
		d.NicId = *nic.Id
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

	log.Info("Starting deleting resources...")

	err := d.client().RemoveNic(d.DatacenterId, d.ServerId, d.NicId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	err = d.client().RemoveVolume(d.DatacenterId, d.VolumeId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	err = d.client().RemoveServer(d.DatacenterId, d.ServerId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	err = d.client().RemoveLan(d.DatacenterId, d.LanId)
	if err != nil {
		result = multierror.Append(result, err)
	}
	// If the DataCenter existed before creating the machine, do not delete it at clean-up
	if !d.DCExists {
		err = d.client().RemoveDatacenter(d.DatacenterId)
		if err != nil {
			result = multierror.Append(result, err)
		}
	}
	err = d.client().RemoveIpBlock(d.IPAddress)
	if err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
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
	// first look if the provided parameter matches an alias, if a match is found we return the image alias
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
	// if no alias matches we do extended search and return the image id
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
