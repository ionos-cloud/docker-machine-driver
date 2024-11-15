package ionoscloud

import (
	"encoding/base64"
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/docker-machine-driver/internal/utils"
	mockutils "github.com/ionos-cloud/docker-machine-driver/internal/utils/mocks"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
)

const (
	defaultHostName  = "default1"
	defaultStorePath = "path"
)

var (
	// Common variables used
	datacenterName         = "datacenter_name"
	datacenterId           = "datacenter_id"
	lanId                  = "1"
	serverName             = "server_name"
	serverId               = "server_id"
	volumeName             = "volume_name"
	volumeId               = "volume_id"
	nicId                  = "nic_id"
	testRegion             = "us/ewr"
	testVar                = "test"
	testImageIdVar         = "test-image-id"
	locationId             = "las"
	imageType              = "HDD"
	imageAlias             = "ubuntu:20.04"
	dcVersion              = int32(1)
	testErr                = fmt.Errorf("errFoo")
	lanId1                 = "2"
	lanId1Int              = 2
	lanName1               = "test_lan1"
	lanId2                 = "5"
	lanName2               = "test2"
	localhost_ip           = "test_local_ip"
	cpuFamily              = "INTEL_ICELAKE"
	imagePassword          = "<testdata>"
	nicDhcp                = false
	nicIps                 = []string{localhost_ip, "127.0.0.3"}
	diskType               = "SSD"
	volumeAvailabilityZone = "ZONE_1"
	serverAvailabilityZone = "ZONE_2"
	cores                  = 4
	cloudInit              = "#cloud-config\nHostname: testdata"
	ram                    = 2048
	diskSize               = 100
)

var (
	// Sdk resources used
	dc = &sdkgo.Datacenter{
		Id: sdkgo.ToPtr(datacenterId),
		Properties: &sdkgo.DatacenterProperties{
			Name:        sdkgo.ToPtr(datacenterName),
			Description: sdkgo.ToPtr("datacenter_description"),
			Location:    sdkgo.ToPtr(testRegion),
			Version:     sdkgo.ToPtr(int32(1)),
		},
	}
	lan_post = &sdkgo.LanPost{
		Id: sdkgo.ToPtr(lanId),
		Properties: &sdkgo.LanPropertiesPost{
			Name:   &lanName1,
			Public: sdkgo.ToPtr(true),
		},
	}
	lan_post_private = &sdkgo.LanPost{
		Id: sdkgo.ToPtr(lanId),
		Properties: &sdkgo.LanPropertiesPost{
			Name:   &lanName1,
			Public: sdkgo.ToPtr(false),
		},
	}
	lan_get = &sdkgo.Lan{
		Id: sdkgo.ToPtr(lanId),
		Properties: &sdkgo.LanProperties{
			Name:   &lanName1,
			Public: sdkgo.ToPtr(true),
		},
	}
	lan_get_private = &sdkgo.Lan{
		Id: sdkgo.ToPtr(lanId),
		Properties: &sdkgo.LanProperties{
			Name:   &lanName1,
			Public: sdkgo.ToPtr(false),
		},
	}

	server = &sdkgo.Server{
		Id: sdkgo.ToPtr(serverId),
		Properties: &sdkgo.ServerProperties{
			Name:             sdkgo.ToPtr(serverName),
			Ram:              sdkgo.ToPtr(int32(2048)),
			Cores:            sdkgo.ToPtr(int32(2)),
			CpuFamily:        sdkgo.ToPtr("AMD_OPTERON"),
			AvailabilityZone: sdkgo.ToPtr("AUTO"),
		},
		Entities: &sdkgo.ServerEntities{
			Volumes: &sdkgo.AttachedVolumes{
				Items: &[]sdkgo.Volume{
					{
						Id: sdkgo.ToPtr(volumeId),
					},
				},
			},
			Nics: &sdkgo.Nics{
				Items: &[]sdkgo.Nic{
					{
						Id: sdkgo.ToPtr(nicId),
						Properties: &sdkgo.NicProperties{
							Name: sdkgo.ToPtr(defaultHostName),
						},
					},
				},
			},
		},
	}

	location = &sdkgo.Location{
		Id: &locationId,
		Properties: &sdkgo.LocationProperties{
			ImageAliases: &[]string{testVar},
		},
	}
	images = sdkgo.Images{
		Items: &[]sdkgo.Image{image, imageFoundById},
	}
	image = sdkgo.Image{
		Id: &testVar,
		Properties: &sdkgo.ImageProperties{
			Name:      &imageAlias,
			ImageType: &imageType,
			Location:  &testRegion,
		},
	}
	imageFoundById = sdkgo.Image{
		Id: &testImageIdVar,
		Properties: &sdkgo.ImageProperties{
			Name:      &imageAlias,
			ImageType: &imageType,
			Location:  &testRegion,
		},
	}
	ipblock = &sdkgo.IpBlock{
		Id: sdkgo.ToPtr("IPBlock_ID"),
		Properties: &sdkgo.IpBlockProperties{
			Name:     sdkgo.ToPtr("ipblock_name"),
			Location: &testRegion,
			Ips:      &[]string{"ip1", "ip2"},
		},
	}
	lan1 = &sdkgo.Lan{
		Id:         &testVar,
		Properties: &sdkgo.LanProperties{Name: &lanName1},
	}
	lan = &sdkgo.LanPost{
		Id:         &testVar,
		Properties: &sdkgo.LanPropertiesPost{Public: sdkgo.ToPtr(false), Name: &lanName1},
	}
	privateLan = &sdkgo.Lan{
		Id:         &testVar,
		Properties: &sdkgo.LanProperties{Public: sdkgo.ToPtr(false)},
	}
	privateLan2 = &sdkgo.Lan{
		Id:         &lanId1,
		Properties: &sdkgo.LanProperties{Public: sdkgo.ToPtr(false), Name: &lanName1},
	}
	privateLan3 = &sdkgo.Lan{
		Id:         &lanId2,
		Properties: &sdkgo.LanProperties{Public: sdkgo.ToPtr(false), Name: &lanName2},
	}
	lans = sdkgo.Lans{
		Items: &[]sdkgo.Lan{},
	}
	additionalLans = sdkgo.Lans{
		Items: &[]sdkgo.Lan{*privateLan2, *privateLan3},
	}

	volume = &sdkgo.Volume{
		Id: &testVar,
	}
	nic = &sdkgo.Nic{
		Id: &testVar,
		Properties: &sdkgo.NicProperties{
			Ips:  &[]string{localhost_ip},
			Name: sdkgo.PtrString("test"),
		},
	}
	lansGateways  = []sdkgo.NatGatewayLanProperties{{Id: sdkgo.PtrInt32(1), GatewayIps: &[]string{"x.x.x.x"}}}
	lansGateways2 = []sdkgo.NatGatewayLanProperties{{Id: sdkgo.PtrInt32(2), GatewayIps: &[]string{"x.x.x.x"}}}
	natName       = "nat-name"
	nat           = &sdkgo.NatGateway{
		Id: &testVar,
		Properties: &sdkgo.NatGatewayProperties{
			PublicIps: &[]string{"x.x.x.x"},
			Lans:      &lansGateways,
			Name:      sdkgo.PtrString(natName),
		},
	}
	nats = &sdkgo.NatGateways{
		Items: &[]sdkgo.NatGateway{*nat},
	}
	dcs = &sdkgo.Datacenters{
		Items: &[]sdkgo.Datacenter{*dc},
	}
	ips = []string{testVar}

	cube_template = &sdkgo.Template{
		Id: sdkgo.PtrString("template-id"),
		Properties: &sdkgo.TemplateProperties{
			Name: sdkgo.PtrString("Basic Cube XS"),
		},
	}

	cube_templates = &sdkgo.Templates{
		Items: &[]sdkgo.Template{*cube_template},
	}
)

