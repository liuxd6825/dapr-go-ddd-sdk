package fs_pkg

import (
	fstest "github.com/liuxd6825/dapr-go-ddd-sdk/fs/test"
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/test"
	"testing"
)

func TestFs_LoadFile(t *testing.T) {
	fsm, err := fstest.NewFsManager("")
	if err != nil {
		t.Fatal(err)
		return
	}

	fsPkg, err := NewFsPkg(test.NewEnvConfig(fsm, ""))
	if err != nil {
		t.Fatal(err)
		return
	}

	fileInfos := fsPkg.ReadAllDir("/testfile")
	for _, fileInfo := range fileInfos {
		t.Log(fileInfo)
	}
}

func TestFs_RemoveFile(t *testing.T) {
	fsm, err := fstest.NewFsManager("")
	if err != nil {
		t.Fatal(err)
		return
	}

	fsPkg, err := NewFsPkg(test.NewEnvConfig(fsm, ""))
	if err != nil {
		t.Fatal(err)
		return
	}
	fsPkg.RemoveFile("/testfile/master/remove.txt")
}
