// Code generated by MockGen. DO NOT EDIT.
// Source: ./utils/clients/reputation/client.go

// Package mock_reputation is a generated GoMock package.
package mock_reputation

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	go_uuid "github.com/satori/go.uuid"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// IsWalletReputable mocks base method
func (m *MockClient) IsWalletReputable(ctx context.Context, id go_uuid.UUID, platform string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsWalletReputable", ctx, id, platform)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsWalletReputable indicates an expected call of IsWalletReputable
func (mr *MockClientMockRecorder) IsWalletReputable(ctx, id, platform interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsWalletReputable", reflect.TypeOf((*MockClient)(nil).IsWalletReputable), ctx, id, platform)
}

// IsWalletOnPlatform mocks base method
func (m *MockClient) IsWalletOnPlatform(ctx context.Context, id go_uuid.UUID, platform string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsWalletOnPlatform", ctx, id, platform)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsWalletOnPlatform indicates an expected call of IsWalletOnPlatform
func (mr *MockClientMockRecorder) IsWalletOnPlatform(ctx, id, platform interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsWalletOnPlatform", reflect.TypeOf((*MockClient)(nil).IsWalletOnPlatform), ctx, id, platform)
}
