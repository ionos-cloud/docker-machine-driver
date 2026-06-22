package ionoscloud

import (
	"encoding/base64"
	"fmt"
	"net"
	"os"
	"reflect"
	"testing"

	"github.com/ionos-cloud/docker-machine-driver/internal/utils"
	mockutils "github.com/ionos-cloud/docker-machine-driver/internal/utils/mocks"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
	"github.com/rancher/machine/libmachine/drivers"
	"github.com/rancher/machine/libmachine/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	defaultHostName  = "default1"
	defaultStorePath = "path"
)

var (
	// Common variables used
	datacenterName         = "datacenter_name"
	datacenterId           = "datacenter_id"
	datacenterId2          = "datacenter_id2"
	lanId                  = "1"
	serverName             = "server_name"
	serverId               = "server_id"
	volumeName             = "volume_name"
	volumeId               = "volume_id"
	nicId                  = "nic_id"
	testRegion             = "us/ewr"
	testRegion2            = "de/fra/2"
	testRegion2Parent      = "de/fra"
	testVar                = "test"
	testImageIdVar         = "test-image-id"
	locationId             = "las"
	imageType              = "HDD"
	imageAlias             = "ubuntu:latest"
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
	dcDeFra2 = &sdkgo.Datacenter{
		Id: sdkgo.ToPtr(datacenterId),
		Properties: &sdkgo.DatacenterProperties{
			Name:        sdkgo.ToPtr(datacenterName),
			Description: sdkgo.ToPtr("datacenter_description"),
			Location:    sdkgo.ToPtr(testRegion2),
			Version:     sdkgo.ToPtr(int32(1)),
		},
	}
	lan_post = &sdkgo.Lan{
		Id: sdkgo.ToPtr(lanId),
		Properties: &sdkgo.LanProperties{
			Name:   &lanName1,
			Public: sdkgo.ToPtr(true),
		},
	}
	lan_post_private = &sdkgo.Lan{
		Id: sdkgo.ToPtr(lanId),
		Properties: &sdkgo.LanProperties{
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
	lan = &sdkgo.Lan{
		Id:         &testVar,
		Properties: &sdkgo.LanProperties{Public: sdkgo.ToPtr(false), Name: &lanName1},
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

func TestSetConfigFromCustomFlagsAdditionalDisks(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, map[string]interface{}{
		flagAdditionalDisks: []string{"HDD:10", "SSD Premium:13"},
	})
	assert.Equal(t, driver.AdditionalDisks, []DiskProperties{{"HDD", 10}, {"SSD Premium", 13}})

	driver, _ = NewTestDriverFlagsSet(t, map[string]interface{}{
		flagAdditionalDisks: []string{"SSD Standard:123", "SSD:124", "SSD Premium:412"},
	})
	assert.Equal(t, driver.AdditionalDisks, []DiskProperties{{"SSD Standard", 123}, {"SSD", 124}, {"SSD Premium", 412}})

	driver, _ = NewTestDriverFlagsSet(t, map[string]interface{}{
		flagAdditionalDisks: []string{},
	})
	assert.Equal(t, driver.AdditionalDisks, []DiskProperties(nil))
}

func TestSetConfigFromCustomFlagsAdditionalDisksError(t *testing.T) {
	driver, _ := NewTestDriver(t, defaultHostName, defaultStorePath)
	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagAdditionalDisks: []string{"HDD:10qdwq::", "SSD Premium:13"},
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.Equal(t, err.Error(), "invalid additional disk configuration: HDD:10qdwq::, must be \"type:size\"")
	assert.Empty(t, checkFlags.InvalidFlags)

	checkFlags = &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagAdditionalDisks: []string{"wrongDiskType:10", "SSD Premium:13"},
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err = driver.SetConfigFromFlags(checkFlags)
	assert.Equal(t, err.Error(), "invalid additional disk type: wrongDiskType, must be one of [\"HDD\" \"SSD\" \"SSD Standard\" \"SSD Premium\"]")
	assert.Empty(t, checkFlags.InvalidFlags)

	checkFlags = &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagAdditionalDisks: []string{"SSD Standard:notInt", "SSD Premium:13"},
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err = driver.SetConfigFromFlags(checkFlags)
	assert.Equal(t, err.Error(), "invalid additional disk size: notInt, must be an integer")
	assert.Empty(t, checkFlags.InvalidFlags)
}

func TestSetConfigFromCustomFlagsAdditionalNicsDhcp(t *testing.T) {
	driver, _ := NewTestDriverFlagsSet(t, map[string]interface{}{
		flagAdditionalNicsDhcp: []string{"5=false", " 7 = true "},
	})
	assert.Equal(t, map[int]bool{5: false, 7: true}, driver.AdditionalNicsDhcp)

	// An unset flag leaves the map nil so additional NICs keep the default DHCP.
	driver, _ = NewTestDriverFlagsSet(t, map[string]interface{}{})
	assert.Nil(t, driver.AdditionalNicsDhcp)
}

func TestSetConfigFromCustomFlagsAdditionalNicsDhcpError(t *testing.T) {
	driver, _ := NewTestDriver(t, defaultHostName, defaultStorePath)

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagAdditionalNicsDhcp: []string{"5"},
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.Equal(t, "invalid value for ionoscloud-additional-nics-dhcp: \"5\" must be in the form lanId=dhcp (e.g. 5=false)", err.Error())
	assert.Empty(t, checkFlags.InvalidFlags)

	checkFlags = &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagAdditionalNicsDhcp: []string{"notInt=false"},
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err = driver.SetConfigFromFlags(checkFlags)
	assert.Equal(t, "invalid value for ionoscloud-additional-nics-dhcp: \"notInt=false\" must be a numeric LAN id followed by =dhcp (e.g. 5=false)", err.Error())
	assert.Empty(t, checkFlags.InvalidFlags)

	checkFlags = &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagAdditionalNicsDhcp: []string{"5=notBool"},
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err = driver.SetConfigFromFlags(checkFlags)
	assert.Equal(t, "invalid value for ionoscloud-additional-nics-dhcp: \"5=notBool\" must have a boolean DHCP value (true/false)", err.Error())
	assert.Empty(t, checkFlags.InvalidFlags)
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
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(nil, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&lans, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
}

func TestPreCreateLocationCutSuccess(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.Location = "de/fra"
	driver.DatacenterId = ""
	clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{}}, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
}

func TestPreCreateLocationCutSuccess2(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.Location = testRegion2
	driver.DatacenterId = ""
	clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{}}, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
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

func TestPreCreateServerTypeErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)
	driver.ServerType = "wrong_value"
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&additionalLans, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.Error(t, err)
}

func TestPreCreateServerTypeNoErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authDcIdFlagsSet)

	driver.ServerType = "ENTERPRISE"
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&additionalLans, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)

	driver.ServerType = "CUBE"
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&additionalLans, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err = driver.PreCreateCheck()
	assert.NoError(t, err)

}

