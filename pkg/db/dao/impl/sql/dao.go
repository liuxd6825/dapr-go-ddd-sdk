package sql

import (
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_sql"
	idao "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/impl"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"gorm.io/gorm"
	gormschema "gorm.io/gorm/schema"
	"sync"
	"time"
)

type Dao[T any] struct {
	*impl.DaoBase[T]
	cfg       *idao.DaoConfig
	db        *gorm.DB
	dao       store.IStore[T]
	gormSch   *gormschema.Schema
	tableName string
}

func NewDao[T any](cfg *idao.DaoConfig, tableNames ...string) idao.Dao[T] {
	cfg.Valid()
	var db *gorm.DB
	//eb := ddd.NewMapEntityBuilder[map[string]any]()
	if cfg.DB != nil {
		if val, ok := cfg.DB.(*gorm.DB); ok {
			db = val
		} else {
			panic("database config error")
		}
	}

	if db == nil {
		item := env.GetDB(cfg.DbKey)
		if item == nil {
			panic("db item not found")
		}
		if val, ok := item.GetDB().(*gorm.DB); ok {
			db = val
		} else {
			panic("database config error")
		}
	}
	tableName := cfg.DBSchema.TableName
	if tableName == "" {
		tableName = cfg.DBSchema.Name
	}
	if len(tableNames) > 0 {
		tableName = tableNames[0]
	}
	var t T
	var gormSch *gormschema.Schema
	if reflectutils.IsMapStringKey[T](t) {
		sch, err := NewGormSchema(cfg.DBSchema)
		if err != nil {
			panic(err)
		}
		gormSch = sch
	} else {
		ins := reflectutils.NewInstance[T]()
		sch, err := gormschema.Parse(ins, &sync.Map{}, gormschema.NamingStrategy{})
		if sch == nil && err != nil {
			panic(err)
		}
		gormSch = sch
	}

	newDaoCfg := &store_sql.NewConfig{
		DbKey:      cfg.DbKey,
		Db:         db,
		TableName:  tableName,
		DBSchema:   cfg.DBSchema,
		GormSchema: gormSch,
	}

	sqlDao := store_sql.NewDao[T](newDaoCfg)

	daoBase := impl.NewDaoBase[T](sqlDao, cfg)

	return &Dao[T]{
		DaoBase:   daoBase,
		dao:       sqlDao,
		cfg:       cfg,
		db:        db,
		gormSch:   gormSch,
		tableName: tableName,
	}
}

func (d *Dao[T]) Table() idao.Table {
	return newTable(d.db, d.tableName, d.dao.NewEntity(), d.cfg.DBSchema, d.gormSch)
}

func init() {
	gorm.GetGoToDbValue = func(db *gorm.DB, field *gormschema.Field, value any) (any, bool) {
		// liuxd lxd
		if field.DataType == dbschema.Json {
			if value == nil {
				return nil, false
			}
			if val, err := json.Marshal(value); err == nil {
				return val, true
			} else {
				panic(err)
			}
		} else if field.DataType == gormschema.Time {
			if val, ok := value.(times.IDate); ok {
				return val.Date(), true
			} else if val, ok := value.(times.ITime); ok {
				return val.Time(), true
			} else {
				return value, false
			}
		} else if field.DataType == gormschema.Date {
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

	gorm.GetDbToGoValue = func(db *gorm.DB, field *gormschema.Field, value any) (any, bool) {
		if field.DataType == gormschema.Json {
			if value == nil {
				return nil, false
			}
			val := getJsonText(value)
			var data any
			if len(val) > 0 {
				err := json.Unmarshal([]byte(val), &data)
				if err != nil {
					panic(err)
				}
			}
			return data, true
		} else if field.DataType == gormschema.Time {
			return newTime(value), true
		} else if field.DataType == gormschema.Date {
			return newDate(value), true
		}
		return value, false
	}
}

func newDate(value any) *times.Date {
	if value != nil {
		if val, ok := value.(**time.Time); ok {
			return times.GetDate(*val)
		} else if val, ok := value.(time.Time); ok {
			return times.GetDate(&val)
		} else if val, ok := value.(*time.Time); ok {
			return times.GetDate(val)
		} else if val, ok := value.(string); ok {
			if tVal, err := times.NewDateWithString(val); err == nil {
				return tVal
			} else {
				panic(err)
			}
		}
	}
	return nil
}

func newTime(value any) *times.Time {
	if value != nil {
		if val, ok := value.(**time.Time); ok {
			return times.GetTime(*val)
		} else if val, ok := value.(time.Time); ok {
			return times.GetTime(&val)
		} else if val, ok := value.(*time.Time); ok {
			return times.GetTime(val)
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
