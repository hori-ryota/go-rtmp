package rtmp

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkMessageHeader(t *testing.T) {
	testCases := []struct {
		name string
		fmt  uint8
		b    []byte
		h    ChunkMessageHeader
	}{
		{
			name: "type 0",
			fmt:  0,
			b: []byte{
				0x1, 0x2, 0x3,
				0x4, 0x5, 0x6,
				0x7,
				0x8, 0x9, 0xa, 0xb, // little endian
			},
			h: NewChunkMessageHeaderType0(
				0x010203,
				0x040506,
				0x07,
				0x0b0a0908,
			),
		},
		{
			name: "type 1",
			fmt:  1,
			b: []byte{
				0x1, 0x2, 0x3,
				0x4, 0x5, 0x6,
				0x7,
			},
			h: NewChunkMessageHeaderType1(
				0x010203,
				0x040506,
				0x07,
			),
		},
		{
			name: "type 2",
			fmt:  2,
			b: []byte{
				0x1, 0x2, 0x3,
			},
			h: NewChunkMessageHeaderType2(
				0x010203,
			),
		},
		{
			name: "type 3",
			fmt:  3,
			b:    []byte{},
			h:    NewChunkMessageHeaderType3(),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			h, err := ReadChunkMessageHeader(
				bufio.NewReader(bytes.NewReader(tt.b)),
				tt.fmt,
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
