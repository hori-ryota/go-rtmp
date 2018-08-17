package handshake

//go:generate go run $DEFDIR/go/cmd/genbinary/genbinary.go -package handshake -toml $DEFDIR/handshake/chunk.toml -toml $DEFDIR/handshake/chunk0.toml -toml $DEFDIR/handshake/chunk1.toml -customInterfaceFunc "Chunk2#Validate(chunk1 Chunk1) error" -toml $DEFDIR/handshake/chunk2.toml -o chunk_gen.go

import (
	"bytes"
	"crypto/rand"

	"github.com/pkg/errors"
)

func GenerateChunk1(time uint32) (Chunk1, error) {
	// Random data (1528 bytes): This field can contain any arbitrary
	// values. Since each endpoint has to distinguish between the
	// response to the handshake it has initiated and the handshake
	// initiated by its peer,this data SHOULD send something sufficiently
	// random. But there is no need for cryptographically-secure
	// randomness, or even dynamic values.
	randomBytes := make([]byte, 1528)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create randomBytes")
	}
	return NewChunk1(
		time,
		randomBytes,
	), nil
}

func GenerateChunk2(chunk1 Chunk1, time2 uint32) Chunk2 {
	return NewChunk2(
		chunk1.Time(),
		time2,
		chunk1.RandomBytes(),
	)
}

func (c chunk2) Validate(chunk1 Chunk1) error {

	// validate time
	if chunk1.Time() != c.time {
		return errors.Errorf("c time is not equal chunk1 time: chunk1 time %d: c time %d", chunk1.Time(), c.time)
	}

	// validate time2
	if c.time2 < chunk1.Time() {
		return errors.Errorf("c time2 is smaller than time: time1 %d: time2 %d", chunk1.Time(), c.time)
	}

	// validate randomEcho
	if !bytes.Equal(c.randomEcho, chunk1.RandomBytes()) {
		return errors.Errorf(
			"randomEcho is not %x: got %x",
			chunk1.RandomBytes(),
			c.randomEcho,
		)
	}
	return nil
}
