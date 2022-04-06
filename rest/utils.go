package rest

import (
	"encoding/json"
	"errors"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"net/http"
	"time"
)

const (
	ContentTypeApplicationJson = "application/json"
	ContentTypeTextPlain       = "text/plain"
)

type Command interface {
	GetCommandId() string
	GetTenantId() string
}

// type GetOneFunc = func(ctx iris.Context) (interface{}, bool, error)
// type GetFunc = func(ctx iris.Context) (interface{}, error)
type CmdFunc func() error
type QueryFunc func() (interface{}, bool, error)

func SetErrorNotFond(ctx iris.Context) error {
	ctx.SetErr(iris.ErrNotFound)
	ctx.StatusCode(http.StatusNotFound)
	ctx.ContentType(ContentTypeTextPlain)
	return iris.ErrNotFound
}

func SetErrorInternalServerError(ctx iris.Context, err error) {
	ctx.SetErr(err)
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.ContentType(ContentTypeTextPlain)
}

func SetErrorVerifyError(ctx iris.Context, err error) {
	bytes, _ := json.Marshal(err)
	_, _ = ctx.Write(bytes)
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.ContentType(ContentTypeTextPlain)
}

func SetError(ctx iris.Context, err error) {
	switch err.(type) {
	case *ddd_errors.NullError:
		_ = SetErrorNotFond(ctx)
		break
	case *ddd_errors.VerifyError:
		SetErrorVerifyError(ctx, err)
		break
	default:
		SetErrorInternalServerError(ctx, err)
		break
	}
}

// DoCmd 执行命令
func DoCmd(ctx iris.Context, cmd Command, fun CmdFunc) error {
	if err := ctx.ReadBody(cmd); err != nil {
		return err
	}

	err := fun()
	if err != nil && !ddd_errors.IsAggregateExistsError(err) {
		SetError(ctx, err)
		return err
	}
	return nil
}

func DoQueryOne(ctx iris.Context, fun QueryFunc) (interface{}, bool, error) {
	data, isFound, err := fun()
	if err != nil {
		SetError(ctx, err)
		return nil, isFound, err
	}
	if data == nil || !isFound {
		return nil, isFound, SetErrorNotFond(ctx)
	}
	_, err = ctx.JSON(data)
	if err != nil {
		return nil, false, err
	}
	return data, isFound, nil
}

func DoQueryList(ctx iris.Context, fun QueryFunc) (interface{}, bool, error) {
	data, isFound, err := fun()
	if err != nil {
		SetError(ctx, err)
		return nil, isFound, err
	}
	_, err = ctx.JSON(data)
	if err != nil {
		return nil, false, err
	}
	return data, isFound, nil
}

type CmdAndQueryOptions struct {
	WaitSecond int
}

type CmdAndQueryOption func(options *CmdAndQueryOptions)

func CmdAndQueryOptionWaitSecond(waitSecond int) CmdAndQueryOption {
	return func(options *CmdAndQueryOptions) {
		options.WaitSecond = waitSecond
	}
}

// DoCmdAndQueryOne 执行命令并返回查询一个数据
func DoCmdAndQueryOne(ctx iris.Context, subAppId string, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	return doCmdAndQuery(ctx, subAppId, true, cmd, cmdFun, queryFun, opts...)
}

// DoCmdAndQueryList 执行命令并返回查询列表
func DoCmdAndQueryList(ctx iris.Context, subAppId string, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	return doCmdAndQuery(ctx, subAppId, false, cmd, cmdFun, queryFun, opts...)
}

func doCmdAndQuery(ctx iris.Context, subAppId string, isGetOne bool, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {

	options := &CmdAndQueryOptions{WaitSecond: 5}
	for _, o := range opts {
		o(options)
	}

	err := DoCmd(ctx, cmd, cmdFun)
	isExists := ddd_errors.IsAggregateExistsError(err)
	if err != nil && !isExists {
		return nil, false, err
	}
	err = nil
	isTimeout := true
	// 循环检查EventLog日志是否存在
	for i := 0; i < options.WaitSecond; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		logs, err := applog.GetEventLogByAppIdAndCommandId(cmd.GetTenantId(), subAppId, cmd.GetCommandId())
		if err != nil {
			return nil, false, err
		}

		// 循环检查EventLog日志是否存在
		if len(*logs) > 0 {
			isTimeout = false
			break
		}
	}

	if isTimeout {
		return nil, false, errors.New("query execution timeout")
	}

	var data interface{}
	var isFound bool
	if isGetOne {
		data, isFound, err = DoQueryOne(ctx, queryFun)
	} else {
		data, isFound, err = DoQueryList(ctx, queryFun)
	}
	return data, isFound, err
}

func SetData(ctx iris.Context, data interface{}) {
	_, err := ctx.JSON(data)
	if err != nil {
		SetError(ctx, err)
		return
	}
}
