package ionoscloud

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	sdkgo "github.com/ionos-cloud/sdk-go/v5"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/log"
	"github.com/rancher/machine/libmachine/mcnflag"
	"github.com/rancher/machine/libmachine/mcnutils"
	"github.com/rancher/machine/libmachine/ssh"
	"github.com/rancher/machine/libmachine/state"
)

const (
	flagEndpoint               = "ionoscloud-endpoint"
	flagUsername               = "ionoscloud-username"
	flagPassword               = "ionoscloud-password"
	flagServerCores            = "ionoscloud-server-cores"
	flagServerRam              = "ionoscloud-server-ram"
	flagServerCpuFamily        = "ionoscloud-server-cpu-family"
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

	driverName = "ionoscloud"
	waitCount  = 1000
)

type Driver struct {
	*drivers.BaseDriver
	URL                    string
	Username               string
	Password               string
	ServerId               string
	Ram                    int
	Cores                  int
	SSHKey                 string
	DatacenterId           string
	VolumeId               string
	NicId                  string
	VolumeAvailabilityZone string
	ServerAvailabilityZone string
	DiskSize               int
	DiskType               string
	Image                  string
	Size                   int
	Location               string
	CpuFamily              string
	DCExists               bool
	UseAlias               bool
	LanId                  string
	Config                 *sdkgo.APIClient
}

// NewDriver returns a new driver instance.
func NewDriver(hostName, storePath string) drivers.Driver {
	return &Driver{
		Size:     defaultSize,
		Location: defaultRegion,
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
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
			EnvVar: "IONOSCLOUD_SERVER_CORES",
			Name:   flagServerCores,
			Value:  4,
			Usage:  "Ionos Cloud Server Cores (2, 3, 4, 5, 6, etc.)",
		},
		mcnflag.IntFlag{
			EnvVar: "IONOSCLOUD_SERVER_RAM",
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
			EnvVar: "IONOSCLOUD_SERVER_CPU_FAMILY",
			Name:   flagServerCpuFamily,
			Value:  defaultCpuFamily,
			Usage:  "Ionos Cloud Server CPU families (AMD_OPTERON,INTEL_XEON)",
		},
		mcnflag.StringFlag{
			Name:  flagDatacenterId,
			Usage: "Ionos Cloud Virtual Data Center Id",
		},
		mcnflag.StringFlag{
			Name:  flagVolumeAvailabilityZone,
			Value: defaultAvailabilityZone,
			Usage: "Ionos Cloud Volume Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)",
		},
		mcnflag.StringFlag{
			Name:  flagServerAvailabilityZone,
			Value: defaultAvailabilityZone,
			Usage: "Ionos Cloud Server Availability Zone (AUTO, ZONE_1, ZONE_2, ZONE_3)",
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
		err := d.getClient()
		if err != nil {
			return err
		}
		dc, resp, err := d.Config.DataCenterApi.DatacentersFindById(context.TODO(), d.DatacenterId).Execute()
		if err != nil {
			return err
		}

		if resp.StatusCode == 404 {
			return fmt.Errorf("DataCenter UUID %s does not exist", d.DatacenterId)
		} else {
			log.Info("Creating machine under " + *dc.Properties.Name + " datacenter")
		}
	}

	if imageId, err := d.getImageId(d.Image); err != nil && imageId == "" {
		return fmt.Errorf("The image/alias  %s does not exist.", d.Image)
	}

	return nil
}

