package fs

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/giteafs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/httpfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/memoryfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/fileutils"
	"github.com/orcaman/concurrent-map"
	"github.com/spf13/afero"
	"io/fs"
	"os"
)

//	Manager
//
// @Description: 文件系统管理器
// @author liuxd6825
type Manager struct {
	fsMap         cmap.ConcurrentMap
	DefaultFsName string
}

type WriteModel fs.FileMode

const (
	WriteModelAllWriteRead       WriteModel = 0666 // 所有用户都可读写
	WriteModelSelfWriteOtherRead WriteModel = 0644 // 当前用户读写，其他用户只读
)

func NewManager() *Manager {
	fsMap := cmap.New()
	return &Manager{fsMap: fsMap}
}

// NewManagerWithConfigs
//
//	@Description: 使用map[string]any配置创建文件系统管理器, 目前支持local和gitea
//	@param maps
//	@return *Manager
//	@return error
func NewManagerWithConfigs(maps []map[string]any, defaultFsName string) (*Manager, error) {
	clist := make(map[string]any)
	i := 0
	for _, m := range maps {
		i++
		o := types.Object(m)
		typeVal := o.GetString("type")
		if typeVal == "" {
			return nil, errors.New("maps[%d]. type not found in config")
		}
		nameVal := o.GetString("name")
		if nameVal == "" {
			return nil, errors.New("maps[%d].name cannot be empty", i)
		}

		var cfg any
		var err error
		switch typeVal {
		case localfs.Name():
			cfg, err = localfs.NewConfig(m)
		case giteafs.Name():
			cfg, err = giteafs.NewConfig(m)
		case httpfs.Name():
			cfg, err = httpfs.NewConfig(m)
		case memoryfs.Name():
			cfg, err = memoryfs.NewConfig(m)
		default:
			return nil, errors.New(" maps[%d].type not support : " + typeVal)
		}
		if err != nil {
			return nil, err
		}
		clist[nameVal] = cfg
	}
	manger := NewManager()
	manger.DefaultFsName = defaultFsName
	for id, cfg := range clist {
		var fs afero.Fs
		var err error
		switch cfg.(type) {
		case *localfs.Config:
			c := *cfg.(*localfs.Config)
			fs, err = localfs.NewFs(c)
		case *giteafs.Config:
			fs, err = giteafs.NewFs(cfg.(*giteafs.Config))
		case *httpfs.Config:
			fs, err = httpfs.NewFs(cfg.(*httpfs.Config))
		case *memoryfs.Config:
			fs, err = memoryfs.NewFs(cfg.(*memoryfs.Config))
		default:
			return nil, errors.New("fs.NewManagerWithConfigs() fs.config not support")
		}
		if err != nil {
			return nil, err
		}
		manger.Add(id, fs)
	}
	return manger, nil
}

func (m *Manager) Add(name string, fs afero.Fs) {
	m.fsMap.Set(name, fs)
}

func (m *Manager) Get(name string) (afero.Fs, bool) {
	fs, ok := m.fsMap.Get(name)
	return fs.(afero.Fs), ok
}

func (m *Manager) Remove(name string) {
	m.fsMap.Remove(name)
}

func (m *Manager) Map() map[string]afero.Fs {
	data := make(map[string]afero.Fs)
	for k, v := range m.fsMap.Items() {
		data[k] = v.(afero.Fs)
	}
	return data
}

func (m *Manager) ReadFile(filename string, opts ...*Options) ([]byte, error) {
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return nil, err
	}
	return ReadFile(afs, fileName, opts...)
}

// WriteFile
//
//	@Description: 写文件
//	@receiver m
//	@param filename 文件名称
//	@param pwd
//	@param bytes
//	@param writeModel
//	@return error
func (m *Manager) WriteFile(filename string, bytes []byte, writeModel WriteModel, opts ...*Options) error {
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return err
	}
	return WriteFile(afs, fileName, bytes, fs.FileMode(writeModel), opts...)
}

// RemoveFile
//
//	@Description: 删除文件
//	@receiver m
//	@param filename 文件名称
//	@return error
func (m *Manager) RemoveFile(filename string, opts ...*Options) error {
	o := NewOptions(opts...)
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return err
	}
	fileName = fileutils.AbsPath(filename, o.BasePath)
	return afs.Remove(fileName)
}

// RemoveAll
//
//	@Description: 删除路径
//	@receiver m
//	@param path 路径名称
//	@return error
func (m *Manager) RemoveAll(path string, basePath ...string) error {
	afs, path, err := m.parse(path)
	if err != nil {
		return err
	}
	path = fileutils.AbsPath(path, basePath...)
	return afs.RemoveAll(path)
}

// Rename
//
//	@Description:  文件重命名
//	@receiver m
//	@param oldname 原文件名
//	@param newname 新文件名
//	@return error
func (m *Manager) Rename(aOldName, aNewName string) error {
	afs, oldName, err := m.parse(aOldName)
	if err != nil {
		return err
	}
	afs, newName, err := m.parse(aNewName)
	if err != nil {
		return err
	}
	return afs.Rename(oldName, newName)
}

func (m *Manager) Mkdir(path string, perm os.FileMode, basePath ...string) error {
	afs, name, err := m.parse(path)
	if err != nil {
		return err
	}
	name = fileutils.AbsPath(name, basePath...)
	return afs.Mkdir(name, perm)
}

func (m *Manager) parse(filename string) (afs afero.Fs, fileName string, err error) {
	var fsName string
	err = ParseFileName(filename, m.DefaultFsName, func(aFsName, aFileName string) {
		fsName = aFsName
		fileName = aFileName
	})
	if err != nil {
		return nil, "", err
	}
	fs, ok := m.Get(fsName)
	if !ok {
		return nil, "", errors.New(fmt.Sprintf("file %s not found in config", fsName))
	}
	return fs, fileName, nil
}
