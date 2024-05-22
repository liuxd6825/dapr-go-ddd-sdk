package giteafs

import (
	"encoding/base64"
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
	"testing"
)

func TestConnect(t *testing.T) {
	opt := &Config{
		Url:      "http://localhost:3000",
		User:     "liuxd",
		Password: "liuxd",
		Repo:     "test",
		Branch:   "main",
	}
	fs, err := NewFs(opt)
	if err != nil {
		t.Error(err)
	}

	t.Run("open", func(t *testing.T) {
		if _, err := open(t, fs, "README.md"); err != nil {
			t.Error(err)
			return
		}
	})

	t.Run("update", func(t *testing.T) {
		if err := update(t, fs, "README.md", "hello world"); err != nil {
			t.Error(err)
			return
		}
	})

	newFile := "newFile.txt"
	t.Run("create", func(t *testing.T) {
		if file, err := fs.Create(newFile); err != nil {
			t.Error(err)
			return
		} else {
			content := base64.StdEncoding.EncodeToString([]byte("newFile.txt"))
			if _, err := file.Write([]byte(content)); err != nil {
				t.Error(err)
			} else {
				if err := file.Close(); err != nil {
					t.Error(err)
				}
			}
		}
	})

	t.Run("rename", func(t *testing.T) {
		if err := fs.Rename("newFile.txt", "newFile2.txt"); err != nil {
			t.Error(err)
		}
	})

	t.Run("delete", func(t *testing.T) {
		if err := fs.Remove("newFile2.txt"); err != nil {
			t.Error(err)
		}
	})

	t.Run("readDir", func(t *testing.T) {
		if dir, err := fs.Open("/"); err != nil {
			t.Error(err)
		} else {
			if list, err := dir.Readdir(100); err != nil {
				t.Error(err)
			} else {
				t.Log(list)
			}
		}
	})

}

func open(t *testing.T, fs afero.Fs, fileName string) (string, error) {
	file, err := fs.Open(fileName)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	content := string(bytes)
	t.Log("content:", content)
	return content, nil
}

func update(t *testing.T, fs afero.Fs, fileName string, content string) error {
	file, err := fs.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	bytes, err := base64.StdEncoding.DecodeString(content)
	if _, err := file.Write(bytes); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}
