package image

import (
	"bufio"
	xioutil "github.com/dmacvicar/terraform-provider-libvirt/libvirt/ioutil"
	"github.com/libvirt/libvirt-go-xml"
	"io"
)

type Format int

const (
	QCOW2 Format = iota
	Raw
)

type Source int

const (
	File = iota
	Vagrant
)

type sized interface {
	Size() (int64, error)
}

type Image struct {
	io.Reader
	io.Closer
	sized
}

func NewImageFromSource(src string) (*Image, error) {
	// network transparent reader
	r, err := xioutil.NewURLReader(src)
	if err != nil {
		return nil, err
	}

	// compression
	a, err := xioutil.NewAnyReader(r)
	if err != nil {
		return nil, err
	}

	return &Image{bufio.NewReader(a), a, a}, nil
}

func Import(src string, vol libvirtxml.StorageVolume) error {
	return nil
}
