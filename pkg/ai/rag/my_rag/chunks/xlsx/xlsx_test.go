package xlsx

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	llm2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/stretchr/testify/assert"
	"testing"
)

var xlsxText = fmt.Sprintf(`
### Sheet1
%scsv
"安徽新安银行银行企业账户交易明细","","","","","","","","",
"账户：","658010100100047060","账户名称:","北京家家财富投资有限公司",
"起始时间：","2019-1-1","截止时间：","2019-3-31","导出时间：","2019-06-04",
"交易时间","我方账号","我方户名","我方账户代号开户行名","对方账号","对方户名","对方账户代号开户行名","汇出金额","汇入金额","余额","摘要","用途","备注",
"2019-03-31  21:59:23","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6212250200008505347","杨孟霏","工商银行","-198.84","0.00","1784765.53","网上汇款","汇款","",
"2019-03-31  20:28:30","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6222080200009966071","谷继红","工商银行","-2000.00","0.00","1784964.37","网上汇款","汇款","",
"2019-03-31  19:49:20","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6217002870072970260","孙晨晨","建设银行","-4800.00","0.00","1786964.37","网上汇款","汇款","",
"2019-03-31  19:48:56","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6217002870072970260","孙晨晨","建设银行","-50000.00","0.00","1791764.37","网上汇款","汇款","",
"2019-03-31  19:46:46","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6217002870072970260","孙晨晨","建设银行","-50000.00","0.00","1841764.37","网上汇款","汇款","",
"2019-03-31  19:46:19","658010100100047060","北京家家财富投资有限公司","新安银行营业部","6217002870072970260","孙晨晨","建设银行","-50000.00","0.00","1891764.37","网上汇款","汇款",""
%s
`, "```", "```")

func Test_SplitXLSX_GetSheets(t *testing.T) {
	ctx := context.Background()
	xlsx := NewXlsxChunk()
	llm := newLLM(ctx)
	sheets, err := xlsx.GetSheets(llm, xlsxText)
	assert.NoError(t, err)
	assert.NotNil(t, sheets)
	t.Log(sheets[0].GetRows())
}

func newLLM(ctx context.Context) llm2.LLM {
	llm, err := llm2.NewOpenAI(ctx, openai.ChatModelConfig{
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "deepseek-r1-distill-llama-70b", // 使用的模型版本
		APIKey:  "sk-4a999651298047efaaf38aea633ba636",
	})
	if err != nil {
		panic(err)
	}
	return llm
}

func newLLM2(ctx context.Context) llm2.LLM {
	llm, err := llm2.NewOllama(ctx, ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434",
		Model:   "modelscope.cn/unsloth/DeepSeek-R1-Distill-Qwen-7B-GGUF:latest",
	})
	if err != nil {
		panic(err)
	}
	return llm
}
