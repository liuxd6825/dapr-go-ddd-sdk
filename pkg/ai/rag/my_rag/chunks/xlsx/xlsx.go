package xlsx

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/utils/llmutils"
	"strings"
)

type XlsxChunk struct {
	llm *llm.LLM
}

type HeadResult struct {
	Head int `json:"head"`
}

func NewXlsxChunk() *XlsxChunk {
	return &XlsxChunk{}
}

func (s *XlsxChunk) GetSheets(llm llm.LLM, text string) ([]*Sheet, error) {
	sheets := []*Sheet{}
	reader := strings.NewReader(text)
	// 创建Scanner来读取文件
	scanner := bufio.NewScanner(reader)
	adding := false
	sheetsMap := map[string][]string{}
	var lines []string
	var sheetName string
	// 逐行读取
	for scanner.Scan() {
		line := scanner.Text() // 获取当前行的文本
		println(line)
		if strings.HasPrefix(line, "### ") {
			sheetName = strings.TrimSpace(strings.TrimPrefix(line, "### "))
			continue
		}
		if line == "```csv" {
			adding = true
			continue
		}
		if line == "```" {
			sheetsMap[sheetName] = lines
			adding = false
			continue
		}
		if adding {
			lines = append(lines, line)
		}
	}
	for name, lines := range sheetsMap {
		count := len(lines)
		if count > 20 {
			count = 20
		}
		text := lines[0:count]
		msg := s.getHeadMessage(text)
		ctx := context.Background()
		res, err := llm.Generate(ctx, msg)
		if err != nil {
			return nil, err
		}
		var headResult HeadResult
		jonsData := llmutils.GetJsonString(res.Content)
		err = json.Unmarshal([]byte(jonsData), &headResult)
		if err != nil {
			return nil, err
		}
		sheet := NewSheet()
		sheet.Init(name, headResult.Head, lines)
		sheets = append(sheets, sheet)
	}

	return sheets, nil
}

func (s *XlsxChunk) getHeadMessage(text []string) []*schema.Message {
	content := s.getHeadPrompt(text)
	msg := &schema.Message{
		Role:    schema.System,
		Content: content,
	}
	return []*schema.Message{msg}
}

func (s *XlsxChunk) getHeadPrompt(text []string) string {
	prompt := `你是一个markdown和csv文件分析器，分析下面csv数据表的字段行，并只以json格式返回结果。
### 数据示例
"交易明细","","","","","","","","",
"账户：","XXXX","账户名称:","XXXX",
"起始时间：","2019-1-1","截止时间：","2019-3-31","导出时间：","2019-06-04",
"交易时间","我方账号","我方户名","我方账户代号开户行名","对方账号","对方户名","对方账户代号开户行名","汇出金额","汇入金额","余额","摘要","用途","备注",
"2019-03-31  21:59:23","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6212250200008505347","杨孟霏","工商银行","-198.84","0.00","1784765.53","网上汇款","汇款","",
"2019-03-31  20:28:30","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6222080200009966071","谷继红","工商银行","-2000.00","0.00","1784964.37","网上汇款","汇款","",

### 结果示例
	{"head":3}

### 数据内容
%s
`
	data := fmt.Sprintf(prompt, "```csv\n"+strings.Join(text, "\n")+"\n```")
	return data
}
