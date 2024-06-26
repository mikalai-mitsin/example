package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func TestNewFromValidationError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewFromValidationError(tt.args.err)
			assert.Equal(t, got, tt.want)
		})
	}
}

func Test_renderErrorMessage(t *testing.T) {
	type args struct {
		object validation.ErrorObject
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := renderErrorMessage(tt.args.object)
			assert.Equal(t, got, tt.want)
		})
	}
}
