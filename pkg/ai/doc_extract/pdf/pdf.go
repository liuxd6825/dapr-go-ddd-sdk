package pdf

import (
	"bytes"
	"fmt"
	"github.com/ledongthuc/pdf"
	"github.com/spf13/afero"
	"sort"
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
	// 创建ReadSeeker
	readSeeker := bytes.NewReader(data)
	// 2. 创建 PDF Reader
	pdfReader, err := pdf.NewReader(readSeeker, readSeeker.Size())
	if err != nil {
		panic(err)
	}

	// 3. 获取 PDF 页数
	numPages := pdfReader.NumPage()
	var markdown bytes.Buffer
	// 4. 逐页提取文本
	for pageNum := 1; pageNum <= numPages; pageNum++ {
		page := pdfReader.Page(pageNum)
		if page.V.IsNull() || page.V.Key("Contents").Kind() == pdf.Null {
			continue
		}

		texts := page.Content().Text
		if len(texts) == 0 {
			continue
		}
		// 按位置排序文本块 (从上到下，从左到右)
		sort.SliceStable(texts, func(i, j int) bool {
			if texts[i].Y != texts[j].Y {
				return texts[i].Y > texts[j].Y
			}
			return texts[i].X < texts[j].X
		})

		// 添加页面标题
		markdown.WriteString(fmt.Sprintf("## 第%d页\n", pageNum))

		// 组织文本块为有结构的格式
		var (
			currentY    float64 = -1
			currentLine strings.Builder
			blocks      []pdf.Text
		)
		// 合并同一行的文本
		for _, text := range texts {
			if currentY != text.Y {
				if currentLine.Len() > 0 {
					blocks = append(blocks, pdf.Text{
						S:        currentLine.String(),
						FontSize: text.FontSize, // 使用最后一个文本块的字体大小
						Y:        currentY,
					})
					currentLine.Reset()
				}
				currentY = text.Y
			}
			currentLine.WriteString(text.S)
		}
		// 添加最后一行
		if currentLine.Len() > 0 {
			blocks = append(blocks, pdf.Text{
				S:        currentLine.String(),
				FontSize: texts[len(texts)-1].FontSize,
				Y:        currentY,
			})
		}

		// 转换为Markdown格式
		for _, block := range blocks {
			text := clearText(block.S)
			if len(text) == 0 {
				continue
			}

			// 根据字体大小判断格式
			switch {
			case block.FontSize > 14: // 大标题
				markdown.WriteString(fmt.Sprintf("# %s\n", text))
			case block.FontSize > 12: // 小标题
				markdown.WriteString(fmt.Sprintf("## %s\n", text))
			//case block.FontSize > 10: // 三级标题
			//markdown.WriteString(fmt.Sprintf("### %s\n", text))
			case strings.HasPrefix(text, "-") || strings.HasPrefix(text, "•"): // 列表项
				markdown.WriteString(fmt.Sprintf("- %s\n", strings.TrimLeft(text, "-• ")))
			case strings.HasPrefix(text, "1.") || strings.HasPrefix(text, "("): // 有序列表
				markdown.WriteString(fmt.Sprintf("1. %s\n", strings.TrimLeft(text, "1234567890). ")))
			default: // 普通段落
				markdown.WriteString(fmt.Sprintf("%s\n", text))
			}
		}
		if pageNum > 1 {
			markdown.WriteString("<!--break-->\n") // 页间分隔
		}
	}
	return markdown.String(), nil
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
		markdown.WriteString(fmt.Sprintf("## %s\n\n", text))
	case fontSize > 10: // 三级标题
		markdown.WriteString(fmt.Sprintf("### %s\n\n", text))
	case strings.HasPrefix(text, "-") || strings.HasPrefix(text, "•"): // 列表项
		markdown.WriteString(fmt.Sprintf("- %s\n", strings.TrimLeft(text, "-• ")))
	case strings.HasPrefix(text, "1.") || strings.HasPrefix(text, "("): // 有序列表
		markdown.WriteString(fmt.Sprintf("1. %s\n", strings.TrimLeft(text, "1234567890). ")))
	default: // 普通段落
		markdown.WriteString(fmt.Sprintf("%s\n\n", text))
	}
}
