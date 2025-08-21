package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
)

func init() {
	xtest.Init(xtest.TestType_MongoLocal)
}
