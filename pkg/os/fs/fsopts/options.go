package fsopts

import (
	"strings"

	fileutils2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/fileutils"
	"github.com/spf13/afero"
)

type Options struct {
	RootPath string // 根目录
	WorkPath string // 当前工作目录
	Fs       afero.Fs
}

func NewOptionsWidthFileName(fileName string, rootPath string) *Options {
	workPath := GetWorkPath(fileName, rootPath)
	return &Options{
		RootPath: rootPath,
		WorkPath: workPath,
	}
}

func NewOptions(o ...*Options) *Options {
	opts := &Options{}
	for _, o := range o {
		if o.RootPath != "" {
			opts.RootPath = o.RootPath
		}
		if o.WorkPath != "" {
			opts.WorkPath = o.WorkPath
		}
		if o.Fs != nil {
			opts.Fs = o.Fs
		}
	}
	return opts
}

// ParseFileName
//
//	@Description: 解析文件名称
//	@param filename
//	@param defaultFsName 默认的fs名称
//	@param init
//	@return error
func ParseFileName(filename, defaultFsName string, init func(fsName, fileName string)) error {
	var fsName, fileName string
	i := strings.Index(filename, ":")
	if i == -1 {
		fsName = defaultFsName
		fileName = filename
	} else {
		fsName = filename[0:i]
		fileName = filename[i+1:]
	}

	if init != nil {
		init(fsName, fileName)
	}
	return nil
}

func GetAbsPath(filename string, opts ...*Options) string {
	opt := newOpts(opts...)
	return fileutils2.AbsPath(filename, opt)
}

func newOpts(opts ...*Options) *fileutils2.ReadOptions {
	opt := &fileutils2.ReadOptions{}
	for _, o := range opts {
		opt.RootPath = o.RootPath
		opt.WorkPath = o.WorkPath
	}
	return opt
}

// GetWorkPath 根据文件名获取文件路径
func GetWorkPath(fileName string, rootPath string) string {
	i := strings.LastIndex(fileName, "/")
	if i == -1 {
		return "/"
	}
	return fileName[:i]
}
