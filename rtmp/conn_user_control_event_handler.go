package rtmp

import (
	"context"

	"go.uber.org/zap"
)

func (conn *defaultConn) OnStreamBegin(ctx context.Context, streamBegin StreamBegin) {
	conn.Logger().Debug(
		"OnStreamBegin",
		zap.Object("streamBegin", streamBegin),
	)
}

func (conn *defaultConn) OnStreamEOF(ctx context.Context, streamEOF StreamEOF) {
	conn.Logger().Debug(
		"OnStreamEOF",
		zap.Object("streamEOF", streamEOF),
	)
}

func (conn *defaultConn) OnStreamDry(ctx context.Context, streamDry StreamDry) {
	conn.Logger().Debug(
		"OnStreamDry",
		zap.Object("streamDry", streamDry),
	)
}

func (conn *defaultConn) OnSetBufferLength(ctx context.Context, setBufferLength SetBufferLength) {
	conn.Logger().Debug(
		"OnSetBufferLength",
		zap.Object("setBufferLength", setBufferLength),
	)
}

func (conn *defaultConn) OnStreamIsRecorded(ctx context.Context, streamIsRecorded StreamIsRecorded) {
	conn.Logger().Debug(
		"OnStreamIsRecorded",
		zap.Object("streamIsRecorded", streamIsRecorded),
	)
}

func (conn *defaultConn) OnPingRequest(ctx context.Context, pingRequest PingRequest) {
	conn.Logger().Debug(
		"OnPingRequest",
		zap.Object("pingRequest", pingRequest),
	)
	if err := conn.PingResponse(ctx, pingRequest.Timestamp()); err != nil {
		conn.Logger().Error(
			"failed to send PingResponse",
			zap.Error(err),
		)
		return
	}
}

func (conn *defaultConn) OnPingResponse(ctx context.Context, pingResponse PingResponse) {
	conn.Logger().Debug(
		"OnPingResponse",
		zap.Object("pingResponse", pingResponse),
	)
}
