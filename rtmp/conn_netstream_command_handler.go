package rtmp

import (
	"context"

	"go.uber.org/zap"
)

func (conn *defaultConn) OnOnStatus(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, onStatus OnStatus) {
	conn.Logger().Info(
		"OnOnStatus",
		zap.Object("onStatus", onStatus),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
	for _, f := range conn.netStreamCommandCallbacks {
		f(onStatus)
	}
	conn.netStreamCommandCallbacks = conn.netStreamCommandCallbacks[:0]
}

func (conn *defaultConn) OnPlay(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, play Play) {
	conn.Logger().Warn(
		"OnPlay is not implemented",
		zap.Object("play", play),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnPlay2(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, play2 Play2) {
	conn.Logger().Warn(
		"OnPlay2 is not implemented",
		zap.Object("play2", play2),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnDeleteStream(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, deleteStream DeleteStream) {
	conn.Logger().Warn(
		"OnDeleteStream is not implemented",
		zap.Object("deleteStream", deleteStream),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnCloseStream(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, closeStream CloseStream) {
	conn.Logger().Warn(
		"OnCloseStream is not implemented",
		zap.Object("closeStream", closeStream),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnReceiveAudio(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, receiveAudio ReceiveAudio) {
	conn.Logger().Warn(
		"OnReceiveAudio is not implemented",
		zap.Object("receiveAudio", receiveAudio),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnReceiveVideo(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, receiveVideo ReceiveVideo) {
	conn.Logger().Warn(
		"OnReceiveVideo is not implemented",
		zap.Object("receiveVideo", receiveVideo),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnPublish(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, publish Publish) {
	conn.Logger().Info(
		"OnPublish",
		zap.Object("publish", publish),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
	if err := conn.StreamBegin(ctx, chunkStreamID); err != nil {
		conn.Logger().Error(
			"failed to StreamBegin",
			zap.Object("publish", publish),
			zap.Error(err),
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
		conn.Logger().Error(
			"failed to OnStatus",
			zap.Object("publish", publish),
			zap.Error(err),
		)
	}
}

func (conn *defaultConn) OnSeek(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, seek Seek) {
	conn.Logger().Warn(
		"OnSeek is not implemented",
		zap.Object("seek", seek),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}

func (conn *defaultConn) OnPause(ctx context.Context, chunkStreamID uint32, messageStreamID uint32, pause Pause) {
	conn.Logger().Warn(
		"OnPause is not implemented",
		zap.Object("pause", pause),
		zap.Uint32("chunkStreamID", chunkStreamID),
		zap.Uint32("messageStreamID", messageStreamID),
	)
}
