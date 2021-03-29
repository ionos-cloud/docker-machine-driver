package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/docker-machine-driver/utils"
	mockutils "github.com/ionos-cloud/docker-machine-driver/utils/mocks"
	sdkgo "github.com/ionos-cloud/sdk-go/v5"
	"github.com/stretchr/testify/assert"
)

const (
	defaultHostName  = "default"
	defaultStorePath = "path"
)

var (
	testRegion = "us/ewr"
	dcVersion  = int32(1)
	dc         = &sdkgo.Datacenter{
		Id: &testVar,
		Properties: &sdkgo.DatacenterProperties{
			Name:        &testVar,
			Description: &testVar,
			Location:    &testRegion,
			Version:     &dcVersion,
		},
	}
	testVar    = "test"
	locationId = "las"
	location   = &sdkgo.Location{
		Id: &locationId,
		Properties: &sdkgo.LocationProperties{
			ImageAliases: &[]string{testVar},
		},
	}
	images = sdkgo.Images{
		Items: &[]sdkgo.Image{
			{
				Properties: &sdkgo.ImageProperties{
					Name: &testVar,
				},
			},
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
	lan = &sdkgo.LanPost{
		Id: &testVar,
	}
	server = &sdkgo.Server{
		Id: &testVar,
	}
	volume = &sdkgo.Volume{
		Id: &testVar,
	}
	nic = &sdkgo.Nic{
		Id: &testVar,
	}
	ips        = []string{testVar}
	properties = &utils.ClientVolumeProperties{
		DiskType:      defaultDiskType,
		Name:          defaultHostName,
		ImageAlias:    "",
		ImagePassword: defaultImagePassword,
		Zone:          defaultAvailabilityZone,
		SshKey:        testVar,
		DiskSize:      float32(50),
	}
)

func NewTestDriver(ctrl *gomock.Controller, hostName, storePath string) (*Driver, *mockutils.MockClientService) {
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, _ := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	sshPort, err := driver.GetSSHPort()
	assert.Equal(t, 22, sshPort)
	assert.NoError(t, err)
	assert.Equal(t, "", driver.Username)
	assert.Equal(t, "", driver.Password)
	assert.Equal(t, defaultApiEndpoint, driver.URL)
	assert.Equal(t, 4, driver.Cores)
	assert.Equal(t, 2048, driver.Ram)
	assert.Equal(t, defaultRegion, driver.Location)
	assert.Equal(t, defaultDiskType, driver.DiskType)
	assert.Equal(t, 50, driver.DiskSize)
	assert.Equal(t, "", driver.DatacenterId)
	assert.Equal(t, defaultAvailabilityZone, driver.VolumeAvailabilityZone)
	assert.Equal(t, defaultAvailabilityZone, driver.ServerAvailabilityZone)
}

func TestSetConfigFromCustomFlags(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, _ := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagServerRam: 1024,
			flagDiskType:  "SSD",
			flagEndpoint:  "",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	sshPort, err := driver.GetSSHPort()
	assert.Equal(t, 22, sshPort)
	assert.NoError(t, err)
	assert.Equal(t, "", driver.Username)
	assert.Equal(t, "", driver.Password)
	assert.Equal(t, defaultApiEndpoint, driver.URL)
	assert.Equal(t, 4, driver.Cores)
	assert.Equal(t, 1024, driver.Ram)
	assert.Equal(t, defaultRegion, driver.Location)
	assert.Equal(t, "SSD", driver.DiskType)
	assert.Equal(t, 50, driver.DiskSize)
	assert.Equal(t, "", driver.DatacenterId)
	assert.Equal(t, defaultAvailabilityZone, driver.VolumeAvailabilityZone)
	assert.Equal(t, defaultAvailabilityZone, driver.ServerAvailabilityZone)
}

func TestDriverName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, _ := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	assert.Equal(t, driverName, driver.DriverName())
}

func TestPreCreateCheckUserNameErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, _ := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	err = driver.PreCreateCheck()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "please provide username as parameter --ionoscloud-username or as environment variable $IONOSCLOUD_USERNAME")
}

func TestPreCreateCheckPasswordErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, _ := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	err = driver.PreCreateCheck()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "please provide password as parameter --ionoscloud-password or as environment variable $IONOSCLOUD_PASSWORD")
}

func TestPreCreateCheckDataCenterIdErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername:     "IONOSCLOUD_USERNAME",
			flagPassword:     "IONOSCLOUD_PASSWORD",
			flagDatacenterId: "IONOSCLOUD_DATACENTER_ID",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	err = driver.PreCreateCheck()
	assert.NoError(t, err)
}

func TestPreCreateCheckDataCenterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername:     "IONOSCLOUD_USERNAME",
			flagPassword:     "IONOSCLOUD_PASSWORD",
			flagDatacenterId: "IONOSCLOUD_DATACENTER_ID",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(nil, fmt.Errorf("error getting datacenter: 404 not found"))
	err = driver.PreCreateCheck()
	assert.Error(t, err)
}

func TestPreCreateImageIdErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername:     "IONOSCLOUD_USERNAME",
			flagPassword:     "IONOSCLOUD_PASSWORD",
			flagDatacenterId: "IONOSCLOUD_DATACENTER_ID",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetLocationById("us", "ewr").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, fmt.Errorf("error getting image: 404 not found"))
	err = driver.PreCreateCheck()
	assert.Error(t, err)
}

func TestPreCreateCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	err = driver.PreCreateCheck()
	assert.NoError(t, err)
}

func TestCreateSSHKeyErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, _ := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = ""
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, _ := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	err = driver.Create()
	assert.Error(t, err)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, driver.MachineName, true).Return(lan, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, driver.Location, driver.MachineName, driver.CpuFamily, driver.ServerAvailabilityZone, int32(driver.Ram), int32(driver.Cores)).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, properties).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(0), &ips).Return(nic, nil)

	err = driver.Create()
	assert.NoError(t, err)
}

func TestCreateIpBlockErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar

	driver.UseAlias = true
	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateGetImageErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateGetDatacenterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateDatacenterErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = ""

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().CreateDatacenter(driver.MachineName, driver.Location).Return(dc, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateLanErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.LanId = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, driver.MachineName, true).Return(lan, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateServerErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, driver.MachineName, true).Return(lan, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, driver.Location, driver.MachineName, driver.CpuFamily, driver.ServerAvailabilityZone, int32(driver.Ram), int32(driver.Cores)).Return(server, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateAttachVolumeErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, driver.MachineName, true).Return(lan, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, driver.Location, driver.MachineName, driver.CpuFamily, driver.ServerAvailabilityZone, int32(driver.Ram), int32(driver.Cores)).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, properties).Return(volume, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateGetIpBlockErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, driver.MachineName, true).Return(lan, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, driver.Location, driver.MachineName, driver.CpuFamily, driver.ServerAvailabilityZone, int32(driver.Ram), int32(driver.Cores)).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, properties).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestCreateAttachNicErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.SSHKey = testVar
	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar

	clientMock.EXPECT().GetLocationById("us", "las").Return(location, nil)
	clientMock.EXPECT().GetImages().Return(images, nil)
	clientMock.EXPECT().CreateIpBlock(int32(1), driver.Location).Return(ipblock, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().CreateLan(driver.DatacenterId, driver.MachineName, true).Return(lan, nil)
	clientMock.EXPECT().CreateServer(driver.DatacenterId, driver.Location, driver.MachineName, driver.CpuFamily, driver.ServerAvailabilityZone, int32(driver.Ram), int32(driver.Cores)).Return(server, nil)
	clientMock.EXPECT().CreateAttachVolume(driver.DatacenterId, driver.ServerId, properties).Return(volume, nil)
	clientMock.EXPECT().GetIpBlockIps(ipblock).Return(&ips, nil)
	clientMock.EXPECT().CreateAttachNIC(driver.DatacenterId, driver.ServerId, driver.MachineName, true, int32(0), &ips).Return(nic, fmt.Errorf("error"))
	err = driver.Create()
	assert.Error(t, err)
}

func TestRemove(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

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
	clientMock.EXPECT().GetIpBlocks().Return(ipblocks, nil)
	clientMock.EXPECT().RemoveIpBlock(ipblocks, driver.IPAddress).Return(nil)
	err = driver.Remove()
	assert.NoError(t, err)
}

func TestRemoveErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.DCExists = false

	errOccured := fmt.Errorf("error occured")
	clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(errOccured)
	clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(errOccured)
	clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(errOccured)
	clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(errOccured)
	clientMock.EXPECT().RemoveDatacenter(driver.DatacenterId).Return(errOccured)
	clientMock.EXPECT().GetIpBlocks().Return(ipblocks, errOccured)
	clientMock.EXPECT().RemoveIpBlock(ipblocks, driver.IPAddress).Return(errOccured)
	err = driver.Remove()
	assert.Error(t, err)
}

func TestStartErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	err = driver.Start()
	assert.Error(t, err)
}

