package ionoscloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/ionos-cloud/docker-machine-driver/internal/pointer"
	"github.com/ionos-cloud/docker-machine-driver/pkg/extflag"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	"github.com/hashicorp/go-multierror"
	"github.com/ionos-cloud/docker-machine-driver/internal/utils"
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
	flagServerType             = "ionoscloud-server-type"
	flagTemplate               = "ionoscloud-template"
	flagImage                  = "ionoscloud-image"
	flagImagePassword          = "ionoscloud-image-password"
	flagLocation               = "ionoscloud-location"
	flagDatacenterId           = "ionoscloud-datacenter-id"
	flagDatacenterName         = "ionoscloud-datacenter-name"
	flagLanId                  = "ionoscloud-lan-id"
	flagNicDhcp                = "ionoscloud-nic-dhcp"
	flagNicIps                 = "ionoscloud-nic-ips"
	flagLanName                = "ionoscloud-lan-name"
	flagVolumeAvailabilityZone = "ionoscloud-volume-availability-zone"
	flagCloudInit              = "ionoscloud-cloud-init"
	flagSSHInCloudInit         = "ionoscloud-ssh-in-cloud-init"
	flagSSHUser                = "ionoscloud-ssh-user"
	flagCloudInitB64           = "ionoscloud-cloud-init-b64"
	flagWaitForIpChange        = "ionoscloud-wait-for-ip-change"
	flagWaitForIpChangeTimeout = "ionoscloud-wait-for-ip-change-timeout"
	flagNatId                  = "ionoscloud-nat-id"
	flagNatName                = "ionoscloud-nat-name"
	flagNatPublicIps           = "ionoscloud-nat-public-ips"
	flagNatFlowlogs            = "ionoscloud-nat-flowlogs"
	flagNatRules               = "ionoscloud-nat-rules"
	flagSkipDefaultNatRules    = "ionoscloud-skip-default-nat-rules"
	flagNatLansToGateways      = "ionoscloud-nat-lans-to-gateways"
	flagPrivateLan             = "ionoscloud-private-lan"
	flagAdditionalLans         = "ionoscloud-additional-lans"
	flagCreateNat              = "ionoscloud-create-nat"
	// ---
)

