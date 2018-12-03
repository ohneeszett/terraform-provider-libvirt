package ioutil

import (
	"os"
)

type PathReader struct {
	path string
	file *os.File
}

func NewPathReader(path string) (*PathReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &PathReader{path: path, file: file}, nil
}

func (r *PathReader) Size() (int64, error) {
	fi, err := r.file.Stat()
	if err != nil {
		return 0, err
	}
	return int64(fi.Size()), nil
}

func (r *PathReader) Close() error {
	if r.file != nil {
		return r.file.Close()
	}
	return nil
}

func (r *PathReader) Read(b []byte) (n int, err error) {
	return r.file.Read(b)
}
