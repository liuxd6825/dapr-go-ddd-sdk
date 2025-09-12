package restapi

import (
	errors2 "errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types/times"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils/webjson"
	"strings"
)

var notFoundError = errors2.New("Not Found")

func NotFoundError() error {
	return notFoundError
}

type WebError struct {
	Error      any         `json:"error"`
	LogId      string      `json:"logId"`
	Time       *times.Time `json:"time"`
	StatusCode int         `json:"statusCode"`
}

type VerifyError struct {
	Error      any         `json:"error"`
	LogId      string      `json:"logId"`
	Time       *times.Time `json:"time"`
	StatusCode int         `json:"statusCode"`
}

func NewWebError(ictx iris.Context, logId string, err error, statusCode int) *WebError {
	ictx.StatusCode(statusCode)
	return &WebError{
		Error:      err.Error(),
		LogId:      logId,
		Time:       times.PNow(),
		StatusCode: statusCode,
	}
}

func SetError(ictx iris.Context, err error) {
	if err != nil {
		req := ictx.Request()
		logId := logs.GetLogId(ictx)
		logs.Error(ictx, logs.Fields{
			"logId":      logId,
			"method":     req.Method,
			"uri":        req.RequestURI,
			"error":      err.Error(),
			"remoteAddr": req.RemoteAddr,
		})
		var data any
		if vErr, ok := err.(*errors.VerifyError); ok {
			data = NewWebError(ictx, logId, vErr, iris.StatusBadRequest)
		} else if isNotFoundError(err) {
			data = NewWebError(ictx, logId, err, iris.StatusNotFound)
		} else {
			data = NewWebError(ictx, logId, err, iris.StatusInternalServerError)
		}
		_ = ictx.JSON(data)
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
func SetOKJsonData(ictx iris.Context, data any) error {
	ictx.StatusCode(iris.StatusOK)
	jsonStr, err := webjson.Marshal(data)
	if err != nil {
		return err
	}
	_, err = ictx.WriteString(jsonStr)
	return err
}

func SetData(ictx iris.Context, data any) error {
	return ictx.JSON(data)
}