const (
	defaultRegion                 = "us/las"
	defaultImageAlias             = "ubuntu:20.04"
	defaultImagePassword          = "" // Must contain both letters and numbers, at least 8 characters
	defaultCpuFamily              = "AMD_OPTERON"
	defaultAvailabilityZone       = "AUTO"
	defaultDiskType               = "HDD"
	defaultServerType             = "ENTERPRISE"
	defaultTemplate               = "CUBES XS"
	defaultSSHUser                = "root"
	defaultDatacenterName         = "docker-machine-data-center"
	defaultLanName                = "docker-machine-lan"
	defaultNatName                = "docker-machine-nat"
	defaultSize                   = 10
	defaultWaitForIpChangeTimeout = 600
	driverName                    = "ionoscloud"
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
	NicDhcp                bool
	NicIps                 []string
	Location               string
	CpuFamily              string
	ServerType             string
	Template               string
	DCExists               bool
	LanExists              bool
	NatExists              bool
	UseAlias               bool
	VolumeAvailabilityZone string
	ServerAvailabilityZone string
	LanId                  string
	LanName                string
	AdditionalLans         []string
	AdditionalLansIds      []int
	AdditionalNicsIds      []string
	DatacenterId           string
	DatacenterName         string
	VolumeId               string
	NicId                  string
	ServerId               string
	IpBlockId              string
	CreateNat              bool
	NatName                string
	NatId                  string
	CloudInit              string
	CloudInitB64           string
	NatPublicIps           []string
	NatFlowlogs            []string
	NatRules               []string
	SkipDefaultNatRules    bool
	NatLansToGateways      map[string][]string
	PrivateLan             bool
	SSHInCloudInit         bool
	WaitForIpChange        bool
	WaitForIpChangeTimeout int

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
			Name:   flagNatName,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNatName),
			Value:  defaultNatName,
			Usage:  "Ionos Cloud NAT Gateway name. Note that setting this will NOT implicitly create a NAT, this flag will only be read if need be",
		},
		mcnflag.StringFlag{
			Name:   flagNatId,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNatId),
			// Value:  nil,
			Usage: "Ionos Cloud existing and configured NAT Gateway",
		},
		mcnflag.BoolFlag{
			Name:   flagCreateNat,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagCreateNat),
			Usage:  "If set, will create a default NAT. Requires private LAN",
		},
		mcnflag.StringSliceFlag{
			Name:   flagNatPublicIps,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNatPublicIps),
			Usage:  "Ionos Cloud NAT Gateway public IPs",
		},
		mcnflag.StringSliceFlag{
			Name:   flagNatFlowlogs,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNatFlowlogs),
			Usage:  "Ionos Cloud NAT Gateway Flowlogs",
		},
		mcnflag.StringSliceFlag{
			Name:   flagNatRules,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNatRules),
			Usage:  "Ionos Cloud NAT Gateway Rules",
		},
		mcnflag.BoolFlag{
			Name:   flagSkipDefaultNatRules,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagSkipDefaultNatRules),
			Usage:  "Should the driver skip creating default nat rules if creating a NAT, creating only the specified rules",
		},
		mcnflag.StringFlag{
			// A string, like "1=10.0.0.1,10.0.0.2:2=10.0.0.10" . Lans MUST be separated by `:`. IPs MUST be separated by `,`
			Name:   flagNatLansToGateways,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNatLansToGateways),
			Usage:  "Ionos Cloud NAT map of LANs to a slice of their Gateway IPs. Example: \"1=10.0.0.1,10.0.0.2:2=10.0.0.10\"",
		},
		mcnflag.BoolFlag{
			Name:   flagPrivateLan,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagPrivateLan),
			Usage:  "Should the created LAN be private? Does nothing if LAN ID is provided",
		},
		mcnflag.StringSliceFlag{
			Name:   flagAdditionalLans,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagAdditionalLans),
			Usage:  "Names of existing IONOS Lans to connect the machine to. Names that are not found are ignored",
		},
		mcnflag.BoolFlag{
			Name:   flagWaitForIpChange,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagWaitForIpChange),
			Usage:  "Should the driver wait for the NIC IP to be set by external sources?",
		},
		mcnflag.IntFlag{
			Name:   flagWaitForIpChangeTimeout,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagWaitForIpChangeTimeout),
			Value:  defaultWaitForIpChangeTimeout,
			Usage:  "Timeout used when waiting for NIC IP changes",
		},
		mcnflag.BoolFlag{
			Name:   flagNicDhcp,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNicDhcp),
			Usage:  "Should the created NIC have DHCP set to true or false?",
		},
		mcnflag.StringSliceFlag{
			Name:   flagNicIps,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagNicIps),
			Usage:  "Ionos Cloud NIC IPs",
		},
		mcnflag.StringFlag{
			Name:   flagEndpoint,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagEndpoint),
			Value:  sdkgo.DefaultIonosServerUrl,
			Usage:  "Ionos Cloud API Endpoint",
		},
		mcnflag.StringFlag{
			Name:   flagUsername,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagUsername),
			Usage:  "Ionos Cloud Username",
		},
		mcnflag.StringFlag{
			Name:   flagPassword,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagPassword),
			Usage:  "Ionos Cloud Password",
		},
		mcnflag.StringFlag{
			Name:   flagToken,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagToken),
			Usage:  "Ionos Cloud Token",
		},
		mcnflag.IntFlag{
			Name:   flagServerCores,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagServerCores),
			Value:  2,
			Usage:  "Ionos Cloud Server Cores (2, 3, 4, 5, 6, etc.)",
		},
		mcnflag.IntFlag{
			Name:   flagServerRam,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagServerRam),
			Value:  2048,
			Usage:  "Ionos Cloud Server Ram in MB(1024, 2048, 3072, 4096, etc.)",
		},
		mcnflag.IntFlag{
			Name:   flagDiskSize,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagDiskSize),
			Value:  50,
			Usage:  "Ionos Cloud Volume Disk-Size in GB(10, 50, 100, 200, 400)",
		},
		mcnflag.StringFlag{
			Name:   flagImage,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagImage),
			Value:  defaultImageAlias,
			Usage:  "Ionos Cloud Image Id or Alias (ubuntu:latest, ubuntu:20.04)",
		},
		mcnflag.StringFlag{
			Name:   flagImagePassword,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagImagePassword),
			Value:  defaultImagePassword,
			Usage:  "Ionos Cloud Image Password to be able to access the server from DCD platform",
		},
		mcnflag.StringFlag{
			Name:   flagLocation,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagLocation),
			Value:  defaultRegion,
			Usage:  "Ionos Cloud Location",
		},
		mcnflag.StringFlag{
			Name:   flagDiskType,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagDiskType),
			Value:  defaultDiskType,
			Usage:  "Ionos Cloud Volume Disk-Type (HDD, SSD)",
		},
		mcnflag.StringFlag{
			Name:   flagServerType,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagServerType),
			Value:  defaultServerType,
			Usage:  "Ionos Cloud Server Type(ENTERPRISE or CUBE). CUBE servers are only available in certain locations.",
		},
		mcnflag.StringFlag{
			Name:   flagTemplate,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagTemplate),
			Value:  defaultTemplate,
			Usage:  "Ionos Cloud CUBE Template, only used for CUBE servers.",
		},
		mcnflag.StringFlag{
			Name:   flagServerCpuFamily,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagServerCpuFamily),
			Value:  defaultCpuFamily,
			Usage:  "Ionos Cloud Server CPU families (AMD_OPTERON, INTEL_XEON, INTEL_SKYLAKE, INTEL_ICELAKE, AMD_EPYC)",
		},
		mcnflag.StringFlag{
			Name:   flagDatacenterId,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagDatacenterId),
			Usage:  "Ionos Cloud Virtual Data Center Id",
		},
		mcnflag.StringFlag{
			Name:   flagDatacenterName,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagDatacenterName),
			Value:  defaultDatacenterName,
			Usage:  "Ionos Cloud Virtual Data Center Name",
		},
		mcnflag.StringFlag{
			Name:   flagLanId,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagLanId),
			Usage:  "Ionos Cloud LAN Id",
		},
		mcnflag.StringFlag{
			EnvVar: "IONOSCLOUD_LAN_Name",
			Name:   flagLanName,
			Value:  defaultLanName,
			Usage:  "Ionos Cloud LAN Name",
		},
		mcnflag.StringFlag{
			Name:   flagVolumeAvailabilityZone,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagVolumeAvailabilityZone),
			Value:  defaultAvailabilityZone,
			Usage:  "Ionos Cloud Volume Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)",
		},
		mcnflag.StringFlag{
			Name:   flagServerAvailabilityZone,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagServerAvailabilityZone),
			Value:  defaultAvailabilityZone,
			Usage:  "Ionos Cloud Server Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)",
		},
		mcnflag.StringFlag{
			Name:   flagCloudInit,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagCloudInit),
			Usage:  "The cloud-init configuration for the volume as a multi-line string",
		},
		mcnflag.StringFlag{
			Name:   flagCloudInitB64,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagCloudInitB64),
			Usage:  "The cloud-init configuration for the volume as base64 encoded string",
		},
		mcnflag.StringFlag{
			Name:   flagSSHUser,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagSSHUser),
			Value:  defaultSSHUser,
			Usage:  "The name of the user the driver will use for ssh",
		},
		mcnflag.BoolFlag{
			Name:   flagSSHInCloudInit,
			EnvVar: extflag.KebabCaseToEnvVarCase(flagSSHInCloudInit),
			Usage:  "Should the driver only add the SSH info in the user data? (required for custom images)",
		},
	}
}

