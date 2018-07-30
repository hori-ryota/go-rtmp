package rtmp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (conn *defaultConn) DefaultNetStreamCommandHandler() NetStreamCommandHandler {
	return NetStreamCommandHandler{
		OnStatusHandlers: []OnStatusHandler{
			OnStatusHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, onStatus OnStatus) ConnError {
				conn.logger.Info(
					"OnOnStatus",
					zap.Object("onStatus", onStatus),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				warnErrors := make([]error, 0, len(conn.netStreamCommandCallbacks))
				for _, f := range conn.netStreamCommandCallbacks {
					if err := f(onStatus); err != nil {
						if IsConnWarnError(err) {
							warnErrors = append(warnErrors, err)
						} else {
							return err
						}
					}
				}
				conn.netStreamCommandCallbacks = conn.netStreamCommandCallbacks[:0]
				if len(warnErrors) == 0 {
					return nil
				}
				return NewConnWarnError(
					errors.New("error on netStreamCommandCallbacks"),
					zap.Errors("callback errors", warnErrors),
				)
			}),
		},

		PlayHandlers: []PlayHandler{
			PlayHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, play Play) ConnError {
				conn.logger.Debug(
					"OnPlay",
					zap.Object("play", play),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},

		Play2Handlers: []Play2Handler{
			Play2HandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, play2 Play2) ConnError {
				conn.logger.Debug(
					"OnPlay2",
					zap.Object("play2", play2),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},

		DeleteStreamHandlers: []DeleteStreamHandler{
			DeleteStreamHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, deleteStream DeleteStream) ConnError {
				conn.logger.Debug(
					"OnDeleteStream",
					zap.Object("deleteStream", deleteStream),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},

		CloseStreamHandlers: []CloseStreamHandler{
			CloseStreamHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, closeStream CloseStream) ConnError {
				conn.logger.Debug(
					"OnCloseStream",
					zap.Object("closeStream", closeStream),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},

		ReceiveAudioHandlers: []ReceiveAudioHandler{
			ReceiveAudioHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, receiveAudio ReceiveAudio) ConnError {
				conn.logger.Debug(
					"OnReceiveAudio",
					zap.Object("receiveAudio", receiveAudio),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},

		ReceiveVideoHandlers: []ReceiveVideoHandler{
			ReceiveVideoHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, receiveVideo ReceiveVideo) ConnError {
				conn.logger.Debug(
					"OnReceiveVideo",
					zap.Object("receiveVideo", receiveVideo),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},

		PublishHandlers: []PublishHandler{
			PublishHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, publish Publish) ConnError {
				conn.logger.Info(
					"OnPublish",
					zap.Object("publish", publish),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				for _, v := range conn.onPublishValidators {
					if onPublishError := v(ctx, publish); onPublishError != nil {
						if err := conn.OnStatus(
							ctx,
							chunkStreamID,
							messageStreamID,
							onPublishError,
						); err != nil {
							return NewConnFatalError(
								errors.Wrap(err, "failed to OnStatus"),
								zap.Object("publish", publish),
							)
						}
						return NewConnRejectedError(
							errors.New("publish request is rejected"),
							zap.Any("publishError", onPublishError),
						)
					}
				}
				if err := conn.StreamBegin(ctx, chunkStreamID); err != nil {
					return NewConnFatalError(
						errors.Wrap(err, "failed to StreamBegin"),
						zap.Object("publish", publish),
					)
				}
				if err := conn.OnStatus(
					ctx,
					chunkStreamID,
					messageStreamID,
					map[string]interface{}{
						"level": "status",
						"code":  "NetStream.Publish.Start",
					},
				); err != nil {
					return NewConnFatalError(
						errors.Wrap(err, "failed to OnStatus"),
						zap.Object("publish", publish),
					)
				}
				return nil
			}),
		},

		SeekHandlers: []SeekHandler{
			SeekHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, seek Seek) ConnError {
				conn.logger.Debug(
					"OnSeek",
					zap.Object("seek", seek),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},

		PauseHandlers: []PauseHandler{
			PauseHandlerFunc(func(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, pause Pause) ConnError {
				conn.logger.Debug(
					"OnPause",
					zap.Object("pause", pause),
					zap.Uint32("chunkStreamID", chunkStreamID),
					zap.Uint32("messageStreamID", messageStreamID),
				)
				return nil
			}),
		},
	}
}
