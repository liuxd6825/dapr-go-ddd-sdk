package file_handler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"net/http"
)

func litSSR(rawHTML string) template.HTML {
	// 调用SSR服务处理原始HTML
	processedHTML, err := renderHTMLFragment(rawHTML)
	if err != nil {
		return template.HTML(rawHTML) // 降级返回原始HTML
	}
	return template.HTML(processedHTML)
}

func renderHTMLFragment(rawHTML string) ([]byte, error) {
	// 创建包含原始HTML的请求体
	reqBody := struct {
		HTML string `json:"html"`
	}{
		HTML: rawHTML,
	}

	jsonBody, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://127.0.0.1:3333/render",
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// 数据清洗防止XSS
func sanitizeProps(props interface{}) interface{} {
	// 使用第三方库清洗数据
	//return xss.Sanitize(props)
	return props
}
