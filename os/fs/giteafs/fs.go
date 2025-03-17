package giteafs

import (
	gitea "code.gitea.io/sdk/gitea"
	"encoding/base64"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/spf13/afero"
	"os"
	"time"
)

type Fs struct {
	cfg   *Config
	gitea *gitea.Client
}

func NewFs(cfg *Config) (afero.Fs, error) {
	client, err := NewGiteaClient(cfg.Url, cfg.User, cfg.Password)
	if err != nil {
		return nil, errors.New("giteafs.NewFs() %s ", err.Error())
	}
	fs := &Fs{cfg: cfg, gitea: client}
	return fs, nil
}

func NewGiteaClient(url, user, password string) (*gitea.Client, error) {
	client, err := gitea.NewClient(url, func(cli *gitea.Client) error {
		cli.SetBasicAuth(user, password)
		return nil
	})
	if err != nil {
		return client, err
	}
	return client, nil
}

// Decode
// 解码
func (s *Fs) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	_, err := base64.StdEncoding.Decode(dst, src)
	if err != nil {
		return nil, errors.New("base64 decode failed")
	}
	return dst, nil
}

// Encode
// 编码
func (s *Fs) Encode(src []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst, nil
}

func (s *Fs) AbsFileName(file string) string {
	return file
}

func (s *Fs) Create(name string) (afero.File, error) {
	return newFile(s, name, os.O_CREATE, 0666), nil
}

func (s *Fs) Mkdir(name string, perm os.FileMode) error {
	return nil
}

func (s *Fs) MkdirAll(path string, perm os.FileMode) error {
	return nil
}

func (s *Fs) Open(name string) (afero.File, error) {
	file := newFile(s, name, os.O_RDONLY, 0666)
	return file, nil
}

func (s *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	file := newFile(s, name, flag, perm)
	return file, nil
}

func (s *Fs) Remove(name string) error {
	return deleteFile(s, name, "")
}

func (s *Fs) RemoveAll(path string) error {
	opts := gitea.DeleteFileOptions{
		FileOptions: gitea.FileOptions{
			BranchName: s.cfg.Branch,
			Message:    "delete file " + path,
		},
	}
	_, err := s.gitea.DeleteFile(s.cfg.User, s.cfg.Repo, path, opts)
	return err
}

func (s *Fs) Rename(oldName, newName string) error {

	opt := gitea.UpdateFileOptions{
		FileOptions: gitea.FileOptions{
			BranchName: s.cfg.Branch,
			Message:    "rename file " + oldName + " to " + newName,
		},
		FromPath: oldName,
	}
	_, _, err := s.gitea.UpdateFile(s.cfg.User, s.cfg.Repo, newName, opt)
	return err
}

func (s *Fs) Stat(name string) (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Fs) Name() string {
	return Name()
}

func (s *Fs) Chmod(name string, mode os.FileMode) error {
	return nil
}

func (s *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	return nil
}
