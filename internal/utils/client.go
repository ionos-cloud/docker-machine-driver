package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ionos-cloud/docker-machine-driver/pkg/sdk_utils"
	"gopkg.in/yaml.v3"

	"github.com/docker/machine/libmachine/log"
	sdkgo "github.com/ionos-cloud/sdk-go/v6"
)

const waitCount = 1000

type Client struct {
	*sdkgo.APIClient
	ctx context.Context
}

type ClientVolumeProperties struct {
	DiskType      string
	Name          string
	ImageId       string
	ImageAlias    string
	ImagePassword string
	Zone          string
	SshKey        string
	UserData      string
	DiskSize      float32
}

func New(ctx context.Context, name, password, token, url, httpUserAgent string) ClientService {
	clientConfig := sdkgo.NewConfiguration(name, password, token, url)
	clientConfig.UserAgent = fmt.Sprintf("%v_%v", httpUserAgent, clientConfig.UserAgent)
	return &Client{
		APIClient: sdkgo.NewAPIClient(clientConfig),
		ctx:       ctx,
	}
}

func (c *Client) UpdateCloudInitFile(cloudInitYAML string, key string, values []interface{}, single_value bool) (string, error) {
	var cf map[string]interface{}
	cf = make(map[string]interface{})
	if err := yaml.Unmarshal([]byte(cloudInitYAML), &cf); err != nil {
		return "", err
	}

	if val, ok := cf[key]; ok {
		u := val.([]interface{})
		cf[key] = append(u, values...)
	} else {
		if single_value {
			cf[key] = values[0]
		} else {
			cf[key] = values
		}
	}

	newCf, err := yaml.Marshal(cf)
	if err != nil {
		return "", err
	}
	cloudInitYAML = "#cloud-config\n" + string(newCf)

	return cloudInitYAML, nil
}

func (c *Client) CreateIpBlock(size int32, location string) (*sdkgo.IpBlock, error) {
	ipBlock, ipBlockResp, err := c.IPBlocksApi.IpblocksPost(c.ctx).Ipblock(sdkgo.IpBlock{
		Properties: &sdkgo.IpBlockProperties{
			Location: &location,
			Size:     &size,
		}}).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}

	if err = sdk_utils.SanitizeStatusCode(ipBlockResp.StatusCode, ipBlockResp.Message); err != nil {
		return nil, err
	}

	log.Info("IPBlock Reserved!")

	err = c.waitTillProvisioned(ipBlockResp.Header.Get("location"))
	if err != nil {
		return &ipBlock, err
	}
	return &ipBlock, nil
}

func (c *Client) GetIpBlockIps(ipBlock *sdkgo.IpBlock) (*[]string, error) {
	if ipBlockProp, ok := ipBlock.GetPropertiesOk(); ok && ipBlockProp != nil {
		if ips, ok := ipBlockProp.GetIpsOk(); ok && ips != nil {
			return ips, nil
		}
	}
	return nil, fmt.Errorf("error: ip block ips have nil properties")
}

func (c *Client) RemoveIpBlock(ipBlockId string) error {
	resp, err := c.IPBlocksApi.IpblocksDelete(c.ctx, ipBlockId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return err
	}
	log.Info("IPBlock Deleted")
	return nil
}

func (c *Client) CreateDatacenter(name, location string) (*sdkgo.Datacenter, error) {
	dc, dcResp, err := c.DataCentersApi.DatacentersPost(c.ctx).Datacenter(sdkgo.Datacenter{
		Properties: &sdkgo.DatacenterProperties{
			Name:     &name,
			Location: &location,
		}}).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}

	if err = sdk_utils.SanitizeStatusCode(dcResp.StatusCode, dcResp.Message); err != nil {
		return nil, err
	}
	log.Info("Datacenter created!")

	err = c.waitTillProvisioned(dcResp.Header.Get("location"))
	if err != nil {
		return &dc, err
	}
	return &dc, nil
}

func (c *Client) GetDatacenter(datacenterId string) (*sdkgo.Datacenter, error) {
	datacenter, resp, err := c.DataCentersApi.DatacentersFindById(c.ctx, datacenterId).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	if err = sdk_utils.SanitizeStatusCodeCustom(resp.StatusCode, resp.Message, sdk_utils.CustomStatusCodeMessages.Set(404, "provided UUID does not match any datacenter")); err != nil {
		return nil, err
	}
	return &datacenter, nil
}

func (c *Client) GetDatacenters() (*sdkgo.Datacenters, error) {
	datacenters, _, err := c.DataCentersApi.DatacentersGet(c.ctx).Depth(1).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting datacenter: %v", err)
	}
	return &datacenters, nil
}

