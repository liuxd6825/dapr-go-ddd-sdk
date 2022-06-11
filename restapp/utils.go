package restapp

import (
	"context"
	"errors"
	"fmt"
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

type CmdFunc func(ctx context.Context) error
type QueryFunc func(ctx context.Context) (interface{}, bool, error)

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

func SetErrorVerifyError(ctx iris.Context, err *ddd_errors.VerifyError) {
	ctx.SetErr(err)
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.ContentType(ContentTypeTextPlain)
}

func SetError(ctx iris.Context, err error) {
	switch err.(type) {
	case *ddd_errors.NullError:
		_ = SetErrorNotFond(ctx)
		break
	case *ddd_errors.VerifyError:
		verr, _ := err.(*ddd_errors.VerifyError)
		SetErrorVerifyError(ctx, verr)
		break
	default:
		SetErrorInternalServerError(ctx, err)
		break
	}
}

//
// DoCmd
// @Description: 执行命令
// @param ctx  上下文
// @param cmd  命令
// @param fun  执行方法
// @return err 错误
//
func DoCmd(ctx iris.Context, fun CmdFunc) (err error) {
	defer func() {
		if e := ddd_errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()

	restCtx := NewContext(ctx)
	err = fun(restCtx)
	if err != nil && !ddd_errors.IsErrorAggregateExists(err) {
		SetError(ctx, err)
		return err
	}
	return err
}

//
// DoQueryOne
// @Description: 单条数据查询，当无数据时返回错误。
// @param ctx 上下文
// @param fun 执行方法
// @return data 返回数据
// @return isFound 是否有数据
// @return err 错误
//
func DoQueryOne(ctx iris.Context, fun QueryFunc) (data interface{}, isFound bool, err error) {
	defer func() {
		if e := ddd_errors.GetRecoverError(recover()); e != nil {
			err = e
		}
	}()
	restCtx := NewContext(ctx)
	data, isFound, err = fun(restCtx)
	if err != nil {
		SetError(ctx, err)
		return nil, false, err
	}
	if data == nil || !isFound {
		return nil, false, SetErrorNotFond(ctx)
	}
	_, err = ctx.JSON(data)
	if err != nil {
		return nil, false, err
	}
	return data, isFound, err
}

//
// DoQuery
// @Description: 多条数据查询
// @param ctx 上下文
// @param fun 执行方法
// @return data 返回数据
// @return isFound 是否有数据
// @return err 错误
//
func DoQuery(ctx iris.Context, fun QueryFunc) (data interface{}, isFound bool, err error) {
	defer func() {
		if e := ddd_errors.GetRecoverError(recover()); e != nil {
			err = e
			SetError(ctx, err)
		}
	}()
	restCtx := NewContext(ctx)
	data, isFound, err = fun(restCtx)
	if err != nil {
		SetError(ctx, err)
		return data, isFound, err
	}

	_, err = ctx.JSON(data)
	if err != nil {
		return nil, false, err
	}
	return data, isFound, err
}

//
// CmdAndQueryOptions
// @Description: 命令执行参数
//
type CmdAndQueryOptions struct {
	WaitSecond int // 超时时间，单位秒
}

type CmdAndQueryOption func(options *CmdAndQueryOptions)

func CmdAndQueryOptionWaitSecond(waitSecond int) CmdAndQueryOption {
	return func(options *CmdAndQueryOptions) {
		options.WaitSecond = waitSecond
	}
}

// DoCmdAndQueryOne 执行命令并返回查询一个数据
//
//  DoCmdAndQueryOne
//  @Description:  执行命令并返回查询一个数据
//  @param ctx 上下文
//  @param queryAppId  查询AppId
//  @param cmd  命令
//  @param cmdFun  命令执行方法
//  @param queryFun 查询执行方法
//  @param opts 参数
//  @return interface{} 返回值
//  @return bool 是否找到数据
//  @return error 错误
//
func DoCmdAndQueryOne(ctx iris.Context, queryAppId string, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	return doCmdAndQuery(ctx, queryAppId, true, cmd, cmdFun, queryFun, opts...)
}

//
// DoCmdAndQueryList
// @Description:  执行命令并返回查询列表
// @param ctx 上下文
// @param queryAppId  查询AppId
// @param cmd  命令
// @param cmdFun  命令执行方法
// @param queryFun 查询执行方法
// @param opts 参数
// @return interface{} 返回值
// @return bool 是否找到数据
// @return error 错误
//
func DoCmdAndQueryList(ctx iris.Context, queryAppId string, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	return doCmdAndQuery(ctx, queryAppId, false, cmd, cmdFun, queryFun, opts...)
}

func doCmdAndQuery(ctx iris.Context, queryAppId string, isGetOne bool, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	options := &CmdAndQueryOptions{WaitSecond: 5}
	for _, o := range opts {
		o(options)
	}

	err := DoCmd(ctx, cmdFun)
	isExists := ddd_errors.IsErrorAggregateExists(err)
	if err != nil && !isExists {
		SetError(ctx, err)
		return nil, false, err
	}
	err = nil
	isTimeout := true
	// 循环检查EventLog日志是否存在
	for i := 0; i < options.WaitSecond; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		logs, err := applog.GetEventLogByAppIdAndCommandId(cmd.GetTenantId(), queryAppId, cmd.GetCommandId())
		if err != nil {
			SetError(ctx, err)
			return nil, false, err
		}

		// 循环检查EventLog日志是否存在
		if len(*logs) > 0 {
			isTimeout = false
			break
		}
	}

	if isTimeout {
		msg := fmt.Sprintf("applog.GetEventLogByAppIdAndCommandId() error: queryAppId=%s, commandId=%s, tenantId=%s  execution timeout", queryAppId, cmd.GetCommandId(), cmd.GetTenantId())
		SetError(ctx, errors.New(msg))
		return nil, false, err
	}

	var data interface{}
	var isFound bool
	if isGetOne {
		data, isFound, err = DoQueryOne(ctx, queryFun)
	} else {
		data, isFound, err = DoQuery(ctx, queryFun)
	}
	if err != nil {
		SetError(ctx, err)
	}
	return data, isFound, err
}

func SetRestData(ctx iris.Context, data interface{}) {
	_, err := ctx.JSON(data)
	if err != nil {
		SetError(ctx, err)
		return
	}
}
