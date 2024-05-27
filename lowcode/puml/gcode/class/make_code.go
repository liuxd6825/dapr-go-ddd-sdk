package class

import (
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/puml/gcode/class/utils"
	"os"
	"path/filepath"
	"text/template"
)

type MackCodeOption struct {
	SetFileName   func(name string) string
	OverwriteFile bool
}
type MakeCodeOptions = func(opt *MackCodeOption)

func MakeCode(ctx context.Context, plantUmlFile string, templatePath string, savePath string, lanType string, opts ...MakeCodeOptions) error {
	opt := &MackCodeOption{OverwriteFile: true}
	for _, optsItem := range opts {
		optsItem(opt)
	}

	plantfile, err := filepath.Abs(plantUmlFile)
	if err != nil {
		return err
	}

	block, err := ParseFile(ctx, plantfile)
	if err != nil {
		return errors.New(plantfile + err.Error())
	}

	// 解析指定文件生成模板对象
	classTemp, err := newTemplate(templatePath + "/class.tpl")
	if err != nil {
		return err
	}

	// 解析指定文件生成模板对象
	enumTemp, err := newTemplate(templatePath + "/enum.tpl")
	if err != nil {
		return err
	}

	// 解析指定文件生成模板对象
	interTemp, err := newTemplate(templatePath + "/interface.tpl")
	if err != nil {
		return err
	}

	for _, ns := range block.Root.Namespaces {
		var temp *template.Template
		for _, node := range ns.Elements {
			temp = nil
			switch node.GetTypeName() {
			case "class":
				temp = classTemp
			case "enum":
				temp = enumTemp
			case "interface":
				temp = interTemp
			}
			path, err := utils.Mkdir(savePath, node.GetNamespaceName())
			if err != nil {
				return err
			}
			fileName := node.GetName()
			if opt.SetFileName != nil {
				fileName = opt.SetFileName(fileName)
			}
			codeFileName := fmt.Sprintf("%s/%s.%s", path, fileName, lanType)
			if !opt.OverwriteFile {
				if ok, err := utils.PathExists(codeFileName); err != nil {
					return err
				} else if ok {
					continue
				}
			}
			w, err := os.OpenFile(codeFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return err
			}
			err = temp.Execute(w, node)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func newTemplate(templateFile string) (*template.Template, error) {
	// 解析指定文件生成模板对象
	temp, err := template.ParseFiles(templateFile)
	if err != nil {
		return nil, err
	}
	return temp, err
}