// SetConfigFromFlags initializes driver values from the command line values.
func (d *Driver) SetConfigFromFlags(opts drivers.DriverOptions) error {
	d.CreateNat = opts.Bool(flagCreateNat)
	d.NatName = opts.String(flagNatName)
	d.NatId = opts.String(flagNatId)
	d.NatPublicIps = opts.StringSlice(flagNatPublicIps)
	d.NatFlowlogs = opts.StringSlice(flagNatFlowlogs)
	d.NatRules = opts.StringSlice(flagNatRules)
	d.NatLansToGateways = extflag.ToMapOfStringToStringSlice(opts.String(flagNatLansToGateways))
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
	d.ServerType = opts.String(flagServerType)
	d.Template = opts.String(flagTemplate)
	d.CpuFamily = opts.String(flagServerCpuFamily)
	d.DatacenterId = opts.String(flagDatacenterId)
	d.DatacenterName = opts.String(flagDatacenterName)
	d.LanId = opts.String(flagLanId)
	d.LanName = opts.String(flagLanName)
	d.NicDhcp = opts.Bool(flagNicDhcp)
	d.WaitForIpChange = opts.Bool(flagWaitForIpChange)
	d.WaitForIpChangeTimeout = opts.Int(flagWaitForIpChangeTimeout)
	d.NicIps = opts.StringSlice(flagNicIps)
	d.VolumeAvailabilityZone = opts.String(flagVolumeAvailabilityZone)
	d.ServerAvailabilityZone = opts.String(flagServerAvailabilityZone)
	d.SkipDefaultNatRules = opts.Bool(flagSkipDefaultNatRules)
	d.CloudInit = opts.String(flagCloudInit)
	d.SSHUser = opts.String(flagSSHUser)
	d.SSHInCloudInit = opts.Bool(flagSSHInCloudInit)
	d.CloudInitB64 = opts.String(flagCloudInitB64)
	d.PrivateLan = opts.Bool(flagPrivateLan)
	d.AdditionalLans = opts.StringSlice(flagAdditionalLans)

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

	d.DCExists = false
	d.LanExists = false
	d.NatExists = false

	if strings.Contains(d.MachineName, "-pool") {
		if !strings.Contains(d.MachineName, "-pool1-") {
			time.Sleep(60 * time.Second)
		}
	} else {
		for i := len(d.MachineName) - 1; i >= 0; i-- {
			if !unicode.IsNumber(rune(d.MachineName[i])) {
				if d.MachineName[i+1:] != "1" {
					time.Sleep(60 * time.Second)
				}
				break
			}
		}
	}
	if d.DatacenterId == "" {
		datacenters, err := d.client().GetDatacenters()
		if err != nil {
			return err
		}

		foundDc := false
		for _, dc := range *datacenters.Items {
			if *dc.Properties.Name == d.DatacenterName {
				if foundDc {
					return fmt.Errorf("multiple Data Centers with name %v found", d.DatacenterName)
				}
				foundDc = true
				if dcId, ok := dc.GetIdOk(); ok && dcId != nil {
					d.DatacenterId = *dcId
				}
			}
		}
	}

	if d.DatacenterId != "" {
		d.DCExists = true

		if d.LanId == "" {
			lans, err := d.client().GetLans(d.DatacenterId)
			if err != nil {
				return err
			}

			foundLan := false
			for _, lan := range *lans.Items {
				if *lan.Properties.Name == d.LanName {
					if foundLan {
						return fmt.Errorf("multiple LANs with name %v found", d.LanName)
					}
					foundLan = true
					if lanId, ok := lan.GetIdOk(); ok && lanId != nil {
						d.LanId = *lanId
					}
				} else if utils.Contains(d.AdditionalLans, *lan.Properties.Name) {
					if lanId, ok := lan.GetIdOk(); ok && lanId != nil {
						lanIdInt, err := strconv.Atoi(*lanId)
						if err != nil {
							return fmt.Errorf("invalid LAN ID found: %v", *lanId)
						}
						d.AdditionalLansIds = append(d.AdditionalLansIds, lanIdInt)
					}
				}
			}

		}
		if d.LanId != "" {
			d.LanExists = true
			lan, err := d.client().GetLan(d.DatacenterId, d.LanId)
			if err != nil {
				return fmt.Errorf("error getting LAN: %w", err)
			}
			log.Info("Creating machine under LAN " + *lan.GetId())
			d.PrivateLan = !*lan.GetProperties().GetPublic()
		}
		dc, err := d.client().GetDatacenter(d.DatacenterId)
		if err != nil {
			return fmt.Errorf("error getting datacenter: %w", err)
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
	}
	if imageId, err := d.getImageId(d.Image); err != nil && imageId == "" {
		return fmt.Errorf("error getting image/alias %s: %w", d.Image, err)
	}

	if d.NatId != "" && d.DatacenterId != "" {
		nats, err := d.client().GetNats(d.DatacenterId)
		if err != nil {
			return err
		}

		foundNat := false
		for _, nat := range *nats.Items {
			if *nat.Properties.Name == d.NatName {
				if foundNat {
					return fmt.Errorf("multiple Nat Gateways with name %v found", d.NatName)
				}
				foundNat = true
				if id, ok := nat.GetIdOk(); ok && id != nil {
					d.NatId = *id
				}
			}
		}
	}

	if d.NatId != "" {
		d.NatExists = true
	}

	if d.NatId != "" && d.CreateNat {
		return fmt.Errorf("trying to create a NAT while also found an existing NAT. Please set only one of: (%s | %s), or try a different NAT name",
			flagNatId, flagCreateNat)
	}

	// d.PrivateLan is set above to false as a side effect if the LAN with the given ID is private. If concerns are separated in this func, be aware of this!
	if !d.PrivateLan && (d.NatId != "" || d.CreateNat) {
		return fmt.Errorf("using a NAT Gateway requires usage of a private LAN. Please enable %s or provide a Private Lan ID for %s", flagPrivateLan, flagLanId)
	}

	return nil
}

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

func getPropertyWithFallback[T comparable](p1 T, p2 T, empty T) T {
	if p1 == empty {
		return p2
	}
	return p1
}

// Create creates the machine.
func (d *Driver) Create() (err error) {
	log.Infof("Creating SSH key...")
	if d.SSHKey == "" {
		d.SSHKey, err = d.createSSHKey()
		if err != nil {
			return fmt.Errorf("error creating SSH keys: %w", err)
		}
		log.Debugf("SSH Key generated in file: %v", d.publicSSHKeyPath())
	}

	result, err := d.getImageId(d.Image)
	if err != nil {
		return fmt.Errorf("error getting image/alias %s: %w", d.Image, err)
	}
	var alias string
	if d.UseAlias {
		alias = result
	}

	// Creating Data Center if one was not provided

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

	// Creating LAN if one was not provided

	if d.LanId == "" {
		lan, err := d.client().CreateLan(d.DatacenterId, d.LanName, !d.PrivateLan)
		if err != nil {
			err = fmt.Errorf("error creating LAN: %w", err)
			// TODO : export below to a func --->
			log.Warn(rollingBackNotice)
			if removeErr := d.Remove(); removeErr != nil {
				return fmt.Errorf("failed to create machine due to error: %w\n Removing created resources: %v", err, removeErr)
			}
			return err
			// TODO: <---
		}
		if lanId, ok := lan.GetIdOk(); ok && lanId != nil {
			d.LanId = *lanId
			log.Debugf("Lan ID: %v", d.LanId)
		}
	}

	lan, err := d.client().GetLan(d.DatacenterId, d.LanId)
	if err != nil {
		return fmt.Errorf("error getting LAN: %w", err)
	}

	if err != nil {
		log.Warn(rollingBackNotice)
		if removeErr := d.Remove(); removeErr != nil {
			return fmt.Errorf("failed to create machine due to error: %w\n Removing created resources: %v", err, removeErr)
		}
		return err
	}

	var isLanPrivate bool
	if lanProp, ok := lan.GetPropertiesOk(); ok && lanProp != nil {
		if public, ok := lanProp.GetPublicOk(); ok && public != nil {
			isLanPrivate = !*public
		}
	}

	// Creating the server with the volume attached

	// User Data for cloud init
	givenB64CloudInit, _ := base64.StdEncoding.DecodeString(d.CloudInitB64)
	if ud := getPropertyWithFallback(string(givenB64CloudInit), d.CloudInit, ""); ud != "" {
		// Provided B64 User Data has priority over UI provided User Data
		d.CloudInit = ud
	}

	if d.SSHUser != "root" || d.SSHInCloudInit {
		d.CloudInit, err = d.addSSHUserToYaml()
		if err != nil {
			return err
		}
	}
	d.CloudInit, err = d.client().UpdateCloudInitFile(
		d.CloudInit, "hostname", []interface{}{d.MachineName}, true, "skip",
	)
	if err != nil {
		return err
	}
	ud := base64.StdEncoding.EncodeToString([]byte(d.CloudInit))

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
		log.Infof("Image Id: %v", result)
		volumeProperties.Image = &result
	} else {
		log.Infof("Image Alias: %v", alias)
		volumeProperties.ImageAlias = &alias
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

	// Add nics to server
	providedNicIps := len(d.NicIps) != 0
	reservedIps := &[]string{}

	// Reserve IP if needed
	if !isLanPrivate && !providedNicIps ||
		d.CreateNat && d.NatPublicIps == nil {
		ipBlock, err := d.client().CreateIpBlock(1, d.Location)
		if err != nil {
			return fmt.Errorf("error creating ipblock: %w", err)
		}
		if ipBlockId, ok := ipBlock.GetIdOk(); ok && ipBlockId != nil {
			d.IpBlockId = *ipBlockId
			log.Debugf("IpBlock ID: %v", d.IpBlockId)
		}
		reservedIps, err = d.client().GetIpBlockIps(ipBlock)
		if err != nil {
			return err
		}
	}

	var ipsForAttachedNic *[]string

	if providedNicIps {
		ipsForAttachedNic = &d.NicIps // If IPs are provided use those
	} else if isLanPrivate {
		ipsForAttachedNic = nil // Let CloudAPI generate an IP, which we can later use for the subnet
	} else {
		ipsForAttachedNic = reservedIps // For public NICs we use the generated IPs
	}

	attachedNics := sdkgo.NewNicsWithDefaults()

	lanId, _ := strconv.Atoi(d.LanId)
	lanId_int32 := int32(lanId)
	nicProperties := &sdkgo.NicProperties{
		Name: &d.MachineName,
		Lan:  &lanId_int32,
		Ips:  ipsForAttachedNic,
		Dhcp: &d.NicDhcp,
	}

	attachedNics.Items = &[]sdkgo.Nic{
		{
			Properties: nicProperties,
		},
	}

	for _, additionalLanId := range d.AdditionalLansIds {
		dhcp_true := true
		additionalLanId_int32 := int32(additionalLanId)
		additionalNicName := d.MachineName + " " + fmt.Sprint(additionalLanId)
		nicProperties = &sdkgo.NicProperties{
			Name: &additionalNicName,
			Lan:  &additionalLanId_int32,
			Ips:  nil,
			Dhcp: &dhcp_true,
		}
		additionalNic := sdkgo.Nic{
			Properties: nicProperties,
		}
		*attachedNics.Items = append(*attachedNics.Items, additionalNic)
	}

	serverToCreate.Entities.SetNics(*attachedNics)

	server, err := d.client().CreateServer(d.DatacenterId, serverToCreate)
	if err != nil {
		// TODO: Export to a func
		log.Warn(rollingBackNotice)
		if removeErr := d.Remove(); removeErr != nil {
			return fmt.Errorf("failed to create server due to error: %w\n Removing created resources: %v", err, removeErr)
		}
		return err
	}
	if serverId, ok := server.GetIdOk(); ok && serverId != nil {
		d.ServerId = *serverId
		log.Debugf("Server ID: %v", d.ServerId)
	}

	server, err = d.client().GetServer(d.DatacenterId, d.ServerId, 2)
	if err != nil {
		return fmt.Errorf("error getting server by id: %w", err)
	}
	volumes, ok := server.Entities.GetVolumesOk()
	if !ok {
		return fmt.Errorf("error getting server: d.ServerId is empty")
	}
	d.VolumeId = *(*volumes.Items)[0].GetId()
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

	nic, err := d.client().GetNic(d.DatacenterId, d.ServerId, d.NicId)
	if err != nil {
		return fmt.Errorf("error getting NIC: %w", err)
	}

	nicIps := &[]string{}
	if nicProp, ok := nic.GetPropertiesOk(); ok && nicProp != nil {
		nicIps = nicProp.GetIps()
	}
	if len(*nicIps) > 0 {
		d.IPAddress = (*nicIps)[0]
		log.Infof(d.IPAddress)
	}

	// --- NAT ---
	if d.CreateNat {
		// TODO: Were CreateNat in a deeper scope, we wouldn't have the need of these variables (they are here to avoid function-wide side-effects)
		natPublicIps := reservedIps
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

	if !d.NatExists && d.DatacenterId != "" && d.NatId != "" {
		log.Debugf("Starting deleting NAT with Id: %v", d.NatId)
		err := d.client().RemoveNat(d.DatacenterId, d.NatId)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("error deleting NAT: %w", err))
		} else {
			d.NatId = ""
		}
	}

	var err error
	if d.DatacenterId != "" && d.ServerId != "" && d.NicId != "" {
		log.Debugf("Starting deleting Nic with Id: %v", d.NicId)
		err := d.client().RemoveNic(d.DatacenterId, d.ServerId, d.NicId)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("error deleting NIC: %w", err))
		} else {
			d.NicId = ""
		}
	}
	if d.DatacenterId != "" && d.ServerId != "" {
		for _, additionalNicId := range d.AdditionalNicsIds {
			log.Debugf("Starting deleting additional Nic with Id: %v", additionalNicId)
			err := d.client().RemoveNic(d.DatacenterId, d.ServerId, additionalNicId)
			if err != nil {
				result = multierror.Append(result, fmt.Errorf("error deleting additional NIC: %w", err))
			}
		}
	}
	if d.DatacenterId != "" && d.VolumeId != "" {
		log.Debugf("Starting deleting Volume with Id: %v", d.VolumeId)
		err = d.client().RemoveVolume(d.DatacenterId, d.VolumeId)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("error removing volume: %w", err))
		} else {
			d.VolumeId = ""
		}
	}
	if d.DatacenterId != "" && d.ServerId != "" {
		log.Debugf("Starting deleting Server with Id: %v", d.ServerId)
		err = d.client().RemoveServer(d.DatacenterId, d.ServerId)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("error deleting server: %w", err))
		} else {
			d.ServerId = ""
		}
	}
	// If the LAN existed before creating the machine, do not delete it at clean-up
	if !d.LanExists {
		if d.DatacenterId != "" && d.LanId != "" {
			log.Debugf("Starting deleting LAN with Id: %v", d.LanId)
			err = d.client().RemoveLan(d.DatacenterId, d.LanId)
			if err != nil {
				result = multierror.Append(result, fmt.Errorf("error deleting LAN: %w", err))
			} else {
				d.LanId = ""
			}
		}
	}
	// If the DataCenter existed before creating the machine, do not delete it at clean-up
	if !d.DCExists {
		if d.DatacenterId != "" {
			log.Debugf("Starting deleting Datacenter with Id: %v", d.DatacenterId)
			err = d.client().RemoveDatacenter(d.DatacenterId)
			if err != nil {
				result = multierror.Append(result, fmt.Errorf("error deleting datacenter: %w", err))
			} else {
				d.DatacenterId = ""
			}
		}
	}

	if d.IpBlockId != "" {
		log.Debugf("Starting deleting IpBlock with Id: %v", d.IpBlockId)
		err = d.client().RemoveIpBlock(d.IpBlockId)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("error deleting ipblock: %w", err))
		} else {
			d.IpBlockId = ""
		}
	}

	return result.ErrorOrNil()
}

