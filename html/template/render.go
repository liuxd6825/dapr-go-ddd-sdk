package template

import (
	"context"
	"errors"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader"
	"io"
)

func Render(ctx context.Context, writer io.Writer, loader loader.Loader, filepath string, data ...func(ctx pongo2.Context)) error {
	//从gitea中加载文件
	fileBytes, err := loader.GetFile(filepath)
	if err != nil {
		return err
	}

	// 解析文件内容
	parse, err := NewParse(loader, filepath, fileBytes).Parse()
	if err != nil {
		return errors.New(fmt.Sprintf("解释 %s 时出错,%s。", filepath, err.Error()))
	}

	// 执行服务端代码
	var vdata any
	if parse.Script.Len() > 0 {
		runTime := NewScriptRuntime()
		vdata, err = runTime.Run(nil, parse.Script.String())
		if err != nil {
			return errors.New(fmt.Sprintf("执行服务脚本时出错,%s。", err.Error()))
		}
	}

	var tplExample = pongo2.Must(pongo2.FromString(parse.HTML))
	tplCtx := pongo2.Context{"vdata": vdata, "schema": parse.Schema, "uischeam": parse.UiSchema}
	for _, d := range data {
		if d != nil {
			d(tplCtx)
		}
	}
	if err = tplExample.ExecuteWriter(tplCtx, writer); err != nil {
		return errors.New(fmt.Sprintf("模板引擎执行错误,%s。", err.Error()))
	}
	return nil
}
