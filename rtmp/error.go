package rtmp

import (
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type ErrorLevel string

const (
	ErrorLevelFatal    ErrorLevel = "fatal"
	ErrorLevelWarn     ErrorLevel = "warn"
	ErrorLevelRejected ErrorLevel = "rejected"
)

type ConnError interface {
	error
	ErrorLevel() ErrorLevel
	Fields() []zapcore.Field
}

type connFatalError struct {
	error
	fields []zapcore.Field
}

func NewConnFatalError(err error, fields ...zapcore.Field) ConnError {
	return &connFatalError{
		error:  err,
		fields: fields,
	}
}

func (e connFatalError) ErrorLevel() ErrorLevel {
	return ErrorLevelFatal
}

func (e connFatalError) Fields() []zapcore.Field {
	return e.fields
}

func IsConnFatalError(err error) bool {
	e, ok := errors.Cause(err).(ConnError)
	return ok && e.ErrorLevel() == ErrorLevelFatal
}

type connWarnError struct {
	error
	fields []zapcore.Field
}

func NewConnWarnError(err error, fields ...zapcore.Field) ConnError {
	return &connWarnError{
		error:  err,
		fields: fields,
	}
}

func (e connWarnError) ErrorLevel() ErrorLevel {
	return ErrorLevelWarn
}

func (e connWarnError) Fields() []zapcore.Field {
	return e.fields
}

func IsConnWarnError(err error) bool {
	e, ok := errors.Cause(err).(ConnError)
	return ok && e.ErrorLevel() == ErrorLevelWarn
}

type connRejectedError struct {
	error
	fields []zapcore.Field
}

func NewConnRejectedError(err error, fields ...zapcore.Field) ConnError {
	return &connRejectedError{
		error:  err,
		fields: fields,
	}
}

func (e connRejectedError) ErrorLevel() ErrorLevel {
	return ErrorLevelRejected
}

func (e connRejectedError) Fields() []zapcore.Field {
	return e.fields
}

func IsConnRejectedError(err error) bool {
	e, ok := errors.Cause(err).(ConnError)
	return ok && e.ErrorLevel() == ErrorLevelRejected
}
