package rtmp

import (
	"bufio"
	"bytes"
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkHeader(t *testing.T) {
	bh := GenerateChunkBasicHeader(2, 0x3)
	mh := NewChunkMessageHeaderType2(0x2)
	mhET := NewChunkMessageHeaderType2(0xffffff)

	toBinary := func(m encoding.BinaryMarshaler) []byte {
		b, err := m.MarshalBinary()
		if err != nil {
			panic(err)
		}
		return b
	}

	appends := func(bs ...[]byte) []byte {
		b := []byte{}
		for i := range bs {
			b = append(b, bs[i]...)
		}
		return b
	}

	testCases := []struct {
		name string
		b    []byte
		h    ChunkHeader
	}{
		{
			name: "without extendTimestamp",
			b: appends(
				toBinary(bh),
				toBinary(mh),
			),
			h: NewChunkHeader(
				bh,
				mh,
				0,
			),
		},
		{
			name: "with extendTimestamp",
			b: appends(
				toBinary(bh),
				toBinary(mhET),
				[]byte{0x1, 0x2, 0x3, 0x4},
			),
			h: NewChunkHeader(
				bh,
				mhET,
				0x01020304,
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			h, err := ReadChunkHeader(
				bufio.NewReader(bytes.NewReader(tt.b)),
			)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.h, h)
			}

			b, err := tt.h.MarshalBinary()
			if assert.NoError(t, err) {
				assert.Equal(t, tt.b, b)
			}
		})
	}
}
