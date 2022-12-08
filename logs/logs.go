package logs

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/sirupsen/logrus"
)

type LogFunction logrus.LogFunction

type Fields map[string]interface{}

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

type loggerKey struct {
}

func NewContext(parentCtx context.Context, logger Logger) context.Context {
	ctx := context.WithValue(parentCtx, loggerKey{}, logger)
	return ctx
}

func GetLogger(ctx context.Context) Logger {
	if v := ctx.Value(loggerKey{}); v != nil {
		if logger, ok := v.(Logger); ok {
			return logger
		}
	}
	return nil
}

func Trace(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Trace(args...)
	}
}

func Debug(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Debug(args...)
	}
}

func Print(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Print(args...)
	}
}

func Info(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Info(args...)
	}
}

func Warn(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Warn(args...)
	}
}

func Warning(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Warning(args...)
	}
}

func Error(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Error(args...)
	}
}

func Panic(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Panic(args...)
	}
}

func Fatal(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Fatal(args...)
	}
}

func Tracef(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Tracef(format, args...)
	}
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Debugf(format, args...)
	}
}

func Printf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Printf(format, args...)
	}
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Infof(format, args...)
	}
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Warnf(format, args...)
	}
}

func Warningf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Warningf(format, args...)
	}
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Errorf(format, args...)
	}
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Panicf(format, args...)
	}
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Fatalf(format, args...)
	}
}

func Traceln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Traceln(args...)
	}
}
func Debugln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Debugln(args...)
	}
}

func Println(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Println(args...)
	}
}

func Infoln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Infoln(args...)
	}
}

func Warnln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Warnln(args...)
	}
}

func Warningln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Warnln(args...)
	}
}

func Errorln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Errorln(args...)
	}
}

func Panicln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Panicln(args...)
	}
}

func Fatalln(ctx context.Context, args ...interface{}) {
	if l := GetLogger(ctx); l != nil {
		l.Fatalln(args...)
	}
}

func DebugEvent(ctx context.Context, event ddd.Event, funcName string) {
	Debugf(ctx, `eventId:'%v',eventType:'%v', tenantId:'%v', commandId:'%v', funcName:'%v'`, event.GetEventId(), event.GetEventType(), event.GetTenantId(), event.GetCommandId(), funcName)
}
