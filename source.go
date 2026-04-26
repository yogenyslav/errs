package errs

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

type source struct {
	file      string
	line      int
	undefined bool
}

func getSource(skip int) source {
	_, file, line, ok := runtime.Caller(skip)
	return source{
		file:      file,
		line:      line,
		undefined: !ok,
	}
}

// String implements fmt.Stringer interface.
func (s source) String() string {
	if s.undefined {
		return ""
	}

	file := s.file
	if trimSourcePref.Load() {
		file = strings.TrimPrefix(s.file, projectRootDir)
	}
	return fmt.Sprintf("%s:%d", file, s.line)
}
