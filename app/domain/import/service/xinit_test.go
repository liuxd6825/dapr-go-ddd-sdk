package service

import (
	"context"
	"os"
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/spark/livyserver"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

const (
	LivyServerHost = "192.168.120.200"
	LivyServerPort = 8998
)

var livyServer *livyserver.LivyServer
var ctx context.Context
var recordService *RecordService

func TestMain(m *testing.M) {
	livyServer = livyserver.NewLivyServer(LivyServerHost, LivyServerPort)
	ctx = xtest2.NewContext()
	xtest2.InitEnv_MongoRemoteMaster()
	recordService = NewRecordService(livyServer)
	os.Exit(m.Run())
}
