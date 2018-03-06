package rtmp

//go:generate go run $DEFDIR/go/cmd/genbinary/genbinary.go -package rtmp -toml $DEFDIR/chunk_header/basic_header.toml -o chunk_basic_header_gen.go

import (
	"io"

	"github.com/pkg/errors"
)

func GenerateChunkBasicHeader(fmt uint8, chunkStreamID uint32) ChunkBasicHeader {
	switch {
	case chunkStreamID < 64:
		return NewChunkBasicHeader1B(fmt, chunkStreamID)
	case chunkStreamID < 320:
		return NewChunkBasicHeader2B(fmt, chunkStreamID)
	default:
		return NewChunkBasicHeader3B(fmt, chunkStreamID)
	}
}

func ReadChunkBasicHeader(r io.Reader) (ChunkBasicHeader, error) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, errors.Wrap(err, "failed to ReadByte")
	}
	b1 := b[0]

	fmt := b1 >> 6
	streamID := uint32(b1 & 0x3f)

	switch streamID {
	case 0:
		// 2 byte form
		if _, err := io.ReadFull(r, b); err != nil {
			return nil, errors.Wrap(err, "failed to ReadByte")
		}
		b2 := b[0]
		c := NewChunkBasicHeader2B(fmt, 0)
		if err := c.UnmarshalBinary([]byte{b1, b2}); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal 2 byte form chunk basic header")
		}
		return c, nil
	case 1:
		// 3 byte form
		c := NewChunkBasicHeader3B(fmt, 0)
		b := make([]byte, 2)
		if _, err := io.ReadFull(r, b); err != nil {
			return nil, errors.Wrap(err, "failed to read bytes")
		}
		if err := c.UnmarshalBinary([]byte{b1, b[0], b[1]}); err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal 2 byte form chunk basic header")
		}
		return c, nil
	default:
		return NewChunkBasicHeader1B(fmt, streamID), nil
	}
}
