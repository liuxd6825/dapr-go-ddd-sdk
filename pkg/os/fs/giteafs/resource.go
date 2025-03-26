package giteafs

import (
	"code.gitea.io/sdk/gitea"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/os/fs/fsopts"
	"os"
)

func getInfo(fs *Fs, name string) (*gitea.ContentsResponse, error) {
	content, _, err := fs.gitea.GetContents(fs.cfg.User, fs.cfg.Repo, fs.cfg.Branch, name)
	if err != nil {
		return nil, err
	}
	return content, err
}

func readFile(fs *Fs, filepath string) ([]byte, error) {
	opts := fs.cfg
	bytes, resp, err := fs.gitea.GetFile(opts.User, opts.Repo, opts.Branch, filepath, true)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", filepath, err.Error()))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("GetFile failed" + resp.Status)
	}
	return bytes, nil
}

func writeFile(fs *Fs, filepath string, data []byte, message string) error {

	info, err := getInfo(fs, filepath)
	if err != nil {
		return err
	}

	content := base64.StdEncoding.EncodeToString(data)

	cfo := gitea.UpdateFileOptions{}
	cfo.Content = content
	cfo.Message = message
	cfo.BranchName = fs.cfg.Branch
	cfo.SHA = info.SHA
	cfo.Committer = gitea.Identity{Name: fs.cfg.User, Email: fs.cfg.Email}
	cfo.Author = gitea.Identity{Name: fs.cfg.User, Email: fs.cfg.Email}

	_, resp, err := fs.gitea.UpdateFile(fs.cfg.User, fs.cfg.Repo, filepath, cfo)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("CreateFile failed" + resp.Status)
	}
	return err
}

func deleteFile(fs *Fs, filepath string, message string) error {
	opts := gitea.DeleteFileOptions{
		FileOptions: gitea.FileOptions{
			BranchName: fs.cfg.Branch,
			Message:    "delete file " + filepath,
		},
	}
	_, err := fs.gitea.DeleteFile(fs.cfg.User, fs.cfg.Repo, filepath, opts)
	return err
}

func create(fs *Fs, filename string, data []byte, message string) error {
	opts := gitea.CreateFileOptions{
		FileOptions: gitea.FileOptions{
			Message:    "create file",
			BranchName: fs.cfg.Branch,
		},
		Content: string(data),
	}

	_, _, err := fs.gitea.CreateFile(fs.cfg.User, fs.cfg.Repo, filename, opts)
	return err
}

func readDir(fs *Fs, filename string, count int) (res []os.FileInfo, err error) {
	list, _, err := fs.gitea.ListContents(fs.cfg.User, fs.cfg.Repo, fs.cfg.Branch, filename)
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		fileInfo := &fsopts.FileInfo{}
		fileInfo.SetName(item.Name)
		fileInfo.SetSize(item.Size)
		if item.Type == "file" {
			fileInfo.SetIsDir(false)
			res = append(res, fileInfo)
		} else if item.Type == "dir" {
			fileInfo.SetIsDir(true)
			res = append(res, fileInfo)
		}
		if len(res) == count {
			break
		}
	}
	return res, nil
}

func rename(fs *Fs, oldName, newName string) error {
	opt := gitea.UpdateFileOptions{
		FileOptions: gitea.FileOptions{
			BranchName: fs.cfg.Branch,
			Message:    "rename file " + oldName + " to " + newName,
		},
		FromPath: oldName,
	}
	_, _, err := fs.gitea.UpdateFile(fs.cfg.User, fs.cfg.Repo, newName, opt)
	return err
}

func mkdir(fs *Fs, path string, message string) error {
	_, _, err := fs.gitea.CreateFile(fs.cfg.User, fs.cfg.Repo, path, gitea.CreateFileOptions{
		FileOptions: gitea.FileOptions{
			Message:    "create directory " + path,
			BranchName: fs.cfg.Branch,
		},
		Content: "",
	})
	if err != nil {
		return err
	}
	return err
}

func Name() string {
	return "gitea"
}
