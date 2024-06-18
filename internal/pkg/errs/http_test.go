package errs

import (
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
