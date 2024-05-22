// Code generated by MockGen. DO NOT EDIT.
// Source: ticks.go

// Package client is a generated GoMock package.
package client

import (
	context "context"
	reflect "reflect"

	client "github.com/lerenn/cryptellation/pkg/client"
	event "github.com/lerenn/cryptellation/pkg/models/event"
	tick "github.com/lerenn/cryptellation/svc/ticks/pkg/tick"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close", ctx)
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close), ctx)
}

// ServiceInfo mocks base method.
func (m *MockClient) ServiceInfo(ctx context.Context) (client.ServiceInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceInfo", ctx)
	ret0, _ := ret[0].(client.ServiceInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServiceInfo indicates an expected call of ServiceInfo.
func (mr *MockClientMockRecorder) ServiceInfo(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceInfo", reflect.TypeOf((*MockClient)(nil).ServiceInfo), ctx)
}

// SubscribeToTicks mocks base method.
func (m *MockClient) SubscribeToTicks(ctx context.Context, sub event.TickSubscription) (<-chan tick.Tick, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToTicks", ctx, sub)
	ret0, _ := ret[0].(<-chan tick.Tick)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeToTicks indicates an expected call of SubscribeToTicks.
func (mr *MockClientMockRecorder) SubscribeToTicks(ctx, sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToTicks", reflect.TypeOf((*MockClient)(nil).SubscribeToTicks), ctx, sub)
}
