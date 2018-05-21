package rtmp

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func (conn *defaultConn) OnOnStatus(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, onStatus OnStatus) ConnError {
	conn.Logger().Info(
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
}

func (conn *defaultConn) OnPlay(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, play Play) ConnError {
	return NewConnWarnError(
		errors.New("OnPlay is not implemented"),
		zap.Object("play", play),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnPlay2(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, play2 Play2) ConnError {
	return NewConnWarnError(
		errors.New("OnPlay2 is not implemented"),
		zap.Object("play2", play2),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnDeleteStream(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, deleteStream DeleteStream) ConnError {
	return NewConnWarnError(
		errors.New("OnDeleteStream is not implemented"),
		zap.Object("deleteStream", deleteStream),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnCloseStream(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, closeStream CloseStream) ConnError {
	return NewConnWarnError(
		errors.New("OnCloseStream is not implemented"),
		zap.Object("closeStream", closeStream),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnReceiveAudio(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, receiveAudio ReceiveAudio) ConnError {
	return NewConnWarnError(
		errors.New("OnReceiveAudio is not implemented"),
		zap.Object("receiveAudio", receiveAudio),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnReceiveVideo(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, receiveVideo ReceiveVideo) ConnError {
	return NewConnWarnError(
		errors.New("OnReceiveVideo is not implemented"),
		zap.Object("receiveVideo", receiveVideo),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnPublish(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, publish Publish) ConnError {
	conn.Logger().Info(
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
}

func (conn *defaultConn) OnSeek(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, seek Seek) ConnError {
	return NewConnWarnError(
		errors.New("OnSeek is not implemented"),
		zap.Object("seek", seek),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnPause(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, pause Pause) ConnError {
	return NewConnWarnError(
		errors.New("OnPause is not implemented"),
		zap.Object("pause", pause),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}