var (
	// Common flags set
	authFlagsSet = map[string]interface{}{
		flagUsername: "IONOSCLOUD_USERNAME",
		flagPassword: "IONOSCLOUD_PASSWORD",
	}
	authDcIdFlagsSet = map[string]interface{}{
		flagUsername:     "IONOSCLOUD_USERNAME",
		flagPassword:     "IONOSCLOUD_PASSWORD",
		flagDatacenterId: "IONOSCLOUD_DATACENTER_ID",
	}
	propertiesImageFoundById = &utils.ClientVolumeProperties{
		DiskType:      defaultDiskType,
		Name:          defaultHostName,
		ImageId:       testImageIdVar,
		ImagePassword: defaultImagePassword,
		Zone:          defaultAvailabilityZone,
		SshKey:        testVar,
		DiskSize:      float32(50),
	}
	propertiesImageAlias = &utils.ClientVolumeProperties{
		DiskType:      defaultDiskType,
		Name:          defaultHostName,
		ImageAlias:    testVar,
		ImagePassword: defaultImagePassword,
		Zone:          defaultAvailabilityZone,
		SshKey:        testVar,
		DiskSize:      float32(50),
	}
)

func NewTestDriverFlagsSet(t *testing.T, flagsSet map[string]interface{}) (*Driver, *mockutils.MockClientService) {
	driver, clientMock := NewTestDriver(t, defaultHostName, defaultStorePath)
	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: flagsSet,
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)
	return driver, clientMock
}

func NewTestDriver(t *testing.T, hostName, storePath string) (*Driver, *mockutils.MockClientService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	clientMock := mockutils.NewMockClientService(ctrl)
	d := NewDerivedDriver(hostName, storePath)
	d.client = func() utils.ClientService {
		return clientMock
	}
	return d, clientMock
}

func TestNewDriver(t *testing.T) {

	NewDriver("test-machine", defaultStorePath)
}

func TestSetConfigFromDefaultFlags(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, map[string]interface{}{})
	sshPort, err := driver.GetSSHPort()
	assert.Equal(t, 22, sshPort)
	assert.NoError(t, err)
	assert.Equal(t, "", driver.Username)
	assert.Equal(t, "", driver.Password)
	assert.Equal(t, "", driver.Token)
	assert.Equal(t, sdkgo.DefaultIonosServerUrl, driver.Endpoint)
	assert.Equal(t, 2, driver.Cores)
	assert.Equal(t, 2048, driver.Ram)
	assert.Equal(t, defaultRegion, driver.Location)
	assert.Equal(t, defaultDiskType, driver.DiskType)
	assert.Equal(t, false, driver.NicDhcp)
	assert.Equal(t, 50, driver.DiskSize)
	assert.Equal(t, "", driver.DatacenterId)
	assert.Equal(t, defaultAvailabilityZone, driver.VolumeAvailabilityZone)
	assert.Equal(t, defaultAvailabilityZone, driver.ServerAvailabilityZone)
}

func TestSetConfigFromCustomFlags(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, map[string]interface{}{
		flagServerRam: 1024,
		flagDiskType:  "SSD",
		flagEndpoint:  "",
		flagNicDhcp:   true,
		flagNicIps:    []string{"127.0.0.4", "127.0.0.54"},
	})
	sshPort, err := driver.GetSSHPort()
	assert.Equal(t, 22, sshPort)
	assert.NoError(t, err)
	assert.Equal(t, "", driver.Username)
	assert.Equal(t, "", driver.Password)
	assert.Equal(t, sdkgo.DefaultIonosServerUrl, driver.Endpoint)
	assert.Equal(t, 2, driver.Cores)
	assert.Equal(t, 1024, driver.Ram)
	assert.Equal(t, defaultRegion, driver.Location)
	assert.Equal(t, "SSD", driver.DiskType)
	assert.Equal(t, true, driver.NicDhcp)
	assert.Equal(t, []string{"127.0.0.4", "127.0.0.54"}, driver.NicIps)
	assert.Equal(t, 50, driver.DiskSize)
	assert.Equal(t, "", driver.DatacenterId)
	assert.Equal(t, defaultAvailabilityZone, driver.VolumeAvailabilityZone)
	assert.Equal(t, defaultAvailabilityZone, driver.ServerAvailabilityZone)
}

