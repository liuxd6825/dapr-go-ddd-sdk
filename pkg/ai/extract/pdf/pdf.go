package pdf

import (
	"bytes"
	"github.com/ledongthuc/pdf"
	"github.com/spf13/afero"
	"strings"
)

type Reader struct {
}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadFile(fs afero.Fs, fileName string) (string, error) {
	sb := new(strings.Builder)
	data, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return "", err
	}
	// 创建ReadSeeker
	readSeeker := bytes.NewReader(data)
	// 2. 创建 PDF Reader
	pdfReader, err := pdf.NewReader(readSeeker, readSeeker.Size())
	if err != nil {
		panic(err)
	}

	// 3. 获取 PDF 页数
	numPages := pdfReader.NumPage()

	// 4. 逐页提取文本
	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() || page.V.Key("Contents").Kind() == pdf.Null {
			continue
		}

		rows, _ := page.GetTextByRow()
		for _, row := range rows {
			for _, word := range row.Content {
				sb.WriteString(word.S)
			}
			sb.WriteString("\n")
		}
	}
	return sb.String(), nil
}
