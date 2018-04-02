package rtmp

import (
	"bufio"
	"context"
	"io"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Reader interface {
	io.Reader
	ReadMessage() (Message, error)
	SetChunkSize(chunkSize uint32)
	AbortMessage(chunkStreamID uint32)
	SetBandwidthLimitType(bandwidthLimitType BandwidthLimitType)
	SetAcknowledgementWindowSize(acknowledgementWindowSize uint32)
	CreateStream() (streamID uint32)
}

type reader Reader

type defaultReader struct {
	conn                         Conn
	r                            *bufio.Reader
	chunkSize                    uint32
	chunkStreams                 map[ /* chunkStreamID */ uint32]chunkStream
	bandwidthLimitType           BandwidthLimitType
	acknowledgementWindowSize    uint32
	preAcknowledgementWindowSize uint32
	sequenceNumber               uint32
}

func NewDefaultReader(conn Conn, r io.Reader) Reader {
	return &defaultReader{
		conn:         conn,
		r:            bufio.NewReader(r),
		chunkSize:    128, /* default RTMP Chunk size */
		chunkStreams: map[uint32]chunkStream{},
	}
}

func (r *defaultReader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (y *defaultReader) ReadMessage() (Message, error) {
	r := y.r
	for {
		h, err := ReadChunkHeader(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to ReadChunkHeader")
		}
		csID := h.BasicHeader().ChunkStreamID()
		cs := y.chunkStreams[csID]

		switch mh := h.MessageHeader().(type) {
		case ChunkMessageHeaderType0:
			cs.messageLength = mh.MessageLength()
			cs.messageTypeID = mh.MessageTypeID()
			cs.messageStreamID = mh.MessageStreamID()
			cs.timestampDelta = 0
			cs.timestamp = mh.Timestamp()
			if uint32(len(cs.buffer)) < cs.messageLength {
				cs.buffer = make([]byte, cs.messageLength)
			}
			if y.acknowledgementWindowSize > 0 {
				y.sequenceNumber += 11
			}
		case ChunkMessageHeaderType1:
			cs.messageLength = mh.MessageLength()
			cs.messageTypeID = mh.MessageTypeID()
			cs.timestampDelta = mh.TimestampDelta()
			cs.timestamp += cs.timestampDelta
			if uint32(len(cs.buffer)) < cs.messageLength {
				cs.buffer = make([]byte, cs.messageLength)
			}
			if y.acknowledgementWindowSize > 0 {
				y.sequenceNumber += 7
			}
		case ChunkMessageHeaderType2:
			cs.timestampDelta = mh.TimestampDelta()
			cs.timestamp += cs.timestampDelta
			if y.acknowledgementWindowSize > 0 {
				y.sequenceNumber += 3
			}
		case ChunkMessageHeaderType3:
			if cs.buffered == 0 {
				cs.timestamp += cs.timestampDelta
			}
		}

		if y.acknowledgementWindowSize > 0 {
			switch h.BasicHeader().(type) {
			case ChunkBasicHeader1B:
				y.sequenceNumber += 1
			case ChunkBasicHeader2B:
				y.sequenceNumber += 2
			default:
				y.sequenceNumber += 3
			}
		}

		if cs.messageLength == 0 {
			y.sendAcknowledgementIfNeeded()
			m := NewMessage(
				csID,
				cs.messageTypeID,
				cs.timestamp,
				cs.messageStreamID,
				cs.buffer[:cs.messageLength],
			)
			cs.buffered = 0
			y.chunkStreams[csID] = cs
			return m, nil
		}

		length := cs.messageLength - cs.buffered
		if length > y.chunkSize {
			length = y.chunkSize
		}
		b := cs.buffer[cs.buffered : cs.buffered+length]
		if _, err := io.ReadFull(r, b); err != nil {
			return nil, errors.Wrapf(err, "failed to ReadBody: chunkStream=%+v: reader=%+v", cs, r)
		}
		cs.buffered += length
		y.sequenceNumber += length
		y.sendAcknowledgementIfNeeded()
		if cs.buffered == cs.messageLength {
			m := NewMessage(
				csID,
				cs.messageTypeID,
				cs.timestamp,
				cs.messageStreamID,
				cs.buffer[:cs.messageLength],
			)
			cs.buffered = 0
			y.chunkStreams[csID] = cs
			return m, nil
		}

		y.chunkStreams[csID] = cs
	}
}

func (r *defaultReader) SetChunkSize(chunkSize uint32) {
	r.chunkSize = chunkSize
}

func (r *defaultReader) AbortMessage(chunkStreamID uint32) {
	delete(r.chunkStreams, chunkStreamID)
}

func (r *defaultReader) SetBandwidthLimitType(bandwidthLimitType BandwidthLimitType) {
	r.bandwidthLimitType = bandwidthLimitType
}

func (r *defaultReader) SetAcknowledgementWindowSize(acknowledgementWindowSize uint32) {
	r.acknowledgementWindowSize = acknowledgementWindowSize
}

func (r *defaultReader) sendAcknowledgementIfNeeded() {
	if r.acknowledgementWindowSize == 0 {
		return
	}
	if r.sequenceNumber >= r.preAcknowledgementWindowSize+r.acknowledgementWindowSize {
		if err := r.conn.Acknowledgement(context.TODO(), r.sequenceNumber); err != nil {
			r.conn.Logger().Error(
				"failed to send Acknowledgement",
				zap.Error(err),
			)
			return
		}
		r.preAcknowledgementWindowSize = r.acknowledgementWindowSize
	}
}

func (r *defaultReader) CreateStream() (streamID uint32) {
	for i := uint32(3); true; i++ {
		if _, ok := r.chunkStreams[i]; !ok {
			r.chunkStreams[i] = chunkStream{ID: i}
			return i
		}
	}
	// dummy
	return 0
}
