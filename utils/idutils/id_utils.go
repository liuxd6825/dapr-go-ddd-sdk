package idutils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

const NanoCode = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const NanoLength = 30

func NewId() string {
	id, _ := NewNId()
	return id
}

func NewNId() (string, error) {
	return gonanoid.Generate(NanoCode, NanoLength)
}
