package rtmp

import (
	"bufio"
	"io"

	"github.com/pkg/errors"
)

type Writer interface {
	io.Writer
	Flush() error
	WriteMessage(message Message) (n int, err error)
	SetChunkSize(chunkSize uint32)
}

type writer Writer

type defaultWriter struct {
	conn         Conn
	w            *bufio.Writer
	chunkSize    uint32
	chunkStreams map[ /* chunkStreamID */ uint32]chunkStream
}

func NewDefaultWriter(conn Conn, w io.Writer) Writer {
	return &defaultWriter{
		conn:         conn,
		w:            bufio.NewWriter(w),
		chunkSize:    128, /* default RTMP Chunk size */
		chunkStreams: map[uint32]chunkStream{},
	}
}

func (w *defaultWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *defaultWriter) Flush() error {
	return w.w.Flush()
}

func (w *defaultWriter) WriteMessage(m Message) (n int, err error) {

	csID := m.ChunkStreamID()
	cs := w.chunkStreams[csID]
	isNotFirst := cs.isNotFirst
	cs.isNotFirst = true

	p := m.Payload()

	var mh ChunkMessageHeader
	var format uint8
	var extendedTimestamp uint32
	messageLength := uint32(len(p))
	timestampDelta := m.Timestamp() - cs.timestamp
	switch {
	case !isNotFirst || cs.messageStreamID != m.StreamID(): // type 0
		cs.messageLength = messageLength
		cs.messageTypeID = m.TypeID()
		cs.messageStreamID = m.StreamID()
		cs.timestampDelta = 0
		cs.timestamp = m.Timestamp()
		tt := cs.timestamp
		if tt >= 0xffffff {
			tt = 0xffffff
			extendedTimestamp = cs.timestamp
		}
		format = 0
		mh = NewChunkMessageHeaderType0(
			tt,
			uint32(len(p)),
			m.TypeID(),
			m.StreamID(),
		)
	case cs.messageTypeID != m.TypeID(), cs.messageLength != messageLength: // type 1
		cs.messageTypeID = m.TypeID()
		cs.messageLength = messageLength
		cs.timestampDelta = timestampDelta
		cs.timestamp += timestampDelta
		tt := cs.timestampDelta
		if tt >= 0xffffff {
			tt = 0xffffff
			extendedTimestamp = cs.timestampDelta
		}
		format = 1
		mh = NewChunkMessageHeaderType1(
			tt,
			uint32(len(p)),
			m.TypeID(),
		)
	case cs.timestampDelta != timestampDelta: // type 2
		cs.timestampDelta = timestampDelta
		cs.timestamp += timestampDelta
		tt := timestampDelta
		if tt >= 0xffffff {
			tt = 0xffffff
			extendedTimestamp = timestampDelta
		}
		format = 2
		mh = NewChunkMessageHeaderType2(
			tt,
		)
	default: //type 3
		cs.timestamp += timestampDelta
		format = 3
		mh = NewChunkMessageHeaderType3()
	}

	bh := GenerateChunkBasicHeader(format, csID)
	h := NewChunkHeader(bh, mh, extendedTimestamp)

	if uint32(len(p)) <= w.chunkSize {
		b, err := NewChunk(h, p).MarshalBinary()
		if err != nil {
			return 0, errors.Wrap(err, "failed to marshal chunk")
		}
		n, err = w.Write(b)
		if err != nil {
			return 0, errors.Wrap(err, "failed to writer chunk")
		}
		w.chunkStreams[csID] = cs
		return n, nil
	}

	b, err := NewChunk(h, p[:w.chunkSize]).MarshalBinary()
	if err != nil {
		return 0, errors.Wrap(err, "failed to marshal chunk")
	}
	nn, err := w.Write(b)
	if err != nil {
		return 0, errors.Wrap(err, "failed to write chunk")
	}
	n += nn

	h = NewChunkHeader(
		GenerateChunkBasicHeader(3, csID),
		NewChunkMessageHeaderType3(),
		0,
	)
	remain := p[w.chunkSize:]

	for len(remain) > 0 {
		l := int(w.chunkSize)
		if l > len(remain) {
			l = len(remain)
		}
		b, err := NewChunk(h, remain[:l]).MarshalBinary()
		if err != nil {
			return 0, errors.Wrap(err, "failed to marshal chunk")
		}
		nn, err := w.Write(b)
		if err != nil {
			return 0, errors.Wrap(err, "failed to write chunk")
		}
		n += nn

		remain = remain[l:]
	}
	w.chunkStreams[csID] = cs

	return n, nil
}

func (w *defaultWriter) SetChunkSize(chunkSize uint32) {
	w.chunkSize = chunkSize
}