func TestDriverName(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, authFlagsSet)
	assert.Equal(t, driverName, driver.DriverName())
}

func TestPreCreateCheckAuthErr(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, map[string]interface{}{})
	err := driver.PreCreateCheck()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "please provide username($IONOSCLOUD_USERNAME) and password($IONOSCLOUD_PASSWORD) or token($IONOSCLOUD_TOKEN) to authenticate")
}

func TestPreCreateCheckUserNameErr(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, map[string]interface{}{
		flagPassword: "IONOSCLOUD_PASSWORD",
	})
	err := driver.PreCreateCheck()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "please provide username as parameter --ionoscloud-username or as environment variable $IONOSCLOUD_USERNAME")
}

func TestPreCreateCheckPasswordErr(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, map[string]interface{}{
		flagUsername: "IONOSCLOUD_USERNAME",
	})
	err := driver.PreCreateCheck()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "please provide password as parameter --ionoscloud-password or as environment variable $IONOSCLOUD_PASSWORD")
}

func TestPreCreateCheckDataCenterIdErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.DatacenterId = datacenterId
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(nil, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
}

func TestPreCreateCheckDataCenterErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.DatacenterId = datacenterId
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(nil, fmt.Errorf("error getting datacenter: 404 not found"))

	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	err := driver.PreCreateCheck()
	assert.Error(t, err)
}

func TestPreCreateImageIdErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.DatacenterId = datacenterId
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, fmt.Errorf("error getting image: 404 not found"))
	err := driver.PreCreateCheck()
	assert.Error(t, err)
}

func TestPreCreateCheck(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
}
func TestPreCreateLans(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = "test"
	driver.AdditionalLans = []string{lanName1, "wrong_value"}
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&additionalLans, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.True(t, reflect.DeepEqual(driver.AdditionalLansIds, []int{lanId1Int}))
	assert.NoError(t, err)
}

