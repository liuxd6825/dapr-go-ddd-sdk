package watcher

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type fileWatcher struct {
	path string
}

func NewFileWatcher(path string) Watcher {
	return &fileWatcher{path}
}

func (w *fileWatcher) Start(opts ...Options) error {
	//1、初始化监控对象watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Fail to create new Watcher[ %s ]\n", err)
	}

	rootPath, err := filepath.Abs(w.path + "/")
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
				extFile := strings.ToLower(e.Name[lenName-3:])
				if extFile != ".js" && extFile != ".ts" {
					break
				}
				eventType := CreateEvent
				if e.Has(fsnotify.Create) {
					eventType = CreateEvent
				} else if e.Has(fsnotify.Rename) {
					eventType = RenameEvent
				} else if e.Has(fsnotify.Write) {
					eventType = WriteEvent
				} else if e.Has(fsnotify.Remove) {
					eventType = RemoveEvent
				} else if e.Has(fsnotify.Chmod) {
					eventType = ChmodEvent
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