// Create creates the machine.
func (d *Driver) Create() error {
	err := d.getClient()
	if err != nil {
		return err
	}

	log.Infof("Creating SSH key...")
	if d.SSHKey == "" {
		d.SSHKey, err = d.createSSHKey()
		if err != nil {
			return fmt.Errorf("error creating SSH keys: %v", err)
		}
	}
	result, err := d.getImageId(d.Image)
	if err != nil {
		return fmt.Errorf("error getting image: %v", err)
	}

	var alias string
	if d.UseAlias {
		alias = result
	}

	ipBlockSize := int32(1)
	ipBlock, ipBlockResp, err := d.Config.IPBlocksApi.IpblocksPost(context.TODO()).Ipblock(sdkgo.IpBlock{
		Properties: &sdkgo.IpBlockProperties{
			Location: &d.Location,
			Size:     &ipBlockSize,
		}}).Execute()
	if err != nil {
		return fmt.Errorf("error creating ipblock: %v", err)
	}

	if ipBlockResp.StatusCode > 299 {
		return fmt.Errorf("error reserving an ipblock: %s", ipBlockResp.Response.Status)
	}

	err = d.waitTillProvisioned(ipBlockResp.Header.Get("location"))
	if err != nil {
		return fmt.Errorf("error waiting untill provisioned: %v", err)
	}

	var dc sdkgo.Datacenter
	if d.DatacenterId == "" {
		d.DCExists = false
		var err error
		var dcResp *sdkgo.APIResponse

		dc, dcResp, err = d.Config.DataCenterApi.DatacentersPost(context.TODO()).Datacenter(sdkgo.Datacenter{
			Properties: &sdkgo.DatacenterProperties{
				Name:     &d.MachineName,
				Location: &d.Location,
			}}).Execute()
		if err != nil {
			return fmt.Errorf("error creating datacenter: %v", err)
		}
		if dcResp.StatusCode == 202 {
			log.Info("DataCenter Created")
		} else {
			return fmt.Errorf("error creating DC: %s", dcResp.Response.Status)
		}
		err = d.waitTillProvisioned(dcResp.Header.Get("location"))
		if err != nil {
			return fmt.Errorf("error waiting untill provisioned: %v", err)
		}
	} else {
		d.DCExists = true
		dc, _, err = d.Config.DataCenterApi.DatacentersFindById(context.TODO(), d.DatacenterId).Execute()
		if err != nil {
			return fmt.Errorf("error getting datacenter: %v", err)
		}
	}
	d.DatacenterId = *dc.Id

	lanPublic := true
	lan, lanResp, err := d.Config.LanApi.DatacentersLansPost(context.TODO(), *dc.Id).Lan(sdkgo.LanPost{
		Properties: &sdkgo.LanPropertiesPost{
			Name:   &d.MachineName,
			Public: &lanPublic,
		}}).Execute()
	if err != nil {
		return fmt.Errorf("error creating LAN: %v", err)
	}
	if lanResp.StatusCode == 202 {
		log.Info("LAN Created")
	} else {
		err := d.Remove()
		if err != nil {
			return fmt.Errorf("error deleting resources after unsuccesfull creation of LAN: %v", err)
		}
		return fmt.Errorf("error creating a LAN: %s", lanResp.Response.Status)
	}
	err = d.waitTillProvisioned(lanResp.Header.Get("location"))
	if err != nil {
		return fmt.Errorf("error waiting untill provisioned: %v", err)
	}
	d.LanId = *lan.Id

	svrRam := int32(d.Ram)
	svrCores := int32(d.Cores)
	diskSize := float32(d.DiskSize)

	server := sdkgo.Server{
		Properties: &sdkgo.ServerProperties{
			Name:             &d.MachineName,
			Ram:              &svrRam,
			Cores:            &svrCores,
			CpuFamily:        &d.CpuFamily,
			AvailabilityZone: &d.ServerAvailabilityZone,
		},
	}

	svr, serverResp, err := d.Config.ServerApi.DatacentersServersPost(context.TODO(), d.DatacenterId).Server(server).Execute()
	if err != nil {
		return fmt.Errorf("error creating server: %v", err)
	}
	if serverResp.StatusCode == 202 {
		log.Info("Server Created")
	} else {
		err := d.Remove()
		if err != nil {
			return fmt.Errorf("error deleting resources after unsuccesfull creation of server: %v", err)
		}
		return fmt.Errorf("error creating a server: %s", serverResp.Status)
	}

	err = d.waitTillProvisioned(serverResp.Header.Get("location"))
	if err != nil {
		return fmt.Errorf("error waiting untill provisioned: %v", err)
	}
	d.ServerId = *svr.Id

	volume := sdkgo.Volume{
		Properties: &sdkgo.VolumeProperties{
			Type:             &d.DiskType,
			Size:             &diskSize,
			Name:             &d.MachineName,
			ImageAlias:       &alias,
			SshKeys:          &[]string{d.SSHKey},
			AvailabilityZone: &d.VolumeAvailabilityZone,
		},
	}
	volume, volumeResp, err := d.Config.ServerApi.DatacentersServersVolumesPost(context.TODO(), d.DatacenterId, d.ServerId).Volume(volume).Execute()
	if err != nil {
		return fmt.Errorf("error attaching volume to server: %v", err)
	}
	if volumeResp.StatusCode == 202 {
		log.Info("Volume Attached to Server")
	} else {
		err := d.Remove()
		if err != nil {
			return fmt.Errorf("error deleting resources after unsuccesfull creation of volume: %v", err)
		}
		return fmt.Errorf("error attaching a volume to a server: %s", volumeResp.Status)
	}
	err = d.waitTillProvisioned(volumeResp.Header.Get("location"))
	if err != nil {
		return fmt.Errorf("error waiting till provisioned: %s", err.Error())
	}
	d.VolumeId = *volume.Id

	l, _ := strconv.Atoi(d.LanId)
	lanId := int32(l)
	nicDhcp := true
	nic := sdkgo.Nic{
		Properties: &sdkgo.NicProperties{
			Name: &d.MachineName,
			Lan:  &lanId,
			Ips:  ipBlock.Properties.Ips,
			Dhcp: &nicDhcp,
		},
	}
	nic, nicResp, err := d.Config.NicApi.DatacentersServersNicsPost(context.TODO(), *dc.Id, d.ServerId).Nic(nic).Execute()
	if err != nil {
		return fmt.Errorf("error attaching NIC to server: %s", err.Error())
	}
	if nicResp.StatusCode == 202 {
		log.Info("NIC Attached to Server")
	} else {
		err := d.Remove()
		if err != nil {
			return fmt.Errorf("error deleting resources after unsuccesfull creation of NIC: %v", err)
		}
		return fmt.Errorf("error creating a NIC: %s", nicResp.Status)
	}
	err = d.waitTillProvisioned(nicResp.Header.Get("location"))
	if err != nil {
		return fmt.Errorf("error waiting till provisioned: %s", err.Error())
	}
	d.NicId = *nic.Id

	ips := *ipBlock.Properties.Ips
	d.IPAddress = ips[0]
	log.Info(d.IPAddress)

	return nil
}

