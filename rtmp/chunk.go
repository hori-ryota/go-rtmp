package rtmp

import (
	"github.com/pkg/errors"
)

//go:generate go run $DEFDIR/go/cmd/genbinary/genbinary.go -package rtmp -toml $DEFDIR/chunk/chunk.toml -o chunk_gen.go

func (c chunk) MarshalBinary() ([]byte, error) {
	h, err := c.Header().MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to header.MarshalBinary")
	}

	return append(h, c.data...), nil
}
