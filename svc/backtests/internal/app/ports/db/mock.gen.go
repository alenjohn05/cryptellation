// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	backtest "cryptellation/svc/backtests/pkg/backtest"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockPort is a mock of Port interface.
type MockPort struct {
	ctrl     *gomock.Controller
	recorder *MockPortMockRecorder
}

// MockPortMockRecorder is the mock recorder for MockPort.
type MockPortMockRecorder struct {
	mock *MockPort
}

// NewMockPort creates a new mock instance.
func NewMockPort(ctrl *gomock.Controller) *MockPort {
	mock := &MockPort{ctrl: ctrl}
	mock.recorder = &MockPortMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPort) EXPECT() *MockPortMockRecorder {
	return m.recorder
}

// CreateBacktest mocks base method.
func (m *MockPort) CreateBacktest(ctx context.Context, bt backtest.Backtest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBacktest", ctx, bt)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateBacktest indicates an expected call of CreateBacktest.
func (mr *MockPortMockRecorder) CreateBacktest(ctx, bt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBacktest", reflect.TypeOf((*MockPort)(nil).CreateBacktest), ctx, bt)
}

// DeleteBacktest mocks base method.
func (m *MockPort) DeleteBacktest(ctx context.Context, bt backtest.Backtest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBacktest", ctx, bt)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBacktest indicates an expected call of DeleteBacktest.
func (mr *MockPortMockRecorder) DeleteBacktest(ctx, bt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBacktest", reflect.TypeOf((*MockPort)(nil).DeleteBacktest), ctx, bt)
}

// LockedBacktest mocks base method.
func (m *MockPort) LockedBacktest(ctx context.Context, id uuid.UUID, fn LockedBacktestCallback) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockedBacktest", ctx, id, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// LockedBacktest indicates an expected call of LockedBacktest.
func (mr *MockPortMockRecorder) LockedBacktest(ctx, id, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockedBacktest", reflect.TypeOf((*MockPort)(nil).LockedBacktest), ctx, id, fn)
}

// ReadBacktest mocks base method.
func (m *MockPort) ReadBacktest(ctx context.Context, id uuid.UUID) (backtest.Backtest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadBacktest", ctx, id)
	ret0, _ := ret[0].(backtest.Backtest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadBacktest indicates an expected call of ReadBacktest.
func (mr *MockPortMockRecorder) ReadBacktest(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadBacktest", reflect.TypeOf((*MockPort)(nil).ReadBacktest), ctx, id)
}

// UpdateBacktest mocks base method.
func (m *MockPort) UpdateBacktest(ctx context.Context, bt backtest.Backtest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBacktest", ctx, bt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBacktest indicates an expected call of UpdateBacktest.
func (mr *MockPortMockRecorder) UpdateBacktest(ctx, bt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBacktest", reflect.TypeOf((*MockPort)(nil).UpdateBacktest), ctx, bt)
}
