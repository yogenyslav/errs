# errs

The easy way to handle errors.

## Example usage

```go
var (
    ErrInvalidParams = errors.New("tried to fetch data with invalid params")
    ErrFetchData     = errors.New("unable to fetch data")
)

func foo() (any, error) {
    // some logic
    resp, err := fetchData()
    if resp.StatuCode >= 400 && resp.StatusCode < 500 {
        return nil, errs.Wrap(err, ErrInvalidParams)
    } else if resp.StatusCode >= 500 {
        return nil, errs.Wrap(err, ErrFetchData)
    }
    return resp, errs.Wrap(err, "fetch data") // if err == nil than errs.Wrap() would also be nil
}

func bar() {
    // ...

    err := foo()
    if errors.Is(err, ErrInvalidParams) {
        // err with description "tried to fetch data with invalid params"
    } else if errors.Is(err, ErrFetchData) {
        // err with description "unable to fetch data"
    } else if err != nil {
        // err with description "fetch data"
        // e.g. "/path/to/the/project/main.go:50 fetch data -> connection reset by peer"
    }

    // ...
}
```

It is possible to turn on trimming the source file prefix up to the level of project root directory.

```go
func main() {
    enableTrim := true
    errs.WithTrimSourcePref(enableTrim)

    // project root: /your/path/to/project

    // errs.Wrap(err, "err description")
    // before: "/your/path/to/project/internal/abc.go:10 err description -> wrapped err"
    // after: "internal/abc.go:10 err description -> wrapped err"
}
```
