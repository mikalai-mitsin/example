// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -package=handlers -destination=interfaces_mock.go
//

// Package handlers is a generated GoMock package.
package handlers

import (
	context "context"
	reflect "reflect"

	entities "github.com/mikalai-mitsin/example/internal/app/auth/entities"
	log "github.com/mikalai-mitsin/example/internal/pkg/log"
	gomock "go.uber.org/mock/gomock"
)

// Mocklogger is a mock of logger interface.
type Mocklogger struct {
	ctrl     *gomock.Controller
	recorder *MockloggerMockRecorder
	isgomock struct{}
}

// MockloggerMockRecorder is the mock recorder for Mocklogger.
type MockloggerMockRecorder struct {
	mock *Mocklogger
}

// NewMocklogger creates a new mock instance.
func NewMocklogger(ctrl *gomock.Controller) *Mocklogger {
	mock := &Mocklogger{ctrl: ctrl}
	mock.recorder = &MockloggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocklogger) EXPECT() *MockloggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *Mocklogger) Debug(msg string, fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockloggerMockRecorder) Debug(msg any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*Mocklogger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *Mocklogger) Error(msg string, fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockloggerMockRecorder) Error(msg any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*Mocklogger)(nil).Error), varargs...)
}

// Fatal mocks base method.
func (m *Mocklogger) Fatal(msg string, fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockloggerMockRecorder) Fatal(msg any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*Mocklogger)(nil).Fatal), varargs...)
}

// Info mocks base method.
func (m *Mocklogger) Info(msg string, fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockloggerMockRecorder) Info(msg any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*Mocklogger)(nil).Info), varargs...)
}

// Panic mocks base method.
func (m *Mocklogger) Panic(msg string, fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Panic", varargs...)
}

// Panic indicates an expected call of Panic.
func (mr *MockloggerMockRecorder) Panic(msg any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Panic", reflect.TypeOf((*Mocklogger)(nil).Panic), varargs...)
}

// Print mocks base method.
func (m *Mocklogger) Print(msg string, fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Print", varargs...)
}

// Print indicates an expected call of Print.
func (mr *MockloggerMockRecorder) Print(msg any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Print", reflect.TypeOf((*Mocklogger)(nil).Print), varargs...)
}

// Warn mocks base method.
func (m *Mocklogger) Warn(msg string, fields ...log.Field) {
	m.ctrl.T.Helper()
	varargs := []any{msg}
	for _, a := range fields {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockloggerMockRecorder) Warn(msg any, fields ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{msg}, fields...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*Mocklogger)(nil).Warn), varargs...)
}

// MockauthUseCase is a mock of authUseCase interface.
type MockauthUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockauthUseCaseMockRecorder
	isgomock struct{}
}

// MockauthUseCaseMockRecorder is the mock recorder for MockauthUseCase.
type MockauthUseCaseMockRecorder struct {
	mock *MockauthUseCase
}

// NewMockauthUseCase creates a new mock instance.
func NewMockauthUseCase(ctrl *gomock.Controller) *MockauthUseCase {
	mock := &MockauthUseCase{ctrl: ctrl}
	mock.recorder = &MockauthUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockauthUseCase) EXPECT() *MockauthUseCaseMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockauthUseCase) CreateToken(ctx context.Context, login entities.Login) (entities.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", ctx, login)
	ret0, _ := ret[0].(entities.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockauthUseCaseMockRecorder) CreateToken(ctx, login any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockauthUseCase)(nil).CreateToken), ctx, login)
}

// RefreshToken mocks base method.
func (m *MockauthUseCase) RefreshToken(ctx context.Context, refresh entities.Token) (entities.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", ctx, refresh)
	ret0, _ := ret[0].(entities.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockauthUseCaseMockRecorder) RefreshToken(ctx, refresh any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockauthUseCase)(nil).RefreshToken), ctx, refresh)
}
