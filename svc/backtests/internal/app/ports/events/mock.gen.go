// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package events is a generated GoMock package.
package events

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	event "github.com/lerenn/cryptellation/pkg/event"
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

// Publish mocks base method.
func (m *MockPort) Publish(ctx context.Context, backtestID uuid.UUID, event event.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, backtestID, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockPortMockRecorder) Publish(ctx, backtestID, event interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPort)(nil).Publish), ctx, backtestID, event)
}

// Subscribe mocks base method.
func (m *MockPort) Subscribe(ctx context.Context, backtestID uuid.UUID) (<-chan event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", ctx, backtestID)
	ret0, _ := ret[0].(<-chan event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockPortMockRecorder) Subscribe(ctx, backtestID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPort)(nil).Subscribe), ctx, backtestID)
}
