package restapp

import (
	"github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-sdk/actor"
	"github.com/liuxd6825/dapr-go-sdk/actor/config"
	actorError "github.com/liuxd6825/dapr-go-sdk/actor/error"
	"github.com/liuxd6825/dapr-go-sdk/actor/runtime"
	"net/http"
)

// Deprecated: Use RegisterActorImplFactoryContext instead.
func (s *HttpServer) RegisterActorImplFactory(f actor.Factory, opts ...config.Option) {
	runtime.GetActorRuntimeInstance().RegisterActorFactory(f, opts...)
}

func (s *HttpServer) RegisterActorImplFactoryContext(f actor.FactoryContext, opts ...config.Option) {
	runtime.GetActorRuntimeInstanceContext().RegisterActorFactory(f, opts...)
}

// register actor method invoke handler
func (s *HttpServer) actorInvokeHandler(ictx *context.Context) {
	const funLog = "restapp.HttpServer.actorInvokeHandler()"
	ctx, _ := NewContextNoAuth(ictx)
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			ictx.StatusCode(http.StatusInternalServerError)
			fields := logs.Fields{
				"func":  funLog,
				"error": err.Error(),
			}
			logs.Info(ctx, "", fields)
		}
	}()

	actorType := ictx.Params().Get("actorType")
	actorId := ictx.Params().Get("actorId")
	methodName := ictx.Params().Get("methodName")
	reqData, _ := ictx.GetBody()
	rspData, actorErr := runtime.GetActorRuntimeInstanceContext().InvokeActorMethod(ctx, actorType, actorId, methodName, reqData)
	if actorErr != actorError.Success {
		fields := newActorFieldError(funLog, actorType, actorId, methodName, ActorErrToError(actorErr))
		logs.Error(ctx, "", fields)
	}
	if actorErr == actorError.ErrActorTypeNotFound || actorErr == actorError.ErrActorIDNotFound {
		ictx.ResponseWriter().WriteHeader(http.StatusNotFound)
		return
	}
	if actorErr != actorError.Success {
		ictx.ResponseWriter().WriteHeader(http.StatusInternalServerError)
		return
	}
	ictx.StatusCode(http.StatusOK)
	_, _ = ictx.Write(rspData)
}

// register actor config handler
func (s *HttpServer) actorConfigHandler(ictx *context.Context) {
	ctx, _ := NewContextNoAuth(ictx)
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			logs.Error(ctx, "", logs.Fields{"func": "restapp.HttpServer.actorConfigHandler()", "error": err.Error()})
		}
	}()
	statusCode := http.StatusOK
	data, err := runtime.GetActorRuntimeInstanceContext().GetJSONSerializedConfig()
	if err != nil {
		statusCode = http.StatusInternalServerError
	} else if _, err = ictx.Write(data); err != nil {
		statusCode = http.StatusInternalServerError
	}

	if statusCode == http.StatusOK {
		logs.Info(ctx, "", logs.Fields{"func": "restapp.HttpServer.actorConfigHandler()", "data": string(data)})
	} else {
		logs.Error(ctx, "", logs.Fields{"func": "restapp.HttpServer.actorConfigHandler()", "error": err.Error()})
	}
	ictx.StatusCode(statusCode)
}

// register actor reminder invoke handler
func (s *HttpServer) actorReminderInvokeHandler(ictx *context.Context) {
	const funLog = "restapp.HttpServer.actorReminderInvokeHandler()"
	ctx, _ := NewContextNoAuth(ictx)
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			ictx.StatusCode(http.StatusInternalServerError)
			logs.Info(ctx, "", logs.Fields{"func": funLog, "error": err.Error()})
		}
	}()
	actorType := ictx.Params().Get("actorType")
	actorId := ictx.Params().Get("actorId")
	reminderName := ictx.Params().Get("reminderName")
	reqData, _ := ictx.GetBody()

	actorErr := runtime.GetActorRuntimeInstanceContext().InvokeReminder(ctx, actorType, actorId, reminderName, reqData)
	if actorErr != actorError.Success {
		fields := newActorFieldError(funLog, actorType, actorId, reminderName, ActorErrToError(actorErr))
		logs.Error(ctx, "", fields)
	}
	if actorErr == actorError.ErrActorTypeNotFound {
		ictx.ResponseWriter().WriteHeader(http.StatusNotFound)
		return
	}
	if actorErr != actorError.Success {
		ictx.ResponseWriter().WriteHeader(http.StatusInternalServerError)
		return
	}
	ictx.StatusCode(actorErrorAsHttpStatus(actorErr))
}

