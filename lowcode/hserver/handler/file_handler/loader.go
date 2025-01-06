package file_handler

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/spf13/afero"
)

// Loader 是基于 Afero 的 Pongo2 加载器
type Loader struct {
	fs afero.Fs // 使用 Afero 文件系统
}

// NewLoader 创建一个新的 Afero 加载器
func NewLoader(fs afero.Fs) *Loader {
	return &Loader{fs: fs}
}

// Abs 返回绝对路径（Afero 直接操作虚拟路径，不需要实际文件路径）
func (l *Loader) Abs(base, name string) string {
	if filepath.IsAbs(name) {
		return name
	}
	return filepath.Join(filepath.Dir(base), name)
}

// Get 读取模板内容
func (l *Loader) Get(path string) (io.Reader, error) {
	content, err := afero.ReadFile(l.fs, path)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(content), nil
}
