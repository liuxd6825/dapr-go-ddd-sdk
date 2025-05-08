package transform

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"strings"
)

type Tsc struct {
	babel *Babel
}

func NewTsc() *Tsc {
	return &Tsc{}
}

type TscTransformTarget = api.Target

const (
	DefaultTarget TscTransformTarget = iota
	ESNext
	ES5
	ES2015
	ES2016
	ES2017
	ES2018
	ES2019
	ES2020
	ES2021
	ES2022
	ES2023
	ES2024
)

func (t *Tsc) TransformEs5(tsCode string, fileName string) ([]byte, error) {
	return t.Transform(tsCode, fileName, ES5)
}

// Transform
//
//	@Description: 使用 ESBuild 将 TypeScript 转为 ES6
//	@receiver t
//	@param tsCode
//	@param target
//	@return []byte
//	@return error
func (t *Tsc) Transform(tsCode string, fileName string, target TscTransformTarget) ([]byte, error) {
	jsTarget := target
	if jsTarget == ES5 {
		jsTarget = ES2015
	}
	esbuildResult := api.Transform(tsCode, api.TransformOptions{
		Loader: api.LoaderTS, // 指定 TypeScript Loader
		Target: jsTarget,     // 转换 api.ES20115 = ES6
	})

	// 检查 ESBuild 编译是否有错误
	if len(esbuildResult.Errors) > 0 {
		var msgList []string
		for _, err := range esbuildResult.Errors {
			msgList = append(msgList, err.Text)
		}
		errMsg := fmt.Sprintf("error transforming TypeScript: %s in %s", strings.Join(msgList, "\n"), fileName)
		return nil, errors.New(errMsg)
	}

	es6Code := string(esbuildResult.Code)
	if target == ES5 {
		if t.babel == nil {
			var err error
			if t.babel, err = NewBabel(); err != nil {
				return nil, err
			}
		}
		date, err := t.babel.Transform(es6Code, fileName, false, nil)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error transforming ES6: %s in %s", err.Error(), fileName))
		}
		return date, nil
	}
	return []byte(es6Code), nil
}
