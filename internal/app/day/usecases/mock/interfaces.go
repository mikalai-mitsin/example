// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/internal/app/day/usecases (interfaces: DayRepository)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/interfaces.go . DayRepository
//

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/example/internal/app/day/models"
	uuid "github.com/018bf/example/internal/pkg/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockDayRepository is a mock of DayRepository interface.
type MockDayRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDayRepositoryMockRecorder
}

// MockDayRepositoryMockRecorder is the mock recorder for MockDayRepository.
type MockDayRepositoryMockRecorder struct {
	mock *MockDayRepository
}

// NewMockDayRepository creates a new mock instance.
func NewMockDayRepository(ctrl *gomock.Controller) *MockDayRepository {
	mock := &MockDayRepository{ctrl: ctrl}
	mock.recorder = &MockDayRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDayRepository) EXPECT() *MockDayRepositoryMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockDayRepository) Count(arg0 context.Context, arg1 *models.DayFilter) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockDayRepositoryMockRecorder) Count(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockDayRepository)(nil).Count), arg0, arg1)
}

// Create mocks base method.
func (m *MockDayRepository) Create(arg0 context.Context, arg1 *models.Day) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDayRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDayRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockDayRepository) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDayRepositoryMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDayRepository)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockDayRepository) Get(arg0 context.Context, arg1 uuid.UUID) (*models.Day, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*models.Day)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDayRepositoryMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDayRepository)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockDayRepository) List(arg0 context.Context, arg1 *models.DayFilter) ([]*models.Day, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*models.Day)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDayRepositoryMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDayRepository)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockDayRepository) Update(arg0 context.Context, arg1 *models.Day) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockDayRepositoryMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockDayRepository)(nil).Update), arg0, arg1)
}
