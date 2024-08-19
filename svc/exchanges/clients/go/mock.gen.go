// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package client is a generated GoMock package.
package client

import (
	context "context"
	client "cryptellation/pkg/client"
	exchange "cryptellation/svc/exchanges/pkg/exchange"
	reflect "reflect"

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

// Read mocks base method.
func (m *MockClient) Read(ctx context.Context, names ...string) ([]exchange.Exchange, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range names {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Read", varargs...)
	ret0, _ := ret[0].([]exchange.Exchange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockClientMockRecorder) Read(ctx interface{}, names ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, names...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockClient)(nil).Read), varargs...)
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
