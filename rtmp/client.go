package rtmp

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Client struct {
	ctx         context.Context
	cancelFunc  context.CancelFunc
	connOptions []ConnOption

	logger *zap.Logger
}

func NewClient(
	ctx context.Context,
	connOps ...ConnOption,
) *Client {
	ctx, cancel := context.WithCancel(ctx)
	return &Client{
		ctx:         ctx,
		cancelFunc:  cancel,
		connOptions: connOps,
		logger:      defaultLogger,
	}
}

func (c *Client) Close() error {
	c.cancelFunc()
	return nil
}

func (c *Client) Connect(addr string) (Conn, error) {
	ctx := c.ctx
	var d net.Dialer
	nc, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial %s", addr)
	}
	if nc == nil {
		return nil, errors.New("conn is nil")
	}

	conn := NewDefaultConn(
		ctx,
		nc,
		false,
		c.Logger(),
		c.connOptions...,
	)

	go func() {
		remoteAddr := nc.RemoteAddr()
		defer func() {
			if err := conn.Close(); err != nil {
				c.Logger().Error(
					"failed to close conn",
					zap.Error(err),
					zap.Stringer("remoteAddr", remoteAddr),
				)
			}
		}()
		if err := conn.Serve(); err != nil {
			if isCanceledErr(err) {
				return
			}
			c.Logger().Error(
				"failed to conn.serve",
				zap.Error(err),
				zap.Stringer("remoteAddr", remoteAddr),
			)
		}
	}()
	return conn, ctx.Err()
}

func (c *Client) Logger() *zap.Logger {
	if c.logger != nil {
		return c.logger
	}
	return defaultLogger
}
