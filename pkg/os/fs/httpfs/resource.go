package httpfs

import (
	"errors"
	"os"
)

func readFile(fs *Fs, filepath string) ([]byte, error) {
	url := fs.absFile(filepath)
	resp, err := fs.client.R().EnableTrace().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.Status())
	}
	return resp.Body(), nil
}

func writeFile(fs *Fs, filepath string, data []byte, message string) error {
	return nil
}

func deleteFile(fs *Fs, filepath string, message string) error {
	return nil
}

func create(fs *Fs, filename string, data []byte, message string) error {
	return nil
}

func readDir(fs *Fs, filename string, count int) (res []os.FileInfo, err error) {
	return nil, nil
}

func rename(fs *Fs, oldName, newName string) error {
	return nil
}

func mkdir(fs *Fs, path string, message string) error {
	return nil
}

func Name() string {
	return "http"
}
