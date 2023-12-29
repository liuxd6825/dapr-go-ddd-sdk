package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/sirupsen/logrus"
	"time"
)

type LogFunction logrus.LogFunction

type Fields map[string]interface{}

// Logger 日志借口类
type Logger interface {
	Trace(args ...interface{})
	Debug(args ...interface{})
	Print(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Traceln(args ...interface{})
	Debugln(args ...interface{})
	Println(args ...interface{})
	Infoln(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Panicln(args ...interface{})
	Fatalln(args ...interface{})

	SetLevel(level Level)
	GetLevel() Level

	//WithField(key string, value interface{}) Entry
}

type Entry interface {
	//WithField(key string, value interface{}) Entry
}

// ArgFunc 异步参数方法，当符合当前日志级别时调用
type ArgFunc = func() any

// Level 日志级别
type Level = logrus.Level

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = logrus.PanicLevel
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = logrus.DebugLevel
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = logrus.TraceLevel
)

type loggerKey struct {
}

func getArgs(args ...any) []any {
	var res []any
	for _, arg := range args {
		if fun, ok := arg.(ArgFunc); ok {
			res = append(res, fun())
		} else {
			res = append(res, arg)
		}
	}
	return res
}

func isLevel(logger Logger, lvl Level) bool {
	if logger.GetLevel() >= lvl {
		return true
	}
	return false
}

// SetLevel sets the logger level.
func SetLevel(ctx context.Context, level Level) {
	l := GetLogger(ctx)
	if l != nil {
		l.SetLevel(level)
	}
}

// GetLevel returns the logger level.
func GetLevel(ctx context.Context) (bool, Level) {
	l := GetLogger(ctx)
	if l != nil {
		return true, l.GetLevel()
	}
	return false, PanicLevel
}

func Trace(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, TraceLevel) {
			l.Trace(getArgs(getArgs(args...)...))
		}
	}
}

func Debug(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, DebugLevel) {
			l.Debug(getArgs(getArgs(args...)...))
		}
	}
}

func Print(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Print(getArgs(getArgs(args...)...))
	}
}

func Info(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, InfoLevel) {
			l.Info(getArgs(getArgs(args...)...))
		}
	}
}

func Warn(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, WarnLevel) {
			l.Warn(getArgs(getArgs(args...)...))
		}
	}
}

func Warning(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, WarnLevel) {
			l.Warning(getArgs(args...)...)
		}
	}
}

func Error(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, ErrorLevel) {
			l.Error(getArgs(args...)...)
		}
	}
}

func Panic(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, PanicLevel) {
			l.Panic(getArgs(args...)...)
		}
	}
}

func Fatal(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, FatalLevel) {
			l.Fatal(getArgs(args...)...)
		}
	}
}

func Tracef(ctx context.Context, format string, args ...any) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, TraceLevel) {
			l.Tracef(format, getArgs(args...)...)
		}
	}
}

func Debugf(ctx context.Context, format string, args ...any) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, DebugLevel) {
			l.Debugf(format, getArgs(args...)...)
		}
	}
}

func DebugStart(ctx context.Context, fun func() error, format string, args ...any) (err error) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, DebugLevel) {
			logId := idutils.NewId()
			funcName := reflectutils.RunFuncName(2)
			debugMsg := fmt.Sprintf("id=%s; type=start; func=%s; info=<<! %s !>>;", logId, funcName, format)
			l.Debugf(debugMsg, getArgs(args...)...)
			startTime := time.Now()

			defer func() {
				var params []any
				params = append(params, args...)

				userTime := time.Now().Sub(startTime)

				debugMsg = fmt.Sprintf("id=%s; type=end;   func=%s; info=<<! %s !>>; useTime=%v", logId, funcName, format, userTime)
				err = errors.GetRecoverError(err, recover())
				if err != nil {
					debugMsg += "; error=%s"
					params = append(params, err.Error())
				}
				debugMsg += ";"

				l.Debugf(debugMsg, getArgs(params...)...)
			}()

			err = fun()
		}
	}
	return err
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Printf(format, getArgs(args...)...)
	}
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, DebugLevel) {
			l.Infof(format, getArgs(args...)...)
		}
	}
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, WarnLevel) {
			l.Warnf(format, getArgs(args...)...)
		}
	}
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, WarnLevel) {
			l.Warningf(format, args...)
		}
	}
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, ErrorLevel) {
			l.Errorf(format, getArgs(args...)...)
		}
	}
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, PanicLevel) {
			l.Panicf(format, getArgs(args...)...)
		}
	}
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, FatalLevel) {
			l.Fatalf(format, getArgs(args...)...)
		}
	}
}

func Traceln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, TraceLevel) {
			l.Traceln(getArgs(args...)...)
		}
	}
}
func Debugln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, DebugLevel) {
			l.Debugln(getArgs(args...)...)
		}
	}
}

func Println(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Println(getArgs(args...)...)
	}
}

func Infoln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, InfoLevel) {
			l.Infoln(getArgs(args...)...)
		}
	}
}

func Warnln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, WarnLevel) {
			l.Warnln(getArgs(args...)...)
		}
	}
}

func Warningln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, WarnLevel) {
			l.Warnln(getArgs(args...)...)
		}
	}
}

func Errorln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, ErrorLevel) {
			l.Errorln(getArgs(args...)...)
		}
	}
}

func Panicln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, PanicLevel) {
			l.Panicln(getArgs(args...)...)
		}
	}
}

func Fatalln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		if isLevel(l, FatalLevel) {
			l.Fatalln(getArgs(args...)...)
		}
	}
}

func DebugEvent(ctx context.Context, event any, funcName string) {
	eventFunc := func() any {
		data, _ := json.Marshal(event)
		return string(data)
	}
	Debugf(ctx, `event:%v, funcName:'%v'`, eventFunc, funcName)
}

func ParseLevel(lvl string) (Level, error) {
	l, err := logrus.ParseLevel(lvl)
	return l, err
}
