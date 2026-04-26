# errs

`errs` helps you wrap errors with source location, readable context, and optional metadata.

## Requirements

- Go `1.26+`

## Install

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

func load() error {
	err := fetchData()
	return errs.Wrap(err, "fetch data", map[string]any{
		"service": "billing",
		"attempt": 2,
	})
}

func main() {
	if err := load(); err != nil {
		fmt.Println(err) // fetch data

		we, ok := errors.AsType[*errs.WrappedErr](err)
		if ok {
			if v, found := we.GetAttr("service"); found {
				fmt.Println("service:", v)
			}
		}
	}
}
```

## API

### `Wrap(e error, desc string, meta ...map[string]any) error`

Wraps an error with:

- source location (`file:line`),
- a human-readable message (`desc`),
- optional metadata (`meta`).

Behavior:

- returns `nil` when `e == nil`,
- appends context if `e` is already `*WrappedErr`,
- merges metadata into the existing wrapped error chain.

### `WrapChain`

Use `WrapChain(we *WrappedErr, e error, desc string, meta ...map[string]any) *WrappedErr`
when you need to append a **new plain error** to an existing wrapped chain.

Behavior:

- if `we == nil`, it starts a new wrapped chain (same as `Wrap`),
- if `we != nil`, it adds the new plain error to `internalErr`,
- it prepends the new description to `msg`,
- it merges metadata into the chain.

```go
base := errs.Wrap(errors.New("dial tcp timeout"), "fetch profile").(*errs.WrappedErr)

// Add a new plain error into the same wrapped chain.
next := errs.WrapChain(base, errors.New("retry failed"), "refresh cache", map[string]any{
	"attempt": 2,
})

fmt.Println(next.Error())
// refresh cache: fetch profile
```

### `WrappedError` (`WrappedErr`)

`WrappedErr` is the public wrapped error structure used by `Wrap` (referred to as wrapped error in this README).

- `Error() string` returns only the human-readable message,
- `GetAttr(key string) (any, bool)` reads metadata values by key.

Use `errors.AsType[*errs.WrappedErr]` (Go 1.26+) to access wrapped metadata:

```go
we, ok := errors.AsType[*errs.WrappedErr](err)
if ok {
	v, found := we.GetAttr("request_id")
	_ = v
	_ = found
}
```

## Source path trimming

By default, source location uses absolute paths. You can trim the project root prefix:

```go
errs.WithTrimSourcePref(true)
```

Example:

- before: `/your/path/to/project/internal/abc.go:10`
- after: `internal/abc.go:10`

## Notes

- The package captures call-site location via `runtime.Caller(1)`.
- Trimming is controlled globally with `WithTrimSourcePref`.

