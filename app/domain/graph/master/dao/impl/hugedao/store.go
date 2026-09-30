package hugedao

import (
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/hugegraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
)

var _store *hugegraph.Store
var _storeOnce sync.Once

func GetStore() (*hugegraph.Store, error) {
	var err error
	_storeOnce.Do(func() {
		e := env.GetEnv()
		db, ok := e.Huge[DB_KEY_HUGE]
		if !ok {
			panic("no huge")
		}
		var pool *hugegraph.ClientPool
		pool, err = hugegraph.NewClientPool(hugegraph.HugeConfig{
			Host:     db.Host,
			Username: db.Username,
			Password: db.Password,
			Port:     db.Port,
		}, hugegraph.GetPoolConfigDefault())
		if err != nil {
			return
		}
		_store = hugegraph.NewStore(pool)
	})
	return _store, nil
}
