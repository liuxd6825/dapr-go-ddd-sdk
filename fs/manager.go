package fs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/giteafs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/httpfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/localfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/memoryfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/spf13/afero"
)

//	Manager
//
// @Description: 文件系统管理器
// @author liuxd6825
type Manager struct {
	fsMap map[string]afero.Fs
}

func NewManager() *Manager {
	return &Manager{fsMap: make(map[string]afero.Fs)}
}

// NewManagerWithConfigs
//
//	@Description: 使用map[string]any配置创建文件系统管理器, 目前支持local和gitea
//	@param maps
//	@return *Manager
//	@return error
func NewManagerWithConfigs(maps []map[string]any) (*Manager, error) {
	clist := make(map[string]any)
	for _, m := range maps {
		o := types.Object(m)
		typeVal := o.GetString("type")
		if typeVal == "" {
			return nil, errors.New("config type not found in config")
		}
		idVal := o.GetString("id")
		if idVal == "" {
			return nil, errors.New("config cfgId not found in config")
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
		clist[idVal] = cfg
	}
	manger := NewManager()
	for id, cfg := range clist {
		var fs afero.Fs
		var err error
		switch cfg.(type) {
		case *localfs.Config:
			fs, err = localfs.NewFs(cfg.(*localfs.Config))
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
	m.fsMap[name] = fs
}

func (m *Manager) Get(name string) (afero.Fs, bool) {
	fs, ok := m.fsMap[name]
	return fs, ok
}

func (m *Manager) Remove(name string) {
	delete(m.fsMap, name)
}

func (m *Manager) Map() map[string]afero.Fs {
	return m.fsMap
}
