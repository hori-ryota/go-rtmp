package rtmp

//go:generate go run $DEFDIR/go/cmd/genbinary/genbinary.go -package rtmp -toml $DEFDIR/chunk_header/chunk_header.toml -customInterfaceFunc "ChunkHeader#Timestamp() uint32" -customInterfaceFunc "ChunkHeader#TimestampDelta() uint32" -o chunk_header_gen.go

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

func ReadChunkHeader(r io.Reader) (ChunkHeader, error) {
	bh, err := ReadChunkBasicHeader(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadChunkBasicHeader")
	}
	mh, err := ReadChunkMessageHeader(r, bh.Fmt())
	if err != nil {
		return nil, errors.Wrap(err, "failed to ReadChunkMessageHeader")
	}
	extendedTimestamp := uint32(0)
	if mh.NeedsExtendedTimestamp() {
		extendedTimestamp, err = ReadExtendedTimestamp(r)
		if err != nil {
			return nil, errors.Wrap(err, "failed to ReadExtendedTimestamp")
		}
	}
	return NewChunkHeader(bh, mh, extendedTimestamp), nil
}

func (h chunkHeader) MarshalBinary() ([]byte, error) {
	bh, err := h.BasicHeader().MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to basicHeader.MarshalBinary")
	}
	mh, err := h.MessageHeader().MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to messageHeader.MarshalBinary")
	}

	if h.ExtendedTimestamp() == 0 {
		return append(bh, mh...), nil
	}

	et := make([]byte, 4)
	binary.BigEndian.PutUint32(et, h.ExtendedTimestamp())

	b := make([]byte, 0, len(bh)+len(mh)+len(et))
	b = append(b, bh...)
	b = append(b, mh...)
	return append(b, et...), nil
}

func ReadExtendedTimestamp(r io.Reader) (uint32, error) {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return 0, errors.Wrap(err, "failed to ReadFull")
	}
	return BigEndianToUint32(b), nil
}

type timestamper interface {
	Timestamp() uint32
}

type timestampDeltaer interface {
	TimestampDelta() uint32
}

func (h chunkHeader) Timestamp() uint32 {
	if h.MessageHeader().NeedsExtendedTimestamp() {
		return h.extendedTimestamp
	}
	if mh, ok := h.MessageHeader().(timestamper); ok {
		return mh.Timestamp()
	}
	return 0
}

func (h chunkHeader) TimestampDelta() uint32 {
	if h.MessageHeader().NeedsExtendedTimestamp() {
		return h.extendedTimestamp
	}
	if mh, ok := h.MessageHeader().(timestampDeltaer); ok {
		return mh.TimestampDelta()
	}
	return 0
}
