package restapp

import (
	"context"
	"github.com/kataras/iris/v12"
	iris_context "github.com/kataras/iris/v12/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/applog"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"net/http"
	"time"
)

const (
	ContentTypeApplicationJson = iris_context.ContentJSONHeaderValue
	ContentTypeTextPlain       = iris_context.ContentTextHeaderValue
)

type JsonTimeSerializer struct {
}

func (j *JsonTimeSerializer) Serialize(v interface{}) ([]byte, error) {
	t, ok := v.(*time.Time)
	if !ok {
		return nil, errors.New("invalid type")
	}
	return []byte(t.Format("2006-01-02 15:04:05")), nil
}

func myHandler(ctx iris.Context) {
	type MyData struct {
		Time time.Time `json:"time" serializer:"customTime"`
		// other fields
	}
	data := MyData{
		Time: time.Now(),
		// other fields
	}
	ctx.JSON(data)
}

type Command interface {
	GetCommandId() string
	GetTenantId() string
}

type CmdFunc func(ctx context.Context) error
type QueryFunc func(ctx context.Context) (interface{}, bool, error)

// CmdAndQueryOptions
// @Description: 命令执行参数
type CmdAndQueryOptions struct {
	WaitSecond int // 超时时间，单位秒
}

type CmdAndQueryOption func(options *CmdAndQueryOptions)

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

func SetErrorVerifyError(ctx iris.Context, err *errors.VerifyError) {
	ctx.SetErr(err)
	ctx.StatusCode(http.StatusInternalServerError)
	ctx.ContentType(ContentTypeTextPlain)
}

func SetError(ctx context.Context, err error) {
	logs.Error(ctx, err)

	ictx, ok := ctx.(iris.Context)
	if !ok {
		ictx = GetIrisContext(ctx)
	}
	if ictx == nil {
		return
	}
	switch err.(type) {
	case *errors.NullError:
		_ = SetErrorNotFond(ictx)
		break
	case *errors.VerifyError:
		verr, _ := err.(*errors.VerifyError)
		SetErrorVerifyError(ictx, verr)
		break
	default:
		SetErrorInternalServerError(ictx, err)
		break
	}
}

// DoCmd
// @Description: 执行命令
// @param ctx  上下文
// @param cmd  命令
// @param fun  执行方法
// @return err 错误
func DoCmd(ictx iris.Context, fun CmdFunc) (err error) {
	ctx := NewContext(ictx)

	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			SetError(ctx, err)
			logs.Error(ctx, err)
		}
	}()

	err = fun(ctx)
	if err != nil && !errors.IsErrorAggregateExists(err) {
		SetError(ctx, err)
		return err
	}
	return err
}

func Do(ictx iris.Context, fun func() error) (err error) {
	ctx := NewContext(ictx)
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			logs.Error(ctx, err)
		}
	}()
	if fun != nil {
		if err = fun(); err != nil {
			SetError(ctx, err)
		}
	}
	return nil
}

func DoDto[T any](ictx iris.Context, fun func(ctx context.Context) (T, error)) (dto T, err error) {
	ctx := NewContext(ictx)
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			SetError(ictx, err)
		}
	}()
	if fun != nil {
		if dto, err = fun(ctx); err != nil {
			SetError(ictx, err)
		}
	}
	return dto, err
}

// DoQueryOne
// @Description: 单条数据查询，当无数据时返回错误。
// @param ctx 上下文
// @param fun 执行方法
// @return data 返回数据
// @return isFound 是否有数据
// @return err 错误
func DoQueryOne(ictx iris.Context, fun QueryFunc) (data interface{}, isFound bool, err error) {
	ctx := NewContext(ictx)
	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			SetError(ctx, err)
		}
	}()

	data, isFound, err = fun(ctx)
	if err != nil {
		SetError(ctx, err)
		return nil, false, err
	}
	if data == nil || !isFound {
		return nil, false, SetErrorNotFond(ictx)
	}
	err = SetJson(ictx, data)
	if err != nil {
		SetError(ctx, err)
		return nil, false, err
	}
	return data, isFound, err
}

