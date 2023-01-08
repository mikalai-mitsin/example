// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/internal/domain/usecases (interfaces: PostUseCase)

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/example/internal/domain/models"
	gomock "github.com/golang/mock/gomock"
)

// MockPostUseCase is a mock of PostUseCase interface.
type MockPostUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockPostUseCaseMockRecorder
}

// MockPostUseCaseMockRecorder is the mock recorder for MockPostUseCase.
type MockPostUseCaseMockRecorder struct {
	mock *MockPostUseCase
}

// NewMockPostUseCase creates a new mock instance.
func NewMockPostUseCase(ctrl *gomock.Controller) *MockPostUseCase {
	mock := &MockPostUseCase{ctrl: ctrl}
	mock.recorder = &MockPostUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostUseCase) EXPECT() *MockPostUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockPostUseCase) Create(arg0 context.Context, arg1 *models.PostCreate) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockPostUseCaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockPostUseCase)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockPostUseCase) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockPostUseCaseMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockPostUseCase)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockPostUseCase) Get(arg0 context.Context, arg1 string) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPostUseCaseMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPostUseCase)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockPostUseCase) List(arg0 context.Context, arg1 *models.PostFilter) ([]*models.Post, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*models.Post)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockPostUseCaseMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPostUseCase)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockPostUseCase) Update(arg0 context.Context, arg1 *models.PostUpdate) (*models.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockPostUseCaseMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPostUseCase)(nil).Update), arg0, arg1)
}
