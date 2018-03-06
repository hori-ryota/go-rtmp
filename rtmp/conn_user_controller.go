package rtmp

import (
	"context"

	"github.com/pkg/errors"
)

func (conn *defaultConn) StreamBegin(ctx context.Context, streamID uint32) error {
	p := NewStreamBegin(streamID)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDUserControlMessages,
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

func (conn *defaultConn) StreamEOF(ctx context.Context, streamID uint32) error {
	p := NewStreamEOF(streamID)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDUserControlMessages,
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

func (conn *defaultConn) StreamDry(ctx context.Context, streamID uint32) error {
	p := NewStreamDry(streamID)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDUserControlMessages,
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

func (conn *defaultConn) SetBufferLength(ctx context.Context, streamID uint32, bufferLength uint32) error {
	p := NewSetBufferLength(streamID, bufferLength)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDUserControlMessages,
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

func (conn *defaultConn) StreamIsRecorded(ctx context.Context, streamID uint32) error {
	p := NewStreamIsRecorded(streamID)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDUserControlMessages,
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

func (conn *defaultConn) PingRequest(ctx context.Context, timestamp uint32) error {
	p := NewPingRequest(timestamp)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDUserControlMessages,
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

func (conn *defaultConn) PingResponse(ctx context.Context, timestamp uint32) error {
	p := NewPingResponse(timestamp)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}

	m := NewMessage(
		2,
		MessageTypeIDUserControlMessages,
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