func TestCreateSSHKeyErr(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = ""
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateErr(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, authFlagsSet)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}
	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{}}, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),

		clientMock.EXPECT().CreateDatacenter(datacenterName, testRegion).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, true).Return(lan_post, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_post.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, *ipblock.Properties.Ips, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreateLanProvided(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}
	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get}}, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, *ipblock.Properties.Ips, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreatePropertiesSet(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterName = datacenterName
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}

	driver.CpuFamily = cpuFamily
	driver.Location = testRegion
	driver.ImagePassword = imagePassword
	driver.NicDhcp = nicDhcp
	driver.NicIps = nicIps
	driver.DiskType = diskType
	driver.VolumeAvailabilityZone = volumeAvailabilityZone
	driver.ServerAvailabilityZone = serverAvailabilityZone
	driver.Cores = cores
	driver.Ram = ram
	driver.CloudInit = cloudInit
	driver.DiskSize = diskSize

	lan_get2 := &sdkgo.Lan{
		Id: sdkgo.ToPtr(lanId2),
		Properties: &sdkgo.LanProperties{
			Name:   &lanName2,
			Public: sdkgo.ToPtr(true),
		},
	}

	server = &sdkgo.Server{
		Id: sdkgo.ToPtr(serverId),
		Properties: &sdkgo.ServerProperties{
			Name:             sdkgo.ToPtr(serverName),
			Ram:              sdkgo.ToPtr(int32(2048)),
			Cores:            sdkgo.ToPtr(int32(2)),
			CpuFamily:        sdkgo.ToPtr("AMD_OPTERON"),
			AvailabilityZone: sdkgo.ToPtr("AUTO"),
		},
		Entities: &sdkgo.ServerEntities{
			Volumes: &sdkgo.AttachedVolumes{
				Items: &[]sdkgo.Volume{
					{
						Id: sdkgo.ToPtr(volumeId),
					},
				},
			},
			Nics: &sdkgo.Nics{
				Items: &[]sdkgo.Nic{
					{
						Id: sdkgo.ToPtr(nicId),
						Properties: &sdkgo.NicProperties{
							Name: sdkgo.ToPtr(defaultHostName),
						},
					},
					{
						Id: sdkgo.ToPtr("nic_id-2"),
						Properties: &sdkgo.NicProperties{
							Name: sdkgo.ToPtr("different_name"),
						},
					},
				},
			},
		},
	}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get, *lan_get2}}, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, cpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, serverAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				assert.Nil(t, serverToCreate.Properties.TemplateUuid)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(diskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 2)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicIps, *nics[0].Properties.Ips)
				assert.Equal(t, nicDhcp, *nics[0].Properties.Dhcp)

				assert.Equal(t, driver.MachineName+" "+lanId2, *nics[1].Properties.Name)
				assert.Equal(t, int32(5), *nics[1].Properties.Lan)
				assert.Nil(t, nics[1].Properties.Ips)
				assert.Equal(t, true, *nics[1].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreateCubePropertiesSet(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterName = datacenterName
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}

	driver.ServerType = "CUBE"
	driver.CpuFamily = cpuFamily
	driver.Location = testRegion
	driver.ImagePassword = imagePassword
	driver.NicDhcp = nicDhcp
	driver.NicIps = nicIps
	driver.DiskType = diskType
	driver.VolumeAvailabilityZone = volumeAvailabilityZone
	driver.ServerAvailabilityZone = serverAvailabilityZone
	driver.Cores = cores
	driver.Ram = ram
	driver.CloudInit = cloudInit
	driver.DiskSize = diskSize
	driver.WaitForIpChange = true
	driver.WaitForIpChangeTimeout = 176

	lan_get2 := &sdkgo.Lan{
		Id: sdkgo.ToPtr(lanId2),
		Properties: &sdkgo.LanProperties{
			Name:   &lanName2,
			Public: sdkgo.ToPtr(true),
		},
	}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get, *lan_get2}}, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().GetTemplates().Return(cube_templates, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Nil(t, serverToCreate.Properties.CpuFamily)
				assert.Nil(t, serverToCreate.Properties.Ram)
				assert.Nil(t, serverToCreate.Properties.Cores)
				assert.Nil(t, serverToCreate.Properties.AvailabilityZone)
				assert.Equal(t, "CUBE", *serverToCreate.Properties.Type)
				assert.Equal(t, *cube_template.Id, *serverToCreate.Properties.TemplateUuid)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, "DAS", *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Nil(t, volumes[0].Properties.Size)
				assert.Nil(t, volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 2)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicIps, *nics[0].Properties.Ips)
				assert.Equal(t, nicDhcp, *nics[0].Properties.Dhcp)

				assert.Equal(t, driver.MachineName+" "+lanId2, *nics[1].Properties.Name)
				assert.Equal(t, int32(5), *nics[1].Properties.Lan)
				assert.Nil(t, nics[1].Properties.Ips)
				assert.Equal(t, true, *nics[1].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().WaitForNicIpChange(*dc.Id, serverId, nicId, 176).Return(nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}
func TestCreateNatPublicIps(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.NicIps = nicIps
	driver.NatPublicIps = []string{"127.0.0.4"}
	driver.CreateNat = true

	driver.NatLansToGateways = map[string][]string{"1": {"127.0.0.3"}}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get_private}}, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicIps, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
		clientMock.EXPECT().CreateNat(
			*dc.Id, "docker-machine-nat", driver.NatPublicIps, driver.NatFlowlogs, driver.NatRules, driver.NatLansToGateways,
			net.ParseIP((driver.NicIps)[0]).Mask(net.CIDRMask(24, 32)).String()+"/24", driver.SkipDefaultNatRules,
		).Return(nat, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreateNat(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.PrivateLan = true
	driver.NicIps = nicIps
	driver.NatPublicIps = []string{"127.0.0.4"}
	driver.SkipDefaultNatRules = true
	driver.CreateNat = true
	driver.NatFlowlogs = []string{"test_name:ACCEPTED:INGRESS:test_bucket", "test_name2:REGECTED:EGRESS:test_bucket"}
	driver.NatRules = []string{
		"name1:SNAT:TCP::10.0.1.0/24:10.0.2.0/24:100:500",
		"name2:SNAT:ALL::10.0.1.0/24::1023:1500",
	}
	driver.NatLansToGateways = map[string][]string{"1": {"127.0.0.3"}}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{}}, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, false).Return(lan_post_private, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicIps, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
		clientMock.EXPECT().CreateNat(
			*dc.Id, "docker-machine-nat", driver.NatPublicIps, driver.NatFlowlogs, driver.NatRules, driver.NatLansToGateways,
			net.ParseIP((driver.NicIps)[0]).Mask(net.CIDRMask(24, 32)).String()+"/24", driver.SkipDefaultNatRules,
		).Return(nat, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreateExistingNatPatch(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.PrivateLan = true
	driver.NicIps = nicIps
	driver.NatPublicIps = []string{"127.0.0.4"}
	driver.SkipDefaultNatRules = true
	driver.NatFlowlogs = []string{"test_name:ACCEPTED:INGRESS:test_bucket", "test_name2:REGECTED:EGRESS:test_bucket"}
	driver.NatRules = []string{
		"name1:SNAT:TCP::10.0.1.0/24:10.0.2.0/24:100:500",
		"name2:SNAT:ALL::10.0.1.0/24::1023:1500",
	}
	driver.NatLansToGateways = map[string][]string{"1": {"127.0.0.3"}}

	nat = &sdkgo.NatGateway{
		Id: &testVar,
		Properties: &sdkgo.NatGatewayProperties{
			PublicIps: &[]string{"x.x.x.x"},
			Lans:      &lansGateways2,
			Name:      sdkgo.PtrString(defaultNatName),
		},
	}

	nats = &sdkgo.NatGateways{
		Items: &[]sdkgo.NatGateway{*nat},
	}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{}}, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, false).Return(lan_post_private, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicIps, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
		clientMock.EXPECT().GetNat(*dc.Id, *nat.Id).Return(nat, nil),
		clientMock.EXPECT().PatchNat(*dc.Id, *nat.Id, *nat.Properties.Name, *nat.Properties.PublicIps, gomock.AssignableToTypeOf([]sdkgo.NatGatewayLanProperties{})).DoAndReturn(
			func(datacenterId, NatId, NatName string, PublicIps []string, NatLans []sdkgo.NatGatewayLanProperties) (*sdkgo.NatGateway, error) {

				return nat, nil
			}),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreateExistingNatNoPatch(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.PrivateLan = true
	driver.NicIps = nicIps
	driver.NatPublicIps = []string{"127.0.0.4"}
	driver.SkipDefaultNatRules = true
	driver.NatFlowlogs = []string{"test_name:ACCEPTED:INGRESS:test_bucket", "test_name2:REGECTED:EGRESS:test_bucket"}
	driver.NatRules = []string{
		"name1:SNAT:TCP::10.0.1.0/24:10.0.2.0/24:100:500",
		"name2:SNAT:ALL::10.0.1.0/24::1023:1500",
	}
	driver.NatLansToGateways = map[string][]string{"1": {"127.0.0.3"}}
	driver.NatName = natName

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{}}, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, false).Return(lan_post_private, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
		clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicIps, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
		clientMock.EXPECT().GetNat(*dc.Id, *nat.Id).Return(nat, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreateImageIdSSHInCloudInit(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.SSHUser = "username"
	driver.DatacenterId = "value"
	driver.LanId = lanId
	driver.Image = testVar
	driver.ImagePassword = imagePassword
	driver.SSHInCloudInit = true
	location.Properties.ImageAliases = &[]string{}

	test := []interface{}{map[interface{}]interface{}{
		"name":                driver.SSHUser,
		"lock_passwd":         true,
		"sudo":                "ALL=(ALL) NOPASSWD:ALL",
		"create_groups":       false,
		"no_user_group":       true,
		"ssh_authorized_keys": []string{driver.SSHKey},
	}}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, lanId).Return(lan1, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil),
		clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil),
		clientMock.EXPECT().GetImageById(driver.Image).Return(&imageFoundById, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "users", test, false, "append").Return("test_string", nil),
		clientMock.EXPECT().UpdateCloudInitFile("test_string", "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Nil(t, volumes[0].Properties.ImageAlias)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, ips, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, *server.Id, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, *server.Id, nicId).Return(nic, nil),
	)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateImageAliasSSHUser(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.SSHUser = "username"
	driver.DatacenterId = "value"
	driver.LanId = lanId
	driver.Image = testVar
	driver.ImagePassword = imagePassword
	location.Properties.ImageAliases = &[]string{testVar}

	test := []interface{}{map[interface{}]interface{}{
		"name":                driver.SSHUser,
		"lock_passwd":         true,
		"sudo":                "ALL=(ALL) NOPASSWD:ALL",
		"create_groups":       false,
		"no_user_group":       true,
		"ssh_authorized_keys": []string{driver.SSHKey},
	}}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, lanId).Return(lan1, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil),
		clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "users", test, false, "append").Return("test_string", nil),
		clientMock.EXPECT().UpdateCloudInitFile("test_string", "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, driver.ServerAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 1)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{driver.SSHKey}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testVar, *volumes[0].Properties.ImageAlias)
				assert.Nil(t, volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)

				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 1)
				assert.Equal(t, driver.MachineName, *nics[0].Properties.Name)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, ips, *nics[0].Properties.Ips)
				assert.Equal(t, driver.NicDhcp, *nics[0].Properties.Dhcp)
				serverToCreate.Id = &serverId

				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, *server.Id, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dc.Id, *server.Id, nicId).Return(nic, nil),
	)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateIpBlockErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
	driver.UseAlias = true
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, "docker-machine-lan", true).Return(lan, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, "test", int32(2)).Return(server, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, testErr)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, *lan1.Id).Return(testErr)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateGetImageErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, testErr)
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	err := driver.PreCreateCheck()
	assert.Error(t, err)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, testErr)
	err = driver.PreCreateCheck()
	assert.Error(t, err)
}

func TestCreateGetDatacenterErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateDatacenterErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = ""
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateDatacenter(driver.DatacenterName, driver.Location).Return(dc, testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateLanErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
	driver.LanId = ""
	driver.IpBlockId = *ipblock.Id
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, "docker-machine-lan", true).Return(lan, testErr)
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(testErr)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(testErr)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(testErr)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(testErr)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateServerErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = *ipblock.Id
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, testErr)
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(testErr)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(testErr)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(testErr)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(testErr)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateServerRemove(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = *ipblock.Id
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, testErr)
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(nil)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(nil)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(nil)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(nil)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(nil)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateGetIpBlockErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(2)).Return(server, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, testErr)
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(testErr)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(testErr)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(testErr)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(testErr)
	clientMock.EXPECT().RemoveIpBlock(*ipblock.Id).Return(testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateServerErr2(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = *ipblock.Id
	driver.AdditionalLansIds = []int{2, 4}
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, fmt.Errorf("error"))
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(testErr)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(testErr)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(testErr)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(testErr)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestRemove(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.DCExists = false
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(nil)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(nil)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(nil)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(nil)
	clientMock.EXPECT().RemoveDatacenter(driver.DatacenterId).Return(nil)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(nil)
	err := driver.Remove()
	assert.NoError(t, err)
}

func TestRemoveErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.DCExists = false
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(testErr)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(testErr)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(testErr)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(testErr)
	clientMock.EXPECT().RemoveDatacenter(driver.DatacenterId).Return(testErr)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(testErr)
	err := driver.Remove()
	assert.Error(t, err)
}

func TestStartErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(server, nil)
	err := driver.Start()
	assert.Error(t, err)
}

func TestStart(t *testing.T) {
	s := serverWithState(testVar, "PAUSED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	clientMock.EXPECT().StartServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Start()
	assert.NoError(t, err)
}

func TestStartServerErr(t *testing.T) {
	s := serverWithState(testVar, "INACTIVE")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	clientMock.EXPECT().StartServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error starting server"))
	err := driver.Start()
	assert.Error(t, err)
}

func TestStartRunningServer(t *testing.T) {
	s := serverWithState(testVar, "RUNNING")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	err := driver.Start()
	assert.NoError(t, err)
}

func TestStopErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	s := &sdkgo.Server{}
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	err := driver.Stop()
	assert.Error(t, err)
}

func TestStop(t *testing.T) {
	s := serverWithState(testVar, "SHUTOFF")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Stop()
	assert.NoError(t, err)
}

func TestStopServerErr(t *testing.T) {
	s := serverWithState(testVar, "PAUSED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error stoping server"))
	err := driver.Stop()
	assert.Error(t, err)
}

func TestStopStoppedServer(t *testing.T) {
	s := serverWithState(testVar, "BLOCKED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	err := driver.Stop()
	assert.NoError(t, err)
}

func TestRestartErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().RestartServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error restarting server"))
	err := driver.Restart()
	assert.Error(t, err)
}

