// Derived from https://github.com/paulcager/aio
// Copyright 2017 Paul Cager
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package io

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"os"
	"os/exec"

	xio "github.com/dmacvicar/terraform-provider-libvirt/libvirt/io"
)

var NotSupported = errors.New("The file type is not yet supported")
var ErrUnknownSize = errors.New("size of stream can't be determined")

// Don't use standard file detection software / libmagic as it requires >= 128 bytes to be read.
// https://en.wikipedia.org/wiki/List_of_file_signatures
//
// ZIP files / tar files - return concat of all contained files.
const (
	compressMagic = "\x1f\x9d"
	gzipMagic     = "\x1f\x8b"
	lzipMagic     = "LZIP"
	bzip2Magic    = "BZh"
	xzMagic       = "\xfd7zXZ\x00"
	zlibMagic     = "\x78\x9c"
)

type AnyReader struct {
	io.Reader
	filter    io.Reader
	sizeKnown bool
}

func NewAnyReader(r io.Reader) (*AnyReader, error) {
	return &AnyReader{r, nil, false}, nil
}

// for Stat() return
type anyFileInfo struct {
	os.FileInfo
	sizeKnown bool
}

// Size is the size of the source if the file is not
// compressed. In that case, we don't know.
func (fi anyFileInfo) Size() int64 {
	if fi.sizeKnown {
		return fi.FileInfo.Size()
	}
	return int64(-1)
}

func (a *AnyReader) Stat() (os.FileInfo, error) {
	f, ok := (a.Reader).(xio.File)
	if ok {
		fi, err := f.Stat()
		if err != nil {
			return nil, err
		}
		afi := anyFileInfo{fi, a.sizeKnown}
		return &afi, nil
	}
	return nil, NotSupported
}

func (a *AnyReader) Close() error {
	f, ok := (a.Reader).(xio.File)
	if ok {
		return f.Close()
	}
	// no op
	return nil
}

func (a *AnyReader) Read(b []byte) (n int, err error) {
	if a.filter == nil {
		err = a.decide()
		if err != nil {
			return 0, err
		}
	}
	return a.filter.Read(b)
}

func (a *AnyReader) decide() error {
	var err error
	if a.filter != nil {
		return nil
	}

	peeker := bufio.NewReader(a.Reader)
	a.filter = peeker

	if b, err := peeker.Peek(len(compressMagic)); err == nil && string(b) == compressMagic {
		// "compress" format. https://en.wikipedia.org/wiki/Lempel-Ziv-Welch
		return NotSupported
	} else if b, err := peeker.Peek(len(gzipMagic)); err == nil && string(b) == gzipMagic {
		// "gzip" format. https://tools.ietf.org/html/rfc1952
		a.filter, err = gzip.NewReader(a.filter)
	} else if b, err := peeker.Peek(len(bzip2Magic)); err == nil && string(b) == bzip2Magic {
		// "bz2" format.
		a.filter = bzip2.NewReader(a.filter)
	} else if b, err := peeker.Peek(len(zlibMagic)); err == nil && string(b) == zlibMagic {
		// "zlib" RFC 1950
		a.filter, err = zlib.NewReader(a.filter)
	} else if b, err := peeker.Peek(len(lzipMagic)); err == nil && string(b) == lzipMagic {
		return NotSupported
	} else if b, err := peeker.Peek(len(xzMagic)); err == nil && string(b) == xzMagic {
		a.filter = NewXZReader(a.filter)
	} else {
		// It is not a known format. Assume no compression.
		a.sizeKnown = true
	}

	return err
}

// NewXZReader creates a reader that decompresses the `xz` format input.
// Note that for convenience this is done by piping the input through an invocation
// of the `xz` command.
func NewXZReader(r io.Reader) io.Reader {
	return NewFilterReader(r, "xzcat")
}

func NewGZIPReader(r io.Reader) io.Reader {
	return NewFilterReader(r, "zcat")
}

func NewFilterReader(r io.Reader, cmdName string, args ...string) io.ReadCloser {
	rfilter, wfilter := io.Pipe()

	cmd := exec.Command(cmdName, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = r
	cmd.Stdout = wfilter

	go func() {
		err := cmd.Run()
		wfilter.CloseWithError(err)
	}()

	return rfilter
}
