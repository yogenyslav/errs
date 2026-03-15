package errs_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yogenyslav/errs"
)

func TestWithTrimSourcePref(t *testing.T) {
	wd, wdErr := os.Getwd()
	require.NoError(t, wdErr)

	tests := []struct {
		name       string
		enableTrim bool
	}{
		{
			name:       "enabed trim",
			enableTrim: true,
		},
		{
			name:       "disabled trim",
			enableTrim: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs.WithTrimSourcePref(tt.enableTrim)

			err := errors.New("some error")
			err = errs.Wrap(err, "asdf")

			if tt.enableTrim {
				assert.NotContains(t, err.Error(), wd)
			} else {
				assert.Contains(t, err.Error(), wd)
			}
		})
	}
}
