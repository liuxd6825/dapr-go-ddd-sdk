package render

import (
	"sync"
	"time"
)

var cache sync.Map // 全局文件缓存
// FileCache
// @Description:
type FileCache struct {
	Content    []byte    // 文件内容
	IsDynamic  bool      // 是否为动态页面
	LastUpdate time.Time // 上次更新时间
}
