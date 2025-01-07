package localfs

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"path/filepath"
	"strings"
)

type fileWatcher struct {
	f *Fs
}

func NewFileWatcher(f *Fs) fs.Watcher {
	return &fileWatcher{f: f}
}

func (w *fileWatcher) Start(opts ...fs.Option) error {
	//1、初始化监控对象watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Fail to create new Watcher[ %s ]\n", err)
	}

	rootPath, err := filepath.Abs(w.f.GetRootPath() + "/")
	if err != nil {
		return err
	}

	//3、启动监听文件对象事件协程
	go func() {
		for {
			select {
			case e := <-watcher.Events:
				lenName := len(e.Name)
				if lenName < 3 {
					break
				}
				/*
					extFile := strings.ToLower(e.Name[lenName-3:])
					if extFile != ".js" && extFile != ".ts" {
						break
					}
				*/
				eventType := fs.CreateEvent
				if e.Has(fsnotify.Create) {
					eventType = fs.CreateEvent
				} else if e.Has(fsnotify.Rename) {
					eventType = fs.RenameEvent
				} else if e.Has(fsnotify.Write) {
					eventType = fs.WriteEvent
				} else if e.Has(fsnotify.Remove) {
					eventType = fs.RemoveEvent
				} else if e.Has(fsnotify.Chmod) {
					eventType = fs.ChmodEvent
				}

				name := strings.ReplaceAll(e.Name, rootPath+"/", "")
				for _, o := range opts {
					o(rootPath, name, eventType)
				}
			case err := <-watcher.Errors:
				fmt.Printf(" %s\n", err.Error())
			}
		}
	}()

	// 2、将需要监听的文件加入到watcher的监听队列中
	err = watcher.Add(rootPath) //将文件加入监听
	if err != nil {
		return errors.New(fmt.Sprintf("Fail to watch directory[ %s ]", err))
	}
	return nil
}
