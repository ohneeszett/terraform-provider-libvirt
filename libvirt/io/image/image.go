package image

import (
	uri "github.com/dmacvicar/terraform-provider-libvirt/libvirt/io/uri"
	xio "github.com/dmacvicar/terraform-provider-libvirt/libvirt/io"
	any "github.com/dmacvicar/terraform-provider-libvirt/libvirt/io/any"

	"github.com/libvirt/libvirt-go-xml"
)

type Format int
const (
	QCOW2 Format = iota
	Raw
)

const (
	qcow2Magic = "QFI\xfb\x00\x00\x00\x03"
)

type Image struct {
	xio.File
	Format Format
}

func Open(src string) (*Image, error) {
	// network transparent reader
	f, err := uri.Open(src)
	if err != nil {
		return nil, err
	}

	// compression
	a, err := any.NewAnyReader(f)
	if err != nil {
		return nil, err
	}

	// figure out format
	format := Raw
	buf, err := xio.NewBuffer(a)
	if err != nil {
		return nil, err
	}

	b, err := buf.Peek(len(qcow2Magic))
	if err != nil {
		return nil, err
	}
	if string(b) == qcow2Magic {
		format = QCOW2
	}
	return &Image{buf, format}, nil
}

func Import(src string, vol libvirtxml.StorageVolume) error {
	return nil
}
