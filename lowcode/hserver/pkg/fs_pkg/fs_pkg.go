package fs_pkg

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg/json_pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/intutils"
	"github.com/spf13/afero"
	"os"
	"sort"
)

// FsPkg
// @Description:  文件系统
type FsPkg struct {
	base        *fsm.Manager
	cfg         common.IEnvConfig
	WriteModels *FsWriteModel
	fsName      string
	fs          afero.Fs
}

var jsonPkg *json_pkg.JsonPkg = json_pkg.NewJsonPkg()

// FsWriteModel
// @Description: 读写权限
type FsWriteModel struct {
	AllWriteRead       fsm.WriteModel
	SelfWriteOtherRead fsm.WriteModel
}

func NewFsWriteModel() *FsWriteModel {
	return &FsWriteModel{
		AllWriteRead:       fsm.WriteModelAllWriteRead,
		SelfWriteOtherRead: fsm.WriteModelSelfWriteOtherRead,
	}
}

// NewFsPkg
//
//	@Description: 创建文件系统
//	@param cfg
//	@return *FsPkg
//	@return error
func NewFsPkg(cfg common.IEnvConfig, fsName string) (*FsPkg, error) {
	fsManager, err := cfg.GetFsManager()
	if err != nil {
		return nil, err
	}

	fsm := &FsPkg{
		cfg:         cfg,
		fsName:      fsName,
		WriteModels: NewFsWriteModel(),
	}
	if fsName != "" {
		fs, ok := fsManager.Get(fsName)
		if ok {
			fsm.fs = fs
		}
	}
	fsm.base = fsManager
	return fsm, nil
}

func (m *FsPkg) NewFs(fsName string) *FsPkg {
	fs, err := NewFsPkg(m.cfg, fsName)
	if err != nil {
		panic(err)
	}
	return fs
}

// Create
//
//	@Description: 创建文件
//	@receiver m
//	@param name
//	@param opts
func (m *FsPkg) Create(name string, opts ...*fsopts.Options) afero.File {

	file, err := m.base.Create(name, opts...)
	if err != nil {
		panic(err)
	}
	return file
}

func (m *FsPkg) newOptions(opts ...*fsopts.Options) *fsopts.Options {
	opt := fsopts.NewOptions(opts...)
	if m.fs != nil {
		opt.Fs = m.fs
	}
	return opt
}

// Rename
//
//	@Description: 文件改名
//	@receiver m
//	@param oldName 原名
//	@param newName 新名
//	@param opts
//	@return error
func (m *FsPkg) Rename(oldName, newName string, opts ...*fsopts.Options) error {
	return m.base.Rename(oldName, newName)
}

// ReadFile
//
//	@Description: 读取文件内容
//	@receiver m
//	@param filename 文件名称
//	@param basePath 当前目录
//	@return []byte
func (m *FsPkg) ReadFile(filename string, opts ...*fsopts.Options) []byte {
	res, err := m.base.ReadFile(filename, opts...)
	if err != nil {
		panic(fmt.Sprintf("fsPkg.readFile() %s %s ", filename, err.Error()))
	}
	return res
}

// WriteJson
//
//	@Description: 写Json文件，并进行格式化。
//	@receiver m
//	@param filename
//	@param data
//	@param opts
func (m *FsPkg) WriteJson(filename string, data any, opts ...*fsopts.Options) {
	var bytes = pkg.ToBytes(data)

	bytes = jsonPkg.Format(bytes)
	err := m.base.WriteFile(filename, bytes, fsm.WriteModelAllWriteRead, opts...)
	if err != nil {
		panic(err)
	}
}

// WriteFile
//
//	@Description: 写文件
//	@receiver m
//	@param filename
//	@param data
//	@param opts
func (m *FsPkg) WriteFile(filename string, data any, opts ...*fsopts.Options) {
	var bytes = pkg.ToBytes(data)
	err := m.base.WriteFile(filename, bytes, fsm.WriteModelAllWriteRead, opts...)
	if err != nil {
		panic(err)
	}
}

// RemoveFile
//
//	@Description: 删除文件
//	@receiver m
//	@param filename
//	@param opts
func (m *FsPkg) RemoveFile(filename string, opts ...*fsopts.Options) {
	err := m.base.RemoveFile(filename, opts...)
	if err != nil {
		panic(err)
	}
}

// RemoveAll
//
//	@Description: 删除所有子目录与在内的所有文件
//	@receiver m
//	@param name
//	@param opts
func (m *FsPkg) RemoveAll(name string, opts ...*fsopts.Options) {
	err := m.base.RemoveAll(name, opts...)
	if err != nil {
		panic(err)
	}
}

// Mkdir
//
//	@Description: 创建目录
//	@receiver m
//	@param name
//	@param perm
//	@param opts
func (m *FsPkg) Mkdir(name string, perm os.FileMode, opts ...*fsopts.Options) {
	err := m.base.Mkdir(name, perm)
	if err != nil {
		panic(err)
	}
}

// ReadDir
//
//	@Description: 读取指定目录中的子目录与文件。
//	@receiver m
//	@param path
//	@param opts
//	@return []*FileInfo
func (m *FsPkg) ReadDir(path string, opts ...*fsopts.Options) []*pkg.FileInfo {
	res := make([]*pkg.FileInfo, 0)
	files, err := m.base.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fileInfo := &pkg.FileInfo{
			IsDir:     file.IsDir(),
			Name:      file.Name(),
			Size:      file.Size(),
			Path:      path,
			SizeTitle: intutils.GetFileSizeTitle(file.Size()),
		}
		res = append(res, fileInfo)
	}
	return res
}

// ReadAllDir
//
//	@Description:  深度读取目录内容，读取所有子目录与文件。
//	@receiver m
//	@param path
//	@param opts
//	@return []*FileInfo
func (m *FsPkg) ReadAllDir(path string, opts ...*fsopts.Options) []*pkg.FileInfo {
	fileInfos := m.ReadDir(path, opts...)
	for _, file := range fileInfos {
		if file.IsDir {
			file.SubFiles = m.ReadAllDir(file.Path+"/"+file.Name, opts...)
		}
	}
	m.SortFileInfos(fileInfos)
	return fileInfos
}

// SortFileInfos
//
//	@Description: sorts a slice of *FileInfo by Name and IsDir, and recursively sorts SubFiles.
//	@param files
func (m *FsPkg) SortFileInfos(files []*pkg.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir // Directories come first
		}
		return files[i].Name < files[j].Name // Then sort by name
	})

	for _, file := range files {
		if file.IsDir && len(file.SubFiles) > 0 {
			m.SortFileInfos(file.SubFiles) // Recursively sort SubFiles
		}
	}
}
