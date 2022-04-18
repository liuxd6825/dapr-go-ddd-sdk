package restapp

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_mongodb"
)

var mongoDB *ddd_mongodb.MongoDB

func initMongo(mongoConfig *MongoConfig) {
	config := &ddd_mongodb.Config{
		Host:         mongoConfig.Host,
		DatabaseName: mongoConfig.Database,
		UserName:     mongoConfig.UserName,
		Password:     mongoConfig.Password,
		MaxPoolSize:  mongoConfig.MaxPoolSize,
	}

	mongoDB = ddd_mongodb.NewMongoDB()
	if err := mongoDB.Init(config); err != nil {
		panic(err)
	}
}

func GetMongoDB() *ddd_mongodb.MongoDB {
	return mongoDB
}

func NewSession() ddd_repository.Session {
	return ddd_mongodb.NewSession(mongoDB)
}
