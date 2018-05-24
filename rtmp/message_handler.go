package rtmp

import (
	"context"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, m Message) ConnError
}

type MessageHandlerFunc func(ctx context.Context, m Message) ConnError

func (f MessageHandlerFunc) HandleMessage(ctx context.Context, m Message) ConnError {
	return f(ctx, m)
}
