package ioutil

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRemoteImageDownloadRetry(t *testing.T) {
	content := []byte("this is a qcow image... well, it is not")

	// returns a server that returns every error from
	// errorList before returning a valid response
	newErrorServer := func(errorList []int) *httptest.Server {
		errorCount := 0
		return httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if errorCount < len(errorList) {
						t.Logf("Server serving retry %d", errorCount)
						http.Error(w, fmt.Sprintf("Error %d", errorCount), errorList[errorCount])
						errorCount = errorCount + 1
					} else {
						t.Logf("Server: success (after %d errors)", errorCount)
						http.ServeContent(w, r, "content", time.Now(), bytes.NewReader(content))
					}
				}))
	}

	server := newErrorServer([]int{503, 503})
	defer server.Close()

	start := time.Now()
	reader, err := NewHTTPReader(server.URL)
	if err != nil {
		t.Errorf("Could not create an HTTP reader: %v", err)
	}
	defer reader.Close()

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatalf("Expected to retry: %v", err)
	}
	if time.Since(start).Seconds() < 4 {
		t.Fatalf("Expected to retry at least 2 times x 2 seconds")
	}
	if !bytes.Equal(data, content) {
		t.Fatalf("Expected reading %v", string(content))
	}

	server = newErrorServer([]int{503, 404})
	defer server.Close()

	start = time.Now()
	reader, err = NewHTTPReader(server.URL)
	if err != nil {
		t.Errorf("Could not create an HTTP reader: %v", err)
	}
	defer reader.Close()

	data, err = ioutil.ReadAll(reader)
	if err == nil {
		t.Fatalf("Expected %s to fail with status 4xx", server.URL)
	}
	if time.Since(start).Seconds() < 2 {
		t.Fatalf("Expected to retry at least 1 times x 2 seconds")
	}
	if len(data) != 0 {
		t.Fatalf("Expected not reading anything")
	}

	server = newErrorServer([]int{304})
	defer server.Close()

	reader, err = NewHTTPReader(server.URL)
	if err != nil {
		t.Errorf("Could not create an HTTP reader: %v", err)
	}
	defer reader.Close()

	data, err = ioutil.ReadAll(reader)
	if err != ErrNotModified {
		t.Fatalf("Expected to fail with %v", ErrNotModified)
	}
	if len(data) != 0 {
		t.Fatalf("Expected not reading anything")
	}

}
