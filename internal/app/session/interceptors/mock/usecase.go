// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mikalai-mitsin/example/internal/app/session/interceptors (interfaces: SessionUseCase)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/usecase.go . SessionUseCase
//

// Package mock_interceptors is a generated GoMock package.
package mock_interceptors

import (
	context "context"
	reflect "reflect"

	models "github.com/mikalai-mitsin/example/internal/app/session/models"
	uuid "github.com/mikalai-mitsin/example/internal/pkg/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionUseCase is a mock of SessionUseCase interface.
type MockSessionUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSessionUseCaseMockRecorder
}

// MockSessionUseCaseMockRecorder is the mock recorder for MockSessionUseCase.
type MockSessionUseCaseMockRecorder struct {
	mock *MockSessionUseCase
}

// NewMockSessionUseCase creates a new mock instance.
func NewMockSessionUseCase(ctrl *gomock.Controller) *MockSessionUseCase {
	mock := &MockSessionUseCase{ctrl: ctrl}
	mock.recorder = &MockSessionUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionUseCase) EXPECT() *MockSessionUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessionUseCase) Create(arg0 context.Context, arg1 *models.SessionCreate) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionUseCaseMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionUseCase)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockSessionUseCase) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionUseCaseMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessionUseCase)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockSessionUseCase) Get(arg0 context.Context, arg1 uuid.UUID) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSessionUseCaseMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSessionUseCase)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockSessionUseCase) List(arg0 context.Context, arg1 *models.SessionFilter) ([]*models.Session, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*models.Session)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockSessionUseCaseMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSessionUseCase)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockSessionUseCase) Update(arg0 context.Context, arg1 *models.SessionUpdate) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSessionUseCaseMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSessionUseCase)(nil).Update), arg0, arg1)
}
