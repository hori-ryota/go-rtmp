package rtmp

import (
	"context"

	"go.uber.org/zap"
)

func (conn *defaultConn) DefaultProtocolControlEventHandler() ProtocolControlEventHandler {
	return ProtocolControlEventHandler{
		SetChunkSizeHandlers: []SetChunkSizeHandler{
			SetChunkSizeHandlerFunc(func(ctx context.Context, setChunkSize SetChunkSize) ConnError {
				conn.Logger().Debug(
					"OnSetChunkSize",
					zap.Object("setChunkSize", setChunkSize),
				)
				conn.Reader().SetChunkSize(setChunkSize.ChunkSize())
				return nil
			}),
		},

		AbortMessageHandlers: []AbortMessageHandler{
			AbortMessageHandlerFunc(func(ctx context.Context, abortMessage AbortMessage) ConnError {
				conn.Logger().Debug(
					"OnAbortMessage",
					zap.Object("abortMessage", abortMessage),
				)
				conn.Reader().AbortMessage(abortMessage.ChunkStreamID())
				return nil
			}),
		},

		AcknowledgementHandlers: []AcknowledgementHandler{
			AcknowledgementHandlerFunc(func(ctx context.Context, acknowledgement Acknowledgement) ConnError {
				conn.Logger().Debug(
					"OnAcknowledgement",
					zap.Object("acknowledgement", acknowledgement),
				)
				return nil
			}),
		},

		WindowAcknowledgementSizeHandlers: []WindowAcknowledgementSizeHandler{
			WindowAcknowledgementSizeHandlerFunc(func(ctx context.Context, windowAcknowledgementSize WindowAcknowledgementSize) ConnError {
				conn.Logger().Debug(
					"OnWindowAcknowledgementSize",
					zap.Object("windowAcknowledgementSize", windowAcknowledgementSize),
				)
				conn.Reader().SetAcknowledgementWindowSize(windowAcknowledgementSize.AcknowledgementWindowSize())
				return nil
			}),
		},

		SetPeerBandwidthHandlers: []SetPeerBandwidthHandler{
			SetPeerBandwidthHandlerFunc(func(ctx context.Context, setPeerBandwidth SetPeerBandwidth) ConnError {
				conn.Logger().Debug(
					"invoke OnSetPeerBandwidth",
					zap.Object("setPeerBandwidth", setPeerBandwidth),
				)
				conn.Reader().SetBandwidthLimitType(setPeerBandwidth.LimitType())
				conn.Reader().SetAcknowledgementWindowSize(setPeerBandwidth.AcknowledgmentWindowSize())
				return nil
			}),
		},
	}
}
