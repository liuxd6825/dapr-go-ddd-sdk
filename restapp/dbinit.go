package restapp

import (
	"context"
)

func InitDbWithConfig(envName string, prefix string, config *Config, tables *Tables) error {
	name := envName
	if len(name) == 0 {
		name = config.Env
	}

	env, err := config.GetEnvConfig(name)
	if err != nil {
		return err
	}

	InitMongo(env.Mongo)
	return InitDb(tables, env, prefix)
}

func InitDb(tables *Tables, env *EnvConfig, prefix string) error {
	mongo := NewMongoManager()
	gorm := NewGormManager()

	ctx := context.Background()
	opt := NewCreateOptions().SetPrefix(prefix)
	for _, table := range tables.items {
		if table.IsMongo {
			mongo.Create(ctx, table, env, opt)
		} else if table.IsGORM {
			gorm.Create(ctx, table, env, opt)
		} else if table.IsElasticSearch {

		}
	}
	return nil
}
