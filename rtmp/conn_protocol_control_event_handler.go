package rtmp

import (
	"context"

	"go.uber.org/zap"
)

func (conn *defaultConn) OnSetChunkSize(ctx context.Context, setChunkSize SetChunkSize) {
	conn.Logger().Debug(
		"invoke OnSetChunkSize",
		zap.Object("setChunkSize", setChunkSize),
	)
	conn.Reader().SetChunkSize(setChunkSize.ChunkSize())
}

func (conn *defaultConn) OnAbortMessage(ctx context.Context, abortMessage AbortMessage) {
	conn.Logger().Debug(
		"invoke OnAbortMessage",
		zap.Object("abortMessage", abortMessage),
	)
	conn.Reader().AbortMessage(abortMessage.ChunkStreamID())
}

func (conn *defaultConn) OnAcknowledgement(ctx context.Context, acknowledgement Acknowledgement) {
	conn.Logger().Debug(
		"invoke OnAcknowledgement",
		zap.Object("acknowledgement", acknowledgement),
	)
}

func (conn *defaultConn) OnWindowAcknowledgementSize(ctx context.Context, windowAcknowledgementSize WindowAcknowledgementSize) {
	conn.Logger().Warn(
		"OnWindowAcknowledgementSize is not implemented",
		zap.Object("windowAcknowledgementSize", windowAcknowledgementSize),
	)
}

func (conn *defaultConn) OnSetPeerBandwidth(ctx context.Context, setPeerBandwidth SetPeerBandwidth) {
	conn.Logger().Debug(
		"invoke OnSetPeerBandwidth",
		zap.Object("setPeerBandwidth", setPeerBandwidth),
	)
	conn.Reader().SetBandwidthLimitType(setPeerBandwidth.LimitType())
	conn.Reader().SetAcknowledgementWindowSize(setPeerBandwidth.AcknowledgmentWindowSize())
}
