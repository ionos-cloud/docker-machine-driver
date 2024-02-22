package ionoscloud

import (
	"fmt"
	"net"
	"reflect"
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/docker-machine-driver/internal/utils"
	mockutils "github.com/ionos-cloud/docker-machine-driver/internal/utils/mocks"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
	"github.com/stretchr/testify/assert"
)

const (
	defaultHostName  = "default1"
	defaultStorePath = "path"
)

var (
	// Common variables used
	testRegion     = "us/ewr"
	testFalse      = false
	testVar        = "test"
	testImageIdVar = "test-image-id"
	locationId     = "las"
	imageType      = "HDD"
	imageName      = defaultImageAlias
	imageLocation  = "us/las"
	dcVersion      = int32(1)
	testErr        = fmt.Errorf("error")
	lanId1         = "2"
	lanId1Int      = 2
	lanName1       = "test"
	lanId2         = "5"
	lanName2       = "test2"
)

var (
	// Sdk resources used
	dc = &sdkgo.Datacenter{
		Id: &testVar,
		Properties: &sdkgo.DatacenterProperties{
			Name:        &testVar,
			Description: &testVar,
			Location:    &testRegion,
			Version:     &dcVersion,
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
			Name:      &imageName,
			ImageType: &imageType,
			Location:  &imageLocation,
		},
	}
	imageFoundById = sdkgo.Image{
		Id: &testImageIdVar,
		Properties: &sdkgo.ImageProperties{
			Name:      &imageName,
			ImageType: &imageType,
			Location:  &imageLocation,
		},
	}
	ipblock = &sdkgo.IpBlock{
		Id: &testVar,
	}
	ipblocks = &sdkgo.IpBlocks{
		Items: &[]sdkgo.IpBlock{
			*ipblock,
		},
	}
	lan1 = &sdkgo.Lan{
		Id: &testVar,
	}
	lan = &sdkgo.LanPost{
		Id: &testVar,
	}
	privateLan = &sdkgo.Lan{
		Id:         &testVar,
		Properties: &sdkgo.LanProperties{Public: &testFalse},
	}
	privateLan2 = &sdkgo.Lan{
		Id:         &lanId1,
		Properties: &sdkgo.LanProperties{Public: &testFalse, Name: &lanName1},
	}
	privateLan3 = &sdkgo.Lan{
		Id:         &lanId2,
		Properties: &sdkgo.LanProperties{Public: &testFalse, Name: &lanName2},
	}
	lans = sdkgo.Lans{
		Items: &[]sdkgo.Lan{},
	}
	additionalLans = sdkgo.Lans{
		Items: &[]sdkgo.Lan{*privateLan2, *privateLan3},
	}

	server = &sdkgo.Server{
		Id: &testVar,
	}
	volume = &sdkgo.Volume{
		Id: &testVar,
	}
	nic = &sdkgo.Nic{
		Id: &testVar,
		Properties: &sdkgo.NicProperties{
			Ips: &[]string{"127.0.0.1"},
		},
	}
	lansGateways = []sdkgo.NatGatewayLanProperties{{Id: &dcVersion, GatewayIps: &[]string{"x.x.x.x"}}}

	nat = &sdkgo.NatGateway{
		Id: &testVar,
		Properties: &sdkgo.NatGatewayProperties{
			PublicIps: &[]string{"x.x.x.x"},
			Lans:      &lansGateways,
			Name:      &testVar,
		},
	}
	dcs = &sdkgo.Datacenters{
		Items: &[]sdkgo.Datacenter{},
	}
	ips = []string{testVar}
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
	// Properties set for volume creation
	propertiesImageId = &utils.ClientVolumeProperties{
		DiskType:      defaultDiskType,
		Name:          defaultHostName,
		ImageId:       testVar,
		ImagePassword: defaultImagePassword,
		Zone:          defaultAvailabilityZone,
		SshKey:        testVar,
		DiskSize:      float32(50),
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
	assert.Equal(t, sdkgo.DefaultIonosServerUrl, driver.URL)
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
	assert.Equal(t, sdkgo.DefaultIonosServerUrl, driver.URL)
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
	driver.DatacenterId = testVar
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
}

func TestPreCreateCheckDataCenterErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.DatacenterId = testVar
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(nil, fmt.Errorf("error getting datacenter: 404 not found"))

	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	err := driver.PreCreateCheck()
	assert.Error(t, err)
}

func TestPreCreateImageIdErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.DatacenterId = testVar
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
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
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
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = ""
	driver.IPAddress = testVar

	attached_volumes := sdkgo.NewAttachedVolumesWithDefaults()
	attached_volumes.Items = &[]sdkgo.Volume{*volume}
	server.Entities = sdkgo.NewServerEntitiesWithDefaults()
	server.Entities.SetVolumes(*attached_volumes)

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, "docker-machine-lan", true).Return(lan, nil).Times(1)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, false, int32(0), &ips).Return(nic, nil)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateLanProvided(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	//driver.NicId = testVar
	//driver.VolumeId = testVar
	driver.LanId = testVar
	driver.Image = "e20d97c2-38ae-11ed-be62-eec5f4d7ee1e"
	//driver.IPAddress = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById("e20d97c2-38ae-11ed-be62-eec5f4d7ee1e").Return(&sdkgo.Image{}, nil)
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, propertiesImageId).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, false, int32(0), &ips).Return(nic, nil)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateNicDhcpIps(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = ""
	driver.IPAddress = testVar
	driver.NicDhcp = true
	driver.NicIps = []string{"127.0.0.1", "127.0.0.3"}
	driver.AdditionalLansIds = []int{1, 2}

	attached_volumes := sdkgo.NewAttachedVolumesWithDefaults()
	attached_volumes.Items = &[]sdkgo.Volume{*volume}
	server.Entities = sdkgo.NewServerEntitiesWithDefaults()
	server.Entities.SetVolumes(*attached_volumes)

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, "docker-machine-lan", true).Return(lan, nil).Times(1)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(0), &driver.NicIps).Return(nic, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(1), nil).Return(nic, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(2), nil).Return(nic, nil)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateNatPublicIps(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = ""
	driver.IPAddress = testVar
	driver.NicDhcp = true
	driver.PrivateLan = true
	driver.NicIps = []string{"127.0.0.1", "127.0.0.3"}
	driver.NatPublicIps = []string{"127.0.0.4"}
	driver.CreateNat = true

	driver.NatLansToGateways = map[string][]string{"1": {"127.0.0.3"}}

	attached_volumes := sdkgo.NewAttachedVolumesWithDefaults()
	attached_volumes.Items = &[]sdkgo.Volume{*volume}
	server.Entities = sdkgo.NewServerEntitiesWithDefaults()
	server.Entities.SetVolumes(*attached_volumes)

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, "docker-machine-lan", false).Return(lan, nil).Times(1)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateIpBlock("SHOULD_NOT_CALL", "SHOULD_NOT_CALL")
	clientMock.EXPECT().GetIpBlockIps("SHOULD_NOT_CALL")
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(0), &driver.NicIps).Return(nic, nil)
	clientMock.EXPECT().CreateNat(
		driver.DatacenterId, "docker-machine-nat", driver.NatPublicIps, driver.NatFlowlogs, driver.NatRules, driver.NatLansToGateways,
		net.ParseIP((driver.NicIps)[0]).Mask(net.CIDRMask(24, 32)).String()+"/24", driver.SkipDefaultNatRules,
	).Return(nat, nil)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateNat(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = ""
	driver.IPAddress = testVar
	driver.NicDhcp = true
	driver.SkipDefaultNatRules = true
	driver.PrivateLan = true
	driver.CreateNat = true
	driver.NatFlowlogs = []string{"test_name:ACCEPTED:INGRESS:test_bucket", "test_name2:REGECTED:EGRESS:test_bucket"}
	driver.NatRules = []string{
		"name1:SNAT:TCP::10.0.1.0/24:10.0.2.0/24:100:500",
		"name2:SNAT:ALL::10.0.1.0/24::1023:1500",
	}

	driver.NatLansToGateways = map[string][]string{"1": {"127.0.0.3"}}

	attached_volumes := sdkgo.NewAttachedVolumesWithDefaults()
	attached_volumes.Items = &[]sdkgo.Volume{*volume}
	server.Entities = sdkgo.NewServerEntitiesWithDefaults()
	server.Entities.SetVolumes(*attached_volumes)

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(privateLan, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, "docker-machine-lan", false).Return(lan, nil).Times(1)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(0), nil).Return(nic, nil)
	clientMock.EXPECT().CreateNat(
		driver.DatacenterId, "docker-machine-nat", ips, driver.NatFlowlogs, driver.NatRules, driver.NatLansToGateways,
		net.ParseIP(([]string{"127.0.0.1"})[0]).Mask(net.CIDRMask(24, 32)).String()+"/24", driver.SkipDefaultNatRules,
	).Return(nat, nil)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateImageId(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.Image = testImageIdVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(testImageIdVar).Return(&imageFoundById, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, propertiesImageFoundById).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, false, int32(0), &ips).Return(nic, nil)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateImageAlias(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.Image = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, propertiesImageAlias).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, false, int32(0), &ips).Return(nic, nil)
	err := driver.Create()
	assert.NoError(t, err)
}

func TestCreateIpBlockErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
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
	clientMock.EXPECT().GetServer(driver.DatacenterId, "test").Return(server, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateGetImageErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, testErr)
	err := driver.Create()
	assert.Error(t, err)
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, testErr)
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateGetDatacenterErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
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
	driver.DatacenterId = testVar
	driver.LanId = ""
	driver.IpBlockId = testVar
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
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
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
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
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
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, propertiesImageId).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateAttachNicErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = testVar
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, propertiesImageId).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, false, int32(0), &ips).Return(nic, testErr)
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(testErr)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(testErr)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(testErr)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(testErr)
	clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(testErr)
	err := driver.Create()
	assert.Error(t, err)
}

func TestCreateAttachAdditionalNicErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.IpBlockId = testVar
	driver.AdditionalLansIds = []int{2, 4}
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLan(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("lan")).Return(lan1, nil)
	clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(driver.CloudInit, nil)
	clientMock.EXPECT().GetNic(gomock.AssignableToTypeOf("dc"), gomock.AssignableToTypeOf("sv"), gomock.AssignableToTypeOf("nic")).Return(nic, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, gomock.AssignableToTypeOf(sdkgo.Server{})).Return(server, nil)
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, propertiesImageId).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, false, int32(0), &ips).Return(nic, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(2), nil).Return(nic, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(4), nil).Return(nic, testErr)
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(testErr)
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
	driver.DatacenterId = testVar
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
	driver.DatacenterId = testVar
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
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	err := driver.Start()
	assert.Error(t, err)
}

func TestStart(t *testing.T) {
	s := serverWithState(testVar, "PAUSED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	clientMock.EXPECT().StartServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Start()
	assert.NoError(t, err)
}

func TestStartServerErr(t *testing.T) {
	s := serverWithState(testVar, "INACTIVE")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	clientMock.EXPECT().StartServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error starting server"))
	err := driver.Start()
	assert.Error(t, err)
}

func TestStartRunningServer(t *testing.T) {
	s := serverWithState(testVar, "RUNNING")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	err := driver.Start()
	assert.NoError(t, err)
}

func TestStopErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	s := &sdkgo.Server{}
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	err := driver.Stop()
	assert.Error(t, err)
}

func TestStop(t *testing.T) {
	s := serverWithState(testVar, "SHUTOFF")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Stop()
	assert.NoError(t, err)
}

func TestStopServerErr(t *testing.T) {
	s := serverWithState(testVar, "PAUSED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error stoping server"))
	err := driver.Stop()
	assert.Error(t, err)
}

func TestStopStoppedServer(t *testing.T) {
	s := serverWithState(testVar, "BLOCKED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	err := driver.Stop()
	assert.NoError(t, err)
}

func TestRestartErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().RestartServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error restarting server"))
	err := driver.Restart()
	assert.Error(t, err)
}

func TestRestart(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().RestartServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Restart()
	assert.NoError(t, err)
}

func TestKillErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error stoping server"))
	err := driver.Kill()
	assert.Error(t, err)
}

func TestKill(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err := driver.Kill()
	assert.NoError(t, err)
}

func TestGetSSHHostnameErr(t *testing.T) {
	s := serverWithState(testVar, "CRASHED")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	_, err := driver.GetSSHHostname()
	assert.Error(t, err)
}

func TestGetURLErr(t *testing.T) {
	s := serverWithState(testVar, "SHUTOFF")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	_, err := driver.GetURL()
	assert.Error(t, err)
}

// Muted because IP is now set during Create
//func TestGetURL(t *testing.T) {
//	s := serverWithNicAttached(testVar, "AVAILABLE", testVar)
//	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
//	driver.DatacenterId = testVar
//	driver.ServerId = testVar
//	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil).Times(2)
//	_, err := driver.GetURL()
//	assert.NoError(t, err)
//}

func TestGetIPErr(t *testing.T) {
	s := serverWithState(testVar, "AVAILABLE")
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, testErr)
	_, err := driver.GetIP()
	assert.Error(t, err)
}

// Muted because IP is now set during Create
//func TestGetIP(t *testing.T) {
//	s := serverWithNicAttached(testVar, "AVAILABLE", testVar)
//	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
//	driver.DatacenterId = testVar
//	driver.ServerId = testVar
//	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
//	_, err := driver.GetIP()
//	assert.NoError(t, err)
//}

func TestGetStateErr(t *testing.T) {
	s := serverWithNicAttached(testVar, "AVAILABLE", testVar)
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, testErr)
	_, err := driver.GetState()
	assert.Error(t, err)
}

func TestGetStateShutDown(t *testing.T) {
	s := serverWithNicAttached(testVar, "SHUTDOWN", testVar)
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
	_, err := driver.GetState()
	assert.NoError(t, err)
}

func TestGetStateCrashed(t *testing.T) {
	s := serverWithNicAttached(testVar, "CRASHED", testVar)
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(s, nil)
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
	_, err := driver.getImageId(imageName)
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
