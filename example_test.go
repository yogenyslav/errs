package errs_test

import (
	"errors"
	"fmt"

	"github.com/yogenyslav/errs"
)

func ExampleWrap() {
	err := errs.Wrap(errors.New("disk full"), "save file")
	fmt.Println(err)
	// Output: save file
}

func ExampleWrappedErr_GetAttr() {
	err := errs.Wrap(errors.New("timeout"), "fetch user", map[string]any{
		"service": "billing",
		"attempt": 2,
	})

	we, ok := errors.AsType[*errs.WrappedErr](err)
	fmt.Println(ok)

	service, serviceFound := we.GetAttr("service")
	fmt.Println(serviceFound, service)

	_, missingFound := we.GetAttr("request_id")
	fmt.Println(missingFound)

	// Output:
	// true
	// true billing
	// false
}

func ExampleWrapChain() {
	base := errs.Wrap(errors.New("dial tcp timeout"), "fetch profile", map[string]any{
		"service": "profiles",
	}).(*errs.WrappedErr)

	next := errs.WrapChain(base, errors.New("retry failed"), "refresh cache", map[string]any{
		"attempt": 2,
	})

	fmt.Println(next.Error())

	service, _ := next.GetAttr("service")
	attempt, _ := next.GetAttr("attempt")
	fmt.Println(service, attempt)

	// Output:
	// refresh cache: fetch profile
	// profiles 2
}

func ExampleWithTrimSourcePref() {
	// Enable trimming so internal source locations are relative to project root.
	errs.WithTrimSourcePref(true)
	defer errs.WithTrimSourcePref(false)

	_ = errs.Wrap(errors.New("disk full"), "save file")

	// If source trimming is enabled, internal source locations look like:
	// internal/file.go:42
}