func TestPreCreateCheck(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	clientMock.EXPECT().GetDatacenters().Return(dcs, nil)
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
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.True(t, reflect.DeepEqual(driver.AdditionalLansIds, []int{lanId1Int}))
	assert.NoError(t, err)
}

// Regression: when the primary NIC is selected by --ionoscloud-lan-id,
// any names listed in --ionoscloud-additional-lans must still be resolved
// to LAN ids instead of being silently dropped.
func TestPreCreateAdditionalLansResolvedWhenLanIdSet(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = "test"
	driver.LanId = "100"
	driver.AdditionalLans = []string{lanName1, lanName2}
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&additionalLans, nil)
	clientMock.EXPECT().GetLan(driver.DatacenterId, driver.LanId).Return(privateLan, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	assert.Equal(t, "100", driver.LanId, "primary LanId must remain untouched")
	assert.ElementsMatch(t, []int{lanId1Int, 5}, driver.AdditionalLansIds)
}

// --ionoscloud-additional-lans-ids must populate AdditionalLansIds directly
// and merge with any ids resolved from --ionoscloud-additional-lans (no dupes).
func TestPreCreateAdditionalLansIdsFromFlag(t *testing.T) {
	flags := map[string]interface{}{
		flagUsername:          "IONOSCLOUD_USERNAME",
		flagPassword:          "IONOSCLOUD_PASSWORD",
		flagAdditionalLans:    []string{lanName1},
		flagAdditionalLansIds: []string{"2", "7"},
	}
	driver, clientMock := NewTestDriverFlagsSet(t, flags)
	driver.DatacenterId = "test"
	driver.LanId = "100"
	clientMock.EXPECT().GetLans(driver.DatacenterId).Return(&additionalLans, nil)
	clientMock.EXPECT().GetLan(driver.DatacenterId, driver.LanId).Return(privateLan, nil)
	clientMock.EXPECT().GetDatacenter(driver.DatacenterId).Return(dc, nil)
	clientMock.EXPECT().GetImageById(defaultImageAlias).Return(&sdkgo.Image{}, fmt.Errorf("no image found with this id"))
	clientMock.EXPECT().GetImages().Return(&images, nil)
	clientMock.EXPECT().GetNats(driver.DatacenterId).Return(nats, nil)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []int{2, 7}, driver.AdditionalLansIds)
}

