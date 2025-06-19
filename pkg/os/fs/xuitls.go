package fs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/spf13/afero"
	"io/fs"
	"os"
)

// IsDir
//
//	@Description:
//	@param afs
//	@param filename
//	@param opts
//	@return bool
//	@return error
func IsDir(afs afero.Fs, filename string, opts ...*fsopts.Options) (bool, error) {
	fileName := fsopts.GetAbsPath(filename, opts...)
	return afero.IsDir(afs, fileName)
}

// Exists 文件是否存在
func Exists(afs afero.Fs, filename string, opts ...*fsopts.Options) (bool, error) {
	fileName := fsopts.GetAbsPath(filename, opts...)
	return afero.Exists(afs, fileName)
}

// IsEmpty
//
//	@Description:
//	@param afs
//	@param filename
//	@param opts
//	@return bool
//	@return error
func IsEmpty(afs afero.Fs, filename string, opts ...*fsopts.Options) (bool, error) {
	fileName := fsopts.GetAbsPath(filename, opts...)
	return afero.IsEmpty(afs, fileName)
}

// Create
//
//	@Description:
//	@param afs
//	@param filename
//	@param opts
//	@return bool
//	@return error
func Create(afs afero.Fs, filename string, opts ...*fsopts.Options) (afero.File, error) {
	fileName := fsopts.GetAbsPath(filename, opts...)
	return afs.Create(fileName)
}

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
	data, err := afero.ReadFile(afs, fileName)
	if err != nil {
		return nil, err
	}
	if decode, ok := afs.(Decode); ok {
		data, err = decode.Decode(data)
	}
	return data, err
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

func Open(afs afero.Fs, filename string, flag int, perm os.FileMode, opts ...*fsopts.Options) (afero.File, error) {
	fileName := fsopts.GetAbsPath(filename, opts...)
	file, err := afs.OpenFile(fileName, flag, perm)
	return file, err
}

func WriteAt(afs afero.Fs, fsFile afero.File, off int64, bytes []byte, opts ...*fsopts.Options) (int, error) {
	var err error
	if encode, ok := afs.(Decode); ok {
		bytes, err = encode.Encode(bytes)
		if err != nil {
			return -1, err
		}
	}
	return fsFile.WriteAt(bytes, off)
}

func ReadAt(afs afero.Fs, fsFile afero.File, off int64, data []byte, opts ...*fsopts.Options) (int, error) {
	count, err := fsFile.ReadAt(data, off)
	if err != nil {
		return -1, err
	}
	if decode, ok := afs.(Decode); ok {
		data, err = decode.Decode(data)
	}
	return count, err
}

func MkdirAll(afs afero.Fs, path string, fileMode fs.FileMode, opts ...*fsopts.Options) error {
	var err error
	fileName := fsopts.GetAbsPath(path, opts...)
	err = afs.MkdirAll(fileName, fileMode)
	return err
}
