package ioutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLReader(t *testing.T) {
	r, err := NewURLReader("testdata/hello-nonexist.txt")
	assert.Error(t, err)

	r, err = NewURLReader("testdata/hello.txt")
	assert.NoError(t, err)
	assert.IsType(t, (*PathReader)(nil), r)

	r, err = NewURLReader("http://www.google.com")
	assert.NoError(t, err)
	assert.IsType(t, (*HTTPReader)(nil), r)
}