// A non-numeric value for --ionoscloud-additional-lans-ids must fail fast
// during flag parsing rather than silently dropping the entry.
func TestSetConfigFromFlagsAdditionalLansIdsInvalid(t *testing.T) {
	driver, _ := NewTestDriver(t, defaultHostName, defaultStorePath)
	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			flagUsername:          "IONOSCLOUD_USERNAME",
			flagPassword:          "IONOSCLOUD_PASSWORD",
			flagAdditionalLansIds: []string{"not-a-number"},
		},
		CreateFlags: driver.GetCreateFlags(),
	}
	err := driver.SetConfigFromFlags(checkFlags)
	assert.Error(t, err)
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),

		clientMock.EXPECT().CreateDatacenter(datacenterName, testRegion).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, true).Return(lan_post, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_post.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
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
				assert.Nil(t, serverToCreate.Properties.NicMultiQueue)
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
	userdataFlag := drivers.DriverUserdataFlag(driver)
	assert.Empty(t, userdataFlag)
}

func TestCreateAdditionalDisks(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}
	driver.AdditionalDisks = []DiskProperties{{"HDD", 10}}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{}}, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),

		clientMock.EXPECT().CreateDatacenter(datacenterName, testRegion).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, true).Return(lan_post, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_post.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
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
				assert.Nil(t, serverToCreate.Properties.NicMultiQueue)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 2)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				assert.Equal(t, "HDD", *volumes[1].Properties.Type)
				assert.Equal(t, fmt.Sprintf("%s-vol-1", driver.MachineName), *volumes[1].Properties.Name)
				assert.Nil(t, volumes[1].Properties.ImagePassword)
				assert.Nil(t, volumes[1].Properties.SshKeys)
				assert.Nil(t, volumes[1].Properties.UserData)
				assert.Nil(t, volumes[1].Properties.Image)
				assert.Equal(t, float32(10), *volumes[1].Properties.Size)
				assert.Nil(t, volumes[1].Properties.AvailabilityZone)
				assert.Nil(t, volumes[1].Properties.ImageAlias)

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
	userdataFlag := drivers.DriverUserdataFlag(driver)
	assert.Empty(t, userdataFlag)
}