// Remove deletes the machine and resources associated to it.
func (d *Driver) Remove() error {
	ctx := context.Background()
	multierr := mcnutils.MultiError{
		Errs: []error{},
	}

	// NOTE:
	//   - if a resource is already gone or errors occur while deleting it, we
	//     continue removing other resources instead of failing

	log.Info("Starting deleting resources...")
	err := d.getClient()
	if err != nil {
		multierr.Errs = append(multierr.Errs, err)
	}

	// If the DataCenter existed before creating the machine, do not delete it at clean-up
	if !d.DCExists {
		err := d.removeServer(d.DatacenterId, d.ServerId, d.LanId)
		if err != nil {
			multierr.Errs = append(multierr.Errs, fmt.Errorf("error deleting server: %v", err))
		}
		_, resp, err := d.Config.DataCenterApi.DatacentersDelete(ctx, d.DatacenterId).Execute()
		if err != nil {
			multierr.Errs = append(multierr.Errs, fmt.Errorf("erro deleting datacenter: %v", err))
		}
		if resp.StatusCode > 299 {
			multierr.Errs = append(multierr.Errs, fmt.Errorf("error deleting datacenter, API Response status: %s", resp.Status))
		}
		err = d.waitTillProvisioned(resp.Header.Get("location"))
		if err != nil {
			multierr.Errs = append(multierr.Errs, fmt.Errorf("error waiting for datacenter to be deleted: %v", err))
		}
		log.Info("DataCenter Deleted")
	} else {
		err := d.removeServer(d.DatacenterId, d.ServerId, d.LanId)
		if err != nil {
			multierr.Errs = append(multierr.Errs, fmt.Errorf("error deleting server: %v", err))
		}
	}

	ipBlocks, _, err := d.Config.IPBlocksApi.IpblocksGet(ctx).Execute()
	if err != nil {
		multierr.Errs = append(multierr.Errs, fmt.Errorf("error getting ipblock: %v", err))
	}
	for _, i := range *ipBlocks.Items {
		for _, v := range *i.Properties.Ips {
			if d.IPAddress == v {
				_, resp, err := d.Config.IPBlocksApi.IpblocksDelete(ctx, *i.Id).Execute()
				if err != nil {
					multierr.Errs = append(multierr.Errs, fmt.Errorf("error deleting ipblock: %v", err))
				}
				if resp.StatusCode > 299 {
					multierr.Errs = append(multierr.Errs, fmt.Errorf("error deleting ipblock, API Response status: %s", resp.Status))
				}
			}
		}
	}
	log.Info("IpBlock Deleted")

	if len(multierr.Errs) == 0 {
		return nil
	}

	return multierr
}

