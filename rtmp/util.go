package rtmp

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

func BigEndianToUint32(b []byte) uint32 {
	result := uint32(0)
	for i := range b {
		result |= uint32(b[i]) << (8 * uint(len(b)-1-i))
	}
	return result
}

func Uint32ToBigEndian(b []byte, v uint32) {
	for i := range b {
		b[i] = byte(v >> (uint(len(b)-1-i) * 8))
	}
}

func LittleEndianToUint32(b []byte) uint32 {
	result := uint32(0)
	for i := range b {
		result |= uint32(b[i]) << (8 * uint(i))
	}
	return result
}

func Uint32ToLittleEndian(b []byte, v uint32) {
	for i := range b {
		b[i] = byte(v >> (uint(i) * 8))
	}
}

func isDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func isCanceledErr(err error) bool {
	return strings.Contains(errors.Cause(err).Error(), "context canceled")
}
