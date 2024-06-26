// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/018bf/example/internal/pkg/grpc (interfaces: AuthInterceptor)
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

	models "github.com/018bf/example/internal/app/auth/models"
	models0 "github.com/018bf/example/internal/app/user/models"
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

// Auth mocks base method.
func (m *MockAuthInterceptor) Auth(arg0 context.Context, arg1 models.Token) (*models0.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Auth", arg0, arg1)
	ret0, _ := ret[0].(*models0.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Auth indicates an expected call of Auth.
func (mr *MockAuthInterceptorMockRecorder) Auth(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Auth", reflect.TypeOf((*MockAuthInterceptor)(nil).Auth), arg0, arg1)
}
