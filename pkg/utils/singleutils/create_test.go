package singleutils

import (
	"fmt"
	"testing"
	"time"
)

type server struct {
	time time.Time
}

type Server interface {
	Print()
}

func Test_Create(t *testing.T) {
	s1 := newServer()
	s1.Print()

	s2 := newServer()
	s2.Print()

	s3 := newServer()
	s3.Print()
}

func Test_GetTypeName(t *testing.T) {
	s1 := getTypeName[*server]()
	t.Log(s1)

	s2 := getTypeName[Server]()
	t.Log(s2)
}

func getTypeName[T any]() string {
	var null T
	return GetTypeName(null)
}

func newServer() *server {
	return CreateObj[*server](func() *server {
		return &server{time: time.Now()}
	})
}

func (s *server) Print() {
	fmt.Printf("%v\n", s.time)
}
