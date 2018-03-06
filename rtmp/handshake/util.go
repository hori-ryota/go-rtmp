package handshake

func maxUint32(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

type flusher interface {
	Flush() error
}
