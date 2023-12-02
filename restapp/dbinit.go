package restapp

import (
	"context"
	"os"
	"strings"
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

func InitDbScript(dbKey string, tables *Tables, env *EnvConfig, prefix string, saveFile string) error {
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
		if sb, err := mongo.GetScript(ctx, dbKey, mTables, env, opt); err != nil {
			return err
		} else if err = writeFile(saveFile, sb); err != nil {
			return err
		}
	}

	if len(gTables) > 0 {
		gorm.GetScript(ctx, dbKey, mTables, env, opt)
	}
	return nil
}

func writeFile(fileName string, sb *strings.Builder) error {
	//在当前路径下，创建test.txt文件，若当前路径存在test.txt，则清空其内容
	fileHandler, err := os.Create(fileName)
	if nil != err {
		panic(err)
	}
	defer fileHandler.Close()

	//使用WriteString
	_, err = fileHandler.WriteString(sb.String())
	return err
}
