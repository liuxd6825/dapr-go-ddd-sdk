package xlsx

import (
	"bytes"
	"fmt"
	"github.com/carmel/gooxml/spreadsheet"
	"github.com/spf13/afero"
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
	f, err := spreadsheet.Read(reader, reader.Size())
	if err != nil {
		return "", err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sb := strings.Builder{}
	for _, sheet := range f.Sheets() {
		sb.WriteString(fmt.Sprintf("*%s\n", sheet.Name()))
		sb.WriteString("```csv\n")
		for _, row := range sheet.Rows() {
			for _, cell := range row.Cells() {
				sb.WriteString(getCellString(cell) + ",")
			}
			sb.WriteString("\n")
		}
		sb.WriteString("```\n")
	}
	return sb.String(), nil
}

func getCellString(cell spreadsheet.Cell) string {
	val := cell.GetString()
	val = strings.Replace(val, "\"", "\"\"", -1)
	return "\"" + val + "\""
}
