// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package events is a generated GoMock package.
package events

import (
	context "context"
	event "cryptellation/pkg/models/event"
	tick "cryptellation/svc/ticks/pkg/tick"
	reflect "reflect"

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

// Close mocks base method.
func (m *MockPort) Close(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close", ctx)
}

// Close indicates an expected call of Close.
func (mr *MockPortMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPort)(nil).Close), ctx)
}

// PublishTick mocks base method.
func (m *MockPort) PublishTick(ctx context.Context, tick tick.Tick) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishTick", ctx, tick)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishTick indicates an expected call of PublishTick.
func (mr *MockPortMockRecorder) PublishTick(ctx, tick interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishTick", reflect.TypeOf((*MockPort)(nil).PublishTick), ctx, tick)
}

// SubscribeToTicks mocks base method.
func (m *MockPort) SubscribeToTicks(ctx context.Context, sub event.TickSubscription) (<-chan tick.Tick, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToTicks", ctx, sub)
	ret0, _ := ret[0].(<-chan tick.Tick)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeToTicks indicates an expected call of SubscribeToTicks.
func (mr *MockPortMockRecorder) SubscribeToTicks(ctx, sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToTicks", reflect.TypeOf((*MockPort)(nil).SubscribeToTicks), ctx, sub)
}
