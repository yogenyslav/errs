# errs

`errs` is a Go 1.26+ package for building structured error chains with:

- source location for the call site,
- a readable message chain for logs and user-facing output,
- optional key-value metadata,
- a public `WrappedErr` structure for inspection and chaining.

## Requirements

- Go `1.26+`

## Installation

```bash
go get github.com/yogenyslav/errs
```

## Quick start

```go
package main

import (
	"errors"
	"fmt"

	"github.com/yogenyslav/errs"
)

func fetchData() error {
	return errors.New("connection reset by peer")
}

func main() {
	err := errs.Wrap(fetchData(), "fetch data", map[string]any{
		"service": "billing",
		"attempt": 2,
	})

	if err != nil {
		fmt.Println(err.Error())

		we, ok := errors.AsType[*errs.WrappedErr](err)
		if ok {
			if v, found := we.GetAttr("service"); found {
				fmt.Println("service:", v)
			}
		}
	}
}
```

## How wrapping works

`WrappedErr` keeps two separate pieces of state:

- `msg` — the human-readable description chain returned by `Error()`,
- `internalErr` — the wrapped internal error chain.

When you call `(*WrappedErr).Wrap(err, desc, meta...)`:

- the new description is prepended to `msg`,
- the new plain error is accumulated into `internalErr`,
- metadata is merged into the existing map,
- the source location of the call is attached automatically.

When you call the package-level `Wrap(err, desc, meta...)` with a plain error, it creates a new `WrappedErr`. When you call it with an existing `*WrappedErr`, it extends the current message chain and keeps the wrapped chain intact.

## API reference

### `Wrap(e error, desc string, meta ...map[string]any) error`

Creates a wrapped error from a plain error, or extends an existing wrapped error.

Behavior:

- returns `nil` when `e == nil`,
- wraps plain errors with source location and description,
- extends an existing `*WrappedErr` by prepending the new description,
- preserves and merges metadata.

### `(*WrappedErr).Wrap(err error, desc string, meta ...map[string]any) error`

Adds a new plain error into an existing wrapped chain.

Behavior:

- if the receiver is `nil`, it falls back to `Wrap(err, desc, meta...)`,
- prepends `desc` to the visible message chain,
- appends the new plain error to the internal error chain,
- merges metadata into the chain.

```go
base := errs.Wrap(errors.New("dial tcp timeout"), "fetch profile").(*errs.WrappedErr)

next := base.Wrap(errors.New("retry failed"), "refresh cache", map[string]any{
	"attempt": 2,
})

fmt.Println(next.Error())
// refresh cache: fetch profile
```

### `WrappedErr`

`WrappedErr` is the public structure returned by the package.

- `Error() string` returns the readable message chain,
- `GetAttr(key string) (any, bool)` reads metadata values by key.

Use `errors.AsType[*errs.WrappedErr]` to inspect wrapped values:

```go
we, ok := errors.AsType[*errs.WrappedErr](err)
if ok {
	v, found := we.GetAttr("request_id")
	_ = v
	_ = found
}
```

## Source trimming

By default, source location uses the full path. If you want shorter paths in wrapped errors, enable prefix trimming:

```go
errs.WithTrimSourcePref(true)
```

Example:

- before: `/your/path/to/project/internal/abc.go:10`
- after: `internal/abc.go:10`

## Notes

- The package captures call-site information via `runtime.Caller(2)`.
- Source trimming is controlled globally with `WithTrimSourcePref`.

