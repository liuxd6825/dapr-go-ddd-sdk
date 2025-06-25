package docx

import (
	"bytes"
	"fmt"
	"github.com/carmel/gooxml/document"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/spf13/afero"
	"io"
	"strings"
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
	reader := bytes.NewReader(data)
	return readFile(reader, reader.Size())
}

func readFile(r io.ReaderAt, size int64) (string, error) {
	sb := new(strings.Builder)
	doc, err := document.Read(r, size)
	if err != nil {
		return "", errors.New("打开文件失败: %v", err)
	}

	// 2. 提取段落文本
	for _, para := range doc.Paragraphs() {
		fmt.Println("=== 段落 ===")
		fmt.Println(para.X()) // 获取段落纯文本

		// 3. 提取带格式的文本（Run级别）
		for _, run := range para.Runs() {
			sb.WriteString(run.Text())
		}
	}

	// 4. 提取表格内容
	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			for _, cell := range row.Cells() {
				for _, para := range cell.Paragraphs() {
					for _, run := range para.Runs() {
						sb.WriteString(run.Text())
					}
				}
			}
		}
	}

	return sb.String(), nil
}
