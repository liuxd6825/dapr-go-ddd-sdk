package xtest

import "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"

const (
	MongoDBKey      = "db"
	MongoDBName     = "test"
	MongoHostLocal  = "127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019"
	MongoHostRemote = "192.168.120.224:27018,192.168.120.224:27019,192.168.120.224:27020"
)

func GetMongoEnv_Remote() *env.Mongo {
	return &env.Mongo{
		DbKey:      MongoDBKey,
		Host:       MongoHostRemote,
		ReplicaSet: "mongors",
		DbName:     MongoDBName,
		User:       "super_admin",
		Pwd:        "123456",
	}
}

func GetMongoEnv_Local() *env.Mongo {
	return &env.Mongo{
		DbKey:      MongoDBKey,
		Host:       MongoHostLocal,
		ReplicaSet: "rs0",
		DbName:     MongoDBName,
		User:       "super_admin",
		Pwd:        "123456",
	}
}
