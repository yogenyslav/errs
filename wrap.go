package errs

import (
	"errors"
	"fmt"
)

const wrapCallerSkip = 2

// Wrap starts a new error chain or fits the message into it, reports source file with human-readable description
// and optionally extracts key-value attributes.
func Wrap(e error, desc string, meta ...map[string]any) error {
	if e == nil {
		return nil
	}

	var m map[string]any
	if len(meta) > 0 {
		m = meta[0]
	}

	src := getSource(wrapCallerSkip)
	we, isWrapped := errors.AsType[*WrappedErr](e)
	if isWrapped {
		we.wrapInChain(nil, src, desc, m)
		return we
	}

	return &WrappedErr{
		meta:        m,
		internalErr: fmt.Errorf("%s %w", src, e),
		msg:         desc,
	}
}

// WrapChain continues the error chain if the provided error is already wrapped, otherwise starts a new one.
func WrapChain(we *WrappedErr, e error, desc string, meta ...map[string]any) *WrappedErr {
	if we == nil {
		return Wrap(e, desc, meta...).(*WrappedErr)
	}

	var m map[string]any
	if len(meta) > 0 {
		m = meta[0]
	}

	src := getSource(wrapCallerSkip)
	we.wrapInChain(e, src, desc, m)
	return we
}
