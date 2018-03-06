package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

	"github.com/hori-ryota/go-rtmp/rtmp"
	amf "github.com/zhangpeihao/goamf"
	"go.uber.org/zap"
)

func main() {
	if err := Main(os.Args); err != nil {
		log.Fatal(err)
	}
}

func Main(args []string) error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}

	loggingHandler := NewMessageLoggingHandler(logger)

	connInit := rtmp.GenerateCommonConnInitializer(
		func(c rtmp.Conn) rtmp.MessageHandler {
			return loggingHandler
		},
	)

	return rtmp.ListenAndServe(
		context.Background(),
		"0.0.0.0:1935",
		connInit,
	)
}

type MessageLoggingHandler struct {
	logger *zap.Logger
}

func NewMessageLoggingHandler(
	logger *zap.Logger,
) rtmp.MessageHandler {
	return &MessageLoggingHandler{
		logger: logger,
	}
}

func (h MessageLoggingHandler) HandleMessage(
	ctx context.Context,
	m rtmp.Message,
) {
	h.logger.Info(
		"message received",
		zap.Object("message", m),
	)

	switch m.TypeID() {
	case rtmp.MessageTypeIDCommandAMF0:
		r := bytes.NewReader(m.Payload())
		for {
			v, err := amf.ReadValue(r)
			if err == io.EOF {
				break
			}
			if err != nil {
				h.logger.Error(
					"failed to ReadValue",
					zap.Error(err),
				)
			}
			zap.Any("value", v)
		}
	case rtmp.MessageTypeIDCommandAMF3:
		r := bytes.NewReader(m.Payload())
		for {
			v, err := amf.AMF3_ReadValue(r)
			if err == io.EOF {
				break
			}
			if err != nil {
				h.logger.Error(
					"failed to ReadValue",
					zap.Error(err),
				)
			}
			zap.Any("value", v)
		}
	}
}
