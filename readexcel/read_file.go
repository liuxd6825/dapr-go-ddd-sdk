package readexcel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/script"
	"github.com/tealeg/xlsx"
	"strings"

	"io"
	"os"
	"strconv"
	"time"
)

type FieldError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func NewFieldError(field string, error string) error {
	return &FieldError{field, error}
}

func (e *FieldError) Error() string {
	return fmt.Sprintf("field:%s, error:%s", e.Field, e.Msg)
}

func ReadFile(ctx context.Context, fileName string, sheetName string, temp *Template) (*DataTable, error) {
	bs, err := readFile(fileName)
	if err != nil {
		return nil, err
	}
	return ReadBytes(ctx, bytes.NewBuffer(bs), sheetName, temp)
}

func ReadBytes(ctx context.Context, buffer *bytes.Buffer, sheetName string, temp *Template) (*DataTable, error) {
	if err := temp.Init(); err != nil {
		return nil, err
	}

	f, err := xlsx.OpenBinary(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	sheet, err := getSheet(f, sheetName)
	if err != nil {
		return nil, err
	}
	if sheet == nil {
		return nil, errors.New("sheetName is null")
	}

	rows := NewRowBySheet(sheet)
	return readBytes(ctx, rows, temp)
}

func ReadByMap(ctx context.Context, list []map[string]any, temp *Template) (*DataTable, error) {
	if err := temp.Init(); err != nil {
		return nil, err
	}

	rows := NewRowByMap(list)
	return readBytes(ctx, rows, temp)
}

func readBytes(ctx context.Context, rows Rows, temp *Template) (*DataTable, error) {
	runtime, err := newRuntime()
	if err != nil {
		return nil, err
	}

	table := NewDataTable(temp)
	countRow := rows.Count()

	iRow := temp.Heads[0].RowNum + 1
	var index int64 = 0
	for {
		if iRow >= countRow {
			break
		}
		dataRow := NewDataRow()
		cells := rows.Row(int(iRow))
		if cells == nil {
			return table, nil
		}

		rowCells := make(map[string]*DataCell)
		for _, rowColumn := range temp.Heads {
			iRow++
			for _, mapColumn := range rowColumn.Columns {
				col := mapColumn.GetColNum()
				count := len(cells)
				if col < int64(count) {
					cell := cells[col]
					if cell != nil {
						cellValue := getCellString(cell.String())

						if isTime(cell) {
							t, err := cell.GetTime(false)
							if err != nil {
								t = getExcelDate(cell.Value())
							}

							switch mapColumn.DataType {
							case Date:
								cellValue = t.Format("2006-01-02")
								break
							case Time:
								cellValue = t.Format("2006-01-02 15:04:05.000")
								break
							default:
								cellValue = t.Format("2006-01-02 15:04:05.000")
								break
							}
						}
						cellValue = getCellString(cellValue)
						rowCells[mapColumn.Key] = &DataCell{
							Key:    mapColumn.Key,
							Value:  cellValue,
							Row:    iRow,
							Col:    col,
							Errors: nil,
						}
					}
				}
			}
		}

		dataRow.RowNum = iRow
		table.AddRow(dataRow)

		getValueOptions := &KeyValueOptions{
			Temp:     temp,
			DataRow:  dataRow,
			RowCells: rowCells,
		}
		runtime.ClearInterrupt()
		for _, field := range temp.Fields {
			var value any = nil
			cellValues := make(map[string]any, 0)
			for _, mapKey := range field.MapKeys {
				dataRow.AddCell(field.Name, rowCells[mapKey])
				cellValue := temp.GetKeyValue(mapKey)
				if cellValue != nil {
					cellValues[mapKey] = ""
					str := getCellString(cellValue.GetValue(getValueOptions))
					if len(str) > 0 {
						cellValues[mapKey] = str
					}
				}
			}

			var v any
			var err error = nil
			if len(field.Script) == 0 { //没有脚本模式，直接拼字符串
				var s = ""
				for _, v := range cellValues {
					if v, ok := v.(string); ok {
						s += v
					}
				}
				value = s
			} else if value, err = script.RunScript(runtime, cellValues, &field.Script); err != nil {
				dataRow.AddValue(field.Name, err.Error())
				dataRow.AddError(field.Name, err)
				logs.Errorf(ctx, "readexcel.ReadBytes() rowNum:%v, fieldName:%s, script:%s; error:%s； ", iRow, field.Name, field.Script, err.Error())
				continue
			}

			switch field.DataType {
			case Date:
				v, err = ValueToDate(value, nil)
			case DateTime:
				v, err = ValueToDate(value, nil)
			case Time:
				v, err = ValueToDate(value, nil)
			case Integer:
				v, err = ValueToInt(value, nil)
			case Money:
				v, err = ValueToFloat(value, nil)
			default:
				v, err = ValueToString(value, "")
			}
			dataRow.AddValue(field.Name, v)
			dataRow.AddError(field.Name, err)
		}

		if dataRow.HasError() {
			table.AddError(index)
		}
		index++
	}
	return table, nil
}

func getCellString(str string) string {
	return strings.ReplaceAll(str, "\t", "")
}

func GetSheetNames(buffer *bytes.Buffer) ([]string, error) {
	f, err := xlsx.OpenBinary(buffer.Bytes())
	if err != nil {
		return nil, err
	}
	var names []string
	for _, s := range f.Sheets {
		names = append(names, s.Name)
	}
	return names, nil
}

func isTime(cell Cell) bool {
	if cell.IsTime() || cell.Type() == xlsx.CellTypeDate {
		return true
	}
	if cell.Type() == xlsx.CellTypeNumeric && strings.Contains(cell.NumFmt(), "yy") {
		return true
	}
	return false
}

func getSheet(file *xlsx.File, sheetName string) (*xlsx.Sheet, error) {
	var sheet *xlsx.Sheet
	if len(sheetName) == 0 {
		if len(file.Sheets) > 0 {
			sheet = file.Sheets[0]
		} else {
			return nil, errors.New("sheetName is empty")
		}
	} else {
		s, ok := file.Sheet[sheetName]
		if !ok {
			return nil, errors.New(fmt.Sprintf("sheet not found '%v'", sheetName))
		}
		sheet = s
	}
	return sheet, nil
}

func readFile(fileName string) ([]byte, error) {
	//获得一个file
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("read fail")
		return nil, err
	}

	//把file读取到缓冲区中
	defer f.Close()
	var chunk []byte
	buf := make([]byte, 1024)

	for {
		//从file读取到buf中
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		//说明读取结束
		if n == 0 {
			break
		}
		//读取到最终的缓冲区中
		chunk = append(chunk, buf[:n]...)
	}

	return chunk, nil
}

