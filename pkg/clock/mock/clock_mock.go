// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/pkg/clock (interfaces: Clock)

// Package mock_clock is a generated GoMock package.
package mock_clock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockClock is a mock of Clock interface.
type MockClock struct {
	ctrl     *gomock.Controller
	recorder *MockClockMockRecorder
}

// MockClockMockRecorder is the mock recorder for MockClock.
type MockClockMockRecorder struct {
	mock *MockClock
}

// NewMockClock creates a new mock instance.
func NewMockClock(ctrl *gomock.Controller) *MockClock {
	mock := &MockClock{ctrl: ctrl}
	mock.recorder = &MockClockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClock) EXPECT() *MockClockMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockClock) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockClockMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockClock)(nil).Now))
}
