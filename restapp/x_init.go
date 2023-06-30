package restapp

import (
	"github.com/sirupsen/logrus"
	"sync"
)

var _logger *logrus.Logger
var _once sync.Once

func init() {
	InitLogs()
}

func InitLogs() *logrus.Logger {
	_once.Do(func() {
		_logger = logrus.New()
		_logger.SetLevel(logrus.DebugLevel)
	})
	return _logger
}
