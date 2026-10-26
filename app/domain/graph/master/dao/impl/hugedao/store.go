package hugedao

import (
	"sync"

	ragconfigs "github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/configs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/hugegraph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
)

var _store *hugegraph.Store
var _storeOnce sync.Once

func GetStore() (*hugegraph.Store, error) {
	var err error
	_storeOnce.Do(func() {
		e := envs.GetEnv()
		cfg := ragconfigs.GetConfig()
		db := e.GetHugeConfig(cfg.HugeDbKey)
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
