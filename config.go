package errs

import (
	"fmt"
	"os"
	"sync/atomic"
)

var (
	trimSourcePref atomic.Bool
	projectRootDir string
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("get working directory: %w", err))
	}
	projectRootDir = wd + "/"
}

// WithTrimSourcePref sets trimming the source file prefix up to the project root dir.
func WithTrimSourcePref(val bool) {
	trimSourcePref.Store(val)
}
