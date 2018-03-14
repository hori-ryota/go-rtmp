package rtmp

type chunkStream struct {
	ID uint32

	isNotFirst bool

	messageLength   uint32
	messageTypeID   MessageTypeID
	messageStreamID uint32
	timestampDelta  uint32

	timestamp uint32

	buffer   []byte
	buffered uint32
}
