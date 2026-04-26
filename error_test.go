package errs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrappedErr_Error(t *testing.T) {
	we := &WrappedErr{
		msg: "test error",
	}
	assert.Equal(t, "test error", we.Error())
}

func TestWrappedErr_GetAttr(t *testing.T) {
	we := &WrappedErr{
		meta: map[string]any{
			"key1": "value1",
			"key2": 42,
		},
	}

	value, ok := we.GetAttr("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	value, ok = we.GetAttr("key2")
	assert.True(t, ok)
	assert.Equal(t, 42, value)

	value, ok = we.GetAttr("nonexistent")
	assert.False(t, ok)
	assert.Nil(t, value)
}
