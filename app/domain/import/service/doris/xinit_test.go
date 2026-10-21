package doris_test

import (
	"context"
	"os"
	"testing"

	envs "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/envs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/miniofs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/boolutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

var env *envs.Env
var fs *miniofs.Fs
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = xtest2.NewContext()
	env = xtest2.NewEnvConfig()
	env.AddMinio(&envs.Minio{
		Name:      "minio",
		Endpoint:  "192.168.120.200",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		Buckets: map[string]string{
			"import-data": "import-data",
		},
	})
	env.AddMySql(&envs.MySql{
		DbKey:     "doris",
		User:      "root",
		Host:      "192.168.120.224",
		Port:      "9030",
		Password:  "",
		DbName:    "master",
		Loc:       "Asia/Shanghai",
		ParseTime: boolutils.PBool(true),
		Charset:   "utf8mb4",
	})
	env.Init()

	envs.SetEnv(env)
	mfs, err := miniofs.NewFs(env, "minio", "import-data")
	if err != nil {
		panic(err)
	}
	fs = mfs

	os.Exit(m.Run())
}
