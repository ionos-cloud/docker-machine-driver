// Code generated by MockGen. DO NOT EDIT.
// Source: internal/utils/client_service.go

// Package mock_utils is a generated GoMock package.
package mock_utils

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// MockClientService is a mock of ClientService interface.
type MockClientService struct {
	ctrl     *gomock.Controller
	recorder *MockClientServiceMockRecorder
}

// MockClientServiceMockRecorder is the mock recorder for MockClientService.
type MockClientServiceMockRecorder struct {
	mock *MockClientService
}

// NewMockClientService creates a new mock instance.
func NewMockClientService(ctrl *gomock.Controller) *MockClientService {
	mock := &MockClientService{ctrl: ctrl}
	mock.recorder = &MockClientServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientService) EXPECT() *MockClientServiceMockRecorder {
	return m.recorder
}

// CreateDatacenter mocks base method.
func (m *MockClientService) CreateDatacenter(name, location string) (*ionoscloud.Datacenter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDatacenter", name, location)
	ret0, _ := ret[0].(*ionoscloud.Datacenter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDatacenter indicates an expected call of CreateDatacenter.
func (mr *MockClientServiceMockRecorder) CreateDatacenter(name, location interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDatacenter", reflect.TypeOf((*MockClientService)(nil).CreateDatacenter), name, location)
}

// CreateIpBlock mocks base method.
func (m *MockClientService) CreateIpBlock(size int32, location string) (*ionoscloud.IpBlock, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIpBlock", size, location)
	ret0, _ := ret[0].(*ionoscloud.IpBlock)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateIpBlock indicates an expected call of CreateIpBlock.
func (mr *MockClientServiceMockRecorder) CreateIpBlock(size, location interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIpBlock", reflect.TypeOf((*MockClientService)(nil).CreateIpBlock), size, location)
}

// CreateLan mocks base method.
func (m *MockClientService) CreateLan(datacenterId, name string, public bool) (*ionoscloud.LanPost, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLan", datacenterId, name, public)
	ret0, _ := ret[0].(*ionoscloud.LanPost)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLan indicates an expected call of CreateLan.
func (mr *MockClientServiceMockRecorder) CreateLan(datacenterId, name, public interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLan", reflect.TypeOf((*MockClientService)(nil).CreateLan), datacenterId, name, public)
}

// CreateNat mocks base method.
func (m *MockClientService) CreateNat(datacenterId, name string, publicIps, flowlogs, natRules []string, lansToGateways map[string][]string, sourceSubnet string, skipDefaultRules bool) (*ionoscloud.NatGateway, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNat", datacenterId, name, publicIps, flowlogs, natRules, lansToGateways, sourceSubnet, skipDefaultRules)
	ret0, _ := ret[0].(*ionoscloud.NatGateway)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNat indicates an expected call of CreateNat.
func (mr *MockClientServiceMockRecorder) CreateNat(datacenterId, name, publicIps, flowlogs, natRules, lansToGateways, sourceSubnet, skipDefaultRules interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNat", reflect.TypeOf((*MockClientService)(nil).CreateNat), datacenterId, name, publicIps, flowlogs, natRules, lansToGateways, sourceSubnet, skipDefaultRules)
}

// CreateServer mocks base method.
func (m *MockClientService) CreateServer(datacenterId string, server ionoscloud.Server) (*ionoscloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateServer", datacenterId, server)
	ret0, _ := ret[0].(*ionoscloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateServer indicates an expected call of CreateServer.
func (mr *MockClientServiceMockRecorder) CreateServer(datacenterId, server interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateServer", reflect.TypeOf((*MockClientService)(nil).CreateServer), datacenterId, server)
}

// GetDatacenter mocks base method.
func (m *MockClientService) GetDatacenter(datacenterId string) (*ionoscloud.Datacenter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatacenter", datacenterId)
	ret0, _ := ret[0].(*ionoscloud.Datacenter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDatacenter indicates an expected call of GetDatacenter.
func (mr *MockClientServiceMockRecorder) GetDatacenter(datacenterId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatacenter", reflect.TypeOf((*MockClientService)(nil).GetDatacenter), datacenterId)
}

// GetDatacenters mocks base method.
func (m *MockClientService) GetDatacenters() (*ionoscloud.Datacenters, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDatacenters")
	ret0, _ := ret[0].(*ionoscloud.Datacenters)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDatacenters indicates an expected call of GetDatacenters.
func (mr *MockClientServiceMockRecorder) GetDatacenters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDatacenters", reflect.TypeOf((*MockClientService)(nil).GetDatacenters))
}

// GetImageById mocks base method.
func (m *MockClientService) GetImageById(imageId string) (*ionoscloud.Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageById", imageId)
	ret0, _ := ret[0].(*ionoscloud.Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageById indicates an expected call of GetImageById.
func (mr *MockClientServiceMockRecorder) GetImageById(imageId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageById", reflect.TypeOf((*MockClientService)(nil).GetImageById), imageId)
}

// GetImages mocks base method.
func (m *MockClientService) GetImages() (*ionoscloud.Images, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImages")
	ret0, _ := ret[0].(*ionoscloud.Images)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImages indicates an expected call of GetImages.
func (mr *MockClientServiceMockRecorder) GetImages() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImages", reflect.TypeOf((*MockClientService)(nil).GetImages))
}

// GetIpBlockIps mocks base method.
func (m *MockClientService) GetIpBlockIps(ipBlock *ionoscloud.IpBlock) (*[]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIpBlockIps", ipBlock)
	ret0, _ := ret[0].(*[]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIpBlockIps indicates an expected call of GetIpBlockIps.
func (mr *MockClientServiceMockRecorder) GetIpBlockIps(ipBlock interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIpBlockIps", reflect.TypeOf((*MockClientService)(nil).GetIpBlockIps), ipBlock)
}

// GetLan mocks base method.
func (m *MockClientService) GetLan(datacenterId, LanId string) (*ionoscloud.Lan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLan", datacenterId, LanId)
	ret0, _ := ret[0].(*ionoscloud.Lan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLan indicates an expected call of GetLan.
func (mr *MockClientServiceMockRecorder) GetLan(datacenterId, LanId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLan", reflect.TypeOf((*MockClientService)(nil).GetLan), datacenterId, LanId)
}

// GetLans mocks base method.
func (m *MockClientService) GetLans(datacenterId string) (*ionoscloud.Lans, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLans", datacenterId)
	ret0, _ := ret[0].(*ionoscloud.Lans)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLans indicates an expected call of GetLans.
func (mr *MockClientServiceMockRecorder) GetLans(datacenterId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLans", reflect.TypeOf((*MockClientService)(nil).GetLans), datacenterId)
}

// GetLocationById mocks base method.
func (m *MockClientService) GetLocationById(regionId, locationId string) (*ionoscloud.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocationById", regionId, locationId)
	ret0, _ := ret[0].(*ionoscloud.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocationById indicates an expected call of GetLocationById.
func (mr *MockClientServiceMockRecorder) GetLocationById(regionId, locationId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocationById", reflect.TypeOf((*MockClientService)(nil).GetLocationById), regionId, locationId)
}

// GetNat mocks base method.
func (m *MockClientService) GetNat(datacenterId, natId string) (*ionoscloud.NatGateway, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNat", datacenterId, natId)
	ret0, _ := ret[0].(*ionoscloud.NatGateway)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNat indicates an expected call of GetNat.
func (mr *MockClientServiceMockRecorder) GetNat(datacenterId, natId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNat", reflect.TypeOf((*MockClientService)(nil).GetNat), datacenterId, natId)
}

// GetNats mocks base method.
func (m *MockClientService) GetNats(datacenterId string) (*ionoscloud.NatGateways, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNats", datacenterId)
	ret0, _ := ret[0].(*ionoscloud.NatGateways)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNats indicates an expected call of GetNats.
func (mr *MockClientServiceMockRecorder) GetNats(datacenterId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNats", reflect.TypeOf((*MockClientService)(nil).GetNats), datacenterId)
}

// GetNic mocks base method.
func (m *MockClientService) GetNic(datacenterId, ServerId, NicId string) (*ionoscloud.Nic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNic", datacenterId, ServerId, NicId)
	ret0, _ := ret[0].(*ionoscloud.Nic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNic indicates an expected call of GetNic.
func (mr *MockClientServiceMockRecorder) GetNic(datacenterId, ServerId, NicId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNic", reflect.TypeOf((*MockClientService)(nil).GetNic), datacenterId, ServerId, NicId)
}

// GetServer mocks base method.
func (m *MockClientService) GetServer(datacenterId, serverId string, depth int32) (*ionoscloud.Server, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServer", datacenterId, serverId, depth)
	ret0, _ := ret[0].(*ionoscloud.Server)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetServer indicates an expected call of GetServer.
func (mr *MockClientServiceMockRecorder) GetServer(datacenterId, serverId, depth interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServer", reflect.TypeOf((*MockClientService)(nil).GetServer), datacenterId, serverId, depth)
}

// GetTemplates mocks base method.
func (m *MockClientService) GetTemplates() (*ionoscloud.Templates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplates")
	ret0, _ := ret[0].(*ionoscloud.Templates)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplates indicates an expected call of GetTemplates.
func (mr *MockClientServiceMockRecorder) GetTemplates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplates", reflect.TypeOf((*MockClientService)(nil).GetTemplates))
}

// RemoveDatacenter mocks base method.
func (m *MockClientService) RemoveDatacenter(datacenterId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveDatacenter", datacenterId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveDatacenter indicates an expected call of RemoveDatacenter.
func (mr *MockClientServiceMockRecorder) RemoveDatacenter(datacenterId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveDatacenter", reflect.TypeOf((*MockClientService)(nil).RemoveDatacenter), datacenterId)
}

// RemoveIpBlock mocks base method.
func (m *MockClientService) RemoveIpBlock(ipBlockId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveIpBlock", ipBlockId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveIpBlock indicates an expected call of RemoveIpBlock.
func (mr *MockClientServiceMockRecorder) RemoveIpBlock(ipBlockId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveIpBlock", reflect.TypeOf((*MockClientService)(nil).RemoveIpBlock), ipBlockId)
}

// RemoveLan mocks base method.
func (m *MockClientService) RemoveLan(datacenterId, lanId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveLan", datacenterId, lanId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveLan indicates an expected call of RemoveLan.
func (mr *MockClientServiceMockRecorder) RemoveLan(datacenterId, lanId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveLan", reflect.TypeOf((*MockClientService)(nil).RemoveLan), datacenterId, lanId)
}

// RemoveNat mocks base method.
func (m *MockClientService) RemoveNat(datacenterId, natId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNat", datacenterId, natId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNat indicates an expected call of RemoveNat.
func (mr *MockClientServiceMockRecorder) RemoveNat(datacenterId, natId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNat", reflect.TypeOf((*MockClientService)(nil).RemoveNat), datacenterId, natId)
}

// RemoveNic mocks base method.
func (m *MockClientService) RemoveNic(datacenterId, serverId, nicId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNic", datacenterId, serverId, nicId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNic indicates an expected call of RemoveNic.
func (mr *MockClientServiceMockRecorder) RemoveNic(datacenterId, serverId, nicId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNic", reflect.TypeOf((*MockClientService)(nil).RemoveNic), datacenterId, serverId, nicId)
}

// RemoveServer mocks base method.
func (m *MockClientService) RemoveServer(datacenterId, serverId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveServer", datacenterId, serverId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveServer indicates an expected call of RemoveServer.
func (mr *MockClientServiceMockRecorder) RemoveServer(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveServer", reflect.TypeOf((*MockClientService)(nil).RemoveServer), datacenterId, serverId)
}

// RemoveVolume mocks base method.
func (m *MockClientService) RemoveVolume(datacenterId, volumeId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveVolume", datacenterId, volumeId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveVolume indicates an expected call of RemoveVolume.
func (mr *MockClientServiceMockRecorder) RemoveVolume(datacenterId, volumeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveVolume", reflect.TypeOf((*MockClientService)(nil).RemoveVolume), datacenterId, volumeId)
}

// RestartServer mocks base method.
func (m *MockClientService) RestartServer(datacenterId, serverId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RestartServer", datacenterId, serverId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestartServer indicates an expected call of RestartServer.
func (mr *MockClientServiceMockRecorder) RestartServer(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestartServer", reflect.TypeOf((*MockClientService)(nil).RestartServer), datacenterId, serverId)
}

// ResumeServer mocks base method.
func (m *MockClientService) ResumeServer(datacenterId, serverId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResumeServer", datacenterId, serverId)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResumeServer indicates an expected call of ResumeServer.
func (mr *MockClientServiceMockRecorder) ResumeServer(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResumeServer", reflect.TypeOf((*MockClientService)(nil).ResumeServer), datacenterId, serverId)
}

// StartServer mocks base method.
func (m *MockClientService) StartServer(datacenterId, serverId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartServer", datacenterId, serverId)
	ret0, _ := ret[0].(error)
	return ret0
}

// StartServer indicates an expected call of StartServer.
func (mr *MockClientServiceMockRecorder) StartServer(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartServer", reflect.TypeOf((*MockClientService)(nil).StartServer), datacenterId, serverId)
}

// StopServer mocks base method.
func (m *MockClientService) StopServer(datacenterId, serverId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopServer", datacenterId, serverId)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopServer indicates an expected call of StopServer.
func (mr *MockClientServiceMockRecorder) StopServer(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopServer", reflect.TypeOf((*MockClientService)(nil).StopServer), datacenterId, serverId)
}

// SuspendServer mocks base method.
func (m *MockClientService) SuspendServer(datacenterId, serverId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SuspendServer", datacenterId, serverId)
	ret0, _ := ret[0].(error)
	return ret0
}

// SuspendServer indicates an expected call of SuspendServer.
func (mr *MockClientServiceMockRecorder) SuspendServer(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuspendServer", reflect.TypeOf((*MockClientService)(nil).SuspendServer), datacenterId, serverId)
}

// UpdateCloudInitFile mocks base method.
func (m *MockClientService) UpdateCloudInitFile(cloudInitYAML, key string, value []interface{}, single_value bool, behaviour string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCloudInitFile", cloudInitYAML, key, value, single_value, behaviour)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateCloudInitFile indicates an expected call of UpdateCloudInitFile.
func (mr *MockClientServiceMockRecorder) UpdateCloudInitFile(cloudInitYAML, key, value, single_value, behaviour interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCloudInitFile", reflect.TypeOf((*MockClientService)(nil).UpdateCloudInitFile), cloudInitYAML, key, value, single_value, behaviour)
}

// WaitForNicIpChange mocks base method.
func (m *MockClientService) WaitForNicIpChange(datacenterId, ServerId, NicId string, timeout int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForNicIpChange", datacenterId, ServerId, NicId, timeout)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitForNicIpChange indicates an expected call of WaitForNicIpChange.
func (mr *MockClientServiceMockRecorder) WaitForNicIpChange(datacenterId, ServerId, NicId, timeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForNicIpChange", reflect.TypeOf((*MockClientService)(nil).WaitForNicIpChange), datacenterId, ServerId, NicId, timeout)
}
