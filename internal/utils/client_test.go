package utils

import (
	"context"
	"testing"

	sdkgo "github.com/ionos-cloud/sdk-go/v6"
	"github.com/stretchr/testify/assert"
)

// These tests check if Client returns error on wrong credentials

var (
	testLocation = "us/las"
	testName     = "test-name"
	testIpAddr   = "x.x.x.x"
	ipBlock      = &sdkgo.IpBlock{
		Id: &testName,
		Properties: &sdkgo.IpBlockProperties{
			Ips:      &[]string{testIpAddr},
			Location: &testLocation,
		},
	}
	testVolumePropertiesImageId = &ClientVolumeProperties{
		DiskType:      testName,
		Name:          testName,
		ImageId:       testName,
		ImagePassword: testName,
		Zone:          testName,
		SshKey:        testName,
		DiskSize:      0,
	}
	ipBlocks = &sdkgo.IpBlocks{
		Items: &[]sdkgo.IpBlock{
			*ipBlock,
		},
	}
)

func TestClientNew(t *testing.T) {
	New(context.Background(), testName, testName, testName, testName, testName)
}

func TestClientCreateIpBlockErr(t *testing.T) {
	_, err := getTestClient().CreateIpBlock(1, testLocation)
	assert.Error(t, err)
}

func TestClientGetIpBlockIpsErr(t *testing.T) {
	_, err := getTestClient().GetIpBlockIps(&sdkgo.IpBlock{})
	assert.Error(t, err)
}

func TestClientGetIpBlockIps(t *testing.T) {
	_, err := getTestClient().GetIpBlockIps(ipBlock)
	assert.NoError(t, err)
}

func TestClientRemoveIpBlockErr(t *testing.T) {
	err := getTestClient().RemoveIpBlock(testName)
	assert.Error(t, err)
}

func TestClientCreateDatacenterErr(t *testing.T) {
	_, err := getTestClient().CreateDatacenter(testName, testLocation)
	assert.Error(t, err)
}

func TestClientGetDatacenterErr(t *testing.T) {
	_, err := getTestClient().GetDatacenter(testName)
	assert.Error(t, err)
}

func TestClientRemoveDatacenterErr(t *testing.T) {
	err := getTestClient().RemoveDatacenter(testName)
	assert.Error(t, err)
}

func TestClientCreateLanErr(t *testing.T) {
	_, err := getTestClient().CreateLan(testName, testLocation, true)
	assert.Error(t, err)
}

func TestClientRemoveLanErr(t *testing.T) {
	err := getTestClient().RemoveLan(testName, testName)
	assert.Error(t, err)
}

func TestClientCreateServerErr(t *testing.T) {
	_, err := getTestClient().CreateServer(testName, testLocation, testName, testName, testName, 2048, 1)
	assert.Error(t, err)
}

func TestClientGetServerErr(t *testing.T) {
	_, err := getTestClient().GetServer(testName, testName)
	assert.Error(t, err)
}

func TestClientGetLanErr(t *testing.T) {
	_, err := getTestClient().GetLan(testName, testName)
	assert.Error(t, err)
}

func TestClientStartServerErr(t *testing.T) {
	err := getTestClient().StartServer(testName, testName)
	assert.Error(t, err)
}

func TestClientStopServerErr(t *testing.T) {
	err := getTestClient().StopServer(testName, testName)
	assert.Error(t, err)
}

func TestClientRemoveServerErr(t *testing.T) {
	err := getTestClient().RemoveServer(testName, testName)
	assert.Error(t, err)
}

func TestClientCreateVolumeErr(t *testing.T) {
	_, err := getTestClient().CreateAttachVolume(testName, testName, &ClientVolumeProperties{})
	assert.Error(t, err)
}

func TestClientCreateVolumeImageIdErr(t *testing.T) {
	_, err := getTestClient().CreateAttachVolume(testName, testName, testVolumePropertiesImageId)
	assert.Error(t, err)
}

func TestClientRemoveVolumeErr(t *testing.T) {
	err := getTestClient().RemoveVolume(testName, testName)
	assert.Error(t, err)
}

func TestClientCreateNicErr(t *testing.T) {
	_, err := getTestClient().CreateAttachNIC(testName, testName, testName, true, 1, &[]string{testIpAddr})
	assert.Error(t, err)
}

func TestClientRemoveNicErr(t *testing.T) {
	err := getTestClient().RemoveNic(testName, testName, testName)
	assert.Error(t, err)
}

func TestClientGetImagesErr(t *testing.T) {
	_, err := getTestClient().GetImages()
	assert.Error(t, err)
}

func TestClientGetImageByIdErr(t *testing.T) {
	_, err := getTestClient().GetImageById(testName)
	assert.Error(t, err)
}

func TestClientGetLocationByIdErr(t *testing.T) {
	_, err := getTestClient().GetLocationById("us", "las")
	assert.Error(t, err)
}

func TestClientWaitTillProvisioned(t *testing.T) {
	err := getTestClient().waitTillProvisioned("https://api.ionos.com/cloudapi/v6/status/requests")
	assert.Error(t, err)
}

func getTestClient() *Client {
	return &Client{
		APIClient: sdkgo.NewAPIClient(&sdkgo.Configuration{
			Username: "test",
			Password: "test@ionos.com",
			Servers: sdkgo.ServerConfigurations{
				sdkgo.ServerConfiguration{
					URL: "https://api.ionos.com/cloudapi/v6",
				},
			}}),
		ctx: context.TODO(),
	}
}
