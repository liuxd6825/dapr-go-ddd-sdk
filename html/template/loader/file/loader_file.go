package file

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/html/template/loader"
	"io/ioutil"
	"os"
)

type Loader interface {
	loader.Loader
	Connect() error
	Init(name string, path string)
}

type loaderFile struct {
	name     string
	userName string
	password string
	email    string
	path     string
}

func NewLoader(name, path string) Loader {
	return &loaderFile{name: name, path: path}
}

func (s *loaderFile) Init(name string, path string) {
	s.path = path
	s.name = name
}

func (s *loaderFile) GetName() string {
	return s.name
}

func (s *loaderFile) GetFile(filepath string) ([]byte, error) {
	return ioutil.ReadFile(s.path + "/" + filepath)
}

func (s *loaderFile) UpdateFile(filepath string, text string, message string) error {
	return nil
}

func (s *loaderFile) Connect() error {
	if ok := s.isDir(s.path); !ok {
		return errors.New("")
	}
	return nil
}

func (s *loaderFile) isDir(path string) bool {
	p, err := os.Stat(path)
	if err != nil {

		return false
	}
	return p.IsDir()
}

func (s *loaderFile) AbsFileName(file string) string {
	return s.path + "/" + file
}
