package rtmp

import (
	"context"
)

type HandleMessage func(ctx context.Context, m Message) ConnError

type MessageHandler interface {
	HandleMessage(ctx context.Context, m Message) ConnError
}

type abstructMessageHandler struct {
	handleMessage HandleMessage
}

func (h abstructMessageHandler) HandleMessage(ctx context.Context, m Message) ConnError {
	return h.handleMessage(ctx, m)
}

func HandleMessageFunc(
	handleFunc HandleMessage,
) MessageHandler {
	return &abstructMessageHandler{
		handleMessage: handleFunc,
	}
}
