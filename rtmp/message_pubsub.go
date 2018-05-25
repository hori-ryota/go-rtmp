package rtmp

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type MessagePubsub interface {
	MessageHandler
	AddMessageHandler(id string, h MessageHandler)
	RemoveMessageHandler(id string)
}

type messagePubsub MessagePubsub

type defaultMessagePubsub struct {
	messageHandlers sync.Map
	handlersLen     uint
}

func NewDefaultMessagePubsub() MessagePubsub {
	return &defaultMessagePubsub{
		messageHandlers: sync.Map{},
	}
}

func (s *defaultMessagePubsub) HandleMessage(ctx context.Context, m Message) ConnError {
	warnErrors := make([]error, 0, int(s.handlersLen))
	var fatalError ConnError
	s.messageHandlers.Range(func(_ interface{}, v interface{}) bool {
		h := v.(MessageHandler)
		if err := h.HandleMessage(ctx, m); err != nil {
			if IsConnWarnError(err) {
				warnErrors = append(warnErrors, err)
			} else {
				fatalError = err
				return false
			}
		}
		return true
	})
	if fatalError != nil {
		return fatalError
	}
	if len(warnErrors) == 0 {
		return nil
	}
	return NewConnWarnError(
		errors.New("error on defaultMessagePubsub.HandleMessage"),
		zap.Errors("handler errors", warnErrors),
	)
}

func (s *defaultMessagePubsub) AddMessageHandler(id string, h MessageHandler) {
	s.messageHandlers.Store(id, h)
	s.handlersLen++
}

func (s *defaultMessagePubsub) RemoveMessageHandler(id string) {
	s.messageHandlers.Delete(id)
	s.handlersLen--
}
