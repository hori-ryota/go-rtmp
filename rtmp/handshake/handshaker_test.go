package handshake

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func Test_defaultHandshaker_Handshake(t *testing.T) {
	client := NewDefaultHandshaker(false)
	server := NewDefaultHandshaker(true)

	csr, csw := io.Pipe()
	scr, scw := io.Pipe()

	ctx := context.Background()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return errors.Wrap(
			client.Handshake(ctx, bufio.NewReader(scr), bufio.NewWriter(csw)),
			"failed to client handshake",
		)
	})

	g.Go(func() error {
		return errors.Wrap(
			server.Handshake(ctx, bufio.NewReader(csr), bufio.NewWriter(scw)),
			"failed to server handshake",
		)
	})

	err := g.Wait()
	assert.NoError(t, err)
}

func Test_defaultHandshaker_Send0And1(t *testing.T) {
	t.Run("flusher: bufio.Writer", func(t *testing.T) {
		b := new(bytes.Buffer)
		w := bufio.NewWriter(b)
		handshaker := NewDefaultHandshaker(false)
		c0, c1, err := handshaker.Send0And1(
			context.Background(),
			w,
			0x03,
			0xaaaaaaaa,
		)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x03), c0.Version())
		assert.Equal(t, uint32(0xaaaaaaaa), c1.Time())
		assert.NotZero(t, c1.RandomBytes())

		written := b.Bytes()

		assert.Equal(t, byte(0x03), written[0])
		assert.Equal(t, []byte{0xaa, 0xaa, 0xaa, 0xaa}, written[1:5])
		assert.Equal(t, []byte{0, 0, 0, 0}, written[5:9])
		assert.Equal(t, c1.RandomBytes(), written[9:1537])

		assert.Zero(t, w.Buffered())
	})

	t.Run("not flusher", func(t *testing.T) {
		w := new(bytes.Buffer)
		handshaker := NewDefaultHandshaker(false)
		c0, c1, err := handshaker.Send0And1(
			context.Background(),
			w,
			0x03,
			0xaaaaaaaa,
		)
		assert.NoError(t, err)
		assert.Equal(t, uint8(0x03), c0.Version())
		assert.Equal(t, uint32(0xaaaaaaaa), c1.Time())
		assert.NotZero(t, c1.RandomBytes())

		written := w.Bytes()

		assert.Equal(t, byte(0x03), written[0])
		assert.Equal(t, []byte{0xaa, 0xaa, 0xaa, 0xaa}, written[1:5])
		assert.Equal(t, []byte{0, 0, 0, 0}, written[5:9])
		assert.Equal(t, c1.RandomBytes(), written[9:1537])
	})
}

func Test_defaultHandshaker_Send2(t *testing.T) {
	c1, err := GenerateChunk1(0x55555555)
	assert.NoError(t, err)

	t.Run("flusher: bufio.Writer", func(t *testing.T) {
		b := new(bytes.Buffer)
		w := bufio.NewWriter(b)
		handshaker := NewDefaultHandshaker(false)
		err := handshaker.Send2(
			context.Background(),
			w,
			c1,
			0xaaaaaaaa,
		)
		assert.NoError(t, err)

		written := b.Bytes()

		assert.Equal(t, []byte{0x55, 0x55, 0x55, 0x55}, written[0:4])
		assert.Equal(t, []byte{0xaa, 0xaa, 0xaa, 0xaa}, written[4:8])
		assert.Equal(t, c1.RandomBytes(), written[8:1536])

		assert.Zero(t, w.Buffered())
	})

	t.Run("not flusher", func(t *testing.T) {
		w := new(bytes.Buffer)
		handshaker := NewDefaultHandshaker(false)
		err := handshaker.Send2(
			context.Background(),
			w,
			c1,
			0xaaaaaaaa,
		)
		assert.NoError(t, err)

		written := w.Bytes()

		assert.Equal(t, []byte{0x55, 0x55, 0x55, 0x55}, written[0:4])
		assert.Equal(t, []byte{0xaa, 0xaa, 0xaa, 0xaa}, written[4:8])
		assert.Equal(t, c1.RandomBytes(), written[8:1536])
	})
}
