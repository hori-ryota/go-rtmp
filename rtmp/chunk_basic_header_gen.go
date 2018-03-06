// Automatically generated by go generate; DO NOT EDIT.

package rtmp

import (
	"encoding"
	"io"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// The Chunk Basic Header encodes the chunk stream ID and the chunk type
// (represented by fmt field in the figure below). Chunk type
// determines the format of the encoded message header. Chunk Basic
// Header field may be 1, 2, or 3 bytes, depending on the chunk stream
// ID.
// An implementation SHOULD use the smallest representation that can
// hold the ID.
// The protocol supports up to 65597 streams with IDs 3-65599. The IDs
// 0, 1, and 2 are reserved. Value 0 indicates the 2 byte form and an
// ID in the range of 64-319 (the second byte + 64). Value 1 indicates
// the 3 byte form and an ID in the range of 64-65599 ((the third
// byte)*256 + the second byte + 64). Values in the range of 3-63
// represent the complete stream ID. Chunk Stream ID with value 2 is
// reserved for low-level protocol control messages and commands.
// The bits 0-5 (least significant) in the chunk basic header represent
// the chunk stream ID.
// Chunk stream IDs 2-63 can be encoded in the 1-byte version of this
// field.
//
//  0 1 2 3 4 5 6 7
// +-+-+-+-+-+-+-+-+
// |fmt|   cs id   |
// +-+-+-+-+-+-+-+-+
//
// Chunk stream IDs 64-319 can be encoded in the 2-byte form of the
// header. ID is computed as (the second byte + 64).
//
//  0                   1
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |fmt|     0     |   cs id - 64  |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// Chunk stream IDs 64-65599 can be encoded in the 3-byte version of
// this field. ID is computed as ((the third byte)*256 + (the second
// byte) + 64).
//
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |fmt|     1     |        cs id - 64             |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
//
// cs id (6 bits): This field contains the chunk stream ID, for values
// from 2-63. Values 0 and 1 are used to indicate the 2- or 3-byte
// versions of this field.
// fmt (2 bits): This field identifies one of four format used by the
// ’chunk message header’. The ’chunk message header’ for each of
// the chunk types is explained in the next section.
// cs id - 64 (8 or 16 bits): This field contains the chunk stream ID
// minus 64. For example, ID 365 would be represented by a 1 in cs
// id, and a 16-bit 301 here.
// Chunk stream IDs with values 64-319 could
type ChunkBasicHeader interface {
	encoding.BinaryMarshaler
	zapcore.ObjectMarshaler
	Fmt() uint8
	ChunkStreamID() uint32
}

type ChunkBasicHeader1B interface {
	ChunkBasicHeader
	encoding.BinaryUnmarshaler
}

type chunkBasicHeader1B struct {
	fmt           uint8
	chunkStreamID uint32
}

func NewChunkBasicHeader1B(
	fmt uint8,
	chunkStreamID uint32,
) ChunkBasicHeader1B {
	return &chunkBasicHeader1B{
		fmt:           fmt,
		chunkStreamID: chunkStreamID,
	}
}

func (m chunkBasicHeader1B) Fmt() uint8 {
	return m.fmt
}
func (m chunkBasicHeader1B) ChunkStreamID() uint32 {
	return m.chunkStreamID
}

func (m chunkBasicHeader1B) MarshalBinary() ([]byte, error) {
	b := make([]byte, 1)
	b[0] |= m.fmt << 6
	b[0] |= byte(m.chunkStreamID)
	return b, nil
}

func (m *chunkBasicHeader1B) UnmarshalBinary(b []byte) error {
	if len(b) != 1 {
		return errors.Errorf("invalid binary size %d: %x", len(b), b)
	}
	m.fmt = b[0] >> 6
	m.chunkStreamID = uint32(b[0]) & 0x3f
	return nil
}

func UnmarshalChunkBasicHeader1BBinary(b []byte) (ChunkBasicHeader1B, error) {
	m := chunkBasicHeader1B{}
	err := m.UnmarshalBinary(b)
	return &m, err
}

func ReadChunkBasicHeader1B(r io.Reader) (ChunkBasicHeader1B, error) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, errors.Wrap(err, "failed to read bytes")
	}
	return UnmarshalChunkBasicHeader1BBinary(b)
}

