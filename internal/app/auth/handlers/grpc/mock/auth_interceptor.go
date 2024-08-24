// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mikalai-mitsin/example/internal/app/auth/handlers/grpc (interfaces: AuthInterceptor)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/auth_interceptor.go . AuthInterceptor
//

// Package mock_grpc is a generated GoMock package.
package mock_grpc

import (
	context "context"
	reflect "reflect"

	models "github.com/mikalai-mitsin/example/internal/app/auth/models"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthInterceptor is a mock of AuthInterceptor interface.
type MockAuthInterceptor struct {
	ctrl     *gomock.Controller
	recorder *MockAuthInterceptorMockRecorder
}

// MockAuthInterceptorMockRecorder is the mock recorder for MockAuthInterceptor.
type MockAuthInterceptorMockRecorder struct {
	mock *MockAuthInterceptor
}

// NewMockAuthInterceptor creates a new mock instance.
func NewMockAuthInterceptor(ctrl *gomock.Controller) *MockAuthInterceptor {
	mock := &MockAuthInterceptor{ctrl: ctrl}
	mock.recorder = &MockAuthInterceptorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthInterceptor) EXPECT() *MockAuthInterceptorMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockAuthInterceptor) CreateToken(arg0 context.Context, arg1 *models.Login) (*models.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", arg0, arg1)
	ret0, _ := ret[0].(*models.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockAuthInterceptorMockRecorder) CreateToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockAuthInterceptor)(nil).CreateToken), arg0, arg1)
}

// RefreshToken mocks base method.
func (m *MockAuthInterceptor) RefreshToken(arg0 context.Context, arg1 models.Token) (*models.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", arg0, arg1)
	ret0, _ := ret[0].(*models.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockAuthInterceptorMockRecorder) RefreshToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuthInterceptor)(nil).RefreshToken), arg0, arg1)
}
