package fs_pkg

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/pkg"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/intutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"os"
	"sort"
)

type FsPkg struct {
	base        *fs.Manager
	cfg         common.IEnvConfig
	WriteModels *FsWriteModel
}

type FsWriteModel struct {
	AllWriteRead       fs.WriteModel
	SelfWriteOtherRead fs.WriteModel
}

func NewFsWriteModel() *FsWriteModel {
	return &FsWriteModel{
		AllWriteRead:       fs.WriteModelAllWriteRead,
		SelfWriteOtherRead: fs.WriteModelSelfWriteOtherRead,
	}
}

func NewFsPkg(cfg common.IEnvConfig) (*FsPkg, error) {
	fsManager, err := cfg.GetFsManager()
	if err != nil {
		return nil, err
	}
	fsm := &FsPkg{cfg: cfg, WriteModels: NewFsWriteModel()}
	fsm.base = fsManager
	return fsm, nil
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
		panic(err)
	}
	return res
}

func (m *FsPkg) WriteFile(filename string, data any, opts ...*fsopts.Options) {
	var bytes []byte = nil
	if str, ok := data.(string); ok {
		bytes = []byte(str)
	} else if bs, ok := data.([]byte); ok {
		bytes = bs
	} else if mapData, ok := data.(map[string]any); ok {
		str, err := jsonutils.Marshal(mapData)
		if err != nil {
			panic(err)
		}
		bytes = []byte(str)
	} else {
		panic("WriteFile() invalid runValues is string or []byte or map[string]any")
	}

	err := m.base.WriteFile(filename, bytes, fs.WriteModelAllWriteRead, opts...)
	if err != nil {
		panic(err)
	}
}

func (m *FsPkg) RemoveFile(filename string, opts ...*fsopts.Options) {
	err := m.base.RemoveFile(filename, opts...)
	if err != nil {
		panic(err)
	}
}

func (m *FsPkg) RemoveAll(name string, opts ...*fsopts.Options) {
	err := m.base.RemoveAll(name, opts...)
	if err != nil {
		panic(err)
	}
}

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
		fileInfo := &pkg.FileInfo{IsDir: file.IsDir(), Name: file.Name(), Size: file.Size(), Path: path, SizeTitle: intutils.GetFileSizeTitle(file.Size())}
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
			file.SubFiles = m.ReadDir(file.Path+"/"+file.Name, opts...)
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
