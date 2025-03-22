// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go
//

// Package usecases is a generated GoMock package.
package usecases

import (
	context "context"
	reflect "reflect"

	entities "github.com/mikalai-mitsin/example/internal/app/comment/entities"
	log "github.com/mikalai-mitsin/example/internal/pkg/log"
	uuid "github.com/mikalai-mitsin/example/internal/pkg/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockcommentService is a mock of commentService interface.
type MockcommentService struct {
	ctrl     *gomock.Controller
	recorder *MockcommentServiceMockRecorder
	isgomock struct{}
}

// MockcommentServiceMockRecorder is the mock recorder for MockcommentService.
type MockcommentServiceMockRecorder struct {
	mock *MockcommentService
}

// NewMockcommentService creates a new mock instance.
func NewMockcommentService(ctrl *gomock.Controller) *MockcommentService {
	mock := &MockcommentService{ctrl: ctrl}
	mock.recorder = &MockcommentServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcommentService) EXPECT() *MockcommentServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockcommentService) Create(arg0 context.Context, arg1 entities.CommentCreate) (entities.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(entities.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockcommentServiceMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockcommentService)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockcommentService) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockcommentServiceMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockcommentService)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockcommentService) Get(arg0 context.Context, arg1 uuid.UUID) (entities.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(entities.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockcommentServiceMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockcommentService)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockcommentService) List(arg0 context.Context, arg1 entities.CommentFilter) ([]entities.Comment, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]entities.Comment)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockcommentServiceMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockcommentService)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockcommentService) Update(arg0 context.Context, arg1 entities.CommentUpdate) (entities.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(entities.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockcommentServiceMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockcommentService)(nil).Update), arg0, arg1)
}

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
