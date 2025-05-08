package fileutils

import "strings"
import "path/filepath"

func GetBasePath(basePath ...string) string {
	if len(basePath) > 0 {
		return basePath[0]
	}
	return ""
}

// AbsPath
//
//	@Description: 取绝对路径
//	@param pwd 当前路径
//	@param filename 待相对路径的文件名， 如：../c.txt 或 ./c.txt
//	@return string 绝对路径
func AbsPath(filename string, opts ...*ReadOptions) string {
	opt := newReadOption(opts...)
	pwd := GetBasePath(opt.WorkPath)
	if strings.HasPrefix(filename, "/") {
		return opt.RootPath + filename
	}

	if strings.HasPrefix(filename, pwd) {
		return filename
	}

	if !strings.Contains(filename, "../") && !strings.Contains(filename, "./") {
		return filepath.Join(pwd, filename)
	}

	var res []string
	paths := splitFilePath(pwd)
	names := splitFilePath(filename)
	ok := false
	start := 0
	for i, name := range names {
		if len(paths) > 0 {
			if name == ".." {
				start = i
				paths = paths[:len(paths)-1]
				ok = true
			} else if name == "." {
				start = i
				ok = true
				continue
			} else if name == "" {
				continue
			} else {
				break
			}
		}
	}
	if ok {
		names = names[start+1:]
		res = append(paths, names...)
	} else {
		res = names
	}

	return strings.Join(res, "/")
}

func splitFilePath(filename string) []string {
	var res []string
	paths := strings.Split(filename, "/")
	if paths[len(paths)-1] == "" {
		paths = paths[:len(paths)-1]
	}
	for _, path := range paths {
		res = append(res, path)
	}
	return res
}
