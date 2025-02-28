package sql

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository/ddd_sql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/db"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/db_pkg/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"gorm.io/gorm"
	dbschema "gorm.io/gorm/schema"
	"time"
)

type Dao struct {
	*impl.DaoBase
	cfg      *db.DaoConfig
	db       *gorm.DB
	dao      *ddd_sql.MapDao
	dbSchema *dbschema.Schema
}

func NewDao(cfg *db.DaoConfig) db.Dao {
	cfg.Valid()
	var database *gorm.DB
	//eb := ddd.NewMapEntityBuilder[map[string]any]()
	if cfg.Database != nil {
		if val, ok := cfg.Database.(*gorm.DB); ok {
			database = val
		} else {
			panic("database config error")
		}
	}

	if database == nil {
		item := restapp.GetDb(cfg.DbKey)
		if item == nil {
			panic("db item not found")
		}
		if val, ok := item.GetDB().(*gorm.DB); ok {
			database = val
		} else {
			panic("database config error")
		}
	}
	dbSchema, err := NewDBSchema(cfg.Schema)
	if err != nil {
		panic(err)
	}

	sqlDao := ddd_sql.NewMapDao(database, cfg.DbKey, cfg.Schema.Name)
	sqlDao.AddMetadata("dbSchema", dbSchema)
	sqlDao.AddMetadata("schema", cfg.Schema)

	daoBase := impl.NewDaoBase(sqlDao, cfg)

	return &Dao{
		DaoBase:  daoBase,
		dao:      sqlDao,
		cfg:      cfg,
		db:       database,
		dbSchema: dbSchema,
	}
}

func (d *Dao) Table() db.Table {
	return newTable(d.db, d.cfg.Schema, d.dbSchema)
}

func init() {
	gorm.GetGoToDbValue = func(db *gorm.DB, field *dbschema.Field, value any) (any, bool) {
		// liuxd lxd
		if field.DataType == dbschema.Object || field.DataType == dbschema.Array {
			if val, err := json.Marshal(value); err == nil {
				return val, true
			} else {
				panic(err)
			}
		} else if field.DataType == dbschema.Time {
			if val, ok := value.(times.IDate); ok {
				return val.Date(), true
			} else if val, ok := value.(times.ITime); ok {
				return val.Time(), true
			} else {
				return value, false
			}
		}
		return value, false
	}

	gorm.GetDbToGoValue = func(db *gorm.DB, field *dbschema.Field, value any) (any, bool) {
		if field.DataType == dbschema.Object || field.DataType == dbschema.Array {
			val := getJsonText(value)
			var data any
			if field.DataType == dbschema.Object {
				data = make(map[string]any)
			} else {
				data = make([]any, 0)
			}
			if len(val) > 0 {
				err := json.Unmarshal([]byte(val), &data)
				if err != nil {
					panic(err)
				}
			}
			return data, true
		} else if field.DataType == dbschema.Time {
			return newTime(value), true
		}
		return value, false
	}
}

func newTime(value any) *times.Time {
	if value != nil {
		if val, ok := value.(**time.Time); ok {
			return times.NewTime(*val)
		} else if val, ok := value.(time.Time); ok {
			return times.NewTime(&val)
		} else if val, ok := value.(*time.Time); ok {
			return times.NewTime(val)
		} else if val, ok := value.(string); ok {
			if tVal, err := times.NewTimeWithString(val); err == nil {
				return tVal
			} else {
				panic(err)
			}
		}
	}
	return nil
}

func getJsonText(value any) string {
	res := ""
	if s, ok := value.(**string); ok {
		if s != nil && *s != nil {
			res = **s
		}

	} else if s, ok := value.(*string); ok {
		if s != nil {
			res = *s
		}
	} else if s, ok := value.(string); ok {
		res = s
	}
	return res
}
