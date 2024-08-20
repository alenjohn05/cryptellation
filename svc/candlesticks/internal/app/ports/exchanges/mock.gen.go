// Code generated by MockGen. DO NOT EDIT.
// Source: port.go

// Package exchanges is a generated GoMock package.
package exchanges

import (
	context "context"
	reflect "reflect"

	candlestick "github.com/lerenn/cryptellation/svc/candlesticks/pkg/candlestick"
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

// GetCandlesticks mocks base method.
func (m *MockPort) GetCandlesticks(ctx context.Context, payload GetCandlesticksPayload) (*candlestick.List, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCandlesticks", ctx, payload)
	ret0, _ := ret[0].(*candlestick.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCandlesticks indicates an expected call of GetCandlesticks.
func (mr *MockPortMockRecorder) GetCandlesticks(ctx, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCandlesticks", reflect.TypeOf((*MockPort)(nil).GetCandlesticks), ctx, payload)
}
