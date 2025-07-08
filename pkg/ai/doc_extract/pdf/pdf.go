package pdf

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"github.com/spf13/afero"
	"sort"
	"strings"
)

/*
# 你是一个golang工程，负责将PDF文件转换为Markdown格式。
### 要求
你的任务是读取PDF文件，提取文本内容，并将其转换为Markdown格式。你需要处理不同的文本块，包括标题、段落、列表等，并确保输出的Markdown格式正确。
*/

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
	readSeeker := bytes.NewReader(data)
	pdfReader, err := pdf.NewReader(readSeeker, readSeeker.Size())
	if err != nil {
		return "", err
	}

	numPages := pdfReader.NumPage()
	var markdown bytes.Buffer
	fontSizes := collectFontSizes(pdfReader, numPages)
	header1, header2 := getHeaderThresholds(fontSizes)

	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() || page.V.Key("Contents").Kind() == pdf.Null {
			continue
		}
		texts := page.Content().Text
		if len(texts) == 0 {
			continue
		}
		// 排序
		sort.SliceStable(texts, func(i, j int) bool {
			if texts[i].Y != texts[j].Y {
				return texts[i].Y > texts[j].Y
			}
			return texts[i].X < texts[j].X
		})
		blocks := mergeTextBlocks(texts)
		if pageNum > 1 {
			markdown.WriteString("<!--break-->\n")
		}
		writeMarkdownBlocks(&markdown, blocks, header1, header2)
	}
	return markdown.String(), nil
}

// 收集所有字体大小
func collectFontSizes(pdfReader *pdf.Reader, numPages int) []float64 {
	var sizes []float64
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() || page.V.Key("Contents").Kind() == pdf.Null {
			continue
		}
		for _, t := range page.Content().Text {
			if t.FontSize > 0 {
				sizes = append(sizes, t.FontSize)
			}
		}
	}
	return sizes
}

// 计算标题阈值
func getHeaderThresholds(sizes []float64) (float64, float64) {
	if len(sizes) == 0 {
		return 14, 12
	}
	sort.Float64s(sizes)
	max := sizes[len(sizes)-1]
	mid := sizes[len(sizes)/3]
	return max, mid
}

// 合并同一行文本
func mergeTextBlocks(texts []pdf.Text) []pdf.Text {
	var blocks []pdf.Text
	var currentY float64 = -1
	var currentLine strings.Builder
	var currentFont float64
	for _, text := range texts {
		if currentY != text.Y {
			if currentLine.Len() > 0 {
				blocks = append(blocks, pdf.Text{S: currentLine.String(), FontSize: currentFont, Y: currentY})
				currentLine.Reset()
			}
			currentY = text.Y
			currentFont = text.FontSize
		}
		currentLine.WriteString(text.S)
	}
	if currentLine.Len() > 0 {
		blocks = append(blocks, pdf.Text{S: currentLine.String(), FontSize: currentFont, Y: currentY})
	}
	return blocks
}

// 写入Markdown，增强表格识别与行合并
func writeMarkdownBlocks2(markdown *bytes.Buffer, blocks []pdf.Text, header1, header2 float64) {
	var para strings.Builder
	var inTable bool
	var tableHeader string
	var tableRows []string
	for i, block := range blocks {
		text := clearText(block.S)
		if len(text) == 0 {
			continue
		}
		// 检测表格头
		if isTableHeader(text) {
			flushPara(markdown, &para)
			if inTable && len(tableRows) > 0 {
				writeMarkdownTable(markdown, tableHeader, tableRows)
				tableRows = nil
			}
			inTable = true
			tableHeader = text
			continue
		}
		if inTable {
			// 判断是否���新表格行（以数字序号开头）
			if isTableRowStart(text) {
				tableRows = append(tableRows, text)
			} else if len(tableRows) > 0 {
				// 换行内容拼接到上一行
				tableRows[len(tableRows)-1] += text
			} else {
				// 表格头下第一行不是序号，直接拼接
				tableRows = append(tableRows, text)
			}
			// 如果下一个block不是表格内容，输出表格
			if i+1 == len(blocks) || (!isTableRowStart(clearText(blocks[i+1].S)) && !isTableHeader(clearText(blocks[i+1].S))) {
				writeMarkdownTable(markdown, tableHeader, tableRows)
				tableRows = nil
				inTable = false
			}
			continue
		}
		if block.FontSize >= header1 {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf("# %s\n", text))
		} else if block.FontSize >= header2 {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf("## %s\n", text))
		} else if isUnorderedList(text) {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf("- %s\n", strings.TrimLeft(text, "-•* ")))
		} else if isOrderedList(text) {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf("1. %s\n", strings.TrimLeft(text, "1234567890). ")))
		} else {
			if para.Len() > 0 {
				para.WriteString(" ")
			}
			para.WriteString(text)
		}
	}
	flushPara(markdown, &para)
	if inTable && len(tableRows) > 0 {
		writeMarkdownTable(markdown, tableHeader, tableRows)
	}
}

// 判断是否为表格头
func isTableHeader(text string) bool {
	return strings.HasPrefix(text, "序号 ") && strings.Contains(text, "企业名称")
}

