package rtmp

import (
	"context"

	"go.uber.org/zap"
)

func (conn *defaultConn) OnConnect(ctx context.Context, connect Connect) {
	conn.Logger().Info(
		"OnConnect",
		zap.Object("connect", connect),
	)
	if err := conn.WindowAcknowledgementSize(ctx, conn.windowAcknowledgementSize); err != nil {
		conn.Logger().Error(
			"failed to WindowAcknowledgementSize",
			zap.Object("connect", connect),
			zap.Error(err),
		)
	}
	if err := conn.SetPeerBandwidth(ctx, conn.windowAcknowledgementSize, conn.bandwidthLimitType); err != nil {
		conn.Logger().Error(
			"failed to SetPeerBandwidth",
			zap.Object("connect", connect),
			zap.Error(err),
		)
	}
	if err := conn.StreamBegin(ctx, 0); err != nil {
		conn.Logger().Error(
			"failed to StreamBegin",
			zap.Object("connect", connect),
			zap.Error(err),
		)
	}
	// TODO: auth check
	if err := conn.ConnectResult(ctx, nil, nil); err != nil {
		conn.Logger().Error(
			"failed to ConnectResult",
			zap.Object("connect", connect),
			zap.Error(err),
		)
	}
}

func (conn *defaultConn) OnConnectResult(ctx context.Context, connectResult ConnectResult) {
	conn.Logger().Warn(
		"OnConnectResult is not implemented",
		zap.Object("connectResult", connectResult),
	)
}

func (conn *defaultConn) OnConnectError(ctx context.Context, connectError ConnectError) {
	conn.Logger().Warn(
		"OnConnectError is not implemented",
		zap.Object("connectError", connectError),
	)
}

func (conn *defaultConn) OnCall(ctx context.Context, call Call) {
	conn.Logger().Info(
		"invoke OnCall",
		zap.Object("call", call),
	)
}

func (conn *defaultConn) OnCallResponse(ctx context.Context, callResponse CallResponse) {
	conn.Logger().Warn(
		"OnCallResponse is not implemented",
		zap.Object("callResponse", callResponse),
	)
}

func (conn *defaultConn) OnClose(ctx context.Context, close Close) {
	conn.Logger().Warn(
		"OnClose is not implemented",
		zap.Object("close", close),
	)
}

func (conn *defaultConn) OnCreateStream(ctx context.Context, createStream CreateStream) {
	chunkStreamID := conn.Reader().CreateStream()
	if err := conn.CreateStreamResult(ctx, createStream.TransactionID(), createStream.CommandObject(), chunkStreamID); err != nil {
		conn.Logger().Error(
			"failed to CreateStreamResult",
			zap.Object("createStream", createStream),
			zap.Error(err),
		)
	}
}

func (conn *defaultConn) OnCreateStreamResult(ctx context.Context, createStreamResult CreateStreamResult) {
	conn.Logger().Warn(
		"OnCreateStreamResult is not implemented",
		zap.Object("createStreamResult", createStreamResult),
	)
}

func (conn *defaultConn) OnCreateStreamError(ctx context.Context, createStreamError CreateStreamError) {
	conn.Logger().Warn(
		"OnCreateStreamError is not implemented",
		zap.Object("createStreamError", createStreamError),
	)
}
