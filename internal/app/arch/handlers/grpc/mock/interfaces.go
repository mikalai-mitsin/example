// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/internal/app/arch/handlers/grpc (interfaces: ArchInterceptor)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/interfaces.go . ArchInterceptor
//

// Package mock_grpc is a generated GoMock package.
package mock_grpc

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/example/internal/app/arch/models"
	uuid "github.com/018bf/example/internal/pkg/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockArchInterceptor is a mock of ArchInterceptor interface.
type MockArchInterceptor struct {
	ctrl     *gomock.Controller
	recorder *MockArchInterceptorMockRecorder
}

// MockArchInterceptorMockRecorder is the mock recorder for MockArchInterceptor.
type MockArchInterceptorMockRecorder struct {
	mock *MockArchInterceptor
}

// NewMockArchInterceptor creates a new mock instance.
func NewMockArchInterceptor(ctrl *gomock.Controller) *MockArchInterceptor {
	mock := &MockArchInterceptor{ctrl: ctrl}
	mock.recorder = &MockArchInterceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArchInterceptor) EXPECT() *MockArchInterceptorMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockArchInterceptor) Create(arg0 context.Context, arg1 *models.ArchCreate) (*models.Arch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.Arch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockArchInterceptorMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockArchInterceptor)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockArchInterceptor) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockArchInterceptorMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockArchInterceptor)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockArchInterceptor) Get(arg0 context.Context, arg1 uuid.UUID) (*models.Arch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*models.Arch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockArchInterceptorMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockArchInterceptor)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockArchInterceptor) List(arg0 context.Context, arg1 *models.ArchFilter) ([]*models.Arch, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*models.Arch)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockArchInterceptorMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockArchInterceptor)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockArchInterceptor) Update(arg0 context.Context, arg1 *models.ArchUpdate) (*models.Arch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*models.Arch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockArchInterceptorMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockArchInterceptor)(nil).Update), arg0, arg1)
}
