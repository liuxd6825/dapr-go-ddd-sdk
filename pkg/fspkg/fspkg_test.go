package fspkg

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"os"
	"testing"
)

func TestFs_LoadFile(t *testing.T) {
	env := newEnv()
	fsPkg, err := NewFsPkg(env, "file")
	if err != nil {
		t.Fatal(err)
		return
	}

	fileInfos := fsPkg.ReadAllPath("/testfile")
	for _, fileInfo := range fileInfos {
		t.Log(fileInfo)
	}
}

func TestFs_RemoveFile(t *testing.T) {
	env := newEnv()
	fsPkg, err := NewFsPkg(env, "file")
	if err != nil {
		t.Fatal(err)
		return
	}
	fsPkg.RemoveFile("/testfile/master/xremove.txt")
}

func TestFs_WriteAt(t *testing.T) {
	env := newEnv()
	fsPkg, err := NewFsPkg(env, "file")
	if err != nil {
		t.Fatal(err)
		return
	}
	fsPkg.Create("001.txt")
	writeFile, err := fsPkg.Open("/001.txt", os.O_WRONLY, 0644)

	if err != nil {
		t.Fatal(err)
		return
	}
	for i := 0; i < 10; i++ {
		str := fmt.Sprintf("<%d>", i)
		if _, err := fsPkg.WriteAt(writeFile, []byte(str), int64(i*3)); err != nil {
			t.Fatal(err)
			return
		}
	}
	writeFile.Close()

	readFile, err := fsPkg.Open("/001.txt", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatal(err)
		return
	}

	fileInfo, err := readFile.Stat()
	if err != nil {
		t.Fatal(err)
		return
	}

	fileSize := fileInfo.Size()
	bufSize := int64(8)
	off := int64(0)
	for {
		fileSize = fileSize - bufSize
		if fileSize <= 0 {
			bufSize = int64(bufSize + fileSize)
		}
		if bufSize <= 0 {
			break
		}

		buf := make([]byte, bufSize)
		_, err := fsPkg.ReadAt(readFile, buf, off)
		if err != nil {
			t.Fatal(err)
			return
		}
		t.Log(string(buf))
		off += bufSize
	}

}

func newEnv() *env.Env {
	return xtest.NewEnvConfig_Fs("file", "/Users/lxd/Projects/liuxd6825/dapr/dapr-go-ddd-sdk/pkg/fspkg/test-file")
}
