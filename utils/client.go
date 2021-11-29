package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/machine/libmachine/log"
	sdkgo "github.com/ionos-cloud/sdk-go/v5"
)

const waitCount = 1000

type Client struct {
	*sdkgo.APIClient
	ctx context.Context
}

type ClientVolumeProperties struct {
	DiskType      string
	Name          string
	ImageAlias    string
	ImagePassword string
	Zone          string
	SshKey        string
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

func (c *Client) CreateIpBlock(size int32, location string) (*sdkgo.IpBlock, error) {
	ipBlock, ipBlockResp, err := c.IPBlocksApi.IpblocksPost(c.ctx).Ipblock(sdkgo.IpBlock{
		Properties: &sdkgo.IpBlockProperties{
			Location: &location,
			Size:     &size,
		}}).Execute()
	if err != nil {
		return nil, fmt.Errorf("error creating ipblock: %v", err)
	}
	if ipBlockResp.StatusCode == 202 {
		log.Info("IPBlock Reserved")
	} else {
		return nil, fmt.Errorf("error reserving an ipblock: %s", ipBlockResp.Response.Status)
	}
	err = c.waitTillProvisioned(ipBlockResp.Header.Get("location"))
	if err != nil {
		return &ipBlock, fmt.Errorf("error waiting until ip block is created: %v", err)
	}
	return &ipBlock, nil
}

func (c *Client) GetIpBlockIps(ipBlock *sdkgo.IpBlock) (*[]string, error) {
	if ipBlockProp, ok := ipBlock.GetPropertiesOk(); ok && ipBlockProp != nil {
		if ips, ok := ipBlockProp.GetIpsOk(); ok && ips != nil {
			return ips, nil
		}
	}
	return nil, fmt.Errorf("error getting ip block ips")
}

func (c *Client) RemoveIpBlock(ipBlockId string) error {
	_, resp, err := c.IPBlocksApi.IpblocksDelete(c.ctx, ipBlockId).Execute()
	if err != nil {
		return fmt.Errorf("error deleting ipblock: %v", err)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("error deleting ipblock, API Response status: %s", resp.Status)
	}
	log.Info("IPBlock Deleted")
	return nil
}

func (c *Client) CreateDatacenter(name, location string) (*sdkgo.Datacenter, error) {
	dc, dcResp, err := c.DataCenterApi.DatacentersPost(c.ctx).Datacenter(sdkgo.Datacenter{
		Properties: &sdkgo.DatacenterProperties{
			Name:     &name,
			Location: &location,
		}}).Execute()
	if err != nil {
		return nil, fmt.Errorf("error creating datacenter: %v", err)
	}
	if dcResp.StatusCode == 202 {
		log.Info("DataCenter Created")
	} else {
		return nil, fmt.Errorf("error creating DC: %s", dcResp.Response.Status)
	}
	err = c.waitTillProvisioned(dcResp.Header.Get("location"))
	if err != nil {
		return &dc, fmt.Errorf("error waiting until data center is created: %v", err)
	}
	return &dc, nil
}

func (c *Client) GetDatacenter(datacenterId string) (*sdkgo.Datacenter, error) {
	datacenter, resp, err := c.DataCenterApi.DatacentersFindById(c.ctx, datacenterId).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting datacenter: %v", err)
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("DataCenter UUID %s does not exist", datacenterId)
	}
	return &datacenter, nil
}

func (c *Client) RemoveDatacenter(datacenterId string) error {
	_, resp, err := c.DataCenterApi.DatacentersDelete(c.ctx, datacenterId).Execute()
	if err != nil {
		return fmt.Errorf("error deleting datacenter: %v", err)
	}
	if resp.StatusCode > 299 {
		if resp.StatusCode == 405 {
			return fmt.Errorf("error deleting datacenter: %v. Please consider to delete it manually", err)
		}
		return fmt.Errorf("error deleting datacenter, API Response status: %s", resp.Status)
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return fmt.Errorf("error waiting for datacenter to be deleted: %v", err)
	}
	log.Info("DataCenter Deleted")
	return nil
}

func (c *Client) CreateLan(datacenterId, name string, public bool) (*sdkgo.LanPost, error) {
	lan, lanResp, err := c.LanApi.DatacentersLansPost(c.ctx, datacenterId).Lan(sdkgo.LanPost{
		Properties: &sdkgo.LanPropertiesPost{
			Name:   &name,
			Public: &public,
		}}).Execute()
	if err != nil {
		return nil, fmt.Errorf("error creating LAN: %v", err)
	}
	if lanResp.StatusCode == 202 {
		log.Info("LAN Created")
	} else {
		return nil, fmt.Errorf("error creating a LAN: %s", lanResp.Response.Status)
	}
	err = c.waitTillProvisioned(lanResp.Header.Get("location"))
	if err != nil {
		return &lan, fmt.Errorf("error waiting until lan is created: %v", err)
	}
	return &lan, nil
}

