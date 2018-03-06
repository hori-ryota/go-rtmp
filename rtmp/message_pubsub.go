package rtmp

import (
	"context"
)

type MessagePubsub interface {
	MessageHandler
	AppendMessageHandler(...MessageHandler)
}

type messagePubsub MessagePubsub

type defaultMessagePubsub struct {
	messageHandlers []MessageHandler
}

func NewDefaultMessagePubsub(hs ...MessageHandler) MessagePubsub {
	s := &defaultMessagePubsub{
		messageHandlers: hs,
	}
	return s
}

func (s *defaultMessagePubsub) HandleMessage(ctx context.Context, m Message) {
	for _, h := range s.messageHandlers {
		h.HandleMessage(ctx, m)
	}
}

func (s *defaultMessagePubsub) AppendMessageHandler(hs ...MessageHandler) {
	s.messageHandlers = append(s.messageHandlers, hs...)
}
