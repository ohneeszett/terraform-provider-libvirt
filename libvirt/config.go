package libvirt

import (
	"log"

	libvirt "github.com/libvirt/libvirt-go"
)

type Config struct {
	Uri string
}

type Client struct {
	libvirt  *libvirt.Connect
	PoolSync *LibVirtPoolSync
}

func (c *Config) Client() (*Client, error) {
	conn, err := libvirt.NewConnect(c.Uri)
	if err != nil {
		return nil, err
	}

	client := &Client{
		libvirt:  conn,
		PoolSync: NewLibVirtPoolSync(),
	}

	log.Println("[INFO] Created libvirt client")

	return client, nil
}
