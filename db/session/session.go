package session

import (
	"context"
	"time"
)

type Func func(ctx context.Context) error
type SessionType int
type doSession func(i int) error

type Options struct {
	mongo        *SessionType
	neo4j        *SessionType
	writeTimeout *time.Duration //写事务超时进行
	readTimeout  *time.Duration //读事务超时进行
}

const (
	NoSession SessionType = iota
	ReadSession
	WriteSession
)

type OptionsBuilder struct {
	mongo        *SessionType
	neo4j        *SessionType
	writeTimeout *time.Duration
	readTimeout  *time.Duration
}

var (
	DefaultWriteTimeout = 5 * time.Minute
	DefaultReadTimeout  = 2 * time.Minute
)

func NewOptions(opts ...*Options) *Options {
	w := WriteSession
	opt := &Options{
		mongo:        &w,
		neo4j:        &w,
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
		if o.mongo != nil {
			opt.mongo = o.mongo
		}
		if o.neo4j != nil {
			opt.neo4j = o.neo4j
		}
	}
	return opt
}

func (o *Options) SetMongo(v SessionType) *Options {
	o.mongo = &v
	return o
}

func (o *Options) GetMongo() SessionType {
	if o.mongo == nil {
		return NoSession
	}
	return *o.mongo
}

func (o *Options) SetNeo4j(v SessionType) *Options {
	o.neo4j = &v
	return o
}

func (o *Options) GetNeo4j() SessionType {
	if o.neo4j == nil {
		return NoSession
	}
	return *o.neo4j
}

/*
func StartSession(ctx context.Context, fun Func, opts ...*Options) error {
	var start doSession
	var sessions []ddd_repository.Session
	var err error

	opt := NewOptions(opts...)
	if opt.GetMongo() != NoSession {
		sessions = append(sessions, mongo_dao.NewSession(opt.GetMongo() == WriteSession))
	}
	if opt.GetNeo4j() != NoSession {
		sessions = append(sessions, neo4j_dao.NewSession(opt.GetNeo4j() == WriteSession))
	}

	start = func(i int) error {
		session := sessions[i]
		return session.UseTransaction(ctx, func(ctx context.Context) error {
			if i < 1 {
				return fun(ctx)
			}
			return start(i - 1)
		})
	}

	length := len(sessions)
	if length > 0 {
		err = start(length - 1)
	} else {
		err = fun(ctx)
	}

	if err != nil {
		logs.Errorln(ctx, "db.StartSession()", err.Error())
	}
	return err
}
*/
