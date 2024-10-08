// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mikalai-mitsin/example/internal/app/auth/usecases (interfaces: AuthRepository)
//
// Generated by this command:
//
//	mockgen -build_flags=-mod=mod -destination mock/auth_repository.go . AuthRepository
//

// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	context "context"
	reflect "reflect"

	models "github.com/mikalai-mitsin/example/internal/app/auth/models"
	models0 "github.com/mikalai-mitsin/example/internal/app/user/models"
	gomock "go.uber.org/mock/gomock"
)

// MockAuthRepository is a mock of AuthRepository interface.
type MockAuthRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRepositoryMockRecorder
}

// MockAuthRepositoryMockRecorder is the mock recorder for MockAuthRepository.
type MockAuthRepositoryMockRecorder struct {
	mock *MockAuthRepository
}

// NewMockAuthRepository creates a new mock instance.
func NewMockAuthRepository(ctrl *gomock.Controller) *MockAuthRepository {
	mock := &MockAuthRepository{ctrl: ctrl}
	mock.recorder = &MockAuthRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthRepository) EXPECT() *MockAuthRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAuthRepository) Create(arg0 context.Context, arg1 *models0.User) (*models.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAuthRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAuthRepository)(nil).Create), arg0, arg1)
}

// GetSubject mocks base method.
func (m *MockAuthRepository) GetSubject(arg0 context.Context, arg1 models.Token) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSubject", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubject indicates an expected call of GetSubject.
func (mr *MockAuthRepositoryMockRecorder) GetSubject(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubject", reflect.TypeOf((*MockAuthRepository)(nil).GetSubject), arg0, arg1)
}

// RefreshToken mocks base method.
func (m *MockAuthRepository) RefreshToken(arg0 context.Context, arg1 models.Token) (*models.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", arg0, arg1)
	ret0, _ := ret[0].(*models.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockAuthRepositoryMockRecorder) RefreshToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuthRepository)(nil).RefreshToken), arg0, arg1)
}

// Validate mocks base method.
func (m *MockAuthRepository) Validate(arg0 context.Context, arg1 models.Token) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Validate", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Validate indicates an expected call of Validate.
func (mr *MockAuthRepositoryMockRecorder) Validate(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Validate", reflect.TypeOf((*MockAuthRepository)(nil).Validate), arg0, arg1)
}
