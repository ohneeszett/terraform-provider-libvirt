package image

import (
	"crypto/sha256"
	"fmt"
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
	img, err := Open("../../testdata/test.qcow2")
	assert.NoError(t, err)
	defer img.Close()

	h := sha256.New()

	if _, err := io.Copy(h, img); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, testQcow2Sha256, fmt.Sprintf("%x", h.Sum(nil)))

	fi, err := img.Stat()
	assert.NoError(t, err)
	assert.Equal(t, int64(testQcow2Size), fi.Size())
	assert.Equal(t, img.Format, QCOW2)
}

func TestImageCompressed(t *testing.T) {
	img, err := Open("../../testdata/gzip/test.qcow2.gz")
	assert.NoError(t, err)
	defer img.Close()

	assert.NoError(t, err)
	h := sha256.New()

	if _, err := io.Copy(h, img); err != nil {
		assert.NoError(t, err)
	}
	assert.Equal(t, testQcow2Sha256, fmt.Sprintf("%x", h.Sum(nil)))

	fi, err := img.Stat()
	assert.NoError(t, err)
	assert.Equal(t, int64(-1), fi.Size())
	assert.Equal(t, img.Format, QCOW2)
}

func TestImageRaw(t *testing.T) {
	img, err := Open("../../testdata/tcl.iso")
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

	img, err := Open(server.URL + "/gzip/test.qcow2.gz")
	assert.NoError(t, err)
	defer img.Close()
	assert.Equal(t, img.Format, QCOW2)
}
