// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mikalai-mitsin/example/internal/app/session/usecases (interfaces: SessionRepository)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/repository.go . SessionRepository
//

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	models "github.com/mikalai-mitsin/example/internal/app/session/models"
	uuid "github.com/mikalai-mitsin/example/internal/pkg/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSessionRepository is a mock of SessionRepository interface.
type MockSessionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryMockRecorder
}

// MockSessionRepositoryMockRecorder is the mock recorder for MockSessionRepository.
type MockSessionRepositoryMockRecorder struct {
	mock *MockSessionRepository
}

// NewMockSessionRepository creates a new mock instance.
func NewMockSessionRepository(ctrl *gomock.Controller) *MockSessionRepository {
	mock := &MockSessionRepository{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepository) EXPECT() *MockSessionRepositoryMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockSessionRepository) Count(arg0 context.Context, arg1 *models.SessionFilter) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockSessionRepositoryMockRecorder) Count(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockSessionRepository)(nil).Count), arg0, arg1)
}

// Create mocks base method.
func (m *MockSessionRepository) Create(arg0 context.Context, arg1 *models.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockSessionRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockSessionRepository) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionRepositoryMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessionRepository)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockSessionRepository) Get(arg0 context.Context, arg1 uuid.UUID) (*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSessionRepositoryMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSessionRepository)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockSessionRepository) List(arg0 context.Context, arg1 *models.SessionFilter) ([]*models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockSessionRepositoryMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSessionRepository)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockSessionRepository) Update(arg0 context.Context, arg1 *models.Session) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockSessionRepositoryMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSessionRepository)(nil).Update), arg0, arg1)
}
