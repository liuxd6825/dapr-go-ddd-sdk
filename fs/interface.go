package fs

import "github.com/liuxd6825/dapr-go-ddd-sdk/fs/fsopts"

type PathFs interface {
	BasePath() string
}

type FS interface {
	GetFsType() string
}

type Decode interface {
	Decode(src []byte) ([]byte, error)
	Encode(src []byte) ([]byte, error)
}

type Reader interface {
	ReadFile(filename string, opts ...*fsopts.Options) ([]byte, error)
}
