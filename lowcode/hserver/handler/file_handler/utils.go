package file_handler

import (
	"context"
	"fmt"
	"github.com/spf13/afero"
	"regexp"
	"sync"
)

// 定义动态标记的正则表达式
var dynamicMetaRegex = regexp.MustCompile(`<meta\s+name="dynamic-page"\s+content="true"\s*/?>`)
var dynamicCommentRegex = regexp.MustCompile(`<!--\s*dynamic-page\s*-->`)

// 缓存文件的动态标记状态
var dynamicCache sync.Map

// IsFileExist 判断文件是否存在
func IsFileExist(fs afero.Fs, path string) (bool, error) {
	exists, err := afero.Exists(fs, path)
	return exists, err
}

// isDynamicPageRegex 使用正则表达式检查 HTML 文件是否是动态页面
func isDynamicPage(ctx context.Context, fs afero.Fs, filePath string) (bool, error) {
	// 检查缓存
	if cached, ok := dynamicCache.Load(filePath); ok {
		return cached.(bool), nil
	}

	// 打开文件
	file, err := fs.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file afero.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 只读取文件前 1024 字节
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		return false, err
	}
	content := buffer[:n]

	// 检查是否匹配动态标记
	isDynamic := dynamicMetaRegex.Match(content) || dynamicCommentRegex.Match(content)

	// 缓存结果
	dynamicCache.Store(filePath, isDynamic)
	return isDynamic, nil
}
