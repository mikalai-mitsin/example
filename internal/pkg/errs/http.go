package errs

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

func (e *Error) GetHTTPStatus() int {
	switch e.Code {
	case ErrorCodeOK:
		return http.StatusOK
	case ErrorCodeUnknown:
		return http.StatusInternalServerError
	case ErrorCodeInvalidArgument:
		return http.StatusBadRequest
	case ErrorCodeDeadlineExceeded:
		return http.StatusInternalServerError
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeAlreadyExists:
		return http.StatusBadRequest
	case ErrorCodePermissionDenied:
		return http.StatusForbidden
	case ErrorCodeResourceExhausted:
		return http.StatusInternalServerError
	case ErrorCodeFailedPrecondition:
		return http.StatusBadRequest
	case ErrorCodeAborted:
		return http.StatusInternalServerError
	case ErrorCodeOutOfRange:
		return http.StatusInternalServerError
	case ErrorCodeUnimplemented:
		return http.StatusMethodNotAllowed
	case ErrorCodeInternal:
		return http.StatusInternalServerError
	case ErrorCodeUnavailable:
		return http.StatusServiceUnavailable
	case ErrorCodeDataLoss:
		return http.StatusInternalServerError
	case ErrorCodeUnauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
func (e *Error) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.GetHTTPStatus())
	render.JSON(w, r, e)
	return nil
}
func RenderToHTTPResponse(err error, w http.ResponseWriter, r *http.Request) {
	var e *Error
	if errors.As(err, &e) {
		if err := render.Render(w, r, e); err != nil {
			render.Status(r, http.StatusInternalServerError)
		}
		return
	}
	render.Status(r, http.StatusInternalServerError)
	render.PlainText(w, r, err.Error())
}