// 判断是否为表格新行（以数字序号开头）
func isTableRowStart(text string) bool {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return false
	}
	if text[0] >= '0' && text[0] <= '9' {
		// 允许前面有序号和空格
		return true
	}
	return false
}

// 输出Markdown表格，支持跨行单元格合并
func writeMarkdownTable(markdown *bytes.Buffer, header string, rows []string) {
	headerFields := splitTableFields(header)
	markdown.WriteString("| " + strings.Join(headerFields, " | ") + " |\n")
	markdown.WriteString("|" + strings.Repeat(" --- |", len(headerFields)) + "\n")
	var prevFields []string
	for _, row := range rows {
		rowFields := splitTableFields(row)
		if len(rowFields) < len(headerFields) && len(prevFields) > 0 {
			// 如果首列为空，合并到上一行
			if strings.TrimSpace(rowFields[0]) == "" {
				for i := 1; i < len(rowFields) && i < len(prevFields); i++ {
					prevFields[i] += rowFields[i]
				}
				continue
			}
		}
		// 补齐列数
		for len(rowFields) < len(headerFields) {
			rowFields = append(rowFields, "")
		}
		markdown.WriteString("| " + strings.Join(rowFields, " | ") + " |\n")
		prevFields = rowFields
	}
}

// 按空格分割表格字段，连续空格视为一列分隔
func splitTableFields(line string) []string {
	fields := strings.Fields(line)
	if len(fields) == 1 && strings.Contains(line, " ") {
		// 处理只有一个字段但有多行的情况
		return append([]string{""}, fields[0])
	}
	return fields
}

// 写入Markdown
func writeMarkdownBlocks(markdown *bytes.Buffer, blocks []pdf.Text, header1, header2 float64) {
	var para strings.Builder
	for _, block := range blocks {
		text := clearText(block.S)
		if len(text) == 0 {
			continue
		}
		if block.FontSize >= header1 {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf("# %s\n", text))
		} else if block.FontSize >= header2 {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf(" %s\n", text))
		} else if isUnorderedList(text) {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf("- %s\n", strings.TrimLeft(text, "-•* ")))
		} else if isOrderedList(text) {
			flushPara(markdown, &para)
			markdown.WriteString(fmt.Sprintf("1. %s\n", strings.TrimLeft(text, "1234567890). ")))
		} else {
			if para.Len() > 0 {
				para.WriteString(" ")
			}
			para.WriteString(text)
		}
	}
	flushPara(markdown, &para)
}

func flushPara(markdown *bytes.Buffer, para *strings.Builder) {
	if para.Len() > 0 {
		markdown.WriteString(para.String() + "\n\n")
		para.Reset()
	}
}

func isUnorderedList(text string) bool {
	return strings.HasPrefix(text, "-") || strings.HasPrefix(text, "•") || strings.HasPrefix(text, "*")
}

func isOrderedList(text string) bool {
	if len(text) > 1 && (text[0] >= '0' && text[0] <= '9') && (text[1] == '.' || text[1] == ')') {
		return true
	}
	return false
}

func clearText(text string) string {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return text
	}
	text = strings.ReplaceAll(text, "�", " ")
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, " ", " ")
	text = strings.ReplaceAll(text, "\u0000", " ")
	text = strings.ReplaceAll(text, "\u001E", " ")
	text = strings.ReplaceAll(text, "\u001F", " ")
	text = strings.ReplaceAll(text, "\u0001", " ")
	return text
}

// 处理非表格文本块
func processNonTableBlocks(markdown *bytes.Buffer, blocks []pdf.Text) {
	var currentY float64 = -1
	var currentLine strings.Builder
	var currentFontSize float64

	for _, block := range blocks {
		if currentY != block.Y {
			if currentLine.Len() > 0 {
				writeMarkdownLine(markdown, currentLine.String(), currentFontSize)
				currentLine.Reset()
			}
			currentY = block.Y
			currentFontSize = block.FontSize
		}
		currentLine.WriteString(block.S + " ")
	}

	if currentLine.Len() > 0 {
		writeMarkdownLine(markdown, currentLine.String(), currentFontSize)
	}
}

// 根据文本特征写入适当的Markdown格式
func writeMarkdownLine(markdown *bytes.Buffer, text string, fontSize float64) {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return
	}

	switch {
	case fontSize > 14: // 大标题
		markdown.WriteString(fmt.Sprintf("# %s\n\n", text))
	case fontSize > 12: // 小标题
		markdown.WriteString(fmt.Sprintf("%s\n\n", text))
	case fontSize > 10: // 三级标题
		markdown.WriteString(fmt.Sprintf("%s\n\n", text))
	case strings.HasPrefix(text, "-") || strings.HasPrefix(text, "•"): // 列表项
		markdown.WriteString(fmt.Sprintf("- %s\n", strings.TrimLeft(text, "-• ")))
	case strings.HasPrefix(text, "1.") || strings.HasPrefix(text, "("): // 有序列表
		markdown.WriteString(fmt.Sprintf("1. %s\n", strings.TrimLeft(text, "1234567890). ")))
	default: // 普通段落
		markdown.WriteString(fmt.Sprintf("%s\n\n", text))
	}
}
