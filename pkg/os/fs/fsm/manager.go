package fsm

import (
	"fmt"
	"io/fs"
	iofs "io/fs"
	"os"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	fs2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/giteafs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/httpfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/localfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/memoryfs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/types"
	"github.com/orcaman/concurrent-map"
	"github.com/spf13/afero"
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

type File struct {
	fs afero.Fs
	afero.File
}

func NewFile(fs afero.Fs, fsFile afero.File) *File {
	return &File{fs: fs, File: fsFile}
}

// SizeTitle 格式化文件大小为易读格式
func (f File) SizeTitle() string {
	const (
		KB = 1 << 10
		MB = 1 << 20
		GB = 1 << 30
	)
	fileInfo, _ := f.File.Stat()
	size := fileInfo.Size()
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

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
		manger.AddFs(id, fs)
	}
	return manger, nil
}

func (m *Manager) NewFsm(fsName string) (*Manager, error) {
	fsm := NewManager()
	fsm.DefaultFsName = fsName
	for _, name := range m.fsMap.Keys() {
		if v, ok := m.GetFs(name); ok {
			fsm.AddFs(name, v)
		}
	}
	return fsm, nil
}

func (m *Manager) HasFs(name string) bool {
	return m.fsMap.Has(name)
}

func (m *Manager) AddFs(name string, fs afero.Fs) {
	m.fsMap.Set(name, fs)
}

func (m *Manager) GetFs(name string) (afero.Fs, bool) {
	fs, ok := m.fsMap.Get(name)
	if !ok {
		return nil, false
	}
	return fs.(afero.Fs), true
}

func (m *Manager) GetFsByTag(tags ...string) (res []afero.Fs) {
	res = make([]afero.Fs, 0)
	for _, key := range m.fsMap.Keys() {
		f, ok := m.fsMap.Get(key)
		if ok {
			fs := f.(afero.Fs)
			if list, ok := f.(fs2.Tags); ok {
				isAdd := false
				for _, tag1 := range list.Tags() {
					for _, tag2 := range tags {
						if tag1 == tag2 {
							isAdd = true
							res = append(res, fs)
							break
						}
					}
					if isAdd {
						break
					}
				}
			}
		}
	}
	return res
}

func (m *Manager) RemoveFs(name string) {
	m.fsMap.Remove(name)
}

func (m *Manager) MapFs() map[string]afero.Fs {
	data := make(map[string]afero.Fs)
	for k, v := range m.fsMap.Items() {
		data[k] = v.(afero.Fs)
	}
	return data
}

func (m *Manager) GetFsByFileUrl(fileUrl string) (afero.Fs, error) {
	afs, _, err := m.parse(fileUrl)
	return afs, err
}

func (m *Manager) CreateFile(filename string, opts ...*fsopts.Options) (afero.File, error) {
	opt := fsopts.NewOptions(opts...)
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return nil, err
	}
	if opt.Fs != nil {
		afs = opt.Fs
	}
	return fs2.Create(afs, fileName, opt)
}

func (m *Manager) MoveDir(aSrcDir, aDestDir string) error {
	fs, srcDir, err := m.parse(aSrcDir)
	if err != nil {
		return err
	}
	_, destDir, err := m.parse(aDestDir)
	if err != nil {
		return err
	}
	// 创建目标目录
	err = m.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// 获取源目录内容
	files, err := afero.ReadDir(fs, srcDir)
	if err != nil {
		return err
	}

	// 移动目录内容
	for _, file := range files {
		srcPath := srcDir + "/" + file.Name()
		destPath := destDir + "/" + file.Name()

		if file.IsDir() {
			err = m.MoveDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			err = fs.Rename(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	// 删除空源目录
	return fs.Remove(srcDir)
}

func (m *Manager) ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error) {
	opt := fsopts.NewOptions(opts...)
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return nil, err
	}
	if opt.Fs != nil {
		afs = opt.Fs
	}
	return fs2.ReadFile(afs, fileName, opt)
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
func (m *Manager) WriteFile(filename string, bytes []byte, writeModel WriteModel, opts ...*fsopts.Options) error {
	opt := fsopts.NewOptions(opts...)
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return err
	}
	if opt.Fs != nil {
		afs = opt.Fs
	}
	return fs2.WriteFile(afs, fileName, bytes, fs.FileMode(writeModel), opt)
}

func (m *Manager) Open(filename string, flag int, perm os.FileMode, opts ...*fsopts.Options) (*File, error) {
	opt := fsopts.NewOptions(opts...)
	afs, fName, err := m.parse(filename)
	if err != nil {
		return nil, err
	}
	fsFile, err := fs2.Open(afs, fName, flag, perm, opt)
	if err != nil {
		return nil, err
	}
	return NewFile(afs, fsFile), nil
}

func (m *Manager) WriteAt(fsFile *File, data []byte, off int64, opts ...*fsopts.Options) (int, error) {
	opt := fsopts.NewOptions(opts...)
	if fsFile == nil {
		return -1, errors.New("fsFile==nil")
	}
	return fs2.WriteAt(fsFile.fs, fsFile, off, data, opt)
}

func (m *Manager) ReadAt(fsFile *File, data []byte, off int64, opts ...*fsopts.Options) (int, error) {
	opt := fsopts.NewOptions(opts...)
	if fsFile == nil {
		return -1, errors.New("fsFile==nil")
	}
	return fs2.ReadAt(fsFile.fs, fsFile, off, data, opt)
}

// RemoveFile
//
//	@Description: removes a file identified by name, returning an error, if any
//	@receiver m
//	@param filename 文件名称
//	@return error
func (m *Manager) RemoveFile(filename string, opts ...*fsopts.Options) error {
	opt := fsopts.NewOptions(opts...)
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return err
	}
	if opt.Fs != nil {
		afs = opt.Fs
	}
	fileName = fsopts.GetAbsPath(filename, opt)
	return afs.Remove(fileName)
}

