package auth

import (
	"github.com/casbin/casbin/v2"
	_ "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

var _enforcer = NewEnforcerEmpty()

/*  Init modelConfFile : "rbac_model.conf" */
func Init(dbKey string, modelConfFile string) {

	dbItem := envs.GetEnv().GetDB(dbKey)
	if dbItem == nil {
		panic(errors.New("db item not found"))
	}

	dbType := dbItem.GetDBType()
	switch dbType {
	case envs.DBType_MongoDB:
		if db := dbItem.GetMongo(); db != nil {
			adapter, err := mongodbadapter.NewAdapterByDB(db.Client(), &mongodbadapter.AdapterConfig{
				DatabaseName:   db.GetDatabase().Name(),
				CollectionName: "casbin_rule",
				IsFiltered:     false,
			})
			if err != nil {
				panic("Failed to create Casbin adapter")
			}
			enforcer, err := casbin.NewEnforcer(modelConfFile, adapter)
			if err != nil {
				panic("Failed to create Casbin enforcer")
			}
			_enforcer = enforcer
		} else {
			panic(errors.New("db type not supported"))
		}
		break
	case envs.DBType_Postgres, envs.DBType_MySQL,
		envs.DBType_Oracle, envs.DBType_MsSQL,
		envs.DBType_Sqlite:
		if db := dbItem.GetGormDB(); db != nil {
			adapter, err := gormadapter.NewAdapterByDB(db)
			if err != nil {
				panic("Failed to create Casbin adapter")
			}
			enforcer, err := casbin.NewEnforcer(modelConfFile, adapter)
			if err != nil {
				panic("Failed to create Casbin enforcer")
			}
			_enforcer = enforcer
		} else {
			panic(errors.New("db type not supported"))
		}
		break
	default:
		panic("auth Database type not supported")
	}

}

func Enforcer() IEnforcer {
	return _enforcer
}
