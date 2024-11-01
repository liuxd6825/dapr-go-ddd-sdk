package server

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/fs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/rs-server/modules/common"
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

func (m *FsManager) ReadFile(filename string, pwd string) []byte {
	res, err := m.base.ReadFile(filename, pwd)
	if err != nil {
		panic(err)
	}
	return res
}

func (m *FsManager) WriteFile(filename string, pwd string, bytes []byte) {
	err := m.base.WriteFile(filename, pwd, bytes, fs.WriteModelAllWriteRead)
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
	err := m.base.RemoveFile(name)
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
