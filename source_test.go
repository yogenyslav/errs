package errs

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTrimSourcePref(t *testing.T) {
	wd, wdErr := os.Getwd()
	require.NoError(t, wdErr)

	tests := []struct {
		name       string
		enableTrim bool
	}{
		{
			name:       "enabled trim",
			enableTrim: true,
		},
		{
			name:       "disabled trim",
			enableTrim: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				WithTrimSourcePref(tt.enableTrim)

				err := Wrap(errors.New("some error"), "asdf")
				we, ok := errors.AsType[*WrappedErr](err)
				require.Truef(t, ok, "error should be of type *WrappedErr, got %T", err)

				if tt.enableTrim {
					assert.NotContains(t, we.internalErr.Error(), wd)
				} else {
					assert.Contains(t, we.internalErr.Error(), wd)
				}
			},
		)
	}
}
