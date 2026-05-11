// Package errs provides error handling utils.
package errs

import (
	"fmt"
	"maps"
)

// WrappedErr separates the internal data from human-readable error message and.
type WrappedErr struct {
	meta        map[string]any
	internalErr error
	msg         string
}

// Error implements error interface.
func (e *WrappedErr) Error() string {
	return e.msg
}

// GetAttr returns a value from wrapped meta by key.
// Second return value indicates whether it was found or not.
func (e *WrappedErr) GetAttr(key string) (any, bool) {
	if e.meta == nil {
		return nil, false
	}

	v, ok := e.meta[key]
	return v, ok
}

// Wrap fits the message into the error chain, reports source file with human-readable description
func (e *WrappedErr) Wrap(err error, desc string, meta ...map[string]any) error {
	if e == nil {
		return Wrap(err, desc, meta...)
	}

	var m map[string]any
	if len(meta) > 0 {
		m = meta[0]
	}

	src := getSource(wrapCallerSkip)
	e.wrap(err, src, desc, m)
	return e
}

func (e *WrappedErr) wrap(err error, src source, desc string, meta map[string]any) {
	if e.meta == nil {
		e.meta = make(map[string]any)
	}
	maps.Insert(e.meta, maps.All(meta))

	if desc != "" {
		e.msg = fmt.Sprintf("%s: %s", desc, e.msg)
	}

	if err != nil {
		e.internalErr = fmt.Errorf("%s %w -> %w", src, err, e.internalErr)
	} else {
		e.internalErr = fmt.Errorf("%s wraps -> %w", src, e.internalErr)
	}
}
