// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/internal/app/equipment/interceptors (interfaces: AuthUseCase)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/auth.go . AuthUseCase
//

// Package mock_interceptors is a generated GoMock package.
package mock_interceptors

import (
	context "context"
	reflect "reflect"

	models "github.com/018bf/example/internal/app/user/models"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthUseCase is a mock of AuthUseCase interface.
type MockAuthUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockAuthUseCaseMockRecorder
}

// MockAuthUseCaseMockRecorder is the mock recorder for MockAuthUseCase.
type MockAuthUseCaseMockRecorder struct {
	mock *MockAuthUseCase
}

// NewMockAuthUseCase creates a new mock instance.
func NewMockAuthUseCase(ctrl *gomock.Controller) *MockAuthUseCase {
	mock := &MockAuthUseCase{ctrl: ctrl}
	mock.recorder = &MockAuthUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthUseCase) EXPECT() *MockAuthUseCaseMockRecorder {
	return m.recorder
}

// GetUser mocks base method.
func (m *MockAuthUseCase) GetUser(arg0 context.Context) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockAuthUseCaseMockRecorder) GetUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockAuthUseCase)(nil).GetUser), arg0)
}

// HasObjectPermission mocks base method.
func (m *MockAuthUseCase) HasObjectPermission(arg0 context.Context, arg1 *models.User, arg2 models.PermissionID, arg3 any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasObjectPermission", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// HasObjectPermission indicates an expected call of HasObjectPermission.
func (mr *MockAuthUseCaseMockRecorder) HasObjectPermission(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasObjectPermission", reflect.TypeOf((*MockAuthUseCase)(nil).HasObjectPermission), arg0, arg1, arg2, arg3)
}

// HasPermission mocks base method.
func (m *MockAuthUseCase) HasPermission(arg0 context.Context, arg1 *models.User, arg2 models.PermissionID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HasPermission", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// HasPermission indicates an expected call of HasPermission.
func (mr *MockAuthUseCaseMockRecorder) HasPermission(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HasPermission", reflect.TypeOf((*MockAuthUseCase)(nil).HasPermission), arg0, arg1, arg2)
}