func TestStart(t *testing.T) {
	var (
		state  = "PAUSED"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().StartServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err = driver.Start()
	assert.NoError(t, err)
}

func TestStartServerErr(t *testing.T) {
	var (
		state  = "INACTIVE"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().StartServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error starting server"))
	err = driver.Start()
	assert.Error(t, err)
}

func TestStartRunningServer(t *testing.T) {
	var (
		state  = "AVAILABLE"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	err = driver.Start()
	assert.NoError(t, err)
}

func TestStopErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	server := &sdkgo.Server{}
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	err = driver.Stop()
	assert.Error(t, err)
}

func TestStop(t *testing.T) {
	var (
		state  = "NOSTATE"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err = driver.Stop()
	assert.NoError(t, err)
}

func TestStopServerErr(t *testing.T) {
	var (
		state  = "PAUSED"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error stoping server"))
	err = driver.Stop()
	assert.Error(t, err)
}

func TestStopStoppedServer(t *testing.T) {
	var (
		state  = "BLOCKED"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	err = driver.Stop()
	assert.NoError(t, err)
}

func TestRestartErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().RestartServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error restarting server"))
	err = driver.Restart()
	assert.Error(t, err)
}

func TestRestart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().RestartServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err = driver.Restart()
	assert.NoError(t, err)
}

func TestKillErr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(fmt.Errorf("error stoping server"))
	err = driver.Kill()
	assert.Error(t, err)
}

func TestKill(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().StopServer(driver.DatacenterId, driver.ServerId).Return(nil)
	err = driver.Kill()
	assert.NoError(t, err)
}

func TestGetSSHHostnameErr(t *testing.T) {
	var (
		state  = "CRASHED"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	_, err = driver.GetSSHHostname()
	assert.Error(t, err)
}

func TestGetURLErr(t *testing.T) {
	var (
		state  = "SHUTOFF"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	_, err = driver.GetURL()
	assert.Error(t, err)
}

func TestGetURL(t *testing.T) {
	var (
		state  = "AVAILABLE"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil).Times(2)
	_, err = driver.GetURL()
	assert.Error(t, err)
}

func TestGetIPErr(t *testing.T) {
	var (
		state  = "AVAILABLE"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, fmt.Errorf("error"))
	_, err = driver.GetIP()
	assert.Error(t, err)
}

func TestGetIP(t *testing.T) {
	var (
		state  = "AVAILABLE"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
			Entities: &sdkgo.ServerEntities{
				Nics: &sdkgo.Nics{
					Items: &[]sdkgo.Nic{
						{
							Properties: &sdkgo.NicProperties{
								Ips: &[]string{testVar},
							},
						},
					},
				},
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	_, err = driver.GetIP()
	assert.NoError(t, err)
}

func TestGetStateErr(t *testing.T) {
	var (
		state  = "AVAILABLE"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
			Entities: &sdkgo.ServerEntities{
				Nics: &sdkgo.Nics{
					Items: &[]sdkgo.Nic{
						{
							Properties: &sdkgo.NicProperties{
								Ips: &[]string{testVar},
							},
						},
					},
				},
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, fmt.Errorf("error"))
	_, err = driver.GetState()
	assert.Error(t, err)
}

func TestGetStateShutDown(t *testing.T) {
	var (
		state  = "SHUTDOWN"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
			Entities: &sdkgo.ServerEntities{
				Nics: &sdkgo.Nics{
					Items: &[]sdkgo.Nic{
						{
							Properties: &sdkgo.NicProperties{
								Ips: &[]string{testVar},
							},
						},
					},
				},
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	_, err = driver.GetState()
	assert.NoError(t, err)
}

func TestGetStateCrashed(t *testing.T) {
	var (
		state  = "CRASHED"
		server = &sdkgo.Server{
			Id: &testVar,
			Metadata: &sdkgo.DatacenterElementMetadata{
				State: &state,
			},
			Entities: &sdkgo.ServerEntities{
				Nics: &sdkgo.Nics{
					Items: &[]sdkgo.Nic{
						{
							Properties: &sdkgo.NicProperties{
								Ips: &[]string{testVar},
							},
						},
					},
				},
			},
		}
	)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	driver, clientMock := NewTestDriver(ctrl, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername: "IONOSCLOUD_USERNAME",
			flagPassword: "IONOSCLOUD_PASSWORD",
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)

	driver.DatacenterId = testVar
	driver.ServerId = testVar
	clientMock.EXPECT().GetServer(driver.DatacenterId, driver.ServerId).Return(server, nil)
	_, err = driver.GetState()
	assert.NoError(t, err)
}
