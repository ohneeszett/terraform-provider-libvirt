package uri

import (
	"os"
	"testing"

	"github.com/dmacvicar/terraform-provider-libvirt/libvirt/io/http"
	"github.com/stretchr/testify/assert"
)

func TestURIOpen(t *testing.T) {
	r, err := Open("../testdata/hello-nonexist.txt")
	assert.Error(t, err)

	r, err = Open("../testdata/hello.txt")
	assert.NoError(t, err)
	assert.IsType(t, (*os.File)(nil), r)

	r, err = Open("http://www.google.com")
	assert.NoError(t, err)
	assert.IsType(t, (*http.File)(nil), r)
}
