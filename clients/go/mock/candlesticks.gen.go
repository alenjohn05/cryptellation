// Code generated by MockGen. DO NOT EDIT.
// Source: candlesticks.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	client "github.com/digital-feather/cryptellation/clients/go"
	candlestick "github.com/digital-feather/cryptellation/pkg/candlestick"
	gomock "github.com/golang/mock/gomock"
)

// MockCandlesticks is a mock of Candlesticks interface.
type MockCandlesticks struct {
	ctrl     *gomock.Controller
	recorder *MockCandlesticksMockRecorder
}

// MockCandlesticksMockRecorder is the mock recorder for MockCandlesticks.
type MockCandlesticksMockRecorder struct {
	mock *MockCandlesticks
}

// NewMockCandlesticks creates a new mock instance.
func NewMockCandlesticks(ctrl *gomock.Controller) *MockCandlesticks {
	mock := &MockCandlesticks{ctrl: ctrl}
	mock.recorder = &MockCandlesticksMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCandlesticks) EXPECT() *MockCandlesticksMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockCandlesticks) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockCandlesticksMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockCandlesticks)(nil).Close))
}

// Read mocks base method.
func (m *MockCandlesticks) Read(ctx context.Context, payload client.ReadCandlesticksPayload) (*candlestick.List, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, payload)
	ret0, _ := ret[0].(*candlestick.List)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockCandlesticksMockRecorder) Read(ctx, payload interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockCandlesticks)(nil).Read), ctx, payload)
}
