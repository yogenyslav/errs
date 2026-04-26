package errs

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func internalErrsEqual(t *testing.T, want, got error) {
	t.Helper()

	if got == nil || want == nil {
		assert.Equal(t, got, want)
		return
	}

	wantErrMsg := removeSourceFromErr(want).Error()
	gotErrMsg := removeSourceFromErr(got).Error()

	assert.Equalf(
		t, wantErrMsg, gotErrMsg,
		"internal error messages should be equal ignoring source file and line number and format type, want %s, got %s",
		wantErrMsg, gotErrMsg,
	)
}

func removeSourceFromErr(err error) error {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	parts := strings.SplitN(errMsg, " ", 2)
	if len(parts) < 2 {
		return err
	}

	if !strings.Contains(parts[0], ":") {
		return err
	}

	return errors.New(parts[1])
}
