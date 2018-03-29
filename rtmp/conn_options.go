package rtmp

import (
	"context"

	"go.uber.org/zap"
)

type connOptions struct {
	logger           *zap.Logger
	connInitializers []func(Conn)

	onConnectValidators []func(
		ctx context.Context,
		connect Connect,
	) ConnectError

	onPublishValidators []func(
		ctx context.Context,
		publish Publish,
	) (errorInfo map[string]interface{})
}

type ConnOption func(*connOptions)

func WithLogger(logger *zap.Logger) ConnOption {
	return func(o *connOptions) {
		o.logger = logger
	}
}

func WithConnInitializers(connInitializers ...func(Conn)) ConnOption {
	return func(o *connOptions) {
		o.connInitializers = append(o.connInitializers, connInitializers...)
	}
}

func WithOnConnectValidators(onConnectValidators ...func(ctx context.Context, connect Connect) ConnectError) ConnOption {
	return func(o *connOptions) {
		o.onConnectValidators = append(o.onConnectValidators, onConnectValidators...)
	}
}

func WithOnPublishValidators(onPublishValidators ...func(ctx context.Context, publish Publish) (errorInfo map[string]interface{})) ConnOption {
	return func(o *connOptions) {
		o.onPublishValidators = append(o.onPublishValidators, onPublishValidators...)
	}
}

func (o connOptions) Apply(c *defaultConn) {
	if o.logger != nil {
		c.logger = o.logger
	}
	if len(o.onConnectValidators) > 0 {
		c.onConnectValidators = o.onConnectValidators
	}
	if len(o.onPublishValidators) > 0 {
		c.onPublishValidators = o.onPublishValidators
	}
	for _, f := range o.connInitializers {
		f(c)
	}
}