// DoQuery
// @Description: 多条数据查询
// @param ctx 上下文
// @param fun 执行方法
// @return data 返回数据
// @return isFound 是否有数据
// @return err 错误
func DoQuery(ictx iris.Context, fun QueryFunc) (data any, isFound bool, err error) {
	ctx := NewContext(ictx)

	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if ok {
				SetError(ctx, err)
			}
		}
	}()

	defer func() {
		if err = errors.GetRecoverError(err, recover()); err != nil {
			SetError(ctx, err)
		}
	}()

	_ = logs.DebugStart(ctx, func() error {
		data, isFound, err = fun(ctx)
		return err
	}, "uri=%v; params=%v;", ictx.FullRequestURI(), ictx.Params())

	if err != nil {
		SetError(ctx, err)
		return data, isFound, err
	}

	err = SetJson(ictx, data)
	if err != nil {
		SetError(ctx, err)
		return nil, false, err
	}
	return data, isFound, err
}

func CmdAndQueryOptionWaitSecond(waitSecond int) CmdAndQueryOption {
	return func(options *CmdAndQueryOptions) {
		options.WaitSecond = waitSecond
	}
}

// DoCmdAndQueryOne 执行命令并返回查询一个数据
//
//	DoCmdAndQueryOne
//	@Description:  执行命令并返回查询一个数据
//	@param ctx 上下文
//	@param queryAppId  查询AppId
//	@param cmd  命令
//	@param cmdFun  命令执行方法
//	@param queryFun 查询执行方法
//	@param opts 参数
//	@return interface{} 返回值
//	@return bool 是否找到数据
//	@return error 错误
func DoCmdAndQueryOne(ictx iris.Context, queryAppId string, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	return doCmdAndQuery(ictx, queryAppId, true, cmd, cmdFun, queryFun, opts...)
}

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
func DoCmdAndQueryList(ictx iris.Context, queryAppId string, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	return doCmdAndQuery(ictx, queryAppId, false, cmd, cmdFun, queryFun, opts...)
}

func doCmdAndQuery(ictx iris.Context, queryAppId string, isGetOne bool, cmd Command, cmdFun CmdFunc, queryFun QueryFunc, opts ...CmdAndQueryOption) (interface{}, bool, error) {
	ctx := NewContext(ictx)
	options := &CmdAndQueryOptions{WaitSecond: 5}
	for _, o := range opts {
		o(options)
	}

	err := DoCmd(ictx, cmdFun)
	isExists := errors.IsErrorAggregateExists(err)
	if err != nil && !isExists {
		SetError(ctx, err)
		return nil, false, err
	}
	err = nil
	//isTimeout := true
	// 循环检查EventLog日志是否存在
	for i := 0; i < options.WaitSecond; i++ {
		time.Sleep(time.Duration(1) * time.Second)
		logs, err := applog.GetEventLogByAppIdAndCommandId(cmd.GetTenantId(), queryAppId, cmd.GetCommandId())
		if err != nil {
			SetError(ctx, err)
			return nil, false, err
		}

		// 循环检查EventLog日志是否存在
		if len(logs) > 0 {
			// isTimeout = false
			break
		}
	}

	/*	if isTimeout {
		msg := fmt.Sprintf("applog.GetEventLogByAppIdAndCommandId() error: queryAppId=%s, commandId=%s, tenantId=%s  execution timeout", queryAppId, cmd.GetCommandId(), cmd.GetTenantId())
		SetError(ctx, errors.New(msg))
		return nil, false, err
	}*/

	var data interface{}
	var isFound bool
	if isGetOne {
		data, isFound, err = DoQueryOne(ictx, queryFun)
	} else {
		data, isFound, err = DoQuery(ictx, queryFun)
	}
	if err != nil {
		SetError(ctx, err)
	}
	return data, isFound, err
}

func SetJson(ictx iris.Context, data interface{}) error {
	bs, err := WriteJSON(data)
	if err != nil {
		SetError(ictx, err)
		return err
	}

	if _, err = ictx.Write(bs); err != nil {
		SetError(ictx, err)
		return err
	}

	ictx.ContentType(iris_context.ContentJSONHeaderValue)
	return nil
}

func ReadJSON(ictx iris.Context, obj any) error {
	data, err := ictx.GetBody()
	if err != nil {
		return err
	}
	return jsonutils.CustomJson.Unmarshal(data, obj)
}

func WriteJSON(data any) ([]byte, error) {
	return jsonutils.CustomJson.Marshal(data)
}
