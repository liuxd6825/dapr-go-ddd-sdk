package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"os"
	"path/filepath"
)

func init() {

}

func initLogs(level logs.Level, saveDays int, rotationHour int) {
	appPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir, appName := filepath.Split(appPath)
	saveFile := fmt.Sprintf("%s/logs/%s", dir, appName)
	logs.Init(saveFile, level, saveDays, rotationHour)
}

func GetLogger() logs.Logger {
	return logs.GetLogrus()
}
