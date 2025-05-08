package runtimeutils

import (
	"bytes"
	"errors"
	"path"
	"runtime"
	"strconv"
)

var (
	goroutinePrefix = []byte("goroutine ")
	errBadStack     = errors.New("invalid runtime.Stack output")
)

// GetFuncName
// @Description: 获取当前运行的方法名称
// @param skip  调用方法的调转级数
// @return string 返回方法名称
func GetFuncName(skip int) string {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+1, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

func GetPackageName(skip int) string {
	_, filename, _, _ := runtime.Caller(skip + 1)
	return path.Base(path.Dir(filename))
}

// GoId This is terrible, slow, and should never be used.
func GoId() (int, error) {
	buf := make([]byte, 32)
	n := runtime.Stack(buf, false)
	buf = buf[:n]
	// goroutine 1 [running]: ...

	buf, ok := bytes.CutPrefix(buf, goroutinePrefix)
	if !ok {
		return 0, errBadStack
	}

	i := bytes.IndexByte(buf, ' ')
	if i < 0 {
		return 0, errBadStack
	}

	return strconv.Atoi(string(buf[:i]))
}
