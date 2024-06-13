package server

import "github.com/dop251/goja"

type Service struct {
	Name    string
	Desc    string
	server  *Server
	Service *goja.Object
}
