package ioutil

import (
	"fmt"
	"strings"
	"net/url"
)

type URLReader interface {
	closerSizedReader
}

func NewURLReader(src string) (URLReader, error) {
	url, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("Can't parse source '%s' as url: %s", src, err)
	}

	if strings.HasPrefix(url.Scheme, "http") {
		return NewHTTPReader(url.String())
	} else if url.Scheme == "file" || url.Scheme == "" {
		return NewPathReader(url.Path)
	} else {
		return nil, fmt.Errorf("Don't know how to read from '%s': %s", url.String(), err)
	}
}