func (c *Client) RemoveLan(datacenterId, lanId string) error {
	_, resp, err := c.LanApi.DatacentersLansDelete(c.ctx, datacenterId, lanId).Execute()
	if err != nil {
		return fmt.Errorf("error deleting LAN: %v", err)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf(resp.Status)
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("LAN Deleted")
	return nil
}

func (c *Client) CreateServer(datacenterId, location, name, cpufamily, zone string, ram, cores int32) (*sdkgo.Server, error) {
	server := sdkgo.Server{
		Properties: &sdkgo.ServerProperties{
			Name:             &name,
			Ram:              &ram,
			Cores:            &cores,
			CpuFamily:        &cpufamily,
			AvailabilityZone: &zone,
		},
	}

	svr, serverResp, err := c.ServerApi.DatacentersServersPost(c.ctx, datacenterId).Server(server).Execute()
	if err != nil {
		return nil, fmt.Errorf("error creating server in location %s err: %v", location, err)
	}
	if serverResp.StatusCode == 202 {
		log.Info("Server Created")
	} else {
		return nil, fmt.Errorf("error creating a server: %s", serverResp.Status)
	}
	err = c.waitTillProvisioned(serverResp.Header.Get("location"))
	if err != nil {
		return &svr, fmt.Errorf("error waiting until server is created: %v", err)
	}
	return &svr, nil
}

func (c *Client) GetServer(datacenterId, serverId string) (*sdkgo.Server, error) {
	server, resp, err := c.ServerApi.DatacentersServersFindById(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting server: %v", err)
	}
	if resp.StatusCode > 299 {
		if resp.StatusCode == 401 {
			return nil, fmt.Errorf("unauthorized: either user name or password are incorrect")

		} else {
			return nil, fmt.Errorf("error occurred fetching a server: %s", resp.Status)
		}
	}
	return &server, nil
}

func (c *Client) StartServer(datacenterId, serverId string) error {
	_, _, err := c.ServerApi.DatacentersServersStartPost(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}
	return nil
}

func (c *Client) StopServer(datacenterId, serverId string) error {
	_, _, err := c.ServerApi.DatacentersServersStopPost(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return fmt.Errorf("error stoping server: %v", err)
	}
	return nil
}

func (c *Client) RestartServer(datacenterId, serverId string) error {
	_, _, err := c.ServerApi.DatacentersServersRebootPost(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return fmt.Errorf("error restarting server: %v", err)
	}
	return nil
}

func (c *Client) RemoveServer(datacenterId, serverId string) error {
	_, resp, err := c.ServerApi.DatacentersServersDelete(c.ctx, datacenterId, serverId).Execute()
	if err != nil {
		return fmt.Errorf("error deleting server: %v", err)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf(resp.Status)
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("Server Deleted")
	return nil
}

func (c *Client) CreateAttachVolume(datacenterId, serverId string, volProperties *ClientVolumeProperties) (*sdkgo.Volume, error) {
	vol := sdkgo.Volume{
		Properties: &sdkgo.VolumeProperties{
			Type:             &volProperties.DiskType,
			Size:             &volProperties.DiskSize,
			Name:             &volProperties.Name,
			ImageAlias:       &volProperties.ImageAlias,
			ImagePassword:    &volProperties.ImagePassword,
			SshKeys:          &[]string{volProperties.SshKey},
			AvailabilityZone: &volProperties.Zone,
		},
	}
	volume, volumeResp, err := c.ServerApi.DatacentersServersVolumesPost(c.ctx, datacenterId, serverId).Volume(vol).Execute()
	if err != nil {
		return nil, fmt.Errorf("error attaching volume to server: %v", err)
	}
	if volumeResp.StatusCode == 202 {
		log.Info("Volume Attached to Server")
	} else {
		return nil, fmt.Errorf("error attaching a volume to a server: %s", volumeResp.Status)
	}
	err = c.waitTillProvisioned(volumeResp.Header.Get("location"))
	if err != nil {
		return &volume, fmt.Errorf("error waiting until volume is created and attached: %s", err.Error())
	}
	return &volume, nil
}

func (c *Client) RemoveVolume(datacenterId, volumeId string) error {
	_, resp, err := c.VolumeApi.DatacentersVolumesDelete(c.ctx, datacenterId, volumeId).Execute()
	if err != nil {
		return fmt.Errorf("error deleting volume: %v", err)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf(resp.Status)
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("Volume Deleted")
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
	nic, nicResp, err := c.NicApi.DatacentersServersNicsPost(c.ctx, datacenterId, serverId).Nic(n).Execute()
	if err != nil {
		return nil, fmt.Errorf("error attaching NIC to server: %s", err.Error())
	}
	if nicResp.StatusCode == 202 {
		log.Info("NIC Attached to Server")
	} else {
		return nil, fmt.Errorf("error creating a NIC: %s", nicResp.Status)
	}
	err = c.waitTillProvisioned(nicResp.Header.Get("location"))
	if err != nil {
		return &nic, fmt.Errorf("error waiting until nic is created and attached: %s", err.Error())
	}
	return &nic, nil
}

func (c *Client) RemoveNic(datacenterId, serverId, nicId string) error {
	_, resp, err := c.NicApi.DatacentersServersNicsDelete(c.ctx, datacenterId, serverId, nicId).Execute()
	if err != nil {
		return fmt.Errorf("error deleting NIC: %v", err)
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf(resp.Status)
	}
	err = c.waitTillProvisioned(resp.Header.Get("location"))
	if err != nil {
		return err
	}
	log.Info("NIC Deleted")
	return nil
}

func (c *Client) GetLocationById(regionId, locationId string) (*sdkgo.Location, error) {
	location, _, err := c.LocationApi.LocationsFindByRegionIdAndId(c.ctx, regionId, locationId).Execute()
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (c *Client) GetImages() (sdkgo.Images, error) {
	images, imagesResp, err := c.ImageApi.ImagesGet(c.ctx).Execute()
	if err != nil {
		return sdkgo.Images{}, err
	}
	if imagesResp.StatusCode == 401 {
		return sdkgo.Images{}, fmt.Errorf("error: authentication failed")
	}
	return images, nil
}

func (c *Client) waitTillProvisioned(path string) error {
	for i := 0; i < waitCount; i++ {
		requestStatus, _, err := c.RequestApi.RequestsStatusGet(c.ctx, getRequestId(path)).Execute()
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

func getRequestId(path string) string {
	str := strings.Split(path, "/")
	return str[len(str)-2]
}
