// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package exchanges is a generated GoMock package.
package exchanges

import (
	reflect "reflect"

	tick "github.com/lerenn/cryptellation/pkg/tick"
	gomock "github.com/golang/mock/gomock"
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

// ListenSymbol mocks base method.
func (m *MockPort) ListenSymbol(exchange, symbol string) (chan tick.Tick, chan struct{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenSymbol", exchange, symbol)
	ret0, _ := ret[0].(chan tick.Tick)
	ret1, _ := ret[1].(chan struct{})
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListenSymbol indicates an expected call of ListenSymbol.
func (mr *MockPortMockRecorder) ListenSymbol(exchange, symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenSymbol", reflect.TypeOf((*MockPort)(nil).ListenSymbol), exchange, symbol)
}
