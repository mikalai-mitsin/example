// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mikalai-mitsin/example/internal/app/session/usecases (interfaces: UUIDGenerator)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/uuid_generator.go . UUIDGenerator
//

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	reflect "reflect"

	uuid "github.com/mikalai-mitsin/example/internal/pkg/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockUUIDGenerator is a mock of UUIDGenerator interface.
type MockUUIDGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDGeneratorMockRecorder
}

// MockUUIDGeneratorMockRecorder is the mock recorder for MockUUIDGenerator.
type MockUUIDGeneratorMockRecorder struct {
	mock *MockUUIDGenerator
}

// NewMockUUIDGenerator creates a new mock instance.
func NewMockUUIDGenerator(ctrl *gomock.Controller) *MockUUIDGenerator {
	mock := &MockUUIDGenerator{ctrl: ctrl}
	mock.recorder = &MockUUIDGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUUIDGenerator) EXPECT() *MockUUIDGeneratorMockRecorder {
	return m.recorder
}

// NewUUID mocks base method.
func (m *MockUUIDGenerator) NewUUID() uuid.UUID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUUID")
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

// NewUUID indicates an expected call of NewUUID.
func (mr *MockUUIDGeneratorMockRecorder) NewUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUUID", reflect.TypeOf((*MockUUIDGenerator)(nil).NewUUID))
}
