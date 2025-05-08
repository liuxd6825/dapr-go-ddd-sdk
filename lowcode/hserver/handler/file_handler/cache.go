package file_handler

import (
	"sync"
	"time"
)

var cache sync.Map // 全局文件缓存
// 定义缓存结构
type FileCache struct {
	Content    []byte    // 文件内容
	IsDynamic  bool      // 是否为动态页面
	LastUpdate time.Time // 上次更新时间
}
