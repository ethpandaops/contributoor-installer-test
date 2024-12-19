// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ethpandaops/contributoor-installer/internal/service (interfaces: GitHubService)
//
// Generated by this command:
//
//	mockgen -package mock -destination mock/github.mock.go github.com/ethpandaops/contributoor-installer/internal/service GitHubService
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockGitHubService is a mock of GitHubService interface.
type MockGitHubService struct {
	ctrl     *gomock.Controller
	recorder *MockGitHubServiceMockRecorder
}

// MockGitHubServiceMockRecorder is the mock recorder for MockGitHubService.
type MockGitHubServiceMockRecorder struct {
	mock *MockGitHubService
}

// NewMockGitHubService creates a new mock instance.
func NewMockGitHubService(ctrl *gomock.Controller) *MockGitHubService {
	mock := &MockGitHubService{ctrl: ctrl}
	mock.recorder = &MockGitHubServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitHubService) EXPECT() *MockGitHubServiceMockRecorder {
	return m.recorder
}

// GetLatestVersion mocks base method.
func (m *MockGitHubService) GetLatestVersion() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestVersion")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestVersion indicates an expected call of GetLatestVersion.
func (mr *MockGitHubServiceMockRecorder) GetLatestVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestVersion", reflect.TypeOf((*MockGitHubService)(nil).GetLatestVersion))
}

// VersionExists mocks base method.
func (m *MockGitHubService) VersionExists(arg0 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VersionExists", arg0)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VersionExists indicates an expected call of VersionExists.
func (mr *MockGitHubServiceMockRecorder) VersionExists(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VersionExists", reflect.TypeOf((*MockGitHubService)(nil).VersionExists), arg0)
}
