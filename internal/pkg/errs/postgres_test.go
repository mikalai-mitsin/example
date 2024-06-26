package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromPostgresError(t *testing.T) {
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
			got := FromPostgresError(tt.args.err)
			assert.Equal(t, got, tt.want)
		})
	}
}