func (c *Client) RemoveDatacenter(datacenterId string) error {
	resp, err := c.DataCentersApi.DatacentersDelete(c.ctx, datacenterId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	if err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message); err != nil {
		return err
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("DataCenter Deleted")
	return nil
}

func (c *Client) CreateLan(datacenterId, name string, public bool) (*sdkgo.LanPost, error) {
	lan, lanResp, err := c.LANsApi.DatacentersLansPost(c.ctx, datacenterId).Lan(sdkgo.LanPost{
		Properties: &sdkgo.LanPropertiesPost{
			Name:   &name,
			Public: &public,
		}}).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}

	err = sdk_utils.SanitizeStatusCode(lanResp.StatusCode, lanResp.Message)
	if err != nil {
		return nil, err
	}
	log.Info("LAN Created")

	err = c.waitTillProvisioned(lanResp.Header.Get("location"))
	if err != nil {
		return &lan, err
	}
	return &lan, nil
}

func (c *Client) RemoveLan(datacenterId, lanId string) error {
	resp, err := c.LANsApi.DatacentersLansDelete(c.ctx, datacenterId, lanId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return err
	}
	log.Info("LAN Deleted")

	return c.waitTillProvisioned(resp.Header.Get("location"))
}

func (c *Client) CreateServer(datacenterId string, server sdkgo.Server) (*sdkgo.Server, error) {
	svr, serverResp, err := c.ServersApi.DatacentersServersPost(c.ctx, datacenterId).Server(server).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}

	err = sdk_utils.SanitizeStatusCode(serverResp.StatusCode, serverResp.Message)
	if err != nil {
		return nil, err
	}
	log.Info("Server created!")

	err = c.waitTillProvisioned(serverResp.Header.Get("location"))
	if err != nil {
		return &svr, fmt.Errorf("error waiting until server is created: %w", err)
	}
	return &svr, nil
}

func (c *Client) GetServer(datacenterId, serverId string) (*sdkgo.Server, error) {
	server, resp, err := c.ServersApi.DatacentersServersFindById(c.ctx, datacenterId, serverId).Depth(1).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (c *Client) GetLan(datacenterId, LanId string) (*sdkgo.Lan, error) {
	lan, resp, err := c.LANsApi.DatacentersLansFindById(c.ctx, datacenterId, LanId).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return nil, err
	}
	log.Info("Got existing LAN!")
	return &lan, nil
}

func (c *Client) GetLans(datacenterId string) (*sdkgo.Lans, error) {
	lans, resp, err := c.LANsApi.DatacentersLansGet(c.ctx, datacenterId).Depth(1).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return nil, err
	}
	return &lans, nil
}

func (c *Client) GetNic(datacenterId, ServerId, NicId string) (*sdkgo.Nic, error) {
	nic, resp, err := c.NetworkInterfacesApi.DatacentersServersNicsFindById(c.ctx, datacenterId, ServerId, NicId).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return nil, err
	}
	log.Info("Got existing NIC!")
	return &nic, nil
}

