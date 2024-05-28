// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package client is a generated GoMock package.
package client

import (
	context "context"
	reflect "reflect"

	client "github.com/lerenn/cryptellation/pkg/client"
	timeserie "github.com/lerenn/cryptellation/pkg/models/timeserie"
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

// SMA mocks base method.
func (m *MockClient) SMA(ctx context.Context, payload SMAPayload) (*timeserie.TimeSerie[float64], error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SMA", ctx, payload)
	ret0, _ := ret[0].(*timeserie.TimeSerie[float64])
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SMA indicates an expected call of SMA.
func (mr *MockClientMockRecorder) SMA(ctx, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SMA", reflect.TypeOf((*MockClient)(nil).SMA), ctx, payload)
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
