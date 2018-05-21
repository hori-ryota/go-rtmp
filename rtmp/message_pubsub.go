package rtmp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
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

func (s *defaultMessagePubsub) HandleMessage(ctx context.Context, m Message) ConnError {
	warnErrors := make([]error, 0, len(s.messageHandlers))
	for _, h := range s.messageHandlers {
		if err := h.HandleMessage(ctx, m); err != nil {
			if IsConnWarnError(err) {
				warnErrors = append(warnErrors, err)
			} else {
				return err
			}
		}
	}
	if len(warnErrors) == 0 {
		return nil
	}
	return NewConnWarnError(
		errors.New("error on defaultMessagePubsub.HandleMessage"),
		zap.Errors("handler errors", warnErrors),
	)
}

func (s *defaultMessagePubsub) AppendMessageHandler(hs ...MessageHandler) {
	s.messageHandlers = append(s.messageHandlers, hs...)
}
