package restapp

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	noExtExeName = "" //当前应用的名称(无扩展名)
	exeName      = "" //当前应用的名称
	pathName     = "" // 当前应用路径
	pid          = "" // 当前进程PID
	envName      = ""
	workPath     = ""
	exeFullName  = ""
	startPath    = "" // 启动目录
)
var (
	AppTitle  = "" // 应用名称
	Version   = "" // 应用版本号
	BuildTime = "" // 编译时间
	GitHead   = "" // Git
	sysPaths  *types.CMap[string]
)

func init() {
	val := int64(os.Getpid())
	pid = strconv.FormatInt(val, 10)

	path, _ := os.Executable()
	pathName, exeName = filepath.Split(path)
	SetExeName(exeName)

	exeFullName = path

	sysPaths = types.NewCMap[string]()
	sysPaths.Set("ExeName", exeName)

	wPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workPath = wPath
	startPath = wPath
	sysPaths.Set("WorkPath", workPath)
}

func GetSysPaths() *types.CMap[string] {
	return sysPaths
}

func GetPathName() string {
	return pathName
}

func GetEnvName() string {
	return envName
}

func SetEnvName(val string) {
	envName = val
}

func GetExeFullName() string {
	return exeFullName
}

func GetStartPath() string {
	return startPath
}

// GetPID
//
//	@Description: 当前进程PID
//	@return string
func GetPID() string {
	return pid
}

// GetExeName
//
//	@Description: 取可执行文件名称
//	@return string
func GetExeName() string {
	return exeName
}

// GetNoExtExeName
//
//	@Description: 取无扩展名的应用程序名称
//	@return string
func GetNoExtExeName() string {
	return noExtExeName
}

func GetWorkPath() string {
	return workPath
}

func SetExeName(name string) {
	exeName = strings.ReplaceAll(name, "___", "")
	noExtExeName = exeName
	idx := strings.Index(exeName, ".")
	if idx > 0 {
		noExtExeName = exeName[:idx]
	}
}

// AbsFileName
//
//	@Description: 取绝对路径
//	@param val
//	@return string
func AbsFileName(filename string) string {
	if filename == "" {
		return ""
	}
	filename = strings.ReplaceAll(filename, "${ExeName}", GetExeName())
	filename = strings.ReplaceAll(filename, "${PID}", GetPID())
	filename = strings.ReplaceAll(filename, "${EnvName}", GetEnvName())
	for _, k := range sysPaths.Keys() {
		if v, ok := sysPaths.Get(k); ok {
			key := fmt.Sprintf("${%s}", k)
			filename = strings.ReplaceAll(filename, key, v)
		}
	}
	filename, _ = filepath.Abs(filename)
	return filename
}

// ReplaceSysValues
//
//	@Description: 取绝对路径
//	@param val
//	@return string
func ReplaceSysValues(filename string) string {
	if filename == "" {
		return ""
	}
	filename = strings.ReplaceAll(filename, "${ExeName}", GetExeName())
	filename = strings.ReplaceAll(filename, "${PID}", GetPID())
	filename = strings.ReplaceAll(filename, "${EnvName}", GetEnvName())
	for _, k := range sysPaths.Keys() {
		if v, ok := sysPaths.Get(k); ok {
			key := fmt.Sprintf("${%s}", k)
			filename = strings.ReplaceAll(filename, key, v)
		}
	}
	return filename
}