func (m chunkBasicHeader1B) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddUint8("fmt", m.Fmt())
	enc.AddUint32("chunkStreamID", m.ChunkStreamID())
	return nil
}

type ChunkBasicHeader2B interface {
	ChunkBasicHeader
	encoding.BinaryUnmarshaler
}

type chunkBasicHeader2B struct {
	fmt           uint8
	chunkStreamID uint32
}

func NewChunkBasicHeader2B(
	fmt uint8,
	chunkStreamID uint32,
) ChunkBasicHeader2B {
	return &chunkBasicHeader2B{
		fmt:           fmt,
		chunkStreamID: chunkStreamID,
	}
}

func (m chunkBasicHeader2B) Fmt() uint8 {
	return m.fmt
}
func (m chunkBasicHeader2B) ChunkStreamID() uint32 {
	return m.chunkStreamID
}

func (m chunkBasicHeader2B) MarshalBinary() ([]byte, error) {
	b := make([]byte, 2)
	b[0] |= m.fmt << 6
	b[1] |= byte(m.chunkStreamID - 64)
	return b, nil
}

func (m *chunkBasicHeader2B) UnmarshalBinary(b []byte) error {
	if len(b) != 2 {
		return errors.Errorf("invalid binary size %d: %x", len(b), b)
	}
	m.fmt = b[0] >> 6
	m.chunkStreamID = uint32(b[1]) + 64
	return nil
}

func UnmarshalChunkBasicHeader2BBinary(b []byte) (ChunkBasicHeader2B, error) {
	m := chunkBasicHeader2B{}
	err := m.UnmarshalBinary(b)
	return &m, err
}

func ReadChunkBasicHeader2B(r io.Reader) (ChunkBasicHeader2B, error) {
	b := make([]byte, 2)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, errors.Wrap(err, "failed to read bytes")
	}
	return UnmarshalChunkBasicHeader2BBinary(b)
}

func (m chunkBasicHeader2B) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddUint8("fmt", m.Fmt())
	enc.AddUint32("chunkStreamID", m.ChunkStreamID())
	return nil
}

type ChunkBasicHeader3B interface {
	ChunkBasicHeader
	encoding.BinaryUnmarshaler
}

type chunkBasicHeader3B struct {
	fmt           uint8
	chunkStreamID uint32
}

func NewChunkBasicHeader3B(
	fmt uint8,
	chunkStreamID uint32,
) ChunkBasicHeader3B {
	return &chunkBasicHeader3B{
		fmt:           fmt,
		chunkStreamID: chunkStreamID,
	}
}

func (m chunkBasicHeader3B) Fmt() uint8 {
	return m.fmt
}
func (m chunkBasicHeader3B) ChunkStreamID() uint32 {
	return m.chunkStreamID
}

func (m chunkBasicHeader3B) MarshalBinary() ([]byte, error) {
	b := make([]byte, 3)
	b[0] |= m.fmt << 6
	b[0] |= 1
	for i := range b[1:] {
		b[1+i] |= byte((m.chunkStreamID - 64) >> (uint(1-i) * 8))
	}
	return b, nil
}

func (m *chunkBasicHeader3B) UnmarshalBinary(b []byte) error {
	if len(b) != 3 {
		return errors.Errorf("invalid binary size %d: %x", len(b), b)
	}
	m.fmt = b[0] >> 6
	for i, bb := range b[1:] {
		m.chunkStreamID |= uint32(bb) << (8 * uint(1-i))
	}
	m.chunkStreamID = m.chunkStreamID + 64

	return nil
}

func UnmarshalChunkBasicHeader3BBinary(b []byte) (ChunkBasicHeader3B, error) {
	m := chunkBasicHeader3B{}
	err := m.UnmarshalBinary(b)
	return &m, err
}

func ReadChunkBasicHeader3B(r io.Reader) (ChunkBasicHeader3B, error) {
	b := make([]byte, 3)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, errors.Wrap(err, "failed to read bytes")
	}
	return UnmarshalChunkBasicHeader3BBinary(b)
}

func (m chunkBasicHeader3B) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddUint8("fmt", m.Fmt())
	enc.AddUint32("chunkStreamID", m.ChunkStreamID())
	return nil
}
