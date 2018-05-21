package rtmp

import (
	"context"

	"go.uber.org/zap"
)

func (conn *defaultConn) OnSetChunkSize(ctx context.Context, setChunkSize SetChunkSize) ConnError {
	conn.Logger().Debug(
		"OnSetChunkSize",
		zap.Object("setChunkSize", setChunkSize),
	)
	conn.Reader().SetChunkSize(setChunkSize.ChunkSize())
	return nil
}

func (conn *defaultConn) OnAbortMessage(ctx context.Context, abortMessage AbortMessage) ConnError {
	conn.Logger().Debug(
		"OnAbortMessage",
		zap.Object("abortMessage", abortMessage),
	)
	conn.Reader().AbortMessage(abortMessage.ChunkStreamID())
	return nil
}

func (conn *defaultConn) OnAcknowledgement(ctx context.Context, acknowledgement Acknowledgement) ConnError {
	conn.Logger().Debug(
		"OnAcknowledgement",
		zap.Object("acknowledgement", acknowledgement),
	)
	return nil
}

func (conn *defaultConn) OnWindowAcknowledgementSize(ctx context.Context, windowAcknowledgementSize WindowAcknowledgementSize) ConnError {
	conn.Logger().Debug(
		"OnWindowAcknowledgementSize",
		zap.Object("windowAcknowledgementSize", windowAcknowledgementSize),
	)
	conn.Reader().SetAcknowledgementWindowSize(windowAcknowledgementSize.AcknowledgementWindowSize())
	return nil
}

func (conn *defaultConn) OnSetPeerBandwidth(ctx context.Context, setPeerBandwidth SetPeerBandwidth) ConnError {
	conn.Logger().Debug(
		"invoke OnSetPeerBandwidth",
		zap.Object("setPeerBandwidth", setPeerBandwidth),
	)
	conn.Reader().SetBandwidthLimitType(setPeerBandwidth.LimitType())
	conn.Reader().SetAcknowledgementWindowSize(setPeerBandwidth.AcknowledgmentWindowSize())
	return nil
}
