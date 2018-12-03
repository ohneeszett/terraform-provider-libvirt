package ioutil

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var ErrNotModified = errors.New("http: Not modified")

// The HTTP reader allows to threat the full request-response cycle
// as a reader.
type HTTPReader struct {
	url             string
	response        *http.Response
	ifModifiedSince *time.Time
}

func (r *HTTPReader) String() string {
	return r.url
}

// The reader will only download data if it was modified
// server-side after the given time.
//
// If the data has not been modified since, reading will
// return a ErrNotModified error.
func (r *HTTPReader) SetIfModifiedSince(t time.Time) {
	*r.ifModifiedSince = t
}

// Returns the length of the data according to server-side
func (r *HTTPReader) Size() (int64, error) {
	response, err := http.Head(r.url)
	if err != nil {
		return 0, err
	}
	if response.StatusCode == 403 {
		// possibly only the HEAD method is forbidden, try a Body-less GET instead
		response, err = http.Get(r.url)
		if err != nil {
			return 0, err
		}

		response.Body.Close()
	}
	if response.StatusCode != 200 {
		return 0,
			fmt.Errorf(
				"Error accessing remote resource: %s - %s",
				r.url,
				response.Status)
	}

	length, err := strconv.Atoi(response.Header.Get("Content-Length"))
	if err != nil {
		err = fmt.Errorf(
			"Error while getting Content-Length of \"%s\": %s - got %s",
			r.url,
			err,
			response.Header.Get("Content-Length"))
		return 0, err
	}
	return int64(length), nil
}

func NewHTTPReader(url string) (*HTTPReader, error) {
	return &HTTPReader{url: url}, nil
}

func (r *HTTPReader) Read(p []byte) (int, error) {
	err := r.doRequest()
	if err != nil {
		return 0, err
	}
	return r.response.Body.Read(p)
}

func (r *HTTPReader) Close() error {
	if r.response != nil {
		return r.response.Body.Close()
	}
	return nil
}

func (r *HTTPReader) doRequest() error {
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

	if r.ifModifiedSince != nil {
		req.Header.Set("If-Modified-Since", r.ifModifiedSince.UTC().Format(http.TimeFormat))
	}

	var response *http.Response
	for retryCount := 0; retryCount < maxHTTPRetries; retryCount++ {
		response, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("Error while downloading %s: %v", r.url, err)
		}

		log.Printf("[DEBUG]: url resp status code %s (retry #%d)\n", response.Status, retryCount)
		if response.StatusCode == http.StatusNotModified {
			response.Body.Close()
			return ErrNotModified
		} else if response.StatusCode == http.StatusOK {
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
