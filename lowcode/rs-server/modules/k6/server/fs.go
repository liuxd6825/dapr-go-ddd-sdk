package server

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"os"
)

type FsManager struct {
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

func NewFsManger(cfg common.IEnvConfig) (*FsManager, error) {
	fsManager, err := cfg.GetFsManager()
	if err != nil {
		return nil, err
	}
	fsm := &FsManager{cfg: cfg, WriteModels: NewFsWriteModel()}
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
func (m *FsManager) ReadFile(filename string, options ...*fsopts.Options) []byte {
	res, err := m.base.ReadFile(filename, options...)
	if err != nil {
		panic(err)
	}
	return res
}

func (m *FsManager) WriteFile(filename string, data any, options ...*fsopts.Options) {
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
		panic("WriteFile() invalid data is string or []byte or map[string]any")
	}

	err := m.base.WriteFile(filename, bytes, fs.WriteModelAllWriteRead, options...)
	if err != nil {
		panic(err)
	}
}

func (m *FsManager) RemoveFile(filename string) {
	err := m.base.RemoveFile(filename)
	if err != nil {
		panic(err)
	}
}

func (m *FsManager) RemoveAll(name string) {
	err := m.base.RemoveAll(name)
	if err != nil {
		panic(err)
	}
}

func (m *FsManager) Mkdir(name string, perm os.FileMode) {
	err := m.base.Mkdir(name, perm)
	if err != nil {
		panic(err)
	}
}
