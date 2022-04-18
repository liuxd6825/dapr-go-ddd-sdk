package restapp

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
)

var _mongodb *ddd_mongodb.MongoDB

func initMongo(appMongoConfig *AppMongoConfig) {
	if err := assert.NotNil(appMongoConfig, assert.WidthOptionsError("appMongoConfig is nil")); err != nil {
		panic(err)
	}

	config := &ddd_mongodb.Config{
		Host:         appMongoConfig.Host,
		DatabaseName: appMongoConfig.Database,
		UserName:     appMongoConfig.UserName,
		Password:     appMongoConfig.Password,
		MaxPoolSize:  appMongoConfig.MaxPoolSize,
		ReplicaSet:   appMongoConfig.ReplicaSet,
		WriteConcern: appMongoConfig.WriteConcern,
		ReadConcern:  appMongoConfig.ReadConcern,
	}

	mongodb, err := ddd_mongodb.NewMongoDB(config)
	if err != nil {
		panic(err)
	}
	_mongodb = mongodb
}

func GetMongoDB() *ddd_mongodb.MongoDB {
	return _mongodb
}

func NewSession() ddd_repository.Session {
	return ddd_mongodb.NewSession(_mongodb)
}
