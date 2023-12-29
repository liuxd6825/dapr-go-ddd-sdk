package logs

import (
	"context"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var _logger *logrus.Logger
var _once sync.Once

type logHook struct {
	level Level
}

func (h *logHook) Levels() []logrus.Level {
	return []Level{DebugLevel, InfoLevel, ErrorLevel, PanicLevel, WarnLevel, TraceLevel}
}

func (h *logHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case DebugLevel:
		logrus.Debug(entry.Message)
	case InfoLevel:
		logrus.Info(entry.Message)
	case ErrorLevel:
		logrus.Error(entry.Message)
	case PanicLevel:
		logrus.Panic(entry.Message)
	case WarnLevel:
		logrus.Warn(entry.Message)
	case TraceLevel:
		logrus.Trace(entry.Message)
	}
	return nil
}

func Init(saveFile string, level Level, saveDays int, rotationHour int) Logger {
	_once.Do(func() {
		// 配置日志每隔 1 小时轮转一个新文件，保留最近 30 天的日志文件，多余的自动清理掉。
		writer, _ := rotatelogs.New(
			saveFile+".%Y-%m-%d-%H.log",
			rotatelogs.WithLinkName(saveFile),
			rotatelogs.WithMaxAge(time.Duration(24*saveDays)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(rotationHour)*time.Hour),
		)
		_logger = logrus.New()
		_logger.Hooks.Add(&logHook{level: level})
		_logger.SetOutput(writer)

		_logger.SetLevel(level)
		_logger.SetReportCaller(true)

		logrus.SetLevel(level)
		logrus.SetReportCaller(true)

	})
	return _logger

}

func NewContext(parentCtx context.Context, logger Logger) context.Context {
	ctx := context.WithValue(parentCtx, loggerKey{}, logger)
	return ctx
}

func GetLogrus() Logger {
	return _logger
}

func GetLogger(ctx context.Context) Logger {
	if v := ctx.Value(loggerKey{}); v != nil {
		if logger, ok := v.(Logger); ok {
			return logger
		}
	}
	return nil
}
