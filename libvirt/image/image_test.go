package image

import (
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

const (
	testQcow2Sha256 = "0f71acdc66da59b04121b939573bec2e5be78a6cdf829b64142cf0a93a7076f5"
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
	assert.Equal(t, int64(256), size)
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
}
