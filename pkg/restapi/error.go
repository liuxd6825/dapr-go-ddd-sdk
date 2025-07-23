package restapi

import (
	errors2 "errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"strings"
)

var notFoundError = errors2.New("Not Found")

func NotFoundError() error {
	return notFoundError
}

type InternalServerError struct {
	Error      any        `json:"error"`
	LogId      string     `json:"logId"`
	Time       times.Time `json:"time"`
	StatusCode int        `json:"statusCode"`
}

type VerifyError struct {
	Error      any        `json:"error"`
	LogId      string     `json:"logId"`
	Time       times.Time `json:"time"`
	StatusCode int        `json:"statusCode"`
}

func NewVerifyError(logId string, err *errors.VerifyError) *VerifyError {
	return &VerifyError{
		Error:      err.GetFieldErrors(),
		LogId:      logId,
		Time:       times.Now(),
		StatusCode: iris.StatusBadRequest,
	}
}
func NewInternalServerError(logId string, err error) *InternalServerError {
	return &InternalServerError{
		Error:      err.Error(),
		LogId:      logId,
		Time:       times.Now(),
		StatusCode: iris.StatusInternalServerError,
	}
}

func SetError(ctx iris.Context, err error) {
	if err != nil {
		req := ctx.Request()
		logId := logs.GetLogId(ctx)
		logs.Error(ctx, logs.Fields{"logId": logId, "method": req.Method, "uri": req.RequestURI, "error": err.Error()})
		if vErr, ok := err.(*errors.VerifyError); ok {
			ctx.StatusCode(iris.StatusBadRequest)
			_ = ctx.JSON(NewVerifyError(logId, vErr))
		} else if isNotFoundError(err) {
			ctx.StatusCode(iris.StatusNotFound)
			_ = ctx.JSON(err)
		} else {
			ctx.StatusCode(iris.StatusInternalServerError)
			_ = ctx.JSON(NewInternalServerError(logId, err))
		}

	}
}

func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if err == notFoundError {
		return true
	}
	if strings.Contains(strings.ToLower(err.Error()), "not found") {
		return true
	}
	return false
}
func SetOKData(ictx iris.Context, data any) error {
	ictx.StatusCode(iris.StatusOK)
	return ictx.JSON(data)
}

func SetData(ictx iris.Context, data any) error {
	return ictx.JSON(data)
}
