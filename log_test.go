package ran

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStdLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewStdLogger(buf, Error)

	logger.Info("aaa")
	assert.Equal(t, "", buf.String())
	buf.Reset()

	logger.Error("aaa")
	assert.Equal(t, "aaa\n", buf.String())
	buf.Reset()

	logger.Error("aaa: %s", "bbb")
	assert.Equal(t, "aaa: bbb\n", buf.String())
	buf.Reset()
}
