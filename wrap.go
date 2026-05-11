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
		we.wrap(nil, src, desc, m)
		return we
	}

	return &WrappedErr{
		meta:        m,
		internalErr: fmt.Errorf("%s %w", src, e),
		msg:         desc,
	}
}
