package gitea

import (
	gitea "code.gitea.io/sdk/gitea"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader"
)

type Loader interface {
	loader.Loader
	Connect(url, userName, password string) error
	Init(name string, repo string, branch string)
}

type giteaLoader struct {
	name     string
	url      string
	userName string
	password string
	email    string
	repoName string
	branch   string
	client   *gitea.Client
}

func NewLoader(name, repo, branch string) Loader {
	return &giteaLoader{name: name, repoName: repo, branch: branch}
}

func (s *giteaLoader) Init(name string, repo string, branch string) {
	s.repoName = repo
	s.branch = branch
	s.name = name
}

func (s *giteaLoader) Connect(url, userName, password string) error {
	client, err := gitea.NewClient(url, func(cli *gitea.Client) error {
		cli.SetBasicAuth(userName, password)
		return nil
	})
	if err != nil {
		return err
	}
	s.client = client
	s.url = url
	s.userName = userName
	s.password = password
	return nil
}

func (s *giteaLoader) GetName() string {
	return s.name
}

func (s *giteaLoader) SetName(val string) {
	s.name = val
}

func (s *giteaLoader) GetFile(filepath string) ([]byte, error) {
	bytes, resp, err := s.client.GetFile(s.userName, s.repoName, s.branch, filepath, true)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("%s %s", filepath, err.Error()))
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("GetFile failed" + resp.Status)
	}
	return bytes, nil
}

func (s *giteaLoader) UpdateFile(filepath string, text string, message string) error {
	content := base64.StdEncoding.EncodeToString([]byte(text))
	mainLicense, _, err := s.client.GetContents(s.userName, s.repoName, s.branch, filepath)
	if err != nil {
		return err
	}

	cfo := gitea.UpdateFileOptions{}
	cfo.Content = content
	cfo.Message = message
	cfo.BranchName = s.branch
	cfo.SHA = mainLicense.SHA
	cfo.Committer = gitea.Identity{Name: s.userName, Email: s.email}
	cfo.Author = gitea.Identity{Name: s.userName, Email: s.email}

	fresp, resp, err := s.client.UpdateFile(s.userName, s.repoName, filepath, cfo)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New("CreateFile failed" + resp.Status)
	}
	fmt.Println(fresp)
	return err
}

func (s *giteaLoader) AbsFileName(file string) string {
	return file
}
