package image

import (
	"crypto/sha256"
	"fmt"
	xioutil "github.com/dmacvicar/terraform-provider-libvirt/libvirt/ioutil"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testQcow2Sha256 = "0f71acdc66da59b04121b939573bec2e5be78a6cdf829b64142cf0a93a7076f5"
	testQcow2Size   = 196616
)

func TestImageUnCompressed(t *testing.T) {
	img, err := NewImageFromSource("../testdata/test.qcow2")
	defer img.Close()

	assert.NoError(t, err)
	h := sha256.New()

	if _, err := io.Copy(h, img); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, testQcow2Sha256, fmt.Sprintf("%x", h.Sum(nil)))

	size, err := img.Size()
	assert.NoError(t, err)
	assert.Equal(t, int64(testQcow2Size), size)
}

func TestImageCompressed(t *testing.T) {
	img, err := NewImageFromSource("../testdata/gzip/test.qcow2.gz")
	assert.NoError(t, err)
	defer img.Close()

	assert.NoError(t, err)
	h := sha256.New()

	if _, err := io.Copy(h, img); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, testQcow2Sha256, fmt.Sprintf("%x", h.Sum(nil)))

	size, err := img.Size()
	assert.Error(t, err)
	assert.Equal(t, size, int64(-1))
	assert.Equal(t, err, xioutil.ErrUnknownSize)

	assert.Equal(t, img.Format, QCOW2)
}

func TestImageRaw(t *testing.T) {
	img, err := NewImageFromSource("../testdata/tcl.iso")
	assert.NoError(t, err)
	defer img.Close()
	assert.Equal(t, img.Format, Raw)
}

func TestImageHTTP(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.ServeFile(w, r, "../testdata/"+r.URL.Path[1:])
			}))
	defer server.Close()

	img, err := NewImageFromSource(server.URL + "/gzip/test.qcow2.gz")
	assert.NoError(t, err)
	defer img.Close()
	assert.Equal(t, img.Format, QCOW2)
}
