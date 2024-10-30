package fs

type PathFs interface {
	BasePath() string
}

type Decode interface {
	Decode(src []byte) ([]byte, error)
	Encode(src []byte) ([]byte, error)
}