// RemoveAll
//
//	@Description:  RemoveAll removes a directory path and any children it contains. It
//	@receiver m
//	@param path 路径名称
//	@return error
func (m *Manager) RemoveAll(path string, opts ...*fsopts.Options) error {
	opt := fsopts.NewOptions(opts...)
	afs, path, err := m.parse(path)
	if err != nil {
		return err
	}
	if opt.Fs != nil {
		afs = opt.Fs
	}
	path = fsopts.GetAbsPath(path, opt)
	return afs.RemoveAll(path)
}

// Rename
//
//	@Description:  文件重命名
//	@receiver m
//	@param aOldName 原文件名
//	@param aNewName 新文件名
//	@return error
func (m *Manager) Rename(aOldName, aNewName string, opts ...*fsopts.Options) error {
	afs, oldName, err := m.parse(aOldName)
	if err != nil {
		return err
	}
	afs, newName, err := m.parse(aNewName)
	if err != nil {
		return err
	}

	opt := fsopts.NewOptions(opts...)
	if opt.Fs != nil {
		afs = opt.Fs
	}
	return afs.Rename(oldName, newName)
}

// Mkdir
//
//	@Description:
//	@receiver m
//	@param path
//	@param perm
//	@param opts
//	@return error
func (m *Manager) Mkdir(path string, perm os.FileMode, opts ...*fsopts.Options) error {
	afs, name, err := m.parse(path)
	if err != nil {
		return err
	}

	opt := fsopts.NewOptions(opts...)
	if opt.Fs != nil {
		afs = opt.Fs
	}

	name = fsopts.GetAbsPath(name, opt)
	return afs.Mkdir(name, perm)
}

// MkdirAll
//
//	@Description:
//	@receiver m
//	@param filename
//	@param fileMode
//	@param opts
//	@return error
func (m *Manager) MkdirAll(filename string, fileMode fs.FileMode, opts ...*fsopts.Options) error {
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return err
	}
	opt := fsopts.NewOptions(opts...)
	if opt.Fs != nil {
		afs = opt.Fs
	}
	return fs2.MkdirAll(afs, fileName, fileMode, opts...)
}

// Exists 文件是否存在
func (m *Manager) Exists(filename string, opts ...*fsopts.Options) (bool, error) {
	afs, fileName, err := m.parse(filename)
	if err != nil {
		return false, err
	}
	opt := fsopts.NewOptions(opts...)
	if opt.Fs != nil {
		afs = opt.Fs
	}

	return fs2.Exists(afs, fileName, opt)
}

// ReadDir
//
//	@Description: 取得所有子目录与文件
//	@receiver m
//	@param path
//	@return []iofs.FileInfo
//	@return error
func (m *Manager) ReadDir(path string, opts ...*fsopts.Options) ([]iofs.FileInfo, error) {
	afs, _, err := m.parse(path)
	if err != nil {
		return nil, err
	}

	opt := fsopts.NewOptions(opts...)
	if opt.Fs != nil {
		afs = opt.Fs
	}

	return afero.ReadDir(afs, path)
}

func (m *Manager) parse(filename string) (afs afero.Fs, fileName string, err error) {
	var fsName string
	err = fsopts.ParseFileName(filename, m.DefaultFsName, func(aFsName, aFileName string) {
		fsName = aFsName
		fileName = aFileName
	})
	if err != nil {
		return nil, "", err
	}
	fs, ok := m.GetFs(fsName)
	if !ok {
		return nil, "", errors.New(fmt.Sprintf("fs name \"%s\" not found", fsName))
	}
	return fs, fileName, nil
}
