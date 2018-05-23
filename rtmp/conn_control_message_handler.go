package rtmp

import (
	"bytes"
	"context"
	"encoding/binary"

	"github.com/pkg/errors"
	amf "github.com/zhangpeihao/goamf"
	"go.uber.org/zap"
)

type ControlMessageHandler struct {
	ProtocolControlEventHandler ProtocolControlEventHandler
	UserControlEventHandler     UserControlEventHandler
	NetConnectionCommandHandler NetConnectionCommandHandler
	NetStreamCommandHandler     NetStreamCommandHandler
}

func NewControlMessageHandler(conn Conn) *ControlMessageHandler {
	return &ControlMessageHandler{
		ProtocolControlEventHandler: conn.DefaultProtocolControlEventHandler(),
		UserControlEventHandler:     conn.DefaultUserControlEventHandler(),
		NetConnectionCommandHandler: conn.DefaultNetConnectionCommandHandler(),
		NetStreamCommandHandler:     conn.DefaultNetStreamCommandHandler(),
	}
}

func (h *ControlMessageHandler) HandleMessage(ctx context.Context, m Message) ConnError {
	switch m.TypeID() {
	case MessageTypeIDSetChunkSize:
		p, err := UnmarshalSetChunkSizeBinary(m.Payload())
		if err != nil {
			return NewConnFatalError(
				errors.Wrap(err, "failed to unmarshal SetChunkSize"),
				zap.Object("message", m),
			)
		}
		return h.ProtocolControlEventHandler.OnSetChunkSize(ctx, p)
	case MessageTypeIDAbortMessage:
		p, err := UnmarshalAbortMessageBinary(m.Payload())
		if err != nil {
			return NewConnFatalError(
				errors.Wrap(err, "failed to unmarshal AbortMessage"),
				zap.Object("message", m),
			)
		}
		return h.ProtocolControlEventHandler.OnAbortMessage(ctx, p)
	case MessageTypeIDAcknowledgement:
		p, err := UnmarshalAcknowledgementBinary(m.Payload())
		if err != nil {
			return NewConnFatalError(
				errors.Wrap(err, "failed to unmarshal Acknowledgement"),
				zap.Object("message", m),
			)
		}
		return h.ProtocolControlEventHandler.OnAcknowledgement(ctx, p)
	case MessageTypeIDUserControlMessages:
		b := m.Payload()
		eventType := EventType(binary.BigEndian.Uint16(b[:2]))
		switch eventType {
		case EventTypeStreamBegin:
			p, err := UnmarshalStreamBeginBinary(m.Payload())
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal StreamBegin"),
					zap.Object("message", m),
				)
			}
			return h.UserControlEventHandler.OnStreamBegin(ctx, p)
		case EventTypeStreamEOF:
			p, err := UnmarshalStreamEOFBinary(m.Payload())
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal StreamEOF"),
					zap.Object("message", m),
				)
			}
			return h.UserControlEventHandler.OnStreamEOF(ctx, p)
		case EventTypeStreamDry:
			p, err := UnmarshalStreamDryBinary(m.Payload())
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal StreamDry"),
					zap.Object("message", m),
				)
			}
			return h.UserControlEventHandler.OnStreamDry(ctx, p)
		case EventTypeSetBufferLength:
			p, err := UnmarshalSetBufferLengthBinary(m.Payload())
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal SetBufferLength"),
					zap.Object("message", m),
				)
			}
			return h.UserControlEventHandler.OnSetBufferLength(ctx, p)
		case EventTypeStreamIsRecorded:
			p, err := UnmarshalStreamIsRecordedBinary(m.Payload())
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal StreamIsRecorded"),
					zap.Object("message", m),
				)
			}
			return h.UserControlEventHandler.OnStreamIsRecorded(ctx, p)
		case EventTypePingRequest:
			p, err := UnmarshalPingRequestBinary(m.Payload())
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal PingRequest"),
					zap.Object("message", m),
				)
			}
			return h.UserControlEventHandler.OnPingRequest(ctx, p)
		case EventTypePingResponse:
			p, err := UnmarshalPingResponseBinary(m.Payload())
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal PingResponse"),
					zap.Object("message", m),
				)
			}
			return h.UserControlEventHandler.OnPingResponse(ctx, p)
		}
	case MessageTypeIDWindowAcknowledgementSize:
		p, err := UnmarshalWindowAcknowledgementSizeBinary(m.Payload())
		if err != nil {
			return NewConnFatalError(
				errors.Wrap(err, "failed to unmarshal WindowAcknowledgementSize"),
				zap.Object("message", m),
			)
		}
		return h.ProtocolControlEventHandler.OnWindowAcknowledgementSize(ctx, p)
	case MessageTypeIDSetPeerBandwidth:
		p, err := UnmarshalSetPeerBandwidthBinary(m.Payload())
		if err != nil {
			return NewConnFatalError(
				errors.Wrap(err, "failed to unmarshal SetPeerBandwidth"),
				zap.Object("message", m),
			)
		}
		return h.ProtocolControlEventHandler.OnSetPeerBandwidth(ctx, p)
	case MessageTypeIDCommandAMF0, MessageTypeIDCommandAMF3:
		var name string
		var encodingAMFType EncodingAMFType
		var err error
		b := m.Payload()
		r := bytes.NewReader(b)
		if m.TypeID() == MessageTypeIDCommandAMF0 {
			encodingAMFType = EncodingAMFTypeAMF0
			name, err = amf.ReadString(r)
		} else {
			encodingAMFType = EncodingAMFTypeAMF3
			name, err = amf.AMF3_ReadString(r)
		}
		if err != nil {
			return NewConnFatalError(
				errors.Wrap(err, "failed to read command name"),
				zap.Object("message", m),
			)
		}
		switch name {
		case "connect":
			p, err := UnmarshalConnectBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Connect"),
					zap.Object("message", m),
				)
			}
			return h.NetConnectionCommandHandler.OnConnect(ctx, p)
		case "createStream":
			p, err := UnmarshalCreateStreamBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal CreateStream"),
					zap.Object("message", m),
				)
			}
			return h.NetConnectionCommandHandler.OnCreateStream(ctx, p)
		case "close":
			p, err := UnmarshalCloseBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Close"),
					zap.Object("message", m),
				)
			}
			return h.NetConnectionCommandHandler.OnClose(ctx, p)
		case "_result", "_error":
			var transactionID float64
			if encodingAMFType == EncodingAMFTypeAMF0 {
				transactionID, err = amf.ReadDouble(r)
			} else {
				transactionID, err = amf.AMF3_ReadDouble(r)
			}
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to read transactionID"),
					zap.Object("message", m),
				)
			}
			if transactionID == 1 {
				// connect response
				if name == "_result" {
					p, err := UnmarshalConnectResultBinary(b, encodingAMFType)
					if err != nil {
						return NewConnFatalError(
							errors.Wrap(err, "failed to unmarshal ConnectResult"),
							zap.Object("message", m),
						)
					}
					return h.NetConnectionCommandHandler.OnConnectResult(ctx, p)
				}
				p, err := UnmarshalConnectErrorBinary(b, encodingAMFType)
				if err != nil {
					return NewConnFatalError(
						errors.Wrap(err, "failed to unmarshal ConnectError"),
						zap.Object("message", m),
					)
				}
				return h.NetConnectionCommandHandler.OnConnectError(ctx, p)
			}
			// create stream response
			if name == "_result" {
				p, err := UnmarshalCreateStreamResultBinary(b, encodingAMFType)
				if err != nil {
					return NewConnFatalError(
						errors.Wrap(err, "failed to unmarshal CreateStreamResult"),
						zap.Object("message", m),
					)
				}
				return h.NetConnectionCommandHandler.OnCreateStreamResult(ctx, p)
			}
			p, err := UnmarshalCreateStreamErrorBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal CreateStreamError"),
					zap.Object("message", m),
				)
			}
			return h.NetConnectionCommandHandler.OnCreateStreamError(ctx, p)
		case "onStatus":
			p, err := UnmarshalOnStatusBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal OnStatus"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnOnStatus(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "play":
			p, err := UnmarshalPlayBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Play"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnPlay(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "play2":
			p, err := UnmarshalPlay2Binary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Play2"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnPlay2(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "deleteStream":
			p, err := UnmarshalDeleteStreamBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal DeleteStream"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnDeleteStream(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "closeStream":
			p, err := UnmarshalCloseStreamBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal CloseStream"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnCloseStream(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "receiveAudio":
			p, err := UnmarshalReceiveAudioBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal ReceiveAudio"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnReceiveAudio(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "receiveVideo":
			p, err := UnmarshalReceiveVideoBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal ReceiveVideo"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnReceiveVideo(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "publish":
			p, err := UnmarshalPublishBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Publish"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnPublish(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "seek":
			p, err := UnmarshalSeekBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Seek"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnSeek(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "pause":
			p, err := UnmarshalPauseBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Pause"),
					zap.Object("message", m),
				)
			}
			return h.NetStreamCommandHandler.OnPause(ctx, m.ChunkStreamID(), m.StreamID(), p)
		default:
			p, err := UnmarshalCallBinary(b, encodingAMFType)
			if err != nil {
				return NewConnFatalError(
					errors.Wrap(err, "failed to unmarshal Call"),
					zap.Object("message", m),
				)
			}
			return h.NetConnectionCommandHandler.OnCall(ctx, p)
		}
	case MessageTypeIDAggregate:
		// TODO
		return NewConnWarnError(
			errors.New("not implemented"),
			zap.Object("message", m),
		)
	}
	return nil
}
