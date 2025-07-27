package xtest

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"

const (
	MongoDBKey      = "db"
	MongoDBName     = "test"
	MongoHostLocal  = "127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019"
	MongoHostRemote = "192.168.120.224:27018,192.168.120.224:27019,192.168.120.224:27020"
	MongoReplicaSet = "mongors"
	MongoUser       = "super_admin"
	MongoPassword   = "123456"
)

type MongoOptions struct {
	DBName     *string
	ReplicaSet *string
	User       *string
	Pwd        *string
}

func NewMongoOptions(opts ...*MongoOptions) *MongoOptions {
	o := &MongoOptions{}
	for _, item := range opts {
		if item.DBName != nil {
			o.DBName = item.DBName
		}
		if item.ReplicaSet != nil {
			o.ReplicaSet = item.ReplicaSet
		}
		if item.User != nil {
			o.User = item.User
		}
		if item.Pwd != nil {
			o.Pwd = item.Pwd
		}
	}
	return o
}

func (o *MongoOptions) GetDBName() string {
	if o.DBName != nil {
		return *o.DBName
	}
	return MongoDBName
}

func (o *MongoOptions) GetReplicaSet() string {
	if o.ReplicaSet != nil {
		return *o.ReplicaSet
	}
	return ""
}

func (o *MongoOptions) GetUser() string {
	if o.ReplicaSet != nil {
		return *o.User
	}
	return MongoUser
}

func (o *MongoOptions) GetPwd() string {
	if o.ReplicaSet != nil {
		return *o.User
	}
	return MongoPassword
}

func (o *MongoOptions) SetDBName(s string) *MongoOptions {
	o.DBName = &s
	return o
}

func GetMongoEnv_Remote(opts ...*MongoOptions) *env.Mongo {
	opt := NewMongoOptions(opts...)
	return &env.Mongo{
		DbKey:      MongoDBKey,
		Host:       MongoHostRemote,
		ReplicaSet: getStr(opt.GetReplicaSet(), MongoReplicaSet),
		DbName:     getStr(opt.GetDBName(), MongoDBName),
		User:       getStr(opt.GetUser(), MongoUser),
		Pwd:        getStr(opt.GetPwd(), MongoPassword),
	}
}

func GetMongoEnv_Local(opts ...*MongoOptions) *env.Mongo {
	opt := NewMongoOptions(opts...)
	return &env.Mongo{
		DbKey:      MongoDBKey,
		Host:       MongoHostLocal,
		ReplicaSet: getStr(opt.GetReplicaSet(), "rs0"),
		DbName:     getStr(opt.GetDBName(), MongoDBName),
		User:       getStr(opt.GetUser(), MongoUser),
		Pwd:        getStr(opt.GetPwd(), MongoPassword),
	}
}

func getStr(val string, def string) string {
	if val == "" {
		return def
	}
	return val
}
