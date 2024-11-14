package fs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/spf13/afero"
	"io/fs"
	"strings"
)

type Options struct {
	BasePath string
}

func NewOptions(o ...*Options) *Options {
	opts := &Options{}
	for _, o := range o {
		opts.BasePath = o.BasePath
	}
	return opts
}

// ReadFile
//
//	@Description: 读取文件内容
//	@param fs afero.Fs
//	@param pwdPath string 当前路径
//	@param filename string 要读取文件名称
//	@return []byte 文件内容
//	@return error
func ReadFile(afs afero.Fs, filename string, opts ...*Options) ([]byte, error) {
	o := NewOptions(opts...)
	fileName := fileutils.AbsPath(filename, o.BasePath)
	context, err := afero.ReadFile(afs, fileName)
	if err != nil {
		return nil, err
	}
	if decode, ok := afs.(Decode); ok {
		context, err = decode.Decode(context)
	}
	return context, err
}

// WriteFile
//
//	@Description: 读取文件内容
//	@param fs afero.Fs
//	@param pwdPath string 当前路径
//	@param filename string 要读取文件名称
//	@return []byte 文件内容
//	@return error
func WriteFile(afs afero.Fs, filename string, bytes []byte, fileMode fs.FileMode, opts ...*Options) error {
	var err error
	o := NewOptions(opts...)
	fileName := fileutils.AbsPath(filename, o.BasePath)
	if encode, ok := afs.(Decode); ok {
		bytes, err = encode.Encode(bytes)
		if err != nil {
			return err
		}
	}
	err = afero.WriteFile(afs, fileName, bytes, fileMode)
	return err
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
