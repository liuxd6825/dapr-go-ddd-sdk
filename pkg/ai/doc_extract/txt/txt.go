package txt

import (
	"bytes"
	"io/ioutil"
	"unicode/utf8"

	"github.com/saintfish/chardet"
	"github.com/spf13/afero"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadFile(fs afero.Fs, fileName string) (string, error) {
	data, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return "", err
	}
	return ConvertToUtf8Auto(data)
}

func ConvertToUtf8Auto(data []byte) (string, error) {
	// 1. 使用 chardet 检测编码
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(data)
	if err != nil {
		return "", err
	}

	// 2. 根据检测到的编码名称获取 Encoding 对象
	// 这里利用 x/net/html/charset 包来通过名字查找编码
	e, _ := charset.Lookup(result.Charset)
	if e == nil {
		// 如果没找到对应编码，默认当做 UTF-8 处理
		return ConvertContent(data), nil
	}

	// 3. 转换编码
	reader := transform.NewReader(bytes.NewReader(data), e.NewDecoder())
	utf8Data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(utf8Data), nil
}

func ConvertContent(data []byte) string {
	// 1. 先尝试判断是否已经是 UTF-8
	if utf8.Valid(data) {
		return string(data)
	}

	// 2. 如果不是 UTF-8，尝试按 GBK (GB18030) 解码
	decoder := simplifiedchinese.GB18030.NewDecoder()
	utf8Data, err := decoder.Bytes(data)
	if err != nil {
		// 如果转换也失败，只好返回原始内容（虽乱码但保留原数据）或报错
		return string(data)
	}

	return string(utf8Data)
}
