package errs

import (
	"encoding/json"
	"errors"
	"reflect"

	"go.uber.org/zap/zapcore"
)

type ErrorCode uint

const (
	ErrorCodeOK ErrorCode = iota
	ErrorCodeCanceled
	ErrorCodeUnknown
	ErrorCodeInvalidArgument
	ErrorCodeDeadlineExceeded
	ErrorCodeNotFound
	ErrorCodeAlreadyExists
	ErrorCodePermissionDenied
	ErrorCodeResourceExhausted
	ErrorCodeFailedPrecondition
	ErrorCodeAborted
	ErrorCodeOutOfRange
	ErrorCodeUnimplemented
	ErrorCodeInternal
	ErrorCodeUnavailable
	ErrorCodeDataLoss
	ErrorCodeUnauthenticated
	ErrorCodeClosedRequest
)

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p Param) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString(p.Key, p.Value)
	return nil
}

type Params []Param

func (p Params) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	for _, param := range p {
		encoder.AddString(param.Key, param.Value)
	}
	return nil
}

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Params  Params    `json:"params"`
	Err     error     `json:"-"`
}

func NewError(code ErrorCode, message string) *Error {
	return &Error{Code: code, Message: message, Params: nil, Err: nil}
}
func NewUnexpectedBehaviorError(details string) *Error {
	err := NewError(ErrorCodeInternal, "Unexpected behavior.")
	err.AddParam("details", details)
	return err
}
func NewInvalidFormError() *Error {
	return NewError(
		ErrorCodeInvalidArgument,
		"The form sent is not valid, please correct the errors below.",
	)
}
func NewInvalidParameter(message string) *Error {
	return NewError(ErrorCodeInvalidArgument, message)
}
func NewEntityNotFoundError() *Error {
	return NewError(ErrorCodeNotFound, "Entity not found.")
}
func NewBadTokenError() *Error {
	return NewError(ErrorCodePermissionDenied, "Bad token.")
}
func NewPermissionDeniedError() *Error {
	return NewError(ErrorCodePermissionDenied, "Permission denied.")
}
func NewSubscriptionAlreadyCancelledError() *Error {
	return NewError(ErrorCodeFailedPrecondition, "Subscription is already cancelled.")
}
func NewInactivePlanError() *Error {
	return NewError(ErrorCodeFailedPrecondition, "This plan is inactive.")
}
func (e *Error) Cause() error {
	return e.Err
}
func (e Error) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("message", e.Message)
	encoder.AddUint("code", uint(e.Code))
	if err := encoder.AddObject("params", e.Params); err != nil {
		return err
	}
	return nil
}
func (e *Error) WithParam(key, value string) *Error {
	e.AddParam(key, value)
	return e
}
func (e *Error) WithCause(err error) *Error {
	e.Err = err
	return e
}
func (e *Error) WithParams(params ...Param) *Error {
	if len(e.Params) == 0 {
		e.Params = params
	} else {
		e.Params = append(e.Params, params...)
	}
	return e
}
func (e Error) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}
func (e *Error) Is(tgt error) bool {
	var target *Error
	if ok := errors.As(tgt, &target); !ok {
		return false
	}
	target.Err = nil
	err := *e
	err.Err = nil
	eq := reflect.DeepEqual(&err, target)
	return eq
}
func (e *Error) SetCode(code ErrorCode) {
	e.Code = code
}
func (e *Error) SetCause(err error) {
	e.Err = err
}
func (e *Error) SetMessage(message string) {
	e.Message = message
}
func (e *Error) SetParams(params Params) {
	e.Params = params
}
func (e *Error) AddParam(key, value string) {
	e.Params = append(e.Params, Param{Key: key, Value: value})
}
