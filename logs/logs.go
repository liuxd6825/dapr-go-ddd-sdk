package logs

import (
	"context"
	"encoding/json"
	"github.com/liuxd6825/dapr-go-ddd-sdk/auth"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/idutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/sirupsen/logrus"
	"time"
)

type LogFunction logrus.LogFunction

type Fields = logrus.Fields

// Logger 日志接口类
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
}

// ArgFunc 异步参数方法，当符合当前日志级别时调用
type ArgFunc = func() any

// Level 日志级别
type Level = logrus.Level

var loggerLevel Level = logrus.DebugLevel

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

type Event interface {
	GetTenantId() string
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
	if loggerLevel >= lvl {
		return true
	}
	return false
}

func getFields(ctx context.Context, tenantId string, fields Fields) Fields {
	f := Fields{}
	if user, err := auth.GetLoginUser(ctx); err == nil {
		f["userId"] = user.GetId()
		f["userName"] = user.GetName()
	}
	for key, val := range fields {
		if fun, ok := val.(ArgFunc); ok {
			f[key] = fun()
		} else {
			f[key] = val
		}
	}
	if tenantId!=""{
		f["tenantId"] = tenantId
	} else 	if id , ok := f["tenantId"]; ok {
		if id==""{
			delete(f, "tenantId")
		}
	}

	return f
}

func write(ctx context.Context, tenantId string, fields Fields, level Level, args []any, fun func(ctx context.Context, l Logger, args ...any)) {
	if !isLevel(logger, level) {
		return
	}
	print(ctx, tenantId, fields, args, fun)
}

func print(ctx context.Context, tenantId string, fields Fields, args []any, fun func(ctx context.Context, l Logger, args ...any)) {
	var entry *logrus.Entry
	arg := getArgs(getArgs(args...)...)
	fs := getFields(ctx, tenantId, fields)
	entry = logger.WithFields(fs)
	fun(ctx, entry, arg...)
}

// SetLevel sets the logger level.
func SetLevel(ctx context.Context, level Level) {
	loggerLevel = level
	logger.SetLevel(level)
}

// GetLevel returns the logger level.
func GetLevel(ctx context.Context) Level {
	return loggerLevel
}

func Trace(ctx context.Context, tenantId string, fields Fields, args ...interface{}) {
	write(ctx, tenantId, fields, TraceLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Trace(args...)
	})
}

func Print(ctx context.Context, tenantId string, fields Fields) {
	print(ctx, tenantId, fields, nil, func(ctx context.Context, l Logger, args ...any) {
		l.Print()
	})
}

func Printf(ctx context.Context, tenantId string, fields Fields, fmt string, args ...any) {
	print(ctx, tenantId, fields, nil, func(ctx context.Context, l Logger, args ...any) {
		l.Printf(fmt, args)
	})
}

func Println(ctx context.Context, tenantId string, fields Fields, fmt string, args ...any) {
	print(ctx, tenantId, fields, nil, func(ctx context.Context, l Logger, args ...any) {
		l.Println()
	})
}

func Debug(ctx context.Context, tenantId string, fields Fields) {
	write(ctx, tenantId, fields, DebugLevel, nil, func(ctx context.Context, l Logger, args ...any) {
		l.Debug()
	})
}

func Debugf(ctx context.Context, tenantId string, fields Fields, fmt string, args ...interface{}) {
	write(ctx, tenantId, fields, DebugLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Debugf(fmt, args...)
	})
}

func DebugEvent(ctx context.Context, event Event, funcName string) {
	eventFunc := func() any {
		data, _ := json.Marshal(event)
		return string(data)
	}
	Debug(ctx, event.GetTenantId(), Fields{"event": eventFunc, "func": funcName})
}

func Debugfmt(ctx context.Context, tenantId string, fmt string, args ...interface{}) {
	write(ctx, tenantId, nil, ErrorLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Debugf(fmt, args...)
	})
}

func Info(ctx context.Context, tenantId string, fields Fields, args ...interface{}) {
	write(ctx, tenantId, fields, InfoLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Info(args...)
	})
}

func Infof(ctx context.Context, tenantId string, fields Fields, fmt string, args ...interface{}) {
	write(ctx, tenantId, fields, InfoLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Infof(fmt, args...)
	})
}

func Warn(ctx context.Context, tenantId string, fields Fields) {
	write(ctx, tenantId, fields, WarnLevel, nil, func(ctx context.Context, l Logger, args ...any) {
		l.Warn()
	})
}

func Warnf(ctx context.Context, tenantId string, fields Fields, fmt string, args ...interface{}) {
	write(ctx, tenantId, fields, WarnLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Infof(fmt, args...)
	})
}

func Warning(ctx context.Context, tenantId string, fields Fields, args ...interface{}) {
	write(ctx, tenantId, fields, WarnLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Warning(args...)
	})
}

func Error(ctx context.Context, tenantId string, fields Fields, args ...interface{}) {
	write(ctx, tenantId, fields, ErrorLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Error(args...)
	})
}

func ErrorErr(ctx context.Context, tenantId string, err error) {
	if err == nil {
		return
	}
	fields := Fields{
		"error": err.Error(),
	}
	write(ctx, tenantId, fields, ErrorLevel, nil, func(ctx context.Context, l Logger, args ...any) {
		l.Error(args...)
	})
}

func Errorf(ctx context.Context, tenantId string, fields Fields, fmt string, args ...interface{}) {
	write(ctx, tenantId, fields, ErrorLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Errorf(fmt, args...)
	})
}

func Errorfmt(ctx context.Context, tenantId string, fmt string, args ...interface{}) {
	write(ctx, tenantId, nil, ErrorLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Errorf(fmt, args...)
	})
}

func Panic(ctx context.Context, tenantId string, fields Fields, args ...interface{}) {
	write(ctx, tenantId, fields, PanicLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Panic(args...)
	})
}

func Fatal(ctx context.Context, tenantId string, fields Fields, args ...interface{}) {
	write(ctx, tenantId, fields, FatalLevel, args, func(ctx context.Context, l Logger, args ...any) {
		l.Panic(args...)
	})
}

// DebugStart
//
//	@Description:
//	@param ctx
//	@param tenantId
//	@param fields
//	@param fun
//	@param format
//	@param args
//	@return err
func DebugStart(ctx context.Context, tenantId string, fields Fields, fun func() error) (err error) {
	if isLevel(logger, DebugLevel) {
		logId := idutils.NewId()
		funcName := reflectutils.RunFuncName(2)
		fs := Fields{
			"logId":   logId,
			"logFunc": funcName,
			"logType": "start",
		}
		for key, val := range fields {
			fs[key] = val
		}

		write(ctx, tenantId, fs, DebugLevel, nil, func(ctx context.Context, l Logger, args ...any) {
			l.Debug()
		})

		startTime := time.Now()

		defer func() {
			var params []any
			useTime := time.Now().Sub(startTime)
			fs["logUseTime"] = useTime
			fs["logType"] = "end"

			err = errors.GetRecoverError(err, recover())
			if err != nil {
				fs["error"] = err.Error()
				params = append(params, err.Error())
			}
			write(ctx, tenantId, fs, DebugLevel, nil, func(ctx context.Context, l Logger, args ...any) {
				l.Debug()
			})
		}()

		err = fun()
	}

	return err
}

func ParseLevel(lvl string) (Level, error) {
	l, err := logrus.ParseLevel(lvl)
	return l, err
}