func (c *Client) StartServer(datacenterId, serverId string) error {
	_, err := c.ServersApi.DatacentersServersStartPost(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	return nil
}

func (c *Client) StopServer(datacenterId, serverId string) error {
	_, err := c.ServersApi.DatacentersServersStopPost(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	return nil
}

func (c *Client) RestartServer(datacenterId, serverId string) error {
	_, err := c.ServersApi.DatacentersServersRebootPost(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	return nil
}

func (c *Client) RemoveServer(datacenterId, serverId string) error {
	resp, err := c.ServersApi.DatacentersServersDelete(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return err
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return fmt.Errorf("error waiting until server removed: %w", err)
	}
	log.Info("Server Deleted")
	return nil
}

func (c *Client) CreateAttachVolume(datacenterId, serverId string, volProperties *ClientVolumeProperties) (*sdkgo.Volume, error) {
	if volProperties == nil {
		return nil, fmt.Errorf("volume properties are nil")
	}
	sshKeys := &[]string{}
	if volProperties.SshKey != "" {
		sshKeys = &[]string{volProperties.SshKey}
	}
	var inputProperties sdkgo.VolumeProperties
	inputProperties.Type = &volProperties.DiskType
	inputProperties.Size = &volProperties.DiskSize
	inputProperties.ImagePassword = &volProperties.ImagePassword
	inputProperties.SshKeys = sshKeys
	inputProperties.AvailabilityZone = &volProperties.Zone

	if volProperties.ImageId != "" {
		inputProperties.Image = &volProperties.ImageId
	} else {
		inputProperties.ImageAlias = &volProperties.ImageAlias
	}
	if volProperties.UserData != "" {
		inputProperties.UserData = &volProperties.UserData
	}
	inputVolume := sdkgo.Volume{Properties: &inputProperties}

	volume, volumeResp, err := c.ServersApi.DatacentersServersVolumesPost(c.ctx, datacenterId, serverId).Volume(inputVolume).Execute()
	if err != nil {
		return nil, sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(volumeResp.StatusCode, volumeResp.Message)
	if err != nil {
		return nil, err
	}
	err = c.waitTillProvisioned(volumeResp.Header.Get("location"))
	if err != nil {
		return &volume, fmt.Errorf("error waiting until volume is created and attached: %w", err)
	}
	log.Info("attached volume to server!")
	return &volume, nil
}

func (c *Client) RemoveVolume(datacenterId, volumeId string) error {
	resp, err := c.VolumesApi.DatacentersVolumesDelete(c.ctx, datacenterId, volumeId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return err
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("Volume removed!")
	return nil
}

func (c *Client) CreateAttachNIC(datacenterId, serverId, name string, dhcp bool, lanId int32, ips *[]string) (*sdkgo.Nic, error) {
	n := sdkgo.Nic{
		Properties: &sdkgo.NicProperties{
			Name: &name,
			Lan:  &lanId,
			Ips:  ips,
			Dhcp: &dhcp,
		},
	}
	nic, nicResp, err := c.NetworkInterfacesApi.DatacentersServersNicsPost(c.ctx, datacenterId, serverId).Nic(n).Execute()
	if err != nil {
		return nil, err
	}
	err = sdk_utils.SanitizeStatusCode(nicResp.StatusCode, nicResp.Message)
	if err != nil {
		return nil, err
	}
	err = c.waitTillProvisioned(nicResp.Header.Get("location"))
	if err != nil {
		return &nic, fmt.Errorf("error waiting until nic is created and attached: %w", err)
	}
	log.Info("NIC attached to datacenter!")
	return &nic, nil
}

func (c *Client) RemoveNic(datacenterId, serverId, nicId string) error {
	resp, err := c.NetworkInterfacesApi.DatacentersServersNicsDelete(c.ctx, datacenterId, serverId, nicId).Execute()
	if err != nil {
		return sdk_utils.ShortenOpenApiErr(err)
	}
	err = sdk_utils.SanitizeStatusCode(resp.StatusCode, resp.Message)
	if err != nil {
		return err
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("NIC Deleted")
	return nil
}

func (c *Client) GetLocationById(regionId, locationId string) (*sdkgo.Location, error) {
	location, _, err := c.LocationsApi.LocationsFindByRegionIdAndId(c.ctx, regionId, locationId).Execute()
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (c *Client) GetTemplates() (*sdkgo.Templates, error) {
	templates, templatesResp, err := c.TemplatesApi.TemplatesGet(c.ctx).Depth(1).Execute()
	if err != nil {
		return nil, err
	}
	err = sdk_utils.SanitizeStatusCode(templatesResp.StatusCode, templatesResp.Message)
	if err != nil {
		return nil, err
	}
	return &templates, nil
}

func (c *Client) GetImages() (*sdkgo.Images, error) {
	images, imagesResp, err := c.ImagesApi.ImagesGet(c.ctx).Execute()
	if err != nil {
		return nil, err
	}
	err = sdk_utils.SanitizeStatusCode(imagesResp.StatusCode, imagesResp.Message)
	if err != nil {
		return nil, err
	}
	return &images, nil
}

func (c *Client) GetImageById(imageId string) (*sdkgo.Image, error) {
	image, imagesResp, err := c.ImagesApi.ImagesFindById(c.ctx, imageId).Execute()
	if err != nil {
		return nil, err
	}
	err = sdk_utils.SanitizeStatusCode(imagesResp.StatusCode, imagesResp.Message)
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (c *Client) waitTillProvisioned(path string) error {
	for i := 0; i < waitCount; i++ {
		requestStatus, _, err := c.RequestsApi.RequestsStatusGet(c.ctx, getRequestId(path)).Execute()
		if err != nil {
			return fmt.Errorf("error getting request status: %w", err)
		}
		if *requestStatus.Metadata.Status == "DONE" {
			return nil
		}
		if *requestStatus.Metadata.Status == "FAILED" {
			return fmt.Errorf("waiting returned a FAILED resource: %s", *requestStatus.Metadata.Message)
		}
		time.Sleep(10 * time.Second)
		i++
	}

	return fmt.Errorf("wait timeout has expired")
}

func getRequestId(path string) string {
	str := strings.Split(path, "/")
	return str[len(str)-2]
}
