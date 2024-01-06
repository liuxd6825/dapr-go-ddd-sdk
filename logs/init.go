package logs

import (
	"context"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type logHook struct {
	level Level
}

type Formatter struct {
}

var (
	logger      *logrus.Logger
	once        sync.Once
	formatter   = &Formatter{}
	startFields = []string{"tenantId", "userId", "userName", "pkg", "logId", "logType", "logFunc", "func", "error", "logUseTime"}
)

// init
//
//	@Description: 初始化基于控制台的日志实例
func init() {
	level := DebugLevel
	logger = logrus.New()
	logger.SetLevel(level)
	logger.SetReportCaller(false)
	logger.SetFormatter(formatter)
}

// Init
//
//	@Description: 初始化基于文件的日志
//	@param saveFile
//	@param level
//	@param saveDays
//	@param rotationHour
//	@return Logger
func Init(saveFile string, level Level, saveDays int, rotationHour int) Logger {
	once.Do(func() {
		logrus.SetLevel(level)
		logrus.SetReportCaller(false)
		logrus.SetFormatter(formatter)

		logFile := saveFile + ".%Y-%m-%d-%H.log"
		// 配置日志每隔 1 小时轮转一个新文件，保留最近 30 天的日志文件，多余的自动清理掉。
		writer, _ := rotatelogs.New(
			logFile,
			rotatelogs.WithLinkName(saveFile),
			rotatelogs.WithMaxAge(time.Duration(24*saveDays)*time.Hour),
			rotatelogs.WithRotationTime(time.Duration(rotationHour)*time.Hour),
		)

		logger = logrus.New()
		logger.Hooks.Add(&logHook{level: level})
		logger.SetOutput(writer)
		logger.SetFormatter(formatter)
		logger.SetLevel(level)
		logger.SetReportCaller(false)
		logger.Infof("ctype=app; logFile=%s", logFile)

	})
	return logger
}

func (h *logHook) Levels() []logrus.Level {
	return []Level{DebugLevel, InfoLevel, ErrorLevel, PanicLevel, WarnLevel, TraceLevel}
}

func (h *logHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case DebugLevel:
		logrus.WithFields(entry.Data).Debug(entry.Message)
	case InfoLevel:
		logrus.WithFields(entry.Data).Info(entry.Message)
	case ErrorLevel:
		logrus.WithFields(entry.Data).Error(entry.Message)
	case PanicLevel:
		logrus.WithFields(entry.Data).Panic(entry.Message)
	case WarnLevel:
		logrus.WithFields(entry.Data).Warn(entry.Message)
	case TraceLevel:
		logrus.WithFields(entry.Data).Trace(entry.Message)
	}
	return nil
}

func (f *Formatter) Format(e *logrus.Entry) ([]byte, error) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(`{"time":"%s", "level":"%s", `, e.Time.Format("2006-01-02 03:04:05.0000"), e.Level.String()))

	for _, name := range startFields {
		if val, ok := e.Data[name]; ok {
			sb.WriteString(fmt.Sprintf(`"%s":"%s", `, name, val))
		}
	}
	isHas := false
	for key, val := range e.Data {
		isHas = false
		for _, name := range startFields {
			if key == name {
				isHas = true
				break
			}
		}
		if !isHas {
			sb.WriteString(fmt.Sprintf(`"%s":"%v", `, key, val))
		}
	}
	if e.Message != "" {
		sb.WriteString(fmt.Sprintf(`"msg":"%s" }`, e.Message))
		sb.WriteString("\r\n")
	} else {
		s := sb.String()
		return []byte(s[0:len(s)-2] + "}\r\n"), nil
	}
	return []byte(sb.String()), nil
}

func NewContext(parentCtx context.Context) context.Context {
	ctx := context.WithValue(parentCtx, loggerKey{}, logger)
	return ctx
}

func GetLogger() Logger {
	return logger
}
