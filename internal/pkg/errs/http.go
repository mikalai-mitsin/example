package errs

import (
	"encoding/json"
	"net/http"
)

const ClientClosedRequest = 499

func GetHTTPStatus(e *Error) int {
	switch e.Code {
	case ErrorCodeOK:
		return http.StatusOK
	case ErrorCodeCanceled:
		return ClientClosedRequest
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
func RenderToHTTPResponse(e *Error, writer http.ResponseWriter) error {
	writer.WriteHeader(GetHTTPStatus(e))
	return json.NewEncoder(writer).Encode(e)
}
