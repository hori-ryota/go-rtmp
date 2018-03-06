package rtmp

import (
	"context"

	"github.com/pkg/errors"
)

func (conn *defaultConn) Connect(ctx context.Context, commandObject map[string]interface{}, optionalUserArguments map[string]interface{}) error {
	panic("not implemented")
}

func (conn *defaultConn) ConnectResult(ctx context.Context, properties map[string]interface{}, information map[string]interface{}) error {
	p := NewConnectResult(properties, information, conn.encodingAMFType)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}
	var msgTypeID MessageTypeID
	if conn.encodingAMFType == EncodingAMFTypeAMF0 {
		msgTypeID = MessageTypeIDCommandAMF0
	} else {
		msgTypeID = MessageTypeIDCommandAMF3
	}

	m := NewMessage(
		3,
		msgTypeID,
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

func (conn *defaultConn) ConnectError(ctx context.Context, properties map[string]interface{}, information map[string]interface{}) error {
	p := NewConnectError(properties, information, conn.encodingAMFType)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}
	var msgTypeID MessageTypeID
	if conn.encodingAMFType == EncodingAMFTypeAMF0 {
		msgTypeID = MessageTypeIDCommandAMF0
	} else {
		msgTypeID = MessageTypeIDCommandAMF3
	}

	m := NewMessage(
		3,
		msgTypeID,
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

func (conn *defaultConn) Call(ctx context.Context, procedureName string, transactionID uint32, commandObject map[string]interface{}, optionalArguments map[string]interface{}) error {
	panic("not implemented")
}

func (conn *defaultConn) CallResponse(ctx context.Context, commandName string, transactionID uint32, commandObject map[string]interface{}, response map[string]interface{}) error {
	panic("not implemented")
}

func (conn *defaultConn) CreateStream(ctx context.Context, transactionID uint32, commandObject map[string]interface{}) error {
	p := NewCreateStream(transactionID, commandObject, conn.encodingAMFType)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}
	var msgTypeID MessageTypeID
	if conn.encodingAMFType == EncodingAMFTypeAMF0 {
		msgTypeID = MessageTypeIDCommandAMF0
	} else {
		msgTypeID = MessageTypeIDCommandAMF3
	}

	m := NewMessage(
		3,
		msgTypeID,
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

func (conn *defaultConn) CreateStreamResult(ctx context.Context, transactionID uint32, commandObject map[string]interface{}, streamID uint32) error {
	p := NewCreateStreamResult(transactionID, commandObject, streamID, conn.encodingAMFType)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}
	var msgTypeID MessageTypeID
	if conn.encodingAMFType == EncodingAMFTypeAMF0 {
		msgTypeID = MessageTypeIDCommandAMF0
	} else {
		msgTypeID = MessageTypeIDCommandAMF3
	}

	m := NewMessage(
		3,
		msgTypeID,
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

func (conn *defaultConn) CreateStreamError(ctx context.Context, transactionID uint32, commandObject map[string]interface{}, streamID uint32) error {
	p := NewCreateStreamError(transactionID, commandObject, streamID, conn.encodingAMFType)
	b, err := p.MarshalBinary()
	if err != nil {
		return errors.Wrap(err, "failed to MarshalBinary")
	}
	var msgTypeID MessageTypeID
	if conn.encodingAMFType == EncodingAMFTypeAMF0 {
		msgTypeID = MessageTypeIDCommandAMF0
	} else {
		msgTypeID = MessageTypeIDCommandAMF3
	}

	m := NewMessage(
		3,
		msgTypeID,
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
