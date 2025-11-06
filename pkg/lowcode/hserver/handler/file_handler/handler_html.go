package file_handler

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/spf13/afero"
)

// 使用 text/template 代替 html/template
import "text/template" // 无自动转义

// 定义动态标记的正则表达式
var dynamicMetaRegex = regexp.MustCompile(`<meta\s+name="ssr"\s+content="true"\s*/?>`)
var dynamicCommentRegex = regexp.MustCompile(`<!--\s*ssr\s*-->`)

// 缓存文件的动态标记状态
var dynamicCache sync.Map

// IsFileExist 判断文件是否存在
func (h *Handler) IsFileExist(fs afero.Fs, path string) (bool, error) {
	exists, err := afero.Exists(fs, path)
	return exists, err
}

// isDynamicPageRegex 使用正则表达式检查 HTML 文件是否是动态页面
func (h *Handler) isDynamicPage(ctx context.Context, ictx iris.Context, fs afero.Fs, filePath string, prodMode bool) (bool, error) {
	// 在dev模式下，默认为后端渲染。
	if !h.prodMode {
		return true, nil
	}

	if ictx != nil && ictx.Params().Exists("ssr") {
		return true, nil
	}
	// 检查缓存
	if cached, ok := dynamicCache.Load(filePath); ok {
		return cached.(bool), nil
	}

	// 打开文件
	file, err := fs.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file afero.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 只读取文件前 1024 字节
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil {
		return false, err
	}
	content := buffer[:n]

	// 检查是否匹配动态标记
	isDynamic := dynamicMetaRegex.Match(content) || dynamicCommentRegex.Match(content)
	// 缓存结果
	dynamicCache.Store(filePath, isDynamic)
	return isDynamic, nil
}

func parserHtml(doc *goquery.Document, envNpm *env.Npm, importFileFun func(scripts *goquery.Selection, importFile *strings.Builder)) ([]byte, error) {
	if envNpm != nil && len(envNpm.Links) > 0 {
		doc.Find("script").Each(func(i int, s *goquery.Selection) {
			parserImport(s, envNpm, importFileFun)
		})
	}
	html, err := doc.Html()
	if err != nil {
		return nil, err
	}
	html = strings.ReplaceAll(html, "&#34;", "\"")
	return []byte(html), err
}

// 模板语法检测
func hasTemplateSyntax(text string) bool {
	return strings.Contains(text, "{{") && strings.Contains(text, "}}")
}

// 原始模板处理器
func processRawTemplate(tpl string) string {
	tmpl := template.Must(template.New("raw").Parse(tpl))
	var buf strings.Builder
	tmpl.Execute(&buf, nil)
	return buf.String()
}

func parserImport(s *goquery.Selection, envNpm *env.Npm, importFileFun func(scripts *goquery.Selection, importFile *strings.Builder)) {
	typeAttr, exists := s.Attr("type")
	if exists && typeAttr == "module" && len(s.Nodes) > 0 && s.Nodes[0].FirstChild != nil {
		sb := transformImports(s.Nodes[0].FirstChild.Data, envNpm)
		if sb != nil && importFileFun != nil {
			importFileFun(s, sb)
		}
	}
}

// 转换导入路径
func transformImports(content string, envNpm *env.Npm) *strings.Builder {
	lines := strings.Split(content, "\n")
	sb := &strings.Builder{}
	for _, line := range lines {
		index := strings.Index(line, "import ")
		if index != -1 {
			line = line[index+7:]
			line = strings.ReplaceAll(line, "\"", "")
			line = strings.TrimSpace(line)
			isImport := false
			for _, link := range envNpm.Links {
				if link != nil && strings.HasPrefix(line, link.Name) {
					isImport = true
					sb.WriteString(fmt.Sprintf("import \"%s/%s\"\n", link.Path, line))
					break
				}
			}
			if !isImport {
				sb.WriteString(fmt.Sprintf("import \"%s\"\n", line))
			}
		}
	}
	return sb
}
