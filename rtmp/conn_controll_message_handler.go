package rtmp

import (
	"bytes"
	"context"
	"encoding/binary"

	amf "github.com/zhangpeihao/goamf"
	"go.uber.org/zap"
)

type ConnControlMessageHandler struct {
	conn Conn
}

func NewConnControllMessageHandler(conn Conn) MessageHandler {
	return &ConnControlMessageHandler{
		conn: conn,
	}
}

func (h ConnControlMessageHandler) HandleMessage(ctx context.Context, m Message) {
	switch m.TypeID() {
	case MessageTypeIDSetChunkSize:
		p, err := UnmarshalSetChunkSizeBinary(m.Payload())
		if err != nil {
			h.Error("failed to unmarshal SetChunkSize", m, err)
			return
		}
		h.conn.OnSetChunkSize(ctx, p)
	case MessageTypeIDAbortMessage:
		p, err := UnmarshalAbortMessageBinary(m.Payload())
		if err != nil {
			h.Error("failed to unmarshal AbortMessage", m, err)
			return
		}
		h.conn.OnAbortMessage(ctx, p)
	case MessageTypeIDAcknowledgement:
		p, err := UnmarshalAcknowledgementBinary(m.Payload())
		if err != nil {
			h.Error("failed to unmarshal Acknowledgement", m, err)
			return
		}
		h.conn.OnAcknowledgement(ctx, p)
	case MessageTypeIDUserControlMessages:
		b := m.Payload()
		eventType := EventType(binary.BigEndian.Uint16(b[:2]))
		switch eventType {
		case EventTypeStreamBegin:
			p, err := UnmarshalStreamBeginBinary(m.Payload())
			if err != nil {
				h.Error("failed to unmarshal StreamBegin", m, err)
				return
			}
			h.conn.OnStreamBegin(ctx, p)
		case EventTypeStreamEOF:
			p, err := UnmarshalStreamEOFBinary(m.Payload())
			if err != nil {
				h.Error("failed to unmarshal StreamEOF", m, err)
				return
			}
			h.conn.OnStreamEOF(ctx, p)
		case EventTypeStreamDry:
			p, err := UnmarshalStreamDryBinary(m.Payload())
			if err != nil {
				h.Error("failed to unmarshal StreamDry", m, err)
				return
			}
			h.conn.OnStreamDry(ctx, p)
		case EventTypeSetBufferLength:
			p, err := UnmarshalSetBufferLengthBinary(m.Payload())
			if err != nil {
				h.Error("failed to unmarshal SetBufferLength", m, err)
				return
			}
			h.conn.OnSetBufferLength(ctx, p)
		case EventTypeStreamIsRecorded:
			p, err := UnmarshalStreamIsRecordedBinary(m.Payload())
			if err != nil {
				h.Error("failed to unmarshal StreamIsRecorded", m, err)
				return
			}
			h.conn.OnStreamIsRecorded(ctx, p)
		case EventTypePingRequest:
			p, err := UnmarshalPingRequestBinary(m.Payload())
			if err != nil {
				h.Error("failed to unmarshal PingRequest", m, err)
				return
			}
			h.conn.OnPingRequest(ctx, p)
		case EventTypePingResponse:
			p, err := UnmarshalPingResponseBinary(m.Payload())
			if err != nil {
				h.Error("failed to unmarshal PingResponse", m, err)
				return
			}
			h.conn.OnPingResponse(ctx, p)
		}
	case MessageTypeIDWindowAcknowledgementSize:
		p, err := UnmarshalWindowAcknowledgementSizeBinary(m.Payload())
		if err != nil {
			h.Error("failed to unmarshal WindowAcknowledgementSize", m, err)
			return
		}
		h.conn.OnWindowAcknowledgementSize(ctx, p)
	case MessageTypeIDSetPeerBandwidth:
		p, err := UnmarshalSetPeerBandwidthBinary(m.Payload())
		if err != nil {
			h.Error("failed to unmarshal SetPeerBandwidth", m, err)
			return
		}
		h.conn.OnSetPeerBandwidth(ctx, p)
	case MessageTypeIDAudio:
		// TODO
	case MessageTypeIDVideo:
		// TODO
	case MessageTypeIDDataAMF3:
		// TODO
	case MessageTypeIDSharedObjectAMF3:
		// TODO
	case MessageTypeIDDataAMF0:
		// TODO
	case MessageTypeIDSharedObjectAMF0:
		// TODO
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
			h.Error("failed to read command name", m, err)
			return
		}
		switch name {
		case "connect":
			p, err := UnmarshalConnectBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Connect", m, err)
				return
			}
			h.conn.OnConnect(ctx, p)
		case "createStream":
			p, err := UnmarshalCreateStreamBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal CreateStream", m, err)
				return
			}
			h.conn.OnCreateStream(ctx, p)
		case "close":
			p, err := UnmarshalCloseBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Close", m, err)
				return
			}
			h.conn.OnClose(ctx, p)
		case "_result", "_error":
			var transactionID float64
			if encodingAMFType == EncodingAMFTypeAMF0 {
				transactionID, err = amf.ReadDouble(r)
			} else {
				transactionID, err = amf.AMF3_ReadDouble(r)
			}
			if err != nil {
				h.Error("failed to read transactionID", m, err)
				return
			}
			if transactionID == 1 {
				// connect response
				if name == "_result" {
					p, err := UnmarshalConnectResultBinary(b, encodingAMFType)
					if err != nil {
						h.Error("failed to unmarshal ConnectResult", m, err)
						return
					}
					h.conn.OnConnectResult(ctx, p)
					return
				}
				p, err := UnmarshalConnectErrorBinary(b, encodingAMFType)
				if err != nil {
					h.Error("failed to unmarshal ConnectError", m, err)
					return
				}
				h.conn.OnConnectError(ctx, p)
				return
			}
			// create stream response
			if name == "_result" {
				p, err := UnmarshalCreateStreamResultBinary(b, encodingAMFType)
				if err != nil {
					h.Error("failed to unmarshal CreateStreamResult", m, err)
					return
				}
				h.conn.OnCreateStreamResult(ctx, p)
				return
			}
			p, err := UnmarshalCreateStreamErrorBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal CreateStreamError", m, err)
				return
			}
			h.conn.OnCreateStreamError(ctx, p)
			return
		case "onStatus":
			p, err := UnmarshalOnStatusBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal OnStatus", m, err)
				return
			}
			h.conn.OnOnStatus(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "play":
			p, err := UnmarshalPlayBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Play", m, err)
				return
			}
			h.conn.OnPlay(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "play2":
			p, err := UnmarshalPlay2Binary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Play2", m, err)
				return
			}
			h.conn.OnPlay2(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "deleteStream":
			p, err := UnmarshalDeleteStreamBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal DeleteStream", m, err)
				return
			}
			h.conn.OnDeleteStream(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "closeStream":
			p, err := UnmarshalCloseStreamBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal CloseStream", m, err)
				return
			}
			h.conn.OnCloseStream(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "receiveAudio":
			p, err := UnmarshalReceiveAudioBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal ReceiveAudio", m, err)
				return
			}
			h.conn.OnReceiveAudio(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "receiveVideo":
			p, err := UnmarshalReceiveVideoBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal ReceiveVideo", m, err)
				return
			}
			h.conn.OnReceiveVideo(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "publish":
			p, err := UnmarshalPublishBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Publish", m, err)
				return
			}
			h.conn.OnPublish(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "seek":
			p, err := UnmarshalSeekBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Seek", m, err)
				return
			}
			h.conn.OnSeek(ctx, m.ChunkStreamID(), m.StreamID(), p)
		case "pause":
			p, err := UnmarshalPauseBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Pause", m, err)
				return
			}
			h.conn.OnPause(ctx, m.ChunkStreamID(), m.StreamID(), p)
		default:
			h.Error("not implemented command", m, err)
			p, err := UnmarshalCallBinary(b, encodingAMFType)
			if err != nil {
				h.Error("failed to unmarshal Call", m, err)
				return
			}
			h.conn.OnCall(ctx, p)
		}
	case MessageTypeIDAggregate:
		// TODO
	}
}

func (h ConnControlMessageHandler) Error(msg string, m Message, err error) {
	if err != nil {
		h.conn.Logger().Error(
			msg,
			zap.Error(err),
			zap.Object("message", m),
		)
	}
}
