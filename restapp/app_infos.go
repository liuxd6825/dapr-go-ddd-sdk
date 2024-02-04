package restapp

import (
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
)
var (
	AppTitle  = "" // 应用名称
	Version   = "" // 应用版本号
	BuildTime = "" // 编译时间
	GitHead   = "" // Git
)

func init() {
	val := int64(os.Getpid())
	pid = strconv.FormatInt(val, 10)

	path, _ := os.Executable()
	pathName, exeName = filepath.Split(path)
	SetExeName(exeName)
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

// GetPID
//
//	@Description: 当前进程PID
//	@return string
func GetPID() string {
	return pid
}

// GetExeName
//
//	@Description: 取应用程序名称
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
func AbsFileName(val string) string {
	if val == "" {
		return ""
	}
	val = strings.ReplaceAll(val, "${AppName}", GetExeName())
	val = strings.ReplaceAll(val, "${PID}", GetPID())
	val = strings.ReplaceAll(val, "${EnvName}", GetEnvName())
	val, _ = filepath.Abs(val)
	return val
}
