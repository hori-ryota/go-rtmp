package rtmp

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/hori-ryota/go-rtmp/rtmp/handshake"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Conn interface {
	Serve() error
	io.Closer

	Reader() Reader
	SetReader(r Reader)

	Writer() Writer
	SetWriter(w Writer)

	MessagePubsub

	ProtocolController
	ProtocolControlEventHandler
	UserController
	UserControlEventHandler
	NetConnectionCommander
	NetConnectionCommandHandler
	NetStreamCommander
	NetStreamCommandHandler

	Logger() *zap.Logger
}

type defaultConn struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	conn       net.Conn

	handshaker handshake.Handshaker

	reader
	writer
	encodingAMFType           EncodingAMFType
	bandwidthLimitType        BandwidthLimitType
	windowAcknowledgementSize uint32

	timestampPoint time.Time

	messagePubsub

	logger *zap.Logger
}

func NewDefaultConn(
	ctx context.Context,
	nc net.Conn,
	isServer bool,
	logger *zap.Logger,
	connInitializers ...func(Conn),
) Conn {
	ctx, cancel := context.WithCancel(ctx)
	conn := &defaultConn{
		ctx:                       ctx,
		cancelFunc:                cancel,
		conn:                      nc,
		handshaker:                handshake.NewDefaultHandshaker(isServer),
		reader:                    NewDefaultReader(nc),
		writer:                    NewDefaultWriter(nc),
		encodingAMFType:           defaultEncodingAMFType,
		bandwidthLimitType:        defaultBandwidthLimitType,
		windowAcknowledgementSize: defaultWindowAcknowledgementSize,
		messagePubsub:             NewDefaultMessagePubsub(),
		logger:                    logger,
	}
	for _, f := range connInitializers {
		f(conn)
	}
	return conn
}

func (conn *defaultConn) Serve() error {
	ctx := conn.ctx

	r := conn.reader
	w := conn.writer

	if err := conn.handshaker.Handshake(
		ctx, r, w,
	); err != nil {
		return errors.Wrap(err, "failed to handshake")
	}

	conn.timestampPoint = time.Now()

	for !isDone(ctx) {
		m, err := r.ReadMessage()
		if err != nil {
			if errors.Cause(err) == io.EOF {
				return nil
			}
			conn.Logger().Error(
				"failed to read message",
				zap.Error(err),
			)
			continue
		}
		ctx = SetChunkStreamIDToContext(ctx, m.ChunkStreamID())
		ctx = SetMessageStreamIDToContext(ctx, m.StreamID())
		conn.HandleMessage(ctx, m)
	}
	return ctx.Err()
}

func (conn *defaultConn) Close() error {
	defer conn.cancelFunc()
	return conn.conn.Close()
}

func (conn *defaultConn) Reader() Reader {
	return conn.reader
}

func (conn *defaultConn) SetReader(r Reader) {
	conn.reader = r
}

func (conn *defaultConn) Writer() Writer {
	return conn.writer
}

func (conn *defaultConn) SetWriter(w Writer) {
	conn.writer = w
}

func (conn *defaultConn) Logger() *zap.Logger {
	return conn.logger
}

func (conn *defaultConn) Timestamp() uint32 {
	return uint32(time.Since(conn.timestampPoint) / time.Millisecond)
}
