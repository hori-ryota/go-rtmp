package rtmp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (conn *defaultConn) DefaultNetConnectionCommandHandler() NetConnectionCommandHandler {
	return NetConnectionCommandHandler{
		ConnectHandlers: []ConnectHandler{
			ConnectHandlerFunc(func(ctx context.Context, connect Connect) ConnError {
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
			}),
		},

		ConnectResultHandlers: []ConnectResultHandler{
			ConnectResultHandlerFunc(func(ctx context.Context, connectResult ConnectResult) ConnError {
				conn.Logger().Debug(
					"OnConnectResult",
					zap.Object("connectResult", connectResult),
				)
				return nil
			}),
		},

		ConnectErrorHandlers: []ConnectErrorHandler{
			ConnectErrorHandlerFunc(func(ctx context.Context, connectError ConnectError) ConnError {
				conn.Logger().Debug(
					"OnConnectError",
					zap.Object("connectError", connectError),
				)
				return nil
			}),
		},

		CallHandlers: []CallHandler{
			CallHandlerFunc(func(ctx context.Context, call Call) ConnError {
				conn.Logger().Debug(
					"OnCall",
					zap.Object("call", call),
				)
				return nil
			}),
		},

		CallResponseHandlers: []CallResponseHandler{
			CallResponseHandlerFunc(func(ctx context.Context, callResponse CallResponse) ConnError {
				conn.Logger().Debug(
					"OnCallResponse",
					zap.Object("callResponse", callResponse),
				)
				return nil
			}),
		},

		CloseHandlers: []CloseHandler{
			CloseHandlerFunc(func(ctx context.Context, close Close) ConnError {
				conn.Logger().Debug(
					"OnClose",
					zap.Object("close", close),
				)
				return nil
			}),
		},

		CreateStreamHandlers: []CreateStreamHandler{
			CreateStreamHandlerFunc(func(ctx context.Context, createStream CreateStream) ConnError {
				chunkStreamID := conn.Reader().CreateStream()
				if err := conn.CreateStreamResult(ctx, createStream.TransactionID(), createStream.CommandObject(), chunkStreamID); err != nil {
					return NewConnFatalError(
						errors.Wrap(err, "failed to CreateStreamResult"),
						zap.Object("createStream", createStream),
					)
				}
				return nil
			}),
		},

		CreateStreamResultHandlers: []CreateStreamResultHandler{
			CreateStreamResultHandlerFunc(func(ctx context.Context, createStreamResult CreateStreamResult) ConnError {
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
			}),
		},

		CreateStreamErrorHandlers: []CreateStreamErrorHandler{
			CreateStreamErrorHandlerFunc(func(ctx context.Context, createStreamError CreateStreamError) ConnError {
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
			}),
		},
	}
}
