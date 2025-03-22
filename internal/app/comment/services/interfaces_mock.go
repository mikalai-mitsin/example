// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go
//
// Generated by this command:
//
//	mockgen -source=interfaces.go -package=services -destination=interfaces_mock.go
//

// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"
	time "time"

	entities "github.com/mikalai-mitsin/example/internal/app/comment/entities"
	log "github.com/mikalai-mitsin/example/internal/pkg/log"
	uuid "github.com/mikalai-mitsin/example/internal/pkg/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockcommentRepository is a mock of commentRepository interface.
type MockcommentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockcommentRepositoryMockRecorder
	isgomock struct{}
}

// MockcommentRepositoryMockRecorder is the mock recorder for MockcommentRepository.
type MockcommentRepositoryMockRecorder struct {
	mock *MockcommentRepository
}

// NewMockcommentRepository creates a new mock instance.
func NewMockcommentRepository(ctrl *gomock.Controller) *MockcommentRepository {
	mock := &MockcommentRepository{ctrl: ctrl}
	mock.recorder = &MockcommentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcommentRepository) EXPECT() *MockcommentRepositoryMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockcommentRepository) Count(arg0 context.Context, arg1 entities.CommentFilter) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockcommentRepositoryMockRecorder) Count(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockcommentRepository)(nil).Count), arg0, arg1)
}

// Create mocks base method.
func (m *MockcommentRepository) Create(arg0 context.Context, arg1 entities.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockcommentRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockcommentRepository)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockcommentRepository) Delete(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockcommentRepositoryMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockcommentRepository)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockcommentRepository) Get(arg0 context.Context, arg1 uuid.UUID) (entities.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(entities.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockcommentRepositoryMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockcommentRepository)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockcommentRepository) List(arg0 context.Context, arg1 entities.CommentFilter) ([]entities.Comment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]entities.Comment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockcommentRepositoryMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockcommentRepository)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockcommentRepository) Update(arg0 context.Context, arg1 entities.Comment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockcommentRepositoryMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockcommentRepository)(nil).Update), arg0, arg1)
}

// Mockclock is a mock of clock interface.
type Mockclock struct {
	ctrl     *gomock.Controller
	recorder *MockclockMockRecorder
	isgomock struct{}
}

// MockclockMockRecorder is the mock recorder for Mockclock.
type MockclockMockRecorder struct {
	mock *Mockclock
}

// NewMockclock creates a new mock instance.
func NewMockclock(ctrl *gomock.Controller) *Mockclock {
	mock := &Mockclock{ctrl: ctrl}
	mock.recorder = &MockclockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockclock) EXPECT() *MockclockMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *Mockclock) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockclockMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*Mockclock)(nil).Now))
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

// MockuuidGenerator is a mock of uuidGenerator interface.
type MockuuidGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockuuidGeneratorMockRecorder
	isgomock struct{}
}

// MockuuidGeneratorMockRecorder is the mock recorder for MockuuidGenerator.
type MockuuidGeneratorMockRecorder struct {
	mock *MockuuidGenerator
}

// NewMockuuidGenerator creates a new mock instance.
func NewMockuuidGenerator(ctrl *gomock.Controller) *MockuuidGenerator {
	mock := &MockuuidGenerator{ctrl: ctrl}
	mock.recorder = &MockuuidGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuuidGenerator) EXPECT() *MockuuidGeneratorMockRecorder {
	return m.recorder
}

// NewUUID mocks base method.
func (m *MockuuidGenerator) NewUUID() uuid.UUID {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewUUID")
	ret0, _ := ret[0].(uuid.UUID)
	return ret0
}

// NewUUID indicates an expected call of NewUUID.
func (mr *MockuuidGeneratorMockRecorder) NewUUID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewUUID", reflect.TypeOf((*MockuuidGenerator)(nil).NewUUID))
}
