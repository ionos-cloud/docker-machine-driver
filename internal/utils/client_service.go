package utils

import ionoscloud "github.com/ionos-cloud/sdk-go/v6"

type ClientService interface {
	CreateIpBlock(size int32, location string) (*ionoscloud.IpBlock, error)
	GetIpBlockIps(ipBlock *ionoscloud.IpBlock) (*[]string, error)
	RemoveIpBlock(ipBlockId string) error
	CreateDatacenter(name, location string) (*ionoscloud.Datacenter, error)
	GetDatacenter(datacenterId string) (*ionoscloud.Datacenter, error)
	RemoveDatacenter(datacenterId string) error
	CreateLan(datacenterId, name string, public bool) (*ionoscloud.LanPost, error)
	RemoveLan(datacenterId, lanId string) error
	CreateServer(datacenterId string, server ionoscloud.Server) (*ionoscloud.Server, error)
	GetServer(datacenterId, serverId string) (*ionoscloud.Server, error)
	GetLan(datacenterId, LanId string) (*ionoscloud.Lan, error)
	GetNic(datacenterId, ServerId, NicId string) (*ionoscloud.Nic, error)
	StartServer(datacenterId, serverId string) error
	StopServer(datacenterId, serverId string) error
	RestartServer(datacenterId, serverId string) error
	RemoveServer(datacenterId, serverId string) error
	CreateAttachVolume(datacenterId, serverId string, properties *ClientVolumeProperties) (*ionoscloud.Volume, error)
	RemoveVolume(datacenterId, volumeId string) error
	CreateAttachNIC(datacenterId, serverId, name string, dhcp bool, lanId int32, ips *[]string) (*ionoscloud.Nic, error)
	RemoveNic(datacenterId, serverId, nicId string) error
	GetLocationById(regionId, locationId string) (*ionoscloud.Location, error)
	GetImages() (*ionoscloud.Images, error)
	GetImageById(imageId string) (*ionoscloud.Image, error)
}
