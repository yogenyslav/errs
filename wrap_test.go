package errs

import (
	"errors"
	"fmt"
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrap(t *testing.T) {
	WithTrimSourcePref(true)
	defer WithTrimSourcePref(false)

	errPlain := errors.New("plain error")
	errWrapped := Wrap(errPlain, "wrapped error").(*WrappedErr)

	sampleMeta := map[string]any{
		"key1": "value1",
		"key2": 42,
	}
	errWrappedWithMeta := Wrap(errPlain, "wrapped error with meta", sampleMeta).(*WrappedErr)

	tests := []struct {
		name string
		err  error
		desc string
		meta []map[string]any
		want *WrappedErr
	}{
		{
			name: "wrap plain error with description",
			err:  errPlain,
			desc: "wrapped error",
			want: &WrappedErr{
				internalErr: errPlain,
				msg:         "wrapped error",
			},
		},
		{
			name: "wrap plain error without description",
			err:  errPlain,
			desc: "",
			want: &WrappedErr{
				internalErr: errPlain,
				msg:         "",
			},
		},
		{
			name: "wrap plain error with meta and description",
			err:  errPlain,
			desc: "wrapped error with meta",
			meta: []map[string]any{sampleMeta},
			want: &WrappedErr{
				internalErr: errPlain,
				msg:         "wrapped error with meta",
				meta:        sampleMeta,
			},
		},
		{
			name: "wrap already wrapped error with description",
			err:  new(*errWrapped),
			desc: "wrapped again",
			want: &WrappedErr{
				internalErr: errWrapped.internalErr,
				msg:         "wrapped again: " + errWrapped.msg,
				meta:        make(map[string]any),
			},
		},
		{
			name: "wrap already wrapped error without description",
			err:  new(*errWrapped),
			desc: "",
			want: &WrappedErr{
				internalErr: errWrapped.internalErr,
				msg:         errWrapped.msg,
				meta:        make(map[string]any),
			},
		},
		{
			name: "wrap already wrapped error with meta and description",
			err:  new(*errWrappedWithMeta),
			desc: "wrapped again with meta",
			meta: []map[string]any{
				{"new_key": "new_value"},
			},
			want: &WrappedErr{
				internalErr: errWrappedWithMeta.internalErr,
				msg:         "wrapped again with meta: " + errWrappedWithMeta.msg,
				meta: func() map[string]any {
					mergedMeta := make(map[string]any)
					maps.Copy(mergedMeta, sampleMeta)
					mergedMeta["new_key"] = "new_value"
					return mergedMeta
				}(),
			},
		},
		{
			name: "wrap nil error",
			err:  nil,
			desc: "should return nil",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				var err error

				wrappedCopy, ok := errors.AsType[*WrappedErr](tt.err)
				if ok {
					err = wrappedCopy
				} else {
					err = tt.err
				}

				e := Wrap(err, tt.desc, tt.meta...)
				if e == nil {
					assert.Nil(t, tt.want)
					return
				}

				we, ok := errors.AsType[*WrappedErr](e)
				require.True(t, ok, "expected error to be of type *WrappedErr")

				assert.Equal(t, tt.want.Error(), we.Error())
				assert.Equal(t, tt.want.meta, we.meta)
				internalErrsEqual(t, tt.want.internalErr, we.internalErr)
			},
		)
	}
}

func TestWrapChain(t *testing.T) {
	WithTrimSourcePref(true)
	defer WithTrimSourcePref(false)

	plainErr := errors.New("plain error")
	wrappedErr := Wrap(plainErr, "wrapped error").(*WrappedErr)

	tests := []struct {
		name string
		we   *WrappedErr
		err  error
		desc string
		meta []map[string]any
		want *WrappedErr
	}{
		{
			name: "wrap chain with nil WrappedErr",
			we:   nil,
			err:  plainErr,
			desc: "wrapped error",
			want: &WrappedErr{
				internalErr: plainErr,
				msg:         "wrapped error",
			},
		},
		{
			name: "wrap chain with existing WrappedErr",
			we:   wrappedErr,
			err:  plainErr,
			desc: "wrapped again",
			meta: []map[string]any{{"new_key": "new_value"}},
			want: &WrappedErr{
				internalErr: fmt.Errorf("%w -> %w", plainErr, wrappedErr.internalErr),
				msg:         "wrapped again: " + wrappedErr.msg,
				meta:        map[string]any{"new_key": "new_value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				we := WrapChain(tt.we, tt.err, tt.desc, tt.meta...)

				assert.Equal(t, tt.want.Error(), we.Error())
				assert.Equal(t, tt.want.meta, we.meta)
				internalErrsEqual(t, tt.want.internalErr, we.internalErr)
			},
		)
	}
}
