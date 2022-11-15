package restapp

import "github.com/sirupsen/logrus"

var _logger *logrus.Logger

func init() {
	_logger = logrus.New()
	_logger.SetLevel(logrus.DebugLevel)
}