func getExcelDate(excelStr string) time.Time {
	eTime := time.Date(1899, time.December, 30, 0, 0, 0, 0, time.UTC)
	var days, _ = strconv.Atoi(excelStr)
	return eTime.AddDate(0, 0, days)
	//return eTime.Add(time.Second * time.Duration(days*86400))
}

func ReviewFile(fileName string, sheetName string, maxRows int64) (*Review, error) {
	bytes, err := readFile(fileName)
	if err != nil {
		return nil, err
	}
	return ReadBytesToMap(bytes, sheetName, maxRows)
}

func ReadBytesToMap(bytes []byte, sheetName string, maxRows int64) (*Review, error) {
	f, err := xlsx.OpenBinary(bytes)
	if err != nil {
		return nil, err
	}

	sheet, err := getSheet(f, sheetName)
	if err != nil {
		return nil, err
	}
	if sheet == nil {
		return nil, errors.New("sheetName is null")
	}
	review := &Review{}
	for c := 0; c < len(sheet.Cols); c++ {
		review.AddColumns(GetCellLabel(c + 1))
	}
	for _, s := range f.Sheets {
		review.AddSheetName(&Sheet{Name: s.Name, MaxCol: s.MaxCol, MaxRow: s.MaxRow})
	}
	review.OpenSheet = sheet.Name
	for rIdx := 0; rIdx < sheet.MaxRow; rIdx++ {
		row := sheet.Row(rIdx)
		item := ReviewItem{}
		item["$row"] = rIdx
		for cIdx := 0; cIdx < len(row.Cells); cIdx++ {
			key := GetCellLabel(cIdx + 1)
			if cell := row.Cells[cIdx]; cell != nil {
				item[key] = row.Cells[cIdx].String()
			}
		}
		review.AddItems(item)
		if maxRows > 0 && int64(rIdx) > maxRows {
			break
		}
	}
	return review, nil

}
