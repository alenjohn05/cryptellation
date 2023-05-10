// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package events is a generated GoMock package.
package events

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

// Close mocks base method.
func (m *MockPort) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockPortMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPort)(nil).Close))
}

// Publish mocks base method.
func (m *MockPort) Publish(tick tick.Tick) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", tick)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockPortMockRecorder) Publish(tick interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPort)(nil).Publish), tick)
}

// Subscribe mocks base method.
func (m *MockPort) Subscribe(symbol string) (<-chan tick.Tick, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", symbol)
	ret0, _ := ret[0].(<-chan tick.Tick)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockPortMockRecorder) Subscribe(symbol interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPort)(nil).Subscribe), symbol)
}
