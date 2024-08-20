// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	forwardtest "github.com/lerenn/cryptellation/forwardtests/pkg/forwardtest"
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

// CreateForwardTest mocks base method.
func (m *MockPort) CreateForwardTest(ctx context.Context, ft forwardtest.ForwardTest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateForwardTest", ctx, ft)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateForwardTest indicates an expected call of CreateForwardTest.
func (mr *MockPortMockRecorder) CreateForwardTest(ctx, ft interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateForwardTest", reflect.TypeOf((*MockPort)(nil).CreateForwardTest), ctx, ft)
}

// DeleteForwardTest mocks base method.
func (m *MockPort) DeleteForwardTest(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteForwardTest", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteForwardTest indicates an expected call of DeleteForwardTest.
func (mr *MockPortMockRecorder) DeleteForwardTest(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteForwardTest", reflect.TypeOf((*MockPort)(nil).DeleteForwardTest), ctx, id)
}

// ListForwardTests mocks base method.
func (m *MockPort) ListForwardTests(ctx context.Context, filters ListFilters) ([]forwardtest.ForwardTest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListForwardTests", ctx, filters)
	ret0, _ := ret[0].([]forwardtest.ForwardTest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListForwardTests indicates an expected call of ListForwardTests.
func (mr *MockPortMockRecorder) ListForwardTests(ctx, filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListForwardTests", reflect.TypeOf((*MockPort)(nil).ListForwardTests), ctx, filters)
}

// ReadForwardTest mocks base method.
func (m *MockPort) ReadForwardTest(ctx context.Context, id uuid.UUID) (forwardtest.ForwardTest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadForwardTest", ctx, id)
	ret0, _ := ret[0].(forwardtest.ForwardTest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadForwardTest indicates an expected call of ReadForwardTest.
func (mr *MockPortMockRecorder) ReadForwardTest(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadForwardTest", reflect.TypeOf((*MockPort)(nil).ReadForwardTest), ctx, id)
}

// UpdateForwardTest mocks base method.
func (m *MockPort) UpdateForwardTest(ctx context.Context, ft forwardtest.ForwardTest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateForwardTest", ctx, ft)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateForwardTest indicates an expected call of UpdateForwardTest.
func (mr *MockPortMockRecorder) UpdateForwardTest(ctx, ft interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateForwardTest", reflect.TypeOf((*MockPort)(nil).UpdateForwardTest), ctx, ft)
}
