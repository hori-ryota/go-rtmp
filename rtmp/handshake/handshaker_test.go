package handshake

import (
	"bufio"
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
