package uri

import (
	"fmt"
	"github.com/dmacvicar/terraform-provider-libvirt/libvirt/io"
	"github.com/dmacvicar/terraform-provider-libvirt/libvirt/io/http"
	"net/url"
	"os"
	"strings"
)

func Open(src string) (io.File, error) {
	url, err := url.Parse(src)
	if err != nil {
		return nil, fmt.Errorf("Can't parse source '%s' as url: %s", src, err)
	}

	if strings.HasPrefix(url.Scheme, "http") {
		r, err := http.Open(url.String())
		if err != nil {
			return nil, err
		}
		return r, nil
	} else if url.Scheme == "file" || url.Scheme == "" {
		r, err := os.Open(url.Path)
		if err != nil {
			return nil, err
		}
		return r, nil
	} else {
		return nil, fmt.Errorf("Don't know how to read from '%s': %s", url.String(), err)
	}
}
