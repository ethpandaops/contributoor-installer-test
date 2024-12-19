// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ethpandaops/contributoor-installer/internal/service (interfaces: ConfigManager)
//
// Generated by this command:
//
//	mockgen -package mock -destination mock/config.mock.go github.com/ethpandaops/contributoor-installer/internal/service ConfigManager
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	service "github.com/ethpandaops/contributoor-installer/internal/service"
	gomock "go.uber.org/mock/gomock"
)

// MockConfigManager is a mock of ConfigManager interface.
type MockConfigManager struct {
	ctrl     *gomock.Controller
	recorder *MockConfigManagerMockRecorder
}

// MockConfigManagerMockRecorder is the mock recorder for MockConfigManager.
type MockConfigManagerMockRecorder struct {
	mock *MockConfigManager
}

// NewMockConfigManager creates a new mock instance.
func NewMockConfigManager(ctrl *gomock.Controller) *MockConfigManager {
	mock := &MockConfigManager{ctrl: ctrl}
	mock.recorder = &MockConfigManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigManager) EXPECT() *MockConfigManagerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockConfigManager) Get() *service.ContributoorConfig {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*service.ContributoorConfig)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockConfigManagerMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfigManager)(nil).Get))
}

// GetConfigDir mocks base method.
func (m *MockConfigManager) GetConfigDir() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfigDir")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetConfigDir indicates an expected call of GetConfigDir.
func (mr *MockConfigManagerMockRecorder) GetConfigDir() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigDir", reflect.TypeOf((*MockConfigManager)(nil).GetConfigDir))
}

// GetConfigPath mocks base method.
func (m *MockConfigManager) GetConfigPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfigPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetConfigPath indicates an expected call of GetConfigPath.
func (mr *MockConfigManagerMockRecorder) GetConfigPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigPath", reflect.TypeOf((*MockConfigManager)(nil).GetConfigPath))
}

// Save mocks base method.
func (m *MockConfigManager) Save() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save")
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockConfigManagerMockRecorder) Save() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockConfigManager)(nil).Save))
}

// Update mocks base method.
func (m *MockConfigManager) Update(arg0 func(*service.ContributoorConfig)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockConfigManagerMockRecorder) Update(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockConfigManager)(nil).Update), arg0)
}