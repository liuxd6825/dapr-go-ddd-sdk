package httpfs

import (
	gitea "code.gitea.io/sdk/gitea"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"os"
	"time"
)

type Fs struct {
	options *GiteaOptions
	client  *gitea.Client
}

type GiteaOptions struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	Repo     string `yaml:"repos"`
	Branch   string `yaml:"branch"`
}

func newFs(options *GiteaOptions) afero.Fs {
	return &Fs{options: options}
}

func (f *Fs) Connect() error {
	client, err := gitea.NewClient(f.options.Url, func(cli *gitea.Client) error {
		cli.SetBasicAuth(f.options.User, f.options.Password)
		return nil
	})
	if err != nil {
		return err
	}
	f.client = client
	return nil
}

func (f *Fs) readFile(filepath string) ([]byte, error) {
	opts := f.options
	bytes, resp, err := f.client.GetFile(opts.User, opts.Repo, opts.Branch, filepath, true)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", filepath, err.Error()))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("GetFile failed" + resp.Status)
	}
	return bytes, nil
}

func (f *Fs) writeFile(filepath string, text string, message string) error {
	content := base64.StdEncoding.EncodeToString([]byte(text))
	mainLicense, _, err := f.client.GetContents(f.options.User, f.options.Repo, f.options.Branch, filepath)
	if err != nil {
		return err
	}

	cfo := gitea.UpdateFileOptions{}
	cfo.Content = content
	cfo.Message = message
	cfo.BranchName = f.options.Branch
	cfo.SHA = mainLicense.SHA
	cfo.Committer = gitea.Identity{Name: f.options.User, Email: f.options.Email}
	cfo.Author = gitea.Identity{Name: f.options.User, Email: f.options.Email}

	fresp, resp, err := f.client.UpdateFile(f.options.User, f.options.Repo, filepath, cfo)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("CreateFile failed" + resp.Status)
	}
	fmt.Println(fresp)
	return err
}

func (f *Fs) GetName() string {
	return f.options.Name
}

func (f *Fs) SetName(val string) {
	f.options.Name = val
}

func (f *Fs) AbsFileName(file string) string {
	return file
}

func (f *Fs) Create(name string) (afero.File, error) {
	return newFile(f, name), nil
}

func (f *Fs) Mkdir(name string, perm os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) MkdirAll(path string, perm os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Open(name string) (afero.File, error) {
	file := newFile(f, name)
	return file, nil
}

func (f *Fs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	file := newFile(f, name)
	return file, nil
}

func (f *Fs) Remove(name string) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) RemoveAll(path string) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Rename(oldname, newname string) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Stat(name string) (os.FileInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Name() string {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Chmod(name string, mode os.FileMode) error {
	//TODO implement me
	panic("implement me")
}

func (f *Fs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	//TODO implement me
	panic("implement me")
}
