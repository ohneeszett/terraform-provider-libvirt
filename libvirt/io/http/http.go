package ioutil

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// The HTTP file allows to threat the full request-response cycle
// as a reader.
type File struct {
	url      string
	response *http.Response
}

func (r *File) String() string {
	return r.url
}

// for Stat() return
type httpFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

// implement os.FileInfo
func (fi httpFileInfo) Name() string       { return fi.name }
func (fi httpFileInfo) Size() int64        { return fi.size }
func (fi httpFileInfo) Mode() os.FileMode  { return fi.mode }
func (fi httpFileInfo) ModTime() time.Time { return fi.modTime }
func (fi httpFileInfo) IsDir() bool        { return fi.isDir }
func (fi httpFileInfo) Sys() interface{}   { return nil }

// Returns file meta-data using a HEAD request
func getHeader(url string) (http.Header, error) {
	response, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == 403 {
		// possibly only the HEAD method is forbidden, try a Body-less GET instead
		response, err = http.Get(url)
		if err != nil {
			return nil, err
		}

		response.Body.Close()
	}
	if response.StatusCode != 200 {
		return nil,
			fmt.Errorf(
				"Error accessing remote resource: %s - %s",
				url,
				response.Status)
	}
	return response.Header, nil
}

func httpStat(urlS string) (os.FileInfo, error) {
	header, err := getHeader(urlS)
	if err != nil {
		return nil, err
	}

	url, err := url.Parse(urlS)
	if err != nil {
		return nil, fmt.Errorf("Can't parse source '%s' as url: %s", urlS, err)
	}

	fi := httpFileInfo{}
	fi.name = url.Path
	fi.size, err = strconv.ParseInt(header.Get("Content-Length"), 10, 0)
	if err != nil {
		err = fmt.Errorf(
			"Error while parsing Content-Length of \"%s\": %s - got %s",
			url,
			err,
			header.Get("Content-Length"))
		return nil, err
	}

	fi.modTime, err = http.ParseTime(header.Get("Last-Modified"))
	if err != nil {
		err = fmt.Errorf(
			"Error while parsing Last-Modified of \"%s\": %s - got %s",
			url,
			err,
			header.Get("Last-Modified"))
		return nil, err
	}
	return &fi, nil
}

func Open(url string) (*File, error) {
	return &File{url: url}, nil
}

func (r *File) Stat() (os.FileInfo, error) {
	return httpStat(r.url)
}

func (r *File) Read(p []byte) (int, error) {
	err := r.doRequest()
	if err != nil {
		return 0, err
	}
	return r.response.Body.Read(p)
}

func (r *File) Close() error {
	if r.response != nil {
		return r.response.Body.Close()
	}
	return nil
}

func (r *File) doRequest() error {
	if r.response != nil {
		return nil
	}

	// number of download retries on non client errors (eg. 5xx)
	const maxHTTPRetries int = 3
	// wait time between retries
	const retryWait time.Duration = 2 * time.Second

	client := &http.Client{}
	req, err := http.NewRequest("GET", r.url, nil)

	if err != nil {
		log.Printf("[DEBUG:] Error creating new request for source url %s: %s", r.url, err)
		return fmt.Errorf("Error while downloading %s: %s", r.url, err)
	}

	var response *http.Response
	for retryCount := 0; retryCount < maxHTTPRetries; retryCount++ {
		response, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("Error while downloading %s: %v", r.url, err)
		}

		log.Printf("[DEBUG]: url resp status code %s (retry #%d)\n", response.Status, retryCount)
		if response.StatusCode == http.StatusOK {
			r.response = response
			return nil
		} else if response.StatusCode < 500 {
			break
		} else {
			// The problem is not client but server side
			// retry a few times after a small wait
			if retryCount < maxHTTPRetries {
				time.Sleep(retryWait)
			}
		}
	}
	response.Body.Close()
	return fmt.Errorf("Error while downloading %s: %v", r.url, response)
}
