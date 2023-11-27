package singleutils

import (
	"fmt"
	"testing"
	"time"
)

type server struct {
	time time.Time
}

func Test_Create(t *testing.T) {
	s1 := newServer()
	s1.print()

	s2 := newServer()
	s2.print()

	s3 := newServer()
	s3.print()
}

func newServer() *server {
	return Create[*server]("10D157C3-0B57-425D-827A-173CCEADBE40", func() *server {
		return &server{time: time.Now()}
	})
}

func (s *server) print() {
	fmt.Printf("%v\n", s.time)
}
