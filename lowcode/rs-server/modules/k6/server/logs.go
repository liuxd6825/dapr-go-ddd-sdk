package server

import "fmt"

type Logs struct {
	Name string
}

func NewLogs() *Logs {
	return &Logs{Name: "k6-logs"}
}
func (c *Logs) Info(args ...interface{}) {
	fmt.Print(args...)
}

func (c *Logs) Logf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}

func (c *Logs) Error(args ...interface{}) {
	fmt.Print(args...)
}
