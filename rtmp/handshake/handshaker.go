package handshake

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	ServerRTMPVersion = 0x03
)

type Handshaker interface {
	Handshake(ctx context.Context, r io.Reader, w io.Writer) error
}

type defaultHandshaker struct {
	isServer bool
}

var DefaultServerHandshaker = NewDefaultHandshaker(true)
var DefaultClientHandshaker = NewDefaultHandshaker(false)

func NewDefaultHandshaker(
	isServer bool,
) defaultHandshaker {
	return defaultHandshaker{
		isServer: isServer,
	}
}

func (y defaultHandshaker) Handshake(ctx context.Context, r io.Reader, w io.Writer) error {
	g, ctx := errgroup.WithContext(ctx)

	var received1, sent1 Chunk1

	g.Go(func() (err error) {
		_, received1, err = y.Receive0And1(ctx, r)
		return err
	})

	if y.isServer {
		// The server MUST wait until C0 has been received before sending sent0 and
		// sent1, and MAY wait until after C1 as well.
		err := g.Wait()
		if err != nil {
			return err
		}
	}

	time := uint32(0)
	g.Go(func() (err error) {
		_, sent1, err = y.Send0And1(ctx, w, ServerRTMPVersion, time)
		return err
	})

	err := g.Wait()
	if err != nil {
		return err
	}

	g.Go(func() error {
		received2, err := y.Receive2(ctx, r)
		if err != nil {
			return errors.Wrap(err, "failed to receive Chunk2")
		}
		if err := received2.Validate(sent1); err != nil {
			return errors.Wrap(err, "invalid Chunk2")
		}
		return nil
	})

	g.Go(func() error {
		time2 := maxUint32(time, received1.Time())
		return y.Send2(ctx, w, received1, time2)
	})

	return g.Wait()
}

func (y defaultHandshaker) Send0And1(
	ctx context.Context,
	w io.Writer,
	rtmpVersion uint8,
	time uint32,
) (c0 Chunk0, c1 Chunk1, err error) {
	// send 0
	c0 = NewChunk0(rtmpVersion)
	b, err := c0.MarshalBinary()
	if err != nil {
		return c0, c1, errors.Wrap(err, "failed to MarshalBinary Chunk0")
	}
	if _, err = w.Write(b); err != nil {
		return c0, c1, errors.Wrap(err, "failed to send Chunk0")
	}

	// send 1
	c1, err = GenerateChunk1(time)
	if err != nil {
		return c0, c1, errors.Wrap(err, "failed to init Chunk1")
	}
	b, err = c1.MarshalBinary()
	if err != nil {
		return c0, c1, errors.Wrap(err, "failed to MarshalBinary Chunk1")
	}
	if _, err = w.Write(b); err != nil {
		return c0, c1, errors.Wrap(err, "failed to send Chunk1")
	}
	if w, ok := w.(flusher); ok {
		if err = w.Flush(); err != nil {
			return c0, c1, errors.Wrap(err, "failed to flush writer for Chunk0 and Chunk1")
		}
	}
	return c0, c1, nil
}

func (y defaultHandshaker) Send2(
	ctx context.Context,
	w io.Writer,
	received1 Chunk1,
	time2 uint32,
) (err error) {
	// send 2
	c2 := GenerateChunk2(received1, time2)
	b, err := c2.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary Chunk2")
	}
	if _, err := w.Write(b); err != nil {
		return errors.Wrap(err, "failed to send Chunk2")
	}
	if w, ok := w.(flusher); ok {
		if err := w.Flush(); err != nil {
			return errors.Wrap(err, "failed to flush writer for Chunk2")
		}
	}
	return nil
}

func (y defaultHandshaker) Receive0And1(
	ctx context.Context,
	r io.Reader,
) (c0 Chunk0, c1 Chunk1, err error) {
	// receive 0
	c0, err = ReadChunk0(r)
	if err != nil {
		return c0, c1, errors.Wrap(err, "failed to read Chunk0")
	}

	// receive 1
	c1, err = ReadChunk1(r)
	return c0, c1, errors.Wrap(err, "failed to read Chunk1")
}

func (y defaultHandshaker) Receive2(
	ctx context.Context,
	r io.Reader,
) (c2 Chunk2, err error) {
	// receive 2
	received2, err := ReadChunk2(r)
	return received2, errors.Wrap(err, "failed to read Chunk2")
}
