package rtmp

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Server struct {
	ctx         context.Context
	cancelFunc  context.CancelFunc
	connOptions []ConnOption

	Addr string

	logger *zap.Logger
}

var defaultLogger *zap.Logger = func() *zap.Logger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return log
}()

func NewServer(
	ctx context.Context,
	connOps ...ConnOption,
) *Server {
	ctx, cancel := context.WithCancel(ctx)
	return &Server{
		ctx:         ctx,
		cancelFunc:  cancel,
		connOptions: connOps,
		logger:      defaultLogger,
	}
}

func (s *Server) Close() error {
	s.cancelFunc()
	return nil
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":1935"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

func (s *Server) Serve(l net.Listener) error {
	ctx := s.ctx
	defer func() {
		if err := l.Close(); err != nil {
			s.Logger().Error(
				"failed to close listener",
				zap.Error(err),
				zap.Stringer("addr", l.Addr()),
			)
		}
	}()

	var tempDelay time.Duration // how long to sleep on accept failure
	for !isDone(ctx) {
		nc, err := l.Accept()
		if err != nil {
			if isCanceledErr(err) {
				return nil
			}
			if ne, ok := errors.Cause(err).(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				s.Logger().Error(
					"Accept error",
					zap.Error(err),
					zap.Stringer("retrying in", tempDelay),
					zap.Stringer("addr", l.Addr()),
				)

				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		c := NewDefaultConn(
			ctx,
			nc,
			true,
			s.Logger(),
			s.connOptions...,
		)

		go func() {
			remoteAddr := nc.RemoteAddr()
			defer func() {
				if err := recover(); err != nil {
					s.Logger().Error(
						"panic: failed to conn.serve",
						zap.Any("error", err),
						zap.Stringer("remoteAddr", remoteAddr),
					)
				}
				defer func() {
					if err := recover(); err != nil {
						s.Logger().Error(
							"panic: failed to conn.close",
							zap.Any("error", err),
							zap.Stringer("remoteAddr", remoteAddr),
						)
					}
				}()
				if err := c.Close(); err != nil {
					s.Logger().Error(
						"failed to close conn",
						zap.Error(err),
						zap.Stringer("remoteAddr", remoteAddr),
					)
				}
			}()
			if err := c.Serve(); err != nil {
				if isCanceledErr(err) {
					return
				}
				if errors.Cause(err) == io.EOF {
					return
				}
				if e, ok := errors.Cause(err).(ConnError); ok {
					s.Logger().Error(
						"failed to conn.serve",
						append(e.Fields(), zap.Error(err), zap.Stringer("remoteAddr", remoteAddr))...,
					)
				} else {
					s.Logger().Error(
						"failed to conn.serve",
						zap.Error(err),
						zap.Stringer("remoteAddr", remoteAddr),
					)
				}
			}
		}()
	}
	return ctx.Err()
}

func (s *Server) Logger() *zap.Logger {
	if s.logger != nil {
		return s.logger
	}
	return defaultLogger
}

func ListenAndServe(ctx context.Context, addr string, connOps ...ConnOption) error {
	s := NewServer(ctx, connOps...)
	s.Addr = addr
	return s.ListenAndServe()
}