// Start issues a power on for the machine instance.
func (d *Driver) Start() error {
	serverState, err := d.GetState()
	if err != nil {
		return fmt.Errorf("error getting state: %v", err)
	}
	err = d.getClient()
	if err != nil {
		return err
	}

	if serverState != state.Running {
		_, _, err = d.Config.ServerApi.DatacentersServersStartPost(context.TODO(), d.DatacenterId, d.ServerId).Execute()
		if err != nil {
			return fmt.Errorf("error starting server: %v", err)
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

	err = d.getClient()
	if err != nil {
		return err
	}
	_, _, err = d.Config.ServerApi.DatacentersServersStopPost(context.TODO(), d.DatacenterId, d.ServerId).Execute()
	if err != nil {
		return fmt.Errorf("error stoping server: %v", err)
	}
	return nil
}

// Restart reboots the machine instance.
func (d *Driver) Restart() error {
	err := d.getClient()
	if err != nil {
		return err
	}
	_, resp, err := d.Config.ServerApi.DatacentersServersRebootPost(context.TODO(), d.DatacenterId, d.ServerId).Execute()
	if err != nil {
		return fmt.Errorf("error restarting server: %v", err)
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf("error restarting server, API Response status: %v", resp.Status)
	}
	return nil
}

// Kill stops the machine instance
func (d *Driver) Kill() error {
	_, resp, err := d.Config.ServerApi.DatacentersServersStopPost(context.TODO(), d.DatacenterId, d.ServerId).Execute()
	if err != nil {
		return fmt.Errorf("error stoping server: %v", err)
	}
	if resp.StatusCode != 202 {
		return fmt.Errorf(resp.Status)
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
	err := d.getClient()
	if err != nil {
		return "", err
	}
	server, _, err := d.Config.ServerApi.DatacentersServersFindById(context.TODO(), d.DatacenterId, d.ServerId).Execute()
	if err != nil {
		return "", fmt.Errorf("error getting server by id: %v", err)
	}

	entitiesNicItems := *server.Entities.Nics.Items
	entityNic := entitiesNicItems[0]
	entityNicIps := *entityNic.Properties.Ips
	d.IPAddress = entityNicIps[0]

	if d.IPAddress == "" {
		return "", fmt.Errorf("IP address is not set")
	}
	return d.IPAddress, nil
}

// GetState returns the state of the machine role instance.
func (d *Driver) GetState() (state.State, error) {
	err := d.getClient()
	if err != nil {
		return state.None, err
	}
	server, serverResp, err := d.Config.ServerApi.DatacentersServersFindById(context.TODO(), d.DatacenterId, d.ServerId).Execute()
	if err != nil {
		return state.None, err
	}
	if serverResp.StatusCode > 299 {
		if serverResp.StatusCode == 401 {
			return state.None, fmt.Errorf("unauthorized: either user name or password are incorrect")

		} else {
			return state.None, fmt.Errorf("error occurred fetching a server: %s", serverResp.Status)
		}
	}

	switch *server.Metadata.State {
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
	case "CHRASHED":
		return state.Error, nil
	case "INACTIVE":
		return state.Stopped, nil
	}
	return state.None, nil
}

/*
	Private helper functions
*/

func (d *Driver) getClient() error {
	if d.Username == "" || d.Password == "" || d.URL == "" {
		return fmt.Errorf("username, password or server-url not provided")
	}
	clientConfig := &sdkgo.Configuration{
		Username: d.Username,
		Password: d.Password,
		Servers: sdkgo.ServerConfigurations{
			sdkgo.ServerConfiguration{
				URL: d.URL,
			},
		},
	}

	d.Config = sdkgo.NewAPIClient(clientConfig)
	return nil
}

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

func (d *Driver) waitTillProvisioned(path string) error {
	err := d.getClient()
	if err != nil {
		return err
	}
	for i := 0; i < waitCount; i++ {
		requestStatus, _, err := d.Config.RequestApi.RequestsStatusGet(context.TODO(), d.getRequestId(path)).Execute()
		if err != nil {
			return fmt.Errorf("error getting request status: %s", err.Error())
		}
		if *requestStatus.Metadata.Status == "DONE" {
			return nil
		}
		if *requestStatus.Metadata.Status == "FAILED" {
			return fmt.Errorf(*requestStatus.Metadata.Message)
		}
		time.Sleep(10 * time.Second)
		i++
	}

	return fmt.Errorf("timeout has expired")
}

func (d *Driver) removeServer(datacenterId string, serverId string, lanId string) error {
	server, serverResp, err := d.Config.ServerApi.DatacentersServersFindById(context.TODO(), datacenterId, serverId).Execute()
	if err != nil {
		return err
	}
	if serverResp.StatusCode > 299 {
		return fmt.Errorf(serverResp.Status)
	}

	if server.Entities != nil && server.Entities.Nics != nil && len(*server.Entities.Nics.Items) > 0 {
		nicItems := *server.Entities.Nics.Items
		nicId := *nicItems[0].Id
		_, resp, err := d.Config.NicApi.DatacentersServersNicsDelete(context.TODO(), d.DatacenterId, serverId, nicId).Execute()
		if err != nil {
			return err
		}
		if resp.StatusCode > 299 {
			return fmt.Errorf(resp.Status)
		}
		err = d.waitTillProvisioned(resp.Header.Get("location"))
		if err != nil {
			return err
		}
		log.Info("NIC Deleted")
	}

	if server.Entities != nil && server.Entities.Volumes != nil && len(*server.Entities.Volumes.Items) > 0 {
		volumesItems := *server.Entities.Volumes.Items
		volumeId := *volumesItems[0].Id
		_, resp, err := d.Config.VolumeApi.DatacentersVolumesDelete(context.TODO(), d.DatacenterId, volumeId).Execute()
		if err != nil {
			return err
		}
		if resp.StatusCode > 299 {
			return fmt.Errorf(resp.Status)
		}
		err = d.waitTillProvisioned(resp.Header.Get("location"))
		if err != nil {
			return err
		}
		log.Info("Volume Deleted")
	}

	_, resp, err := d.Config.ServerApi.DatacentersServersDelete(context.TODO(), datacenterId, serverId).Execute()
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf(resp.Status)
	}
	err = d.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("Server Deleted")

	_, resp, err = d.Config.LanApi.DatacentersLansDelete(context.TODO(), datacenterId, lanId).Execute()
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf(resp.Status)
	}
	err = d.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("LAN Deleted")

	return nil
}

func (d *Driver) getRequestId(path string) string {
	if !strings.Contains(path, d.URL) {
		fmt.Errorf("path does not contain %s", d.URL)
		return ""
	}
	str := strings.Split(path, "/")
	return str[len(str)-2]
}

func (d *Driver) getImageId(imageName string) (string, error) {
	err := d.getClient()
	if err != nil {
		return "", err
	}
	d.UseAlias = false
	// first look if the provided parameter matches an alias, if a match is found we return the image alias
	regionId, locationId := d.getRegionIdAndLocationId()
	location, _, err := d.Config.LocationApi.LocationsFindByRegionIdAndId(context.TODO(), regionId, locationId).Execute()
	if err != nil {
		return "", err
	}
	for _, alias := range *location.Properties.ImageAliases {
		if alias == imageName {
			d.UseAlias = true
			return imageName, nil
		}
	}

	// if no alias matches we do extended search and return the image id
	images, imagesResp, err := d.Config.ImageApi.ImagesGet(context.TODO()).Execute()
	if err != nil {
		return "", err
	}

	if imagesResp.StatusCode == 401 {
		return "", fmt.Errorf("error: authentication failed")
	}

	for _, image := range *images.Items {
		imgName := ""
		if *image.Properties.Name != "" {
			imgName = *image.Properties.Name
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
	return "", nil
}

func (d *Driver) getRegionIdAndLocationId() (regionId, locationId string) {
	ids := strings.Split(d.Location, "/")
	// location has standard format: {regionId}/{locationId}
	if len(ids) != 2 {
		fmt.Errorf("error getting Region Id and Location Id from %s", d.Location)
		return "", ""
	}
	return ids[0], ids[1]
}