func TestCreateAdditionalDisks2(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}
	driver.AdditionalDisks = []DiskProperties{{"HDD", 10}, {"SSD Premium", 11}, {"SSD", 13}}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{}}, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),

		clientMock.EXPECT().CreateDatacenter(datacenterName, testRegion).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, true).Return(lan_post, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_post.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
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
				assert.Nil(t, serverToCreate.Properties.NicMultiQueue)
				volumes := *serverToCreate.Entities.Volumes.Items
				assert.Len(t, volumes, 4)
				assert.Equal(t, driver.DiskType, *volumes[0].Properties.Type)
				assert.Equal(t, driver.MachineName, *volumes[0].Properties.Name)
				assert.Equal(t, driver.ImagePassword, *volumes[0].Properties.ImagePassword)
				assert.Equal(t, []string{testVar}, *volumes[0].Properties.SshKeys)
				assert.Equal(t, base64.StdEncoding.EncodeToString([]byte(cloudInit)), *volumes[0].Properties.UserData)
				assert.Equal(t, testImageIdVar, *volumes[0].Properties.Image)
				assert.Equal(t, float32(driver.DiskSize), *volumes[0].Properties.Size)
				assert.Equal(t, driver.VolumeAvailabilityZone, *volumes[0].Properties.AvailabilityZone)
				assert.Nil(t, volumes[0].Properties.ImageAlias)

				assert.Equal(t, "HDD", *volumes[1].Properties.Type)
				assert.Equal(t, fmt.Sprintf("%s-vol-1", driver.MachineName), *volumes[1].Properties.Name)
				assert.Nil(t, volumes[1].Properties.ImagePassword)
				assert.Nil(t, volumes[1].Properties.SshKeys)
				assert.Nil(t, volumes[1].Properties.UserData)
				assert.Nil(t, volumes[1].Properties.Image)
				assert.Equal(t, float32(10), *volumes[1].Properties.Size)
				assert.Nil(t, volumes[1].Properties.AvailabilityZone)
				assert.Nil(t, volumes[1].Properties.ImageAlias)

				assert.Equal(t, "SSD Premium", *volumes[2].Properties.Type)
				assert.Equal(t, fmt.Sprintf("%s-vol-2", driver.MachineName), *volumes[2].Properties.Name)
				assert.Nil(t, volumes[2].Properties.ImagePassword)
				assert.Nil(t, volumes[2].Properties.SshKeys)
				assert.Nil(t, volumes[2].Properties.UserData)
				assert.Nil(t, volumes[2].Properties.Image)
				assert.Equal(t, float32(11), *volumes[2].Properties.Size)
				assert.Nil(t, volumes[2].Properties.AvailabilityZone)
				assert.Nil(t, volumes[2].Properties.ImageAlias)

				assert.Equal(t, "SSD", *volumes[3].Properties.Type)
				assert.Equal(t, fmt.Sprintf("%s-vol-3", driver.MachineName), *volumes[3].Properties.Name)
				assert.Nil(t, volumes[3].Properties.ImagePassword)
				assert.Nil(t, volumes[3].Properties.SshKeys)
				assert.Nil(t, volumes[3].Properties.UserData)
				assert.Nil(t, volumes[3].Properties.Image)
				assert.Equal(t, float32(13), *volumes[3].Properties.Size)
				assert.Nil(t, volumes[3].Properties.AvailabilityZone)
				assert.Nil(t, volumes[3].Properties.ImageAlias)

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
	userdataFlag := drivers.DriverUserdataFlag(driver)
	assert.Empty(t, userdataFlag)
}
func TestCreateAppendRkeProvisioning(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}
	driver.AppendRKECloudInit = true
	tmpFile, err := os.CreateTemp("", "rke-provision-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	rkeContent := `#cloud-config
runcmd:
- sh /etc/rke.sh
write_files:
- path: /etc/rke.sh
  content: some install content
`
	if _, err := tmpFile.WriteString(rkeContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	driver.RKEProvisionUserData = tmpFile.Name()

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{}}, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),

		clientMock.EXPECT().CreateDatacenter(datacenterName, testRegion).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, true).Return(lan_post, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_post.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return("", nil),
		clientMock.EXPECT().UpdateCloudInitFile("", "runcmd", []interface{}{"sh /etc/rke.sh"}, false, "append").Return("", nil),
		clientMock.EXPECT().UpdateCloudInitFile(
			"",
			"write_files",
			[]interface{}{map[string]interface{}{
				"path": "/etc/rke.sh", "content": "some install content",
			}},
			false,
			"append",
		).Return(cloudInit, nil),

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
	err = driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
	userdataFlag := drivers.DriverUserdataFlag(driver)
	assert.Equal(t, "ionoscloud-rancher-provision-user-data", userdataFlag)
}

func TestCreateDeFra2(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.CpuFamily = "INTEL_SKYLAKE"
	driver.Location = testRegion2
	driver.DatacenterName = datacenterName
	driver.ImagePassword = "<testdata>"
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}
	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{}}, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),

		clientMock.EXPECT().CreateDatacenter(datacenterName, testRegion2).Return(dcDeFra2, nil),
		clientMock.EXPECT().CreateLan(*dcDeFra2.Id, lanName1, true).Return(lan_post, nil),
		clientMock.EXPECT().GetLan(*dcDeFra2.Id, *lan_post.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion2Parent).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dcDeFra2.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, driver.CpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(driver.Ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(driver.Cores), *serverToCreate.Properties.Cores)
				assert.Nil(t, serverToCreate.Properties.NicMultiQueue)
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
		clientMock.EXPECT().GetServer(*dcDeFra2.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dcDeFra2.Id, serverId, nicId).Return(nic, nil),
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().CreateIpBlock(int32(1), testRegion).Return(ipblock, nil),
		clientMock.EXPECT().GetIpBlockIps(ipblock).Return(ipblock.Properties.Ips, nil),
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
	driver.NicMultiQueue = true
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
			NicMultiQueue:    sdkgo.ToPtr(true),
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
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
				assert.True(t, *serverToCreate.Properties.NicMultiQueue)
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

func TestCreateAdditionalNicDhcpDisabled(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterName = datacenterName
	driver.LanName = lanName1
	// Two additional NICs: LAN 5 is overridden to DHCP off while LAN 7 is left to
	// the default (on), so this exercises both the override and the fallback paths.
	driver.AdditionalLansIds = []int{5, 7}
	driver.AdditionalNicsDhcp = map[int]bool{5: false}

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

	srv := &sdkgo.Server{
		Id: sdkgo.ToPtr(serverId),
		Entities: &sdkgo.ServerEntities{
			Volumes: &sdkgo.AttachedVolumes{
				Items: &[]sdkgo.Volume{{Id: sdkgo.ToPtr(volumeId)}},
			},
			Nics: &sdkgo.Nics{
				Items: &[]sdkgo.Nic{
					{
						Id:         sdkgo.ToPtr(nicId),
						Properties: &sdkgo.NicProperties{Name: sdkgo.ToPtr(defaultHostName)},
					},
					{
						Id:         sdkgo.ToPtr("nic_id-2"),
						Properties: &sdkgo.NicProperties{Name: sdkgo.ToPtr("different_name")},
					},
				},
			},
		},
	}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get}}, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 3)
				// Primary NIC keeps following --ionoscloud-nic-dhcp.
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicDhcp, *nics[0].Properties.Dhcp)
				// Additional NICs follow the order of AdditionalLansIds ([5, 7]).
				// LAN 5 picks up the per-NIC override (off)...
				assert.Equal(t, int32(5), *nics[1].Properties.Lan)
				assert.Equal(t, false, *nics[1].Properties.Dhcp)
				// ...while LAN 7, absent from the map, keeps the default (on).
				assert.Equal(t, int32(7), *nics[2].Properties.Lan)
				assert.Equal(t, true, *nics[2].Properties.Dhcp)
				serverToCreate.Id = &serverId
				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(srv, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)
}

// A DHCP override whose LAN id is not among the attached additional NICs must be
// ignored (no NIC is created for it) and surfaced as a warning, while the valid
// overrides still apply. This drives the warn-and-ignore branch in CreateIonosServer.
func TestCreateAdditionalNicDhcpOverrideIgnoredForUnattachedLan(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterName = datacenterName
	driver.LanName = lanName1
	// LAN 5 is overridden off, LAN 7 keeps the default on, and LAN 99 is an override
	// for a LAN that is not attached: it must yield no NIC and a warning.
	driver.AdditionalLansIds = []int{5, 7}
	driver.AdditionalNicsDhcp = map[int]bool{5: false, 99: true}

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

	srv := &sdkgo.Server{
		Id: sdkgo.ToPtr(serverId),
		Entities: &sdkgo.ServerEntities{
			Volumes: &sdkgo.AttachedVolumes{
				Items: &[]sdkgo.Volume{{Id: sdkgo.ToPtr(volumeId)}},
			},
			Nics: &sdkgo.Nics{
				Items: &[]sdkgo.Nic{
					{
						Id:         sdkgo.ToPtr(nicId),
						Properties: &sdkgo.NicProperties{Name: sdkgo.ToPtr(defaultHostName)},
					},
					{
						Id:         sdkgo.ToPtr("nic_id-2"),
						Properties: &sdkgo.NicProperties{Name: sdkgo.ToPtr("different_name")},
					},
				},
			},
		},
	}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get}}, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				nics := *serverToCreate.Entities.Nics.Items
				// Only the primary plus the two attached additional LANs; LAN 99 produces no NIC.
				assert.Len(t, nics, 3)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicDhcp, *nics[0].Properties.Dhcp)
				assert.Equal(t, int32(5), *nics[1].Properties.Lan)
				assert.Equal(t, false, *nics[1].Properties.Dhcp)
				assert.Equal(t, int32(7), *nics[2].Properties.Lan)
				assert.Equal(t, true, *nics[2].Properties.Dhcp)
				for _, n := range nics {
					assert.NotEqual(t, int32(99), *n.Properties.Lan, "override for unattached LAN 99 must not create a NIC")
				}
				serverToCreate.Id = &serverId
				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(srv, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)

	// The override for the unattached LAN 99 is reported and dropped.
	assert.Contains(t, log.History(),
		fmt.Sprintf("%s: ignoring DHCP override for LAN id %d, no additional NIC is attached to that LAN", flagAdditionalNicsDhcp, 99))
}

// End-to-end: an additional LAN attached by name via --ionoscloud-additional-lans is
// resolved to its numeric id during PreCreateCheck, and a DHCP override keyed by that
// resolved id (--ionoscloud-additional-nics-dhcp=5=false) is applied to the resulting
// NIC. Exercises the full flag-parse -> name-resolution -> NIC-build chain that the
// docs promise ("a LAN attached by name must be keyed by its resolved ID").
func TestCreateAdditionalNicDhcpOverrideForLanAttachedByName(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, map[string]interface{}{
		flagUsername:           "IONOSCLOUD_USERNAME",
		flagPassword:           "IONOSCLOUD_PASSWORD",
		flagAdditionalLans:     []string{lanName2},
		flagAdditionalNicsDhcp: []string{lanId2 + "=false"},
	})
	// Parsing populates the override map keyed by the numeric LAN id.
	assert.Equal(t, map[int]bool{5: false}, driver.AdditionalNicsDhcp)

	driver.SSHKey = testVar
	driver.DatacenterName = datacenterName
	driver.LanName = lanName1
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

	srv := &sdkgo.Server{
		Id: sdkgo.ToPtr(serverId),
		Entities: &sdkgo.ServerEntities{
			Volumes: &sdkgo.AttachedVolumes{
				Items: &[]sdkgo.Volume{{Id: sdkgo.ToPtr(volumeId)}},
			},
			Nics: &sdkgo.Nics{
				Items: &[]sdkgo.Nic{
					{
						Id:         sdkgo.ToPtr(nicId),
						Properties: &sdkgo.NicProperties{Name: sdkgo.ToPtr(defaultHostName)},
					},
					{
						Id:         sdkgo.ToPtr("nic_id-2"),
						Properties: &sdkgo.NicProperties{Name: sdkgo.ToPtr("different_name")},
					},
				},
			},
		},
	}

	gomock.InOrder(
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dc}}, nil),
		// lanName1 resolves to the primary LAN 1; lanName2 resolves to additional LAN 5.
		clientMock.EXPECT().GetLans(*dc.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get, *privateLan3}}, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				nics := *serverToCreate.Entities.Nics.Items
				assert.Len(t, nics, 2)
				assert.Equal(t, int32(1), *nics[0].Properties.Lan)
				assert.Equal(t, nicDhcp, *nics[0].Properties.Dhcp)
				// The LAN attached by name resolved to id 5 and picked up the 5=false override.
				assert.Equal(t, int32(5), *nics[1].Properties.Lan)
				assert.Equal(t, false, *nics[1].Properties.Dhcp)
				serverToCreate.Id = &serverId
				return &serverToCreate, nil
			}),
		clientMock.EXPECT().GetServer(*dc.Id, serverId, int32(2)).Return(srv, nil),
		clientMock.EXPECT().GetNic(*dc.Id, serverId, nicId).Return(nic, nil),
	)
	err := driver.PreCreateCheck()
	assert.NoError(t, err)
	err = driver.Create()
	assert.NoError(t, err)

	// PreCreateCheck resolved the by-name LAN into the numeric id the override is keyed by.
	assert.Contains(t, driver.AdditionalLansIds, 5)
}

func TestCreatePropertiesSetDeFra2(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterName = datacenterName
	driver.LanName = lanName1
	driver.AdditionalLans = []string{lanName1, lanName2}

	driver.CpuFamily = cpuFamily
	driver.Location = testRegion2
	driver.ImagePassword = imagePassword
	driver.NicDhcp = nicDhcp
	driver.NicIps = nicIps
	driver.DiskType = diskType
	driver.NicMultiQueue = false
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
			NicMultiQueue:    sdkgo.ToPtr(false),
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
		clientMock.EXPECT().GetDatacenters().Return(&sdkgo.Datacenters{Items: &[]sdkgo.Datacenter{*dcDeFra2}}, nil),
		clientMock.EXPECT().GetLans(*dcDeFra2.Id).Return(&sdkgo.Lans{Items: &[]sdkgo.Lan{*lan_get, *lan_get2}}, nil),
		clientMock.EXPECT().GetLan(*dcDeFra2.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetDatacenter(*dcDeFra2.Id).Return(dcDeFra2, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dcDeFra2.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dcDeFra2.Id).Return(dcDeFra2, nil),
		clientMock.EXPECT().GetLan(*dcDeFra2.Id, *lan_get.Id).Return(lan_get, nil),
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dcDeFra2.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Equal(t, cpuFamily, *serverToCreate.Properties.CpuFamily)
				assert.Equal(t, int32(ram), *serverToCreate.Properties.Ram)
				assert.Equal(t, int32(cores), *serverToCreate.Properties.Cores)
				assert.Equal(t, serverAvailabilityZone, *serverToCreate.Properties.AvailabilityZone)
				assert.Nil(t, serverToCreate.Properties.Type)
				assert.Nil(t, serverToCreate.Properties.TemplateUuid)
				assert.Nil(t, serverToCreate.Properties.NicMultiQueue)
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
		clientMock.EXPECT().GetServer(*dcDeFra2.Id, serverId, int32(2)).Return(server, nil),
		clientMock.EXPECT().GetNic(*dcDeFra2.Id, serverId, nicId).Return(nic, nil),
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get.Id).Return(lan_get, nil),
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, false).Return(lan_post_private, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, false).Return(lan_post_private, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
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
		clientMock.EXPECT().GetImageById(imageAlias).Return(&sdkgo.Image{Id: sdkgo.ToPtr(testImageIdVar)}, nil),
		clientMock.EXPECT().GetNats(*dc.Id).Return(nats, nil),

		clientMock.EXPECT().GetDatacenter(*dc.Id).Return(dc, nil),
		clientMock.EXPECT().CreateLan(*dc.Id, lanName1, false).Return(lan_post_private, nil),
		clientMock.EXPECT().GetLan(*dc.Id, *lan_get_private.Id).Return(lan_get_private, nil),
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
		clientMock.EXPECT().GetImageById(driver.Image).Return(&imageFoundById, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "users", test, false, "append").Return("test_string", nil),
		clientMock.EXPECT().UpdateCloudInitFile("test_string", "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)

				assert.Nil(t, serverToCreate.Properties.CpuFamily)
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
		clientMock.EXPECT().GetImageById(driver.Image).Return(nil, fmt.Errorf("no image found with this id")),
		clientMock.EXPECT().GetImages().Return(&sdkgo.Images{Items: &[]sdkgo.Image{}}, nil),
		clientMock.EXPECT().UpdateCloudInitFile(driver.CloudInit, "users", test, false, "append").Return("test_string", nil),
		clientMock.EXPECT().UpdateCloudInitFile("test_string", "hostname", []interface{}{driver.MachineName}, true, "skip").Return(cloudInit, nil),
		clientMock.EXPECT().CreateServer(*dc.Id, gomock.AssignableToTypeOf(sdkgo.Server{})).DoAndReturn(
			func(datacenterId string, serverToCreate sdkgo.Server) (*sdkgo.Server, error) {
				assert.Equal(t, driver.MachineName, *serverToCreate.Properties.Name)
				assert.Nil(t, serverToCreate.Properties.CpuFamily)
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

func TestCreateGetDatacenterErr(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.SSHKey = testVar
	driver.DatacenterId = datacenterId
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

	gomock.InOrder(
		clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(nil),
		clientMock.EXPECT().RemoveVolume(driver.DatacenterId, driver.VolumeId).Return(nil),
		clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(nil),
		clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(nil),
		clientMock.EXPECT().RemoveDatacenter(driver.DatacenterId).Return(nil),
		clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(nil),
	)
	err := driver.Remove()
	assert.NoError(t, err)
}

func TestRemoveCube(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
	driver.DatacenterId = datacenterId
	driver.ServerId = testVar
	driver.NicId = testVar
	driver.VolumeId = testVar
	driver.LanId = testVar
	driver.IPAddress = testVar
	driver.ServerType = "CUBE"
	driver.DCExists = false

	gomock.InOrder(
		clientMock.EXPECT().RemoveNic(driver.DatacenterId, driver.ServerId, driver.NicId).Return(nil),
		clientMock.EXPECT().RemoveServer(driver.DatacenterId, driver.ServerId).Return(nil),
		clientMock.EXPECT().RemoveLan(driver.DatacenterId, driver.LanId).Return(nil),
		clientMock.EXPECT().RemoveDatacenter(driver.DatacenterId).Return(nil),
		clientMock.EXPECT().RemoveIpBlock(driver.IpBlockId).Return(nil),
	)
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

func TestGetImageId(t *testing.T) {
	driver, clientMock := NewTestDriverFlagsSet(t, authFlagsSet)
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
