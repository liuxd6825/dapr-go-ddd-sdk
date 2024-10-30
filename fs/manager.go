package fs

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/giteafs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/httpfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/memoryfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
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
	for _, m := range maps {
		o := types.Object(m)
		typeVal := o.GetString("type")
		if typeVal == "" {
			return nil, errors.New("config type not found in config")
		}
		nameVal := o.GetString("name")
		if nameVal == "" {
			return nil, errors.New("config name not found in config")
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
			return nil, errors.New("config type not support : " + typeVal)
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
			return nil, errors.New("config type not support")
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

func (m *Manager) ReadFile(filename string, pwd string) ([]byte, error) {
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return nil, err
	}
	return ReadFile(afs, pwd, fileName)
}

func (m *Manager) WriteFile(filename string, pwd string, bytes []byte, writeModel WriteModel) error {
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return err
	}
	return WriteFile(afs, pwd, fileName, bytes, fs.FileMode(writeModel))
}

func (m *Manager) RemoveFile(filename string) error {
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return err
	}
	return afs.Remove(fileName)
}

func (m *Manager) RemoveAll(name string) error {
	afs, name, err := m.parse(name)
	if err != nil {
		return err
	}
	return afs.RemoveAll(name)
}

func (m *Manager) Mkdir(name string, perm os.FileMode) error {
	afs, name, err := m.parse(name)
	if err != nil {
		return err
	}
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