// register actor timer invoke handler
func (s *HttpServer) actorTimerInvokeHandler(ictx *context.Context) {
	const funLog = "restapp.HttpServer.actorTimerInvokeHandler()"
	ctx, _ := NewContextNoAuth(ictx)
	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			ictx.StatusCode(http.StatusInternalServerError)
			logs.Error(ctx, "", logs.Fields{"func": funLog, "error": err.Error()})
		}
	}()
	actorType := ictx.Params().Get("actorType")
	actorId := ictx.Params().Get("actorId")
	timerName := ictx.Params().Get("timerName")
	reqData, _ := ictx.GetBody()
	actorErr := runtime.GetActorRuntimeInstanceContext().InvokeTimer(ctx, actorType, actorId, timerName, reqData)
	if actorErr != actorError.Success {
		fields := logs.Fields{
			"func":       funLog,
			"actorType":  actorType,
			"actorId":    actorId,
			"timerName":  timerName,
			"reqData":    reqData,
			"actorError": ActorErrToError(actorErr).Error(),
		}
		logs.Error(ctx, "", fields)
	}
	if actorErr == actorError.ErrActorTypeNotFound {
		ictx.ResponseWriter().WriteHeader(http.StatusNotFound)
		return
	}
	if actorErr != actorError.Success {
		ictx.ResponseWriter().WriteHeader(http.StatusInternalServerError)
		return
	}
	ictx.StatusCode(actorErrorAsHttpStatus(actorErr))
}

// register deactivate actor handler
func (s *HttpServer) actorDeactivateHandler(ictx *context.Context) {
	const funLog = "restapp.HttpServer.actorDeactivateHandler()"
	ctx, _ := NewContextNoAuth(ictx)

	defer func() {
		if err := errors.GetRecoverError(nil, recover()); err != nil {
			ictx.StatusCode(http.StatusInternalServerError)
			fields := logs.Fields{
				"func":  funLog,
				"error": err.Error(),
			}
			logs.Error(ctx, "", fields)
		}
	}()

	actorType := ictx.Params().Get("actorType")
	actorId := ictx.Params().Get("actorId")
	actorErr := runtime.GetActorRuntimeInstanceContext().Deactivate(ctx, actorType, actorId)

	if actorErr != actorError.Success && actorErr != actorError.ErrActorIDNotFound && actorErr != actorError.ErrActorTypeNotFound {
		ictx.ResponseWriter().WriteHeader(http.StatusInternalServerError)
		return
	}
	ictx.StatusCode(http.StatusOK)
}

func actorErrorAsHttpStatus(err actorError.ActorErr) int {
	statusCode := http.StatusOK
	if err == actorError.ErrActorTypeNotFound || err == actorError.ErrActorIDNotFound {
		statusCode = http.StatusNotFound
	} else if err != actorError.Success {
		statusCode = http.StatusInternalServerError
	}
	return statusCode
}

func ActorErrToError(actorErr actorError.ActorErr) error {
	msg := ""
	switch actorErr {
	case actorError.ErrActorTypeNotFound:
		msg = "ErrActorTypeNotFound"
		break
	case actorError.ErrRemindersParamsInvalid:
		msg = "ErrRemindersParamsInvalid"
		break
	case actorError.ErrActorMethodNoFound:
		msg = "ErrActorMethodNoFound"
		break
	case actorError.ErrActorInvokeFailed:
		msg = "ErrActorInvokeFailed"
		break
	case actorError.ErrReminderFuncUndefined:
		msg = "ErrReminderFuncUndefined"
		break
	case actorError.ErrActorMethodSerializeFailed:
		msg = "ErrActorMethodSerializeFailed"
		break
	case actorError.ErrActorSerializeNoFound:
		msg = "ErrActorSerializeNoFound"
		break
	case actorError.ErrActorIDNotFound:
		msg = "ErrActorIDNotFound"
		break
	case actorError.ErrActorFactoryNotSet:
		msg = "ErrActorFactoryNotSet"
		break
	case actorError.ErrTimerParamsInvalid:
		msg = "ErrTimerParamsInvalid"
		break
	case actorError.ErrSaveStateFailed:
		msg = "ErrSaveStateFailed"
		break
	case actorError.ErrActorServerInvalid:
		msg = "ErrActorServerInvalid"
		break
	default:
		msg = "unknown"
		break
	}
	if len(msg) == 0 {
		return nil
	}
	return errors.New(msg)
}

func newActorFieldError(funLog, actorType, actorId, methodName string, err error) logs.Fields {
	fields := logs.Fields{
		"funLog":     funLog,
		"actorType":  actorType,
		"actorId":    actorId,
		"methodName": methodName,
		"error":      err.Error(),
	}
	return fields
}
