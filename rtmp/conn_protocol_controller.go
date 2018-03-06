package rtmp

import (
	"context"

	"github.com/pkg/errors"
)

func (conn *defaultConn) SetChunkSize(ctx context.Context, chunkSize uint32) error {
	p := NewSetChunkSize(chunkSize)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDSetChunkSize,
		conn.Timestamp(),
		0,
		b,
	)
	_, err = conn.Writer().WriteMessage(m)
	if err != nil {
		return errors.Wrap(err, "failed to WriteMessage")
	}
	if err := conn.Writer().Flush(); err != nil {
		return errors.Wrap(err, "failed to Flush Writer")
	}
	conn.Writer().SetChunkSize(chunkSize)
	return nil
}

func (conn *defaultConn) AbortMessage(ctx context.Context, chunkStreamID uint32) error {
	p := NewAbortMessage(chunkStreamID)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDAbortMessage,
		conn.Timestamp(),
		0,
		b,
	)
	_, err = conn.Writer().WriteMessage(m)
	if err != nil {
		return errors.Wrap(err, "failed to WriteMessage")
	}
	if err := conn.Writer().Flush(); err != nil {
		return errors.Wrap(err, "failed to Flush Writer")
	}
	return nil
}

func (conn *defaultConn) Acknowledgement(ctx context.Context, sequenceNumber uint32) error {
	p := NewAcknowledgement(sequenceNumber)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDAcknowledgement,
		conn.Timestamp(),
		0,
		b,
	)
	_, err = conn.Writer().WriteMessage(m)
	if err != nil {
		return errors.Wrap(err, "failed to WriteMessage")
	}
	if err := conn.Writer().Flush(); err != nil {
		return errors.Wrap(err, "failed to Flush Writer")
	}
	return nil
}

func (conn *defaultConn) WindowAcknowledgementSize(ctx context.Context, acknowledgementWindowSize uint32) error {
	p := NewWindowAcknowledgementSize(acknowledgementWindowSize)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDWindowAcknowledgementSize,
		conn.Timestamp(),
		0,
		b,
	)
	_, err = conn.Writer().WriteMessage(m)
	if err != nil {
		return errors.Wrap(err, "failed to WriteMessage")
	}
	if err := conn.Writer().Flush(); err != nil {
		return errors.Wrap(err, "failed to Flush Writer")
	}
	return nil
}

func (conn *defaultConn) SetPeerBandwidth(ctx context.Context, acknowledgementWindowSize uint32, limitType BandwidthLimitType) error {
	p := NewSetPeerBandwidth(acknowledgementWindowSize, limitType)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDSetPeerBandwidth,
		conn.Timestamp(),
		0,
		b,
	)
	_, err = conn.Writer().WriteMessage(m)
	if err != nil {
		return errors.Wrap(err, "failed to WriteMessage")
	}
	if err := conn.Writer().Flush(); err != nil {
		return errors.Wrap(err, "failed to Flush Writer")
	}
	return nil
}
