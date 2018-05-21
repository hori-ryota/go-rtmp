package rtmp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (conn *defaultConn) DefaultUserControlEventHandler() UserControlEventHandler {
	return UserControlEventHandler{
		StreamBeginHandlers: []StreamBeginHandler{
			StreamBeginHandlerFunc(func(ctx context.Context, streamBegin StreamBegin) ConnError {
				conn.Logger().Debug(
					"OnStreamBegin",
					zap.Object("streamBegin", streamBegin),
				)
				return nil
			}),
		},

		StreamEOFHandlers: []StreamEOFHandler{
			StreamEOFHandlerFunc(func(ctx context.Context, streamEOF StreamEOF) ConnError {
				conn.Logger().Debug(
					"OnStreamEOF",
					zap.Object("streamEOF", streamEOF),
				)
				return nil
			}),
		},

		StreamDryHandlers: []StreamDryHandler{
			StreamDryHandlerFunc(func(ctx context.Context, streamDry StreamDry) ConnError {
				conn.Logger().Debug(
					"OnStreamDry",
					zap.Object("streamDry", streamDry),
				)
				return nil
			}),
		},

		SetBufferLengthHandlers: []SetBufferLengthHandler{
			SetBufferLengthHandlerFunc(func(ctx context.Context, setBufferLength SetBufferLength) ConnError {
				conn.Logger().Debug(
					"OnSetBufferLength",
					zap.Object("setBufferLength", setBufferLength),
				)
				return nil
			}),
		},

		StreamIsRecordedHandlers: []StreamIsRecordedHandler{
			StreamIsRecordedHandlerFunc(func(ctx context.Context, streamIsRecorded StreamIsRecorded) ConnError {
				conn.Logger().Debug(
					"OnStreamIsRecorded",
					zap.Object("streamIsRecorded", streamIsRecorded),
				)
				return nil
			}),
		},

		PingRequestHandlers: []PingRequestHandler{
			PingRequestHandlerFunc(func(ctx context.Context, pingRequest PingRequest) ConnError {
				conn.Logger().Debug(
					"OnPingRequest",
					zap.Object("pingRequest", pingRequest),
				)
				if err := conn.PingResponse(ctx, pingRequest.Timestamp()); err != nil {
					return NewConnWarnError(
						errors.Wrap(err, "failed to send PingResponse"),
					)
				}
				return nil
			}),
		},

		PingResponseHandlers: []PingResponseHandler{
			PingResponseHandlerFunc(func(ctx context.Context, pingResponse PingResponse) ConnError {
				conn.Logger().Debug(
					"OnPingResponse",
					zap.Object("pingResponse", pingResponse),
				)
				return nil
			}),
		},
	}
}
