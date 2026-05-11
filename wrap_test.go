package errs

import (
	"errors"
	"maps"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type wrappedSnapshot struct {
	msg      string
	meta     map[string]any
	internal string
}

func assertWrappedErr(t *testing.T, err error, want wrappedSnapshot) {
	t.Helper()

	got, ok := errors.AsType[*WrappedErr](err)
	require.True(t, ok, "expected error to be of type *WrappedErr")

	assert.Equal(t, want, wrappedSnapshot{
		msg:      got.msg,
		meta:     got.meta,
		internal: normalizeInternalErr(got.internalErr),
	})
}

func normalizeInternalErr(err error) string {
	if err == nil {
		return ""
	}

	segments := strings.Split(err.Error(), " -> ")
	for i, segment := range segments {
		segments[i] = removeSourcePrefix(segment)
	}

	return strings.Join(segments, " -> ")
}

func removeSourcePrefix(text string) string {
	parts := strings.SplitN(text, " ", 2)
	if len(parts) < 2 || !strings.Contains(parts[0], ":") {
		return text
	}

	return parts[1]
}

func TestWrap(t *testing.T) {
	WithTrimSourcePref(true)
	t.Cleanup(func() { WithTrimSourcePref(false) })

	sampleMeta := map[string]any{
		"key1": "value1",
		"key2": 42,
	}

	tests := []struct {
		name string
		in   func() error
		desc string
		meta []map[string]any
		want *wrappedSnapshot
	}{
		{
			name: "wrap plain error with description",
			in:   func() error { return errors.New("plain error") },
			desc: "wrapped error",
			want: &wrappedSnapshot{msg: "wrapped error", internal: "plain error"},
		},
		{
			name: "wrap plain error without description",
			in:   func() error { return errors.New("plain error") },
			desc: "",
			want: &wrappedSnapshot{msg: "", internal: "plain error"},
		},
		{
			name: "wrap plain error with meta and description",
			in:   func() error { return errors.New("plain error") },
			desc: "wrapped error with meta",
			meta: []map[string]any{sampleMeta},
			want: &wrappedSnapshot{msg: "wrapped error with meta", meta: sampleMeta, internal: "plain error"},
		},
		{
			name: "wrap already wrapped error with description",
			in: func() error {
				return Wrap(errors.New("plain error"), "wrapped error")
			},
			desc: "wrapped again",
			want: &wrappedSnapshot{
				msg:      "wrapped again: wrapped error",
				meta:     map[string]any{},
				internal: "wraps -> plain error",
			},
		},
		{
			name: "wrap already wrapped error without description",
			in: func() error {
				return Wrap(errors.New("plain error"), "wrapped error")
			},
			desc: "",
			want: &wrappedSnapshot{
				msg:      "wrapped error",
				meta:     map[string]any{},
				internal: "wraps -> plain error",
			},
		},
		{
			name: "wrap already wrapped error with meta and description",
			in: func() error {
				return Wrap(errors.New("plain error"), "wrapped error with meta", sampleMeta)
			},
			desc: "wrapped again with meta",
			meta: []map[string]any{{"new_key": "new_value"}},
			want: &wrappedSnapshot{
				msg: "wrapped again with meta: wrapped error with meta",
				meta: func() map[string]any {
					mergedMeta := maps.Clone(sampleMeta)
					mergedMeta["new_key"] = "new_value"
					return mergedMeta
				}(),
				internal: "wraps -> plain error",
			},
		},
		{
			name: "wrap nil error",
			in:   func() error { return nil },
			desc: "should return nil",
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.in(), tt.desc, tt.meta...)
			if tt.want == nil {
				require.Nil(t, got)
				return
			}

			assertWrappedErr(t, got, *tt.want)
		})
	}
}

func TestWrappedErr_Wrap(t *testing.T) {
	WithTrimSourcePref(true)
	t.Cleanup(func() { WithTrimSourcePref(false) })

	sampleMeta := map[string]any{
		"key1": "value1",
		"key2": 42,
	}

	tests := []struct {
		name string
		we   func() *WrappedErr
		err  func() error
		desc string
		meta []map[string]any
		want *wrappedSnapshot
	}{
		{
			name: "nil receiver falls back to package Wrap",
			we:   func() *WrappedErr { return nil },
			err:  func() error { return errors.New("plain error") },
			desc: "wrapped via method",
			want: &wrappedSnapshot{msg: "wrapped via method", internal: "plain error"},
		},
		{
			name: "wrap plain error into existing wrapped chain",
			we: func() *WrappedErr {
				return Wrap(errors.New("plain error"), "wrapped error").(*WrappedErr)
			},
			err:  func() error { return errors.New("retry failed") },
			desc: "wrapped again",
			meta: []map[string]any{{"new_key": "new_value"}},
			want: &wrappedSnapshot{
				msg:      "wrapped again: wrapped error",
				meta:     map[string]any{"new_key": "new_value"},
				internal: "retry failed -> plain error",
			},
		},
		{
			name: "wrap plain error into existing wrapped chain without description",
			we: func() *WrappedErr {
				return Wrap(errors.New("plain error"), "wrapped error").(*WrappedErr)
			},
			err:  func() error { return errors.New("retry failed") },
			desc: "",
			want: &wrappedSnapshot{
				msg:      "wrapped error",
				meta:     map[string]any{},
				internal: "retry failed -> plain error",
			},
		},
		{
			name: "wrap plain error into wrapped chain with existing metadata",
			we: func() *WrappedErr {
				return Wrap(errors.New("plain error"), "wrapped error with meta", sampleMeta).(*WrappedErr)
			},
			err:  func() error { return errors.New("retry failed") },
			desc: "wrapped again with meta",
			meta: []map[string]any{{"new_key": "new_value"}},
			want: &wrappedSnapshot{
				msg: "wrapped again with meta: wrapped error with meta",
				meta: func() map[string]any {
					mergedMeta := maps.Clone(sampleMeta)
					mergedMeta["new_key"] = "new_value"
					return mergedMeta
				}(),
				internal: "retry failed -> plain error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.we().Wrap(tt.err(), tt.desc, tt.meta...)
			assertWrappedErr(t, got, *tt.want)
		})
	}
}
