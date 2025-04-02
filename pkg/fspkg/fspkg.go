package fspkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/intutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/spf13/afero"
	"os"
	"sort"
	"strings"
)

type IFsPkg interface {
	NewFs(fsName string) *FsPkg
	Exists(fileName string, opts ...*fsopts.Options) bool
	Create(name string, opts ...*fsopts.Options) afero.File
	Rename(oldName, newName string, opts ...*fsopts.Options) error
	ReadFile(filename string, opts ...*fsopts.Options) []byte
	WriteJson(filename string, data any, opts ...*fsopts.Options)
	WriteFile(filename string, data any, opts ...*fsopts.Options)
	RemoveFile(filename string, opts ...*fsopts.Options)
	RemoveAll(name string, opts ...*fsopts.Options)
	Mkdir(name string, perm os.FileMode, opts ...*fsopts.Options)
	ReadPath(path string, opts ...*fsopts.Options) []*FileInfo
	ReadAllPath(path string, opts ...*fsopts.Options) []*FileInfo
	SortFileInfos(files []*FileInfo)
}

// FsPkg
// @Description:  文件系统
type FsPkg struct {
	base        *fsm.Manager
	cfg         env.IEnvConfig
	WriteModels *FsWriteModel
	fsName      string
}

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
func NewFsPkg(cfg env.IEnvConfig, fsName string) (*FsPkg, error) {
	var err error
	fsM := cfg.GetFsManager()

	if fsName != "" {
		fsM, err = fsM.NewFsm(fsName)
		if err != nil {
			panic(err)
		}
	}

	fsPkg := &FsPkg{
		cfg:         cfg,
		fsName:      fsName,
		WriteModels: NewFsWriteModel(),
		base:        fsM,
	}

	return fsPkg, nil
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
	name = getFileName(name)
	file, err := m.base.CreateFile(name, opts...)
	if err != nil {
		panic(err)
	}
	return file
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
	oldName = getFileName(oldName)
	newName = getFileName(newName)
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
	filename = getFileName(filename)
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
	filename = getFileName(filename)
	var toBytes = ToBytes(data)
	toBytes = m.JsonFormat(toBytes)
	err := m.base.WriteFile(filename, toBytes, fsm.WriteModelAllWriteRead, opts...)
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
	filename = getFileName(filename)
	var toBytes = ToBytes(data)
	err := m.base.WriteFile(filename, toBytes, fsm.WriteModelAllWriteRead, opts...)
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
	filename = getFileName(filename)
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
	name = getFileName(name)
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
	name = getFileName(name)
	err := m.base.Mkdir(name, perm)
	if err != nil {
		panic(err)
	}
}

func (m *FsPkg) Exists(fileName string, opts ...*fsopts.Options) bool {
	fileName = getFileName(fileName)
	v, err := m.base.Exists(fileName, opts...)
	if err != nil {
		panic(err)
	}
	return v
}

// ReadPath
//
//	@Description: 读取指定目录中的子目录与文件。
//	@receiver m
//	@param path
//	@param opts
//	@return []*FileInfo
func (m *FsPkg) ReadPath(path string, opts ...*fsopts.Options) []*FileInfo {
	path = getFileName(path)
	res := make([]*FileInfo, 0)
	files, err := m.base.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fileInfo := &FileInfo{
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

// ReadAllPath
//
//	@Description:  深度读取目录内容，读取所有子目录与文件。
//	@receiver m
//	@param path
//	@param opts
//	@return []*FileInfo
func (m *FsPkg) ReadAllPath(path string, opts ...*fsopts.Options) []*FileInfo {
	path = getFileName(path)
	fileInfos := m.ReadPath(path, opts...)
	for _, file := range fileInfos {
		if file.IsDir {
			file.SubFiles = m.ReadAllPath(file.Path+"/"+file.Name, opts...)
		}
	}
	m.SortFileInfos(fileInfos)
	return fileInfos
}

// SortFileInfos
//
//	@Description: sorts a slice of *FileInfo by Name and IsDir, and recursively sorts SubFiles.
//	@param files
func (m *FsPkg) SortFileInfos(files []*FileInfo) {
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

// JsonFormat
//
//	@Description: 对json字符串进行格式化
//	@receiver j
//	@param rawJSON
//	@return []byte
func (m *FsPkg) JsonFormat(rawJSON any) []byte {
	data := ToBytes(rawJSON)
	// 创建一个缓冲区来存储格式化后的 JSON
	var formattedJSON bytes.Buffer
	// 使用 json.Indent 进行格式化
	err := json.Indent(&formattedJSON, data, "", "  ") // 第二个参数是缩进字符串
	if err != nil {
		panic(errors.New("json format error: %s", err.Error()))
	}
	return formattedJSON.Bytes()
}

// ToBytes
//
//	@Description: 将any转换成byte数组，支持string, []byte, map[string]any类型的转换
//	@param data
//	@return []byte
func ToBytes(data any) []byte {
	var b []byte = nil
	if str, ok := data.(string); ok {
		b = []byte(str)
	} else if bs, ok := data.([]byte); ok {
		b = bs
	} else if mapData, ok := data.(map[string]any); ok {
		str, err := jsonutils.Marshal(mapData)
		if err != nil {
			panic(err)
		}
		b = []byte(str)
	} else {
		panic("WriteFile() invalid runValues is string or []byte or map[string]any")
	}
	return b
}

func getFileName(name string) string {
	if strings.HasPrefix(name, "/") || strings.HasPrefix(name, "./") || strings.HasPrefix(name, "../") {
		return name
	}
	return "/" + name
}