func TestRestart(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().RestartServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Restart()
	assert.NoError(t, err)
}

func TestKillErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error stoping server"))
	err := driver.Kill()
	assert.Error(t, err)
}

func TestKill(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Kill()
	assert.NoError(t, err)
}

func TestGetSSHHostnameErr(t *testing.T) {
	s := serverWithState(testVar, "CRASHED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	_, err := driver.GetSSHHostname()
	assert.Error(t, err)
}

func TestGetURLErr(t *testing.T) {
	s := serverWithState(testVar, "SHUTOFF")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	_, err := driver.GetURL()
	assert.Error(t, err)
}

func TestGetIPErr(t *testing.T) {
	s := serverWithState(testVar, "AVAILABLE")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(2)).Return(s, testErr)
	_, err := driver.GetIP()
	assert.Error(t, err)
}

func TestGetStateErr(t *testing.T) {
	s := serverWithNicAttached(testVar, "AVAILABLE", testVar)
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, testErr)
	_, err := driver.GetState()
	assert.Error(t, err)
}

func TestGetStateShutDown(t *testing.T) {
	s := serverWithNicAttached(testVar, "SHUTDOWN", testVar)
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	_, err := driver.GetState()
	assert.NoError(t, err)
}

func TestGetStateCrashed(t *testing.T) {
	s := serverWithNicAttached(testVar, "CRASHED", testVar)
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId, int32(1)).Return(s, nil)
	_, err := driver.GetState()
	assert.NoError(t, err)
}

func TestPublicSSHKeyPath(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.publicSSHKeyPath()
}

func TestIsSwarmMaster(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.isSwarmMaster()
}

func TestGetRegionIdAndLocationId(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.Location = "test/test/test/test"
	driver.getRegionIdAndLocationId()
}

func TestGetImageId(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	driver.Location = defaultRegion
	driver.DiskType = "SSD"
	_, err := driver.getImageIdOrAlias(imageAlias)
	assert.NoError(t, err)
}

func serverWithState(serverId, serverState string) *sdkgo.Server {
	return &sdkgo.Server{
		Id: &serverId,
		Properties: &sdkgo.ServerProperties{
			VmState: &serverState,
		},
	}
}

func serverWithNicAttached(serverId, serverState, nicId string) *sdkgo.Server {
	return &sdkgo.Server{
		Id: &serverId,
		Properties: &sdkgo.ServerProperties{
			VmState: &serverState,
		},
		Entities: &sdkgo.ServerEntities{
			Nics: &sdkgo.Nics{
				Items: &[]sdkgo.Nic{
					{
						Properties: &sdkgo.NicProperties{
							Ips: &[]string{nicId},
						},
					},
				},
			},
		},
	}
}
