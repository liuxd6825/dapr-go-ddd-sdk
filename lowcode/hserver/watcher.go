package hserver

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/element"
	"github.com/spf13/afero"
)

// Watcher
// @Description: 监控服务目录，当文件修改时重新启动服务
type Watcher struct {
	server element.Server
	fs     afero.Fs
}

type WatcherHandler = func(rootPath, fileName string, eventType fs.WatcherEventType) error

func NewWatcher(server element.Server, afs afero.Fs, handler WatcherHandler) *Watcher {
	w := &Watcher{server: server, fs: afs}
	if wfs, ok := afs.(fs.WatcherFs); ok {
		watcher, err := wfs.NewWatcher()
		if err != nil {
			panic(err)
		}
		if err = watcher.Start(handler); err != nil {
			panic(err)
		}
	}
	return w
}
