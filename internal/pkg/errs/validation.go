package errs

import (
	"bytes"
	"errors"
	"text/template"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func NewFromValidationError(err error) *Error {
	var validationErrors validation.Errors
	var validationErrorObject validation.ErrorObject
	if errors.As(err, &validationErrors) {
		e := NewError(
			ErrorCodeInvalidArgument,
			"The form sent is not valid, please correct the errors below.",
		)
		for key, value := range validationErrors {
			switch t := value.(type) {
			case validation.ErrorObject:
				e.AddParam(key, renderErrorMessage(t))
			case *Error:
				e.AddParam(key, t.Message)
			default:
				e.AddParam(key, value.Error())
			}
		}
		return e.WithCause(validationErrors)
	}
	if errors.As(err, &validationErrorObject) {
		return NewInvalidParameter(
			renderErrorMessage(validationErrorObject),
		).WithCause(validationErrorObject)
	}
	return nil
}
func renderErrorMessage(object validation.ErrorObject) string {
	parse, err := template.New("message").Parse(object.Message())
	if err != nil {
		return ""
	}
	var tpl bytes.Buffer
	_ = parse.Execute(&tpl, object.Params())
	return tpl.String()
}
