package rtmp

//go:generate go run $DEFDIR/go/cmd/genbinary/genbinary.go -package rtmp -toml $DEFDIR/chunk_header/message_header.toml -customInterfaceFunc "ChunkMessageHeader#NeedsExtendedTimestamp() bool" -o chunk_message_header_gen.go

import (
	"io"

	"github.com/pkg/errors"
)

func ReadChunkMessageHeader(r io.Reader, fmt uint8) (ChunkMessageHeader, error) {
	switch fmt {
	case 0:
		return ReadChunkMessageHeaderType0(r)
	case 1:
		return ReadChunkMessageHeaderType1(r)
	case 2:
		return ReadChunkMessageHeaderType2(r)
	case 3:
		return ReadChunkMessageHeaderType3(r)
	default:
		return nil, errors.Errorf("unknown fmt %d", fmt)
	}
}

func (m chunkMessageHeaderType0) NeedsExtendedTimestamp() bool {
	return m.Timestamp() == 0xffffff
}
func (m chunkMessageHeaderType1) NeedsExtendedTimestamp() bool {
	return m.TimestampDelta() == 0xffffff
}
func (m chunkMessageHeaderType2) NeedsExtendedTimestamp() bool {
	return m.TimestampDelta() == 0xffffff
}
func (m chunkMessageHeaderType3) NeedsExtendedTimestamp() bool {
	return false
}
