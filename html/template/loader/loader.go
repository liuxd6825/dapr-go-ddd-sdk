package loader

import (
	"errors"
	"fmt"
)

type Loader interface {
	GetName() string
	AbsFileName(file string) string
	GetFile(filepath string) ([]byte, error)
	UpdateFile(filepath string, text string, message string) error
}

var loaders map[string]Loader

func init() {
	loaders = make(map[string]Loader)
}

func GetLoader(name string) (Loader, bool) {
	item, ok := loaders[name]
	return item, ok
}

func AddLoader(loader Loader) error {
	name := loader.GetName()
	if _, ok := loaders[name]; ok {
		return errors.New(fmt.Sprintf("file server name %s ", name))
	}
	loaders[name] = loader
	return nil
}
