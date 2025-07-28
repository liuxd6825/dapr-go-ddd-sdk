package idutils

// import gonanoid "github.com/matoous/go-nanoid/v2"
// import github.com/oklog/ulid/v2
import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	ulid "github.com/oklog/ulid/v2"
)

const NanoCode = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const NanoLength = 24

func NewId() string {
	id, _ := NewNId()
	return id
}

func NewNId() (string, error) {
	return gonanoid.Generate(NanoCode, NanoLength)
}

func NewUlid() (string, error) {
	return ulid.Make().String(), nil
}
