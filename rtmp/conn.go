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
	Context() context.Context

	Reader() Reader
	SetReader(r Reader)

	Writer() Writer
	SetWriter(w Writer)

	MessagePubsub

	ProtocolController
	DefaultProtocolControlEventHandler() ProtocolControlEventHandler
	UserController
	DefaultUserControlEventHandler() UserControlEventHandler
	NetConnectionCommander
	DefaultNetConnectionCommandHandler() NetConnectionCommandHandler
	NetStreamCommander
	DefaultNetStreamCommandHandler() NetStreamCommandHandler

	SetCreateStreamCallbacks(transactionID uint32, f func(CreateStreamResponse) ConnError)
	AddNetstreamCommandCallbacks(func(OnStatus) ConnError)
	TransactionID() uint32

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

	createStreamCallbacks     map[uint32] /* transactionID */ func(CreateStreamResponse) ConnError
	netStreamCommandCallbacks []func(onStatus OnStatus) ConnError

	onConnectValidators []func(
		ctx context.Context,
		connect Connect,
	) ConnectError

	onPublishValidators []func(
		ctx context.Context,
		publish Publish,
	) (errorInfo map[string]interface{})

	logger *zap.Logger
}

func NewDefaultConn(
	ctx context.Context,
	nc net.Conn,
	isServer bool,
	logger *zap.Logger,
	connOps ...ConnOption,
) Conn {
	logger = logger.With(zap.Bool("isServer", isServer))
	logger = logger.With(zap.Stringer("remoteAddr", nc.RemoteAddr()))
	ctx, cancel := context.WithCancel(ctx)
	conn := &defaultConn{
		ctx:                       ctx,
		cancelFunc:                cancel,
		conn:                      nc,
		handshaker:                handshake.NewDefaultHandshaker(isServer),
		encodingAMFType:           defaultEncodingAMFType,
		bandwidthLimitType:        defaultBandwidthLimitType,
		windowAcknowledgementSize: defaultWindowAcknowledgementSize,
		messagePubsub:             NewDefaultMessagePubsub(),

		createStreamCallbacks: map[uint32]func(CreateStreamResponse) ConnError{},

		logger: logger,
	}
	conn.reader = NewDefaultReader(conn, nc, conn.logger)
	conn.writer = NewDefaultWriter(conn, nc)
	ops := &connOptions{}
	for _, o := range connOps {
		o(ops)
	}
	ops.Apply(conn)
	return conn
}

func (conn *defaultConn) Serve() error {
	ctx := conn.ctx

	r := conn.reader
	w := conn.writer

	if err := conn.handshaker.Handshake(
		ctx, r, w,
	); err != nil {
		if errors.Cause(err) == io.EOF || isDone(ctx) {
			return nil
		}
		return errors.Wrap(err, "failed to handshake")
	}

	conn.timestampPoint = time.Now()

	for !isDone(ctx) {
		m, err := r.ReadMessage()
		if err != nil {
			if errors.Cause(err) == io.EOF || isDone(ctx) {
				return nil
			}
			conn.logger.Error(
				"failed to read message",
				zap.Error(err),
			)
			continue
		}
		if err := conn.HandleMessage(ctx, m); err != nil {
			switch {
			case IsConnWarnError(err):
				conn.logger.Warn(
					"caught error",
					append(err.Fields(), zap.Error(err))...,
				)
			case IsConnRejectedError(err):
				conn.logger.Info(
					"caught rejected error",
					append(err.Fields(), zap.Error(err))...,
				)
				return nil
			default:
				return err
			}
		}
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

func (conn *defaultConn) Timestamp() uint32 {
	return uint32(time.Since(conn.timestampPoint) / time.Millisecond)
}

func (conn *defaultConn) Context() context.Context {
	return conn.ctx
}

func (conn *defaultConn) SetCreateStreamCallbacks(transactionID uint32, f func(CreateStreamResponse) ConnError) {
	conn.createStreamCallbacks[transactionID] = f
}

func (conn *defaultConn) AddNetstreamCommandCallbacks(f func(OnStatus) ConnError) {
	conn.netStreamCommandCallbacks = append(conn.netStreamCommandCallbacks, f)
}

func (conn *defaultConn) TransactionID() uint32 {
	for i := uint32(2); true; i++ {
		if _, ok := conn.createStreamCallbacks[i]; !ok {
			conn.createStreamCallbacks[i] = func(_ CreateStreamResponse) ConnError {
				// noop
				return nil
			}
			return i
		}
	}
	// dummy
	return 0
}

func (conn *defaultConn) Logger() *zap.Logger {
	return conn.logger
}
