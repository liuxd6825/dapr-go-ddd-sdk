package xls

import (
	"bytes"
	"fmt"
	"github.com/extrame/xls"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"strings"
)

type Reader struct {
	logger *logrus.Logger
}

func NewReader(logger *logrus.Logger) *Reader {
	return &Reader{
		logger: logger,
	}
}

func (r *Reader) ReadFile(fs afero.Fs, fileName string) (string, error) {
	data, err := afero.ReadFile(fs, fileName)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(data)
	sb := strings.Builder{}
	if xlFile, err := xls.OpenReader(reader, "utf-8"); err == nil {
		fmt.Println(xlFile.Author)
		//第一个sheet
		sheet := xlFile.GetSheet(0)
		if sheet.MaxRow != 0 {
			for i := 0; i < int(sheet.MaxRow); i++ {
				row := sheet.Row(i)
				if row.LastCol() > 0 {
					for j := 0; j < row.LastCol(); j++ {
						col := row.Col(j)
						sb.WriteString(col + ",")
					}
					sb.WriteString("\n")
				}
			}
		}
	} else {
		r.logger.Errorf("open %s error:%s", fileName, err)
	}
	return sb.String(), err
}
