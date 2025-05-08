package ddd_repository

import (
	"context"
	"time"
)

type Session interface {
	UseTransaction(context.Context, SessionFunc, ...*SessionOptions) error
}

type SessionFunc func(ctx context.Context) error
type SessionType int

type SessionOptions struct {
	writeTimeout *time.Duration //写事务超时进行
	readTimeout  *time.Duration //读事务超时进行
}

const (
	NoSession SessionType = iota
	ReadSession
	WriteSession
)

type SessionOptionsBuilder struct {
	writeTimeout *time.Duration
	readTimeout  *time.Duration
}

var (
	DefaultWriteTimeout = 5 * time.Minute
	DefaultReadTimeout  = 2 * time.Minute
)

func NewSessionOptionsBuilder() *SessionOptionsBuilder {
	return &SessionOptionsBuilder{}
}

func NewSessionOptions(opts ...*SessionOptions) *SessionOptions {
	opt := &SessionOptions{
		writeTimeout: &DefaultWriteTimeout,
		readTimeout:  &DefaultReadTimeout,
	}
	for _, o := range opts {
		if o == nil {
			continue
		}
		if o.writeTimeout != nil {
			opt.writeTimeout = o.writeTimeout
		}
		if o.readTimeout != nil {
			opt.readTimeout = o.readTimeout
		}
	}
	return opt
}

func (o *SessionOptions) GetWriteTime() time.Duration {
	if o == nil || o.writeTimeout == nil {
		return DefaultWriteTimeout
	}
	return *o.writeTimeout
}

func (o *SessionOptions) GetReadTime() time.Duration {
	if o == nil || o.readTimeout == nil {
		return DefaultReadTimeout
	}
	return *o.readTimeout
}

func (b *SessionOptionsBuilder) SetWriteTime(v time.Duration) *SessionOptionsBuilder {
	b.writeTimeout = &v
	return b
}

func (b *SessionOptionsBuilder) SetReadTime(v time.Duration) *SessionOptionsBuilder {
	b.readTimeout = &v
	return b
}

func (b *SessionOptionsBuilder) Build() *SessionOptions {
	options := &SessionOptions{
		writeTimeout: b.writeTimeout,
		readTimeout:  b.readTimeout,
	}
	return options
}

func StartSession(ctx context.Context, session Session, dbFunc SessionFunc, opts ...*SessionOptions) error {
	return session.UseTransaction(ctx, dbFunc, opts...)
}
