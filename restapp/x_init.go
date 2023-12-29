package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"os"
	"path/filepath"
)

func init() {

}

func initLogs(level logs.Level) {
	appPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir, exec := filepath.Split(appPath)
	logPath := fmt.Sprintf("%s/logs/%s.log", dir, exec)
	logs.Init(logPath, level)
}

func GetLogger() logs.Logger {
	return logs.GetLogrus()
}
