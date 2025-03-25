package mongodb

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/core/restapp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/store/store_mongodb"
)

type RepositoryOptions struct {
	MongoDB         *store_mongodb.MongoDB
	GetCollCallback GetCollectionCallback
}

func NewRepositoryOptions(opts ...*RepositoryOptions) *RepositoryOptions {
	o := &RepositoryOptions{}
	for _, item := range opts {
		if item.MongoDB != nil {
			o.MongoDB = item.MongoDB
		}
		if item.GetCollCallback != nil {
			o.GetCollCallback = item.GetCollCallback
		}
	}
	if o.MongoDB == nil {
		o.MongoDB = _mongodb
	}
	if o.MongoDB == nil {
		o.MongoDB = GetDB()
	}
	return o
}

func GetDB() *store_mongodb.MongoDB {
	if _mongodb != nil {
		return _mongodb
	}
	return restapp.GetMongoDB()
}
