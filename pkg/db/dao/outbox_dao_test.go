package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	gormschema "gorm.io/gorm/schema"
	"sync"
	"testing"
)

func Test_NewOutboxDao(t *testing.T) {
	data := &dbevent.Outbox{}
	gSch, err := gormschema.ParseWithSpecialTableName(data, &sync.Map{}, gormschema.NamingStrategy{}, "sys_outbox")
	if gSch == nil && err != nil {
		panic(err)
	}
	t.Log(gSch)
}
