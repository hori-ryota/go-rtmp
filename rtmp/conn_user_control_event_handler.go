package rtmp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (conn *defaultConn) OnStreamBegin(ctx context.Context, streamBegin StreamBegin) ConnError {
	conn.Logger().Debug(
		"OnStreamBegin",
		zap.Object("streamBegin", streamBegin),
	)
	return nil
}

func (conn *defaultConn) OnStreamEOF(ctx context.Context, streamEOF StreamEOF) ConnError {
	conn.Logger().Debug(
		"OnStreamEOF",
		zap.Object("streamEOF", streamEOF),
	)
	return nil
}

func (conn *defaultConn) OnStreamDry(ctx context.Context, streamDry StreamDry) ConnError {
	conn.Logger().Debug(
		"OnStreamDry",
		zap.Object("streamDry", streamDry),
	)
	return nil
}

func (conn *defaultConn) OnSetBufferLength(ctx context.Context, setBufferLength SetBufferLength) ConnError {
	conn.Logger().Debug(
		"OnSetBufferLength",
		zap.Object("setBufferLength", setBufferLength),
	)
	return nil
}

func (conn *defaultConn) OnStreamIsRecorded(ctx context.Context, streamIsRecorded StreamIsRecorded) ConnError {
	conn.Logger().Debug(
		"OnStreamIsRecorded",
		zap.Object("streamIsRecorded", streamIsRecorded),
	)
	return nil
}

func (conn *defaultConn) OnPingRequest(ctx context.Context, pingRequest PingRequest) ConnError {
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
}

func (conn *defaultConn) OnPingResponse(ctx context.Context, pingResponse PingResponse) ConnError {
	conn.Logger().Debug(
		"OnPingResponse",
		zap.Object("pingResponse", pingResponse),
	)
	return nil
}
