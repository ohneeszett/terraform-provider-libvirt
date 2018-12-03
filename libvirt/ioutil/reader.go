package ioutil

import (
	"io"
)

type sizedReader interface {
	io.Reader
	Size() (int64, error)
}

type closerSizedReader interface {
	io.Reader
	Close() error
	Size() (int64, error)
}