// Start issues a power on for the machine instance.
func (d *Driver) Start() error {
	serverState, err := d.GetState()
	if err != nil {
		return fmt.Errorf("error getting state: %w", err)
	}
	if serverState != state.Running {
		if d.ServerType == "ENTERPRISE" {
			err = d.client().StartServer(d.DatacenterId, d.ServerId)
		} else if d.ServerType == "CUBE" {
			err = d.client().ResumeServer(d.DatacenterId, d.ServerId)
		} else {
			err = fmt.Errorf("wrong server type: %s", d.ServerType)
		}

		if err != nil {
			return fmt.Errorf("error starting server: %w", err)
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
		return fmt.Errorf("error getting state: %w", err)
	}
	if vmState == state.Stopped {
		log.Infof("Host is already stopped")
		return nil
	}
	if d.ServerType == "ENTERPRISE" {
		err = d.client().StopServer(d.DatacenterId, d.ServerId)
	} else if d.ServerType == "CUBE" {
		err = d.client().SuspendServer(d.DatacenterId, d.ServerId)
	} else {
		err = fmt.Errorf("wrong server type: %s", d.ServerType)
	}
	if err != nil {
		return fmt.Errorf("error stoping server: %w", err)
	}
	return nil
}

// Restart reboots the machine instance.
func (d *Driver) Restart() error {
	err := d.client().RestartServer(d.DatacenterId, d.ServerId)
	if err != nil {
		return fmt.Errorf("error restarting server: %w", err)
	}
	return nil
}

// Kill stops the machine instance
func (d *Driver) Kill() error {
	err := d.client().StopServer(d.DatacenterId, d.ServerId)
	if err != nil {
		return fmt.Errorf("error stopping server: %w", err)
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
	return fmt.Sprintf("tcp://%s:2376", ip), nil // TODO: Perhaps we can allow customization of the Docker Port: https://github.com/rancher/machine/blob/master/drivers/azure/azure.go#L619
}

// GetIP returns public IP address or hostname of the machine instance.
func (d *Driver) GetIP() (string, error) {
	if d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return d.IPAddress, nil
}

// GetState returns the state of the machine role instance.
func (d *Driver) GetState() (state.State, error) {
	if d.ServerId == "" {
		return state.None, fmt.Errorf("error getting server: d.ServerID is empty")
	}
	server, err := d.client().GetServer(d.DatacenterId, d.ServerId, 1)
	if err != nil {
		return state.None, fmt.Errorf("error getting server: %w", err)
	}

	if serverProperties, ok := server.GetPropertiesOk(); ok && serverProperties != nil {
		if vmState, ok := serverProperties.GetVmStateOk(); ok && vmState != nil {
			switch *vmState {
			case "NOSTATE":
				return state.None, nil
			case "RUNNING":
				return state.Running, nil
			case "PAUSED":
				return state.Paused, nil
			case "SUSPENDED":
				return state.Stopped, nil
			case "BLOCKED":
				return state.Stopped, nil
			case "SHUTDOWN":
				return state.Stopped, nil
			case "SHUTOFF":
				return state.Stopped, nil
			case "CRASHED":
				return state.Error, nil
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
