// Code generated by MockGen. DO NOT EDIT.
// Source: ./eyeshade/datastore.go

// Package eyeshade is a generated GoMock package.
package eyeshade

import (
	context "context"
	v4 "github.com/golang-migrate/migrate/v4"
	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
	reflect "reflect"
)

// MockDatastore is a mock of Datastore interface
type MockDatastore struct {
	ctrl     *gomock.Controller
	recorder *MockDatastoreMockRecorder
}

// MockDatastoreMockRecorder is the mock recorder for MockDatastore
type MockDatastoreMockRecorder struct {
	mock *MockDatastore
}

// NewMockDatastore creates a new mock instance
func NewMockDatastore(ctrl *gomock.Controller) *MockDatastore {
	mock := &MockDatastore{ctrl: ctrl}
	mock.recorder = &MockDatastoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatastore) EXPECT() *MockDatastoreMockRecorder {
	return m.recorder
}

// RawDB mocks base method
func (m *MockDatastore) RawDB() *sqlx.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RawDB")
	ret0, _ := ret[0].(*sqlx.DB)
	return ret0
}

// RawDB indicates an expected call of RawDB
func (mr *MockDatastoreMockRecorder) RawDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RawDB", reflect.TypeOf((*MockDatastore)(nil).RawDB))
}

// NewMigrate mocks base method
func (m *MockDatastore) NewMigrate() (*v4.Migrate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMigrate")
	ret0, _ := ret[0].(*v4.Migrate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMigrate indicates an expected call of NewMigrate
func (mr *MockDatastoreMockRecorder) NewMigrate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMigrate", reflect.TypeOf((*MockDatastore)(nil).NewMigrate))
}

// Migrate mocks base method
func (m *MockDatastore) Migrate(currentMigrationVersion uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate", currentMigrationVersion)
	ret0, _ := ret[0].(error)
	return ret0
}

// Migrate indicates an expected call of Migrate
func (mr *MockDatastoreMockRecorder) Migrate(currentMigrationVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockDatastore)(nil).Migrate), currentMigrationVersion)
}

// RollbackTxAndHandle mocks base method
func (m *MockDatastore) RollbackTxAndHandle(tx *sqlx.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackTxAndHandle", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RollbackTxAndHandle indicates an expected call of RollbackTxAndHandle
func (mr *MockDatastoreMockRecorder) RollbackTxAndHandle(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTxAndHandle", reflect.TypeOf((*MockDatastore)(nil).RollbackTxAndHandle), tx)
}

// RollbackTx mocks base method
func (m *MockDatastore) RollbackTx(tx *sqlx.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RollbackTx", tx)
}

// RollbackTx indicates an expected call of RollbackTx
func (mr *MockDatastoreMockRecorder) RollbackTx(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTx", reflect.TypeOf((*MockDatastore)(nil).RollbackTx), tx)
}

// GetAccountEarnings mocks base method
func (m *MockDatastore) GetAccountEarnings(ctx context.Context, options AccountEarningsOptions) (*[]AccountEarnings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountEarnings", ctx, options)
	ret0, _ := ret[0].(*[]AccountEarnings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountEarnings indicates an expected call of GetAccountEarnings
func (mr *MockDatastoreMockRecorder) GetAccountEarnings(ctx, options interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountEarnings", reflect.TypeOf((*MockDatastore)(nil).GetAccountEarnings), ctx, options)
}