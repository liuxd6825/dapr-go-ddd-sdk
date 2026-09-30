package service_test

import (
	"context"
	_ "embed"
	"os"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

var envInv *env.Env
var ctx context.Context

func TestMain(m *testing.M) {
	envInv = xtest2.NewEnvConfigHuge(xtest2.HugeOptions{
		Addr:     "192.168.120.200",
		Port:     18080,
		DBKey:    "huge",
		Database: "dev",
		UserName: "admin",
		Password: "admin",
	})
	env.SetEnv(envInv)
	ctx = xtest2.NewContext()
	os.Exit(m.Run())
}
