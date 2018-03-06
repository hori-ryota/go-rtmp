package rtmp

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunkBasicHeader(t *testing.T) {
	testCases := []struct {
		name string
		b    []byte
		h    ChunkBasicHeader
	}{
		{
			name: "1 byte format: fmt=0, min",
			b:    []byte{3},
			h: GenerateChunkBasicHeader(
				0,
				3,
			),
		},
		{
			name: "1 byte format: fmt=0, max",
			b:    []byte{0x3f},
			h: GenerateChunkBasicHeader(
				0,
				63,
			),
		},
		{
			name: "1 byte format: fmt=3, min",
			b:    []byte{(3 << 6) + 3},
			h: GenerateChunkBasicHeader(
				3,
				3,
			),
		},
		{
			name: "2 byte format: fmt=0, min",
			b:    []byte{0, 0},
			h: GenerateChunkBasicHeader(
				0,
				64,
			),
		},
		{
			name: "2 byte format: fmt=0, max",
			b:    []byte{0, 0xff},
			h: GenerateChunkBasicHeader(
				0,
				319,
			),
		},
		{
			name: "2 byte format: fmt=3, min",
			b:    []byte{(3 << 6), 0},
			h: GenerateChunkBasicHeader(
				3,
				64,
			),
		},
		{
			name: "3 byte format: fmt=0, min",
			b:    []byte{1, 1, 0},
			h: GenerateChunkBasicHeader(
				0,
				320,
			),
		},
		{
			name: "3 byte format: fmt=0, max",
			b:    []byte{1, 0xff, 0xff},
			h: GenerateChunkBasicHeader(
				0,
				65599,
			),
		},
		{
			name: "3 byte format: fmt=3, min",
			b:    []byte{(3 << 6) + 1, 1, 0},
			h: GenerateChunkBasicHeader(
				3,
				320,
			),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			h, err := ReadChunkBasicHeader(
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

	t.Run("3 byte format: 2 byte range data", func(t *testing.T) {
		h, err := ReadChunkBasicHeader(
			bufio.NewReader(bytes.NewReader([]byte{1, 0, 0})),
		)
		if assert.NoError(t, err) {
			assert.Equal(t, NewChunkBasicHeader3B(
				0,
				64,
			), h)
		}
	})
}
