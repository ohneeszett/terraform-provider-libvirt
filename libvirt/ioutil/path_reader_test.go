package ioutil

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestPathReader(t *testing.T) {
	r, err := NewPathReader("testdata/hello-nonexist.txt")
	assert.Error(t, err)

	r, err = NewPathReader("testdata/hello.txt")
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	assert.Equal(t, "Hello\n\n", string(b))
	size, err := r.Size()
	assert.NoError(t, err)
	assert.Equal(t, int64(7), size)
}
