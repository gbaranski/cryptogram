package cryptogram

import (
	"context"
	"net"
)

type Client struct {
	conn net.Conn
}

type Options struct {
	MDNS bool
}

func Connect(ctx context.Context, opts Options) {


}
