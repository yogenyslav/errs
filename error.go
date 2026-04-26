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

func (e *WrappedErr) wrapInChain(err error, src source, desc string, meta map[string]any) {
	if e.meta == nil {
		e.meta = make(map[string]any)
	}
	maps.Insert(e.meta, maps.All(meta))
	if desc != "" {
		e.msg = fmt.Sprintf("%s: %s", desc, e.msg)
	}
	if err != nil {
		e.internalErr = fmt.Errorf("%s %w -> %w", src, err, e.internalErr)
	}
}
