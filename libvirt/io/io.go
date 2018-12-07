package io

import (
	"io"
	"os"
)

// os.File subset
type File interface {
	io.Closer
	io.Reader
	Stat() (os.FileInfo, error)
}
