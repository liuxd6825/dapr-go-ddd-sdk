package txt

import (
	"github.com/spf13/afero"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadFile(fs afero.Fs, fileName string) (string, error) {
	data, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
