package client

import (
	"bytes"
	"context"
	"net"

	"github.com/tidwall/resp"
)

type client struct {
	addr string
}

func New(addr string) *client {
	return &client{}
}

func (c *client) Set(ctx context.Context, key string, val string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})
	_, err = conn.Write(buf.Bytes())
	return err
}
