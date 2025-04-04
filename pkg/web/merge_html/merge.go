package merge_html

import (
	"bytes"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"
)

type fileInfo struct {
	fullName string
	fileInfo os.FileInfo
}

func MergeFile(fs afero.Fs, root string, searchFile string, saveFileName string) error {
	files, err := getFiles(fs, searchFile, root)
	if err != nil {
		return err
	}
	buf := bytes.Buffer{}
	buf.WriteString("var w=window;if(!w.$h){w.$h={}}\n")
	for _, file := range files {
		data, err := afero.ReadFile(fs, file.fullName)
		if err != nil {
			return err
		}
		fileName := strings.TrimPrefix(file.fullName, root)
		buf.WriteString("w.$h[\"")
		buf.WriteString(fileName)
		buf.WriteString("\"]=`")
		buf.Write(escapeJSONString(data))
		buf.WriteString("`\n")
	}
	return afero.WriteFile(fs, saveFileName, buf.Bytes(), os.ModePerm)
}

// removeWhitespace 删除所有 \n 并将多个空格压缩成一个
func removeWhitespace(s string) string {
	var builder strings.Builder
	builder.Grow(len(s)) // 预分配内存，提高性能

	prevSpace := false // 标记前一个字符是否是空格

	for _, ch := range s {
		if ch == '\n' || ch == '\t' {
			continue // 直接跳过换行符
		}

		if ch == '`' {
			builder.Write([]byte("\\`"))
			continue
		}

		if ch == '$' {
			builder.Write([]byte("\\$"))
			continue
		}

		if ch == ' ' {
			if !prevSpace { // 如果前一个字符不是空格，才写入
				builder.WriteRune(' ')
				prevSpace = true
			}
			continue
		}

		// 普通字符直接写入
		builder.WriteRune(ch)
		prevSpace = false
	}

	return builder.String()
}

// escapeJSONString 处理字符串以确保JSON格式正确
func escapeJSONString(s []byte) []byte {
	// 使用json.Marshal来正确转义字符串
	str := removeWhitespace(string(s))
	return []byte(str)
	/*
		b, err := json.Marshal(str)
		if err != nil {
			panic(err) // 如果出错返回空字符串
		}
		// 去除首尾的引号，因为Marshal会在字符串外添加引号
		return b[1 : len(b)-1]

	*/
}

func getFiles(fs afero.Fs, searchFile, root string) (files []fileInfo, err error) {
	list, err := afero.ReadDir(fs, root)
	if err != nil {
		return nil, err
	}
	for _, file := range list {
		if file.IsDir() {
			fl, er := getFiles(fs, searchFile, filepath.Join(root, file.Name()))
			if er != nil {
				return nil, er
			}
			files = append(files, fl...)
		} else if strings.HasSuffix(file.Name(), searchFile) {
			f := fileInfo{fullName: filepath.Join(root, file.Name()), fileInfo: file}
			files = append(files, f)
		}
	}
	return files, err
}
