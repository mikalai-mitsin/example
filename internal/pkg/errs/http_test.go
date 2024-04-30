package errs

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHTTPStatus(t *testing.T) {
	type args struct {
		e *Error
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := GetHTTPStatus(tt.args.e)
			assert.Equal(t, got, tt.want)
		})
	}
}

func TestRenderToHTTPResponse(t *testing.T) {
	type args struct {
		e      *Error
		writer *httptest.ResponseRecorder
	}
	tests := []struct {
		name       string
		setup      func()
		args       args
		wantStatus int
		wantBody   string
		wantError  error
	}{
		{
			name: "ok",
			setup: func() {

			},
			args: args{
				e: NewInvalidFormError().
					WithCause(errors.New("validator error")).
					WithParams(Param{Key: "name", Value: "to short"}).
					WithParam("test", "12"),
				writer: httptest.NewRecorder(),
			},
			wantStatus: http.StatusBadRequest,
			wantBody: `{"meta":{"code":3,"message":"The form sent is not valid, please correct the errors below.","details":[{"fields":["name"],"message":"to short"},{"fields":["test"],"message":"12"}]}}
`,
			wantError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := RenderToHTTPResponse(tt.args.e, tt.args.writer)
			assert.ErrorIs(t, err, tt.wantError)
			assert.Equal(t, tt.wantStatus, tt.args.writer.Code)
			body, _ := io.ReadAll(tt.args.writer.Body)
			assert.Equal(t, tt.wantBody, string(body))
		})
	}
}