package ioutil

import (
	"fmt"
	"io"
	"net/url"
	"strings"
)

type URLReader interface {
	io.Reader
	io.Closer
	Sized
}

func NewURLReader(src string) (URLReader, error) {
	url, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("Can't parse source '%s' as url: %s", src, err)
	}

	if strings.HasPrefix(url.Scheme, "http") {
		r, err := NewHTTPReader(url.String())
		if err != nil {
			return nil, err
		}
		return r, nil
	} else if url.Scheme == "file" || url.Scheme == "" {
		r, err := NewPathReader(url.Path)
		if err != nil {
			return nil, err
		}
		return r, nil
	} else {
		return nil, fmt.Errorf("Don't know how to read from '%s': %s", url.String(), err)
	}
}
