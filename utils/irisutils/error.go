package irisutils

import "github.com/kataras/iris/v12"

func SetError(ctx iris.Context, err error) {
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.SetErr(err)
	}
}
