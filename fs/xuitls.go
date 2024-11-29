package fs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/spf13/afero"
	"io/fs"
)

// ReadFile
//
//	@Description: 读取文件内容
//	@param fs afero.Fs
//	@param pwdPath string 当前路径
//	@param filename string 要读取文件名称
//	@return []byte 文件内容
//	@return error
func ReadFile(afs afero.Fs, filename string, opts ...*fsopts.Options) ([]byte, error) {
	fileName := fsopts.GetAbsPath(filename, opts...)
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
func WriteFile(afs afero.Fs, filename string, bytes []byte, fileMode fs.FileMode, opts ...*fsopts.Options) error {
	var err error
	fileName := fsopts.GetAbsPath(filename, opts...)
	if encode, ok := afs.(Decode); ok {
		bytes, err = encode.Encode(bytes)
		if err != nil {
			return err
		}
	}
	err = afero.WriteFile(afs, fileName, bytes, fileMode)
	return err
}
