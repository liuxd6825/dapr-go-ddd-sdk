package template

import (
	"context"
	"errors"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/spf13/afero"
	"io"
)

// Render
//
//	@Description: 渲染HTML模板，并执行服务端代码
//	@param ctx 上下文
//	@param writer 输出流
//	@param loader 模板文件加载器
//	@param filepath 模板文件路径
//	@param data 模板数据
//	@return error 错误信息
func Render(ctx context.Context, writer io.Writer, fs afero.Fs, filepath string, options ...func(ctx pongo2.Context)) error {
	//从gitea中加载文件
	fileBytes, err := readFile(fs, "", filepath)
	if err != nil {
		return err
	}

	// 解析文件内容
	parse, err := NewParse(fs, filepath, fileBytes).Parse()
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
	for _, opt := range options {
		if opt != nil {
			opt(tplCtx)
		}
	}
	if err = tplExample.ExecuteWriter(tplCtx, writer); err != nil {
		return errors.New(fmt.Sprintf("模板引擎执行错误,%s。", err.Error()))
	}
	return nil
}
