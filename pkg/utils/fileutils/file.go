package fileutils

import (
	"io/ioutil"
	"os"
	"strings"
)

func IsExist(pathOrFile string) bool {
	_, err := os.Stat(pathOrFile)
	return err == nil || os.IsExist(err)
}

func GetFileInfos(path string, extFileName string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var resList []os.FileInfo
	extName := strings.ToLower(extFileName)
	for _, f := range files {
		fileName := strings.ToLower(f.Name())
		if !f.Mode().IsDir() {
			if extName == "" || strings.HasSuffix(fileName, extName) {
				resList = append(resList, f)
			}
		}
	}
	return resList, nil
}

func GetFileCount(path string, extFileName string) int {
	files, err := GetFileInfos(path, extFileName)
	if err != nil {
		return 0
	}
	return len(files)
}
