package rtmp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (conn *defaultConn) OnConnect(ctx context.Context, connect Connect) ConnError {
	conn.Logger().Info(
		"OnConnect",
		zap.Object("connect", connect),
	)
	for _, v := range conn.onConnectValidators {
		if onConnectError := v(ctx, connect); onConnectError != nil {
			if err := conn.ConnectError(ctx, onConnectError.Properties(), onConnectError.Information()); err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to ConnectResult"),
					zap.Object("connect", connect),
				)
			}
			return NewConnRejectedError(
				errors.New("connect request is rejected"),
				zap.Object("connectError", onConnectError),
			)
		}
	}
	if err := conn.WindowAcknowledgementSize(ctx, conn.windowAcknowledgementSize); err != nil {
		return NewConnFatalError(
			errors.Wrap(err, "failed to WindowAcknowledgementSize"),
			zap.Object("connect", connect),
		)
	}
	if err := conn.SetPeerBandwidth(ctx, conn.windowAcknowledgementSize, conn.bandwidthLimitType); err != nil {
		return NewConnFatalError(
			errors.Wrap(err, "failed to SetPeerBandwidth"),
			zap.Object("connect", connect),
		)
	}
	if err := conn.StreamBegin(ctx, 0); err != nil {
		return NewConnFatalError(
			errors.Wrap(err, "failed to StreamBegin"),
			zap.Object("connect", connect),
		)
	}
	if err := conn.ConnectResult(ctx, nil, nil); err != nil {
		return NewConnFatalError(
			errors.Wrap(err, "failed to ConnectResult"),
			zap.Object("connect", connect),
		)
	}
	return nil
}

func (conn *defaultConn) OnConnectResult(ctx context.Context, connectResult ConnectResult) ConnError {
	return NewConnWarnError(
		errors.New("OnConnectResult is not implemented"),
		zap.Object("connectResult", connectResult),
	)
}

func (conn *defaultConn) OnConnectError(ctx context.Context, connectError ConnectError) ConnError {
	return NewConnWarnError(
		errors.New("OnConnectError is not implemented"),
		zap.Object("connectError", connectError),
	)
}

func (conn *defaultConn) OnCall(ctx context.Context, call Call) ConnError {
	return NewConnWarnError(
		errors.New("OnCall is not implemented"),
		zap.Object("call", call),
	)
}

func (conn *defaultConn) OnCallResponse(ctx context.Context, callResponse CallResponse) ConnError {
	return NewConnWarnError(
		errors.New("OnCallResponse is not implemented"),
		zap.Object("callResponse", callResponse),
	)
}

func (conn *defaultConn) OnClose(ctx context.Context, close Close) ConnError {
	return NewConnWarnError(
		errors.New("OnClose is not implemented"),
		zap.Object("close", close),
	)
}

func (conn *defaultConn) OnCreateStream(ctx context.Context, createStream CreateStream) ConnError {
	chunkStreamID := conn.Reader().CreateStream()
	if err := conn.CreateStreamResult(ctx, createStream.TransactionID(), createStream.CommandObject(), chunkStreamID); err != nil {
		return NewConnFatalError(
			errors.Wrap(err, "failed to CreateStreamResult"),
			zap.Object("createStream", createStream),
		)
	}
	return nil
}

func (conn *defaultConn) OnCreateStreamResult(ctx context.Context, createStreamResult CreateStreamResult) ConnError {
	conn.Logger().Info(
		"OnCreateStreamResult",
		zap.Object("createStreamResult", createStreamResult),
	)
	var warnError ConnError
	if f, ok := conn.createStreamCallbacks[createStreamResult.TransactionID()]; ok {
		if err := f(createStreamResult); err != nil {
			if IsConnWarnError(err) {
				warnError = err
			} else {
				return err
			}
		}
	}
	delete(conn.createStreamCallbacks, createStreamResult.TransactionID())
	return warnError
}

func (conn *defaultConn) OnCreateStreamError(ctx context.Context, createStreamError CreateStreamError) ConnError {
	conn.Logger().Info(
		"OnCreateStreamError",
		zap.Object("createStreamError", createStreamError),
	)
	var warnError ConnError
	if f, ok := conn.createStreamCallbacks[createStreamError.TransactionID()]; ok {
		if err := f(createStreamError); err != nil {
			if IsConnWarnError(err) {
				warnError = err
			} else {
				return err
			}
		}
	}
	delete(conn.createStreamCallbacks, createStreamError.TransactionID())
	return warnError
}
