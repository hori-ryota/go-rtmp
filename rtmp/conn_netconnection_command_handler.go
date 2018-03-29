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
	for _, v := range conn.onConnectValidators {
		if onConnectError := v(ctx, connect); onConnectError != nil {
			if err := conn.ConnectError(ctx, onConnectError.Properties(), onConnectError.Information()); err != nil {
				conn.Logger().Error(
					"failed to ConnectResult",
					zap.Object("connect", connect),
					zap.Error(err),
				)
			}
			if err := conn.Close(); err != nil {
				conn.Logger().Error(
					"failed to Close conn when OnConnect",
					zap.Object("connect", connect),
					zap.Error(err),
				)
			}
			return
		}
	}
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
	if err := conn.ConnectResult(ctx, nil, nil); err != nil {
		conn.Logger().Error(
			"failed to ConnectResult",
			zap.Object("connect", connect),
			zap.Error(err),
		)
	}
}

func (conn *defaultConn) OnConnectResult(ctx context.Context, connectResult ConnectResult) {
	conn.Logger().Debug(
		"OnConnectResult is not implemented",
		zap.Object("connectResult", connectResult),
	)
}

func (conn *defaultConn) OnConnectError(ctx context.Context, connectError ConnectError) {
	conn.Logger().Debug(
		"OnConnectError is not implemented",
		zap.Object("connectError", connectError),
	)
}

func (conn *defaultConn) OnCall(ctx context.Context, call Call) {
	conn.Logger().Info(
		"OnCall",
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
	conn.Logger().Info(
		"OnCreateStreamResult",
		zap.Object("createStreamResult", createStreamResult),
	)
	if f, ok := conn.createStreamCallbacks[createStreamResult.TransactionID()]; ok {
		f(createStreamResult)
	}
	delete(conn.createStreamCallbacks, createStreamResult.TransactionID())
}

func (conn *defaultConn) OnCreateStreamError(ctx context.Context, createStreamError CreateStreamError) {
	conn.Logger().Info(
		"OnCreateStreamError",
		zap.Object("createStreamError", createStreamError),
	)
	if f, ok := conn.createStreamCallbacks[createStreamError.TransactionID()]; ok {
		f(createStreamError)
	}
	delete(conn.createStreamCallbacks, createStreamError.TransactionID())
}
