package errs_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yogenyslav/errs"
)

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("Wrap error with context", func(t *testing.T) {
		rawErr := errors.New("original error")
		wrappedErr := errs.Wrap(rawErr, "additional context")

		assert.ErrorIs(t, wrappedErr, rawErr)
		assert.Contains(t, wrappedErr.Error(), "additional context")
		assert.Contains(t, wrappedErr.Error(), "original error")
	})

	t.Run("Wrap nil error returns nil", func(t *testing.T) {
		var rawErr error
		wrappedErr := errs.Wrap(rawErr, "additional context")

		assert.Nil(t, wrappedErr)
	})

	t.Run("Multiple wraps preserve original error", func(t *testing.T) {
		rawErr := errors.New("original error")
		firstWrap := errs.Wrap(rawErr, "first context")
		secondWrap := errs.Wrap(firstWrap, "second context")

		assert.ErrorIs(t, secondWrap, rawErr)
		assert.Contains(t, secondWrap.Error(), "first context")
		assert.Contains(t, secondWrap.Error(), "second context")
		assert.Contains(t, secondWrap.Error(), "original error")
	})

	t.Run("Wrap with empty context", func(t *testing.T) {
		rawErr := errors.New("original error")
		wrappedErr := errs.Wrap(rawErr, "")

		assert.ErrorIs(t, wrappedErr, rawErr)
	})
}
