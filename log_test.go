package workflow

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewStdLogger(buf)

	logger.Info("aaa")
	assert.Equal(t, "aaa", buf.String())
	buf.Reset()

	logger.Info("aaa: %s", "bbb")
	assert.Equal(t, "aaa: bbb", buf.String())
	buf.Reset()

	logger.Error("aaa")
	assert.Equal(t, "aaa", buf.String())
	buf.Reset()

	logger.Error("aaa: %s", "bbb")
	assert.Equal(t, "aaa: bbb", buf.String())
	buf.Reset()
}
