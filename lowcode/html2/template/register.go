package template

import (
	"errors"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/spf13/afero"
)

type Handler struct {
	fs afero.Fs
}
type RenderHandle func(ictx iris.Context, filepath string)

// NewHandler
//
//	@Description: 添加路由
//	@param app *iris.Application
//	@param url  relativePath + "/{repoName}/{filepath:path}"
//	@return error 错误信息
func NewHandler(fs afero.Fs) (RenderHandle, error) {
	if fs == nil {
		return nil, errors.New("loaders is nil")
	}
	h := &Handler{fs: fs}
	return h.handle, nil
}

func (h *Handler) handle(ictx iris.Context, filepath string) {
	do(ictx, func(ctx iris.Context) error {
		return Render(ictx, ictx.ResponseWriter(), h.fs, filepath, func(data pongo2.Context) {
			data["params"] = ctx.Params()
		})
	})
}

// do
//
//	@Description:
//	@param ctx
//	@param fun
func do(ctx iris.Context, fun func(ctx iris.Context) error) {
	var err error
	defer func() {
		_ = CatchError(ctx, err, recover())
	}()
	err = fun(ctx)
	if err != nil {
		SetError(ctx, err)
	}
}

// CatchError
//
//	@Description:
//	@param ctx
//	@param e
//	@param recover
//	@return error
func CatchError(ctx iris.Context, e error, recover any) error {
	var err error
	if e != nil {
		err = e
	} else if recover != nil {
		if ve := recover.(error); ve != nil {
			err = ve
		} else {
			err = fmt.Errorf("unknown error %v", recover)
		}
	}
	if err != nil && ctx != nil {
		SetError(ctx, err)
	}
	return nil
}

func SetError(ctx iris.Context, err error) {
	if err != nil && ctx != nil {
		ctx.SetErr(err)
		ctx.StatusCode(httptest.StatusInternalServerError)
	}
}
