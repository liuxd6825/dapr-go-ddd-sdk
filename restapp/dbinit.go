package restapp

import (
	"context"
)

func InitDbWithConfig(envName string, prefix string, dbKey string, config *Config, tables *Tables) error {
	if dbKey == "" {
		panic("dbkey is not emtiy")
	}
	name := envName
	if len(name) == 0 {
		name = config.Env
	}

	env, err := config.GetEnvConfig(name)
	if err != nil {
		return err
	}

	InitMongo(env.App.AppId, env.Mongo)
	return InitDb(dbKey, tables, env, prefix)
}

func InitDb(dbKey string, tables *Tables, env *EnvConfig, prefix string) error {
	if dbKey == "" {
		panic("dbkey is not emtiy")
	}
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

func InitDbScript(dbKey string, tables *Tables, env *EnvConfig, prefix string) error {
	if dbKey == "" {
		panic("dbkey is not emtiy")
	}
	var mTables []*Table
	var gTables []*Table

	mongo := NewMongoScriptManager()
	gorm := NewGormScriptManager()

	ctx := context.Background()
	opt := NewCreateOptions().SetPrefix(prefix)

	for _, table := range tables.items {
		if table.IsMongo {
			mTables = append(mTables, table)
		} else if table.IsGORM {
			gTables = append(gTables, table)
		}
	}

	if len(mTables) > 0 {
		mongo.GetInitScript(ctx, dbKey, mTables, env, opt)
	}

	if len(gTables) > 0 {
		gorm.GetInitScript(ctx, dbKey, mTables, env, opt)
	}
	return nil
}
