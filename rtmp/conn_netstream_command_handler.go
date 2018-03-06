package rtmp

import (
	"context"

	"go.uber.org/zap"
)

func (conn *defaultConn) OnOnStatus(ctx context.Context, onStatus OnStatus) {
	conn.Logger().Warn(
		"OnOnStatus is not implemented",
		zap.Object("onStatus", onStatus),
	)
}

func (conn *defaultConn) OnPlay(ctx context.Context, play Play) {
	conn.Logger().Warn(
		"OnPlay is not implemented",
		zap.Object("play", play),
	)
}

func (conn *defaultConn) OnPlay2(ctx context.Context, play2 Play2) {
	conn.Logger().Warn(
		"OnPlay2 is not implemented",
		zap.Object("play2", play2),
	)
}

func (conn *defaultConn) OnDeleteStream(ctx context.Context, deleteStream DeleteStream) {
	conn.Logger().Warn(
		"OnDeleteStream is not implemented",
		zap.Object("deleteStream", deleteStream),
	)
}

func (conn *defaultConn) OnCloseStream(ctx context.Context, closeStream CloseStream) {
	conn.Logger().Warn(
		"OnCloseStream is not implemented",
		zap.Object("closeStream", closeStream),
	)
}

func (conn *defaultConn) OnReceiveAudio(ctx context.Context, receiveAudio ReceiveAudio) {
	conn.Logger().Warn(
		"OnReceiveAudio is not implemented",
		zap.Object("receiveAudio", receiveAudio),
	)
}

func (conn *defaultConn) OnReceiveVideo(ctx context.Context, receiveVideo ReceiveVideo) {
	conn.Logger().Warn(
		"OnReceiveVideo is not implemented",
		zap.Object("receiveVideo", receiveVideo),
	)
}

func (conn *defaultConn) OnPublish(ctx context.Context, publish Publish) {
	conn.Logger().Info(
		"OnPublish",
		zap.Object("publish", publish),
	)
	csID := ChunkStreamIDFromContext(ctx)
	if err := conn.StreamBegin(ctx, csID); err != nil {
		conn.Logger().Error(
			"failed to StreamBegin",
			zap.Object("publish", publish),
			zap.Error(err),
		)
	}
	if err := conn.OnStatus(
		ctx,
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

func (conn *defaultConn) OnSeek(ctx context.Context, seek Seek) {
	conn.Logger().Warn(
		"OnSeek is not implemented",
		zap.Object("seek", seek),
	)
}

func (conn *defaultConn) OnPause(ctx context.Context, pause Pause) {
	conn.Logger().Warn(
		"OnPause is not implemented",
		zap.Object("pause", pause),
	)
}
