package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"sync"
)

type RagService struct {
	chatModel model.ToolCallingChatModel
}

var ragServiceOnce sync.Once

func NewRagService() *RagService {
	var ragService *RagService
	ragServiceOnce.Do(func() {
		ragService = newRagService()
	})
	return ragService
}

func newRagService() *RagService {
	ctx := context.Background()
	e := env.GetEnv()
	ragMeta := e.App.Meta["rag"]
	if ragMeta == nil {
		panic("rag not found in env.app")
	}

	ragCfg, err := config.ReadRagConfig(ragMeta)
	if err != nil {
		panic("read RagConfig error" + err.Error())
	}

	//chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
	//	BaseURL: ragCfg.LLM.BaseUrl,
	//	Model:   ragCfg.LLM.Model, // 使用的模型版本
	//	APIKey:  ragCfg.LLM.APIKey,
	//})

	var llmModel model.ToolCallingChatModel
	switch ragCfg.LLM.Type {
	case "ollama":
		m, err := llm.NewOllama(ctx, ollama.ChatModelConfig{
			BaseURL: ragCfg.LLM.BaseUrl,
			Model:   ragCfg.LLM.Model, // 使用的模型版本
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		llmModel = m
	case "openai":
		m, err := llm.NewOpenAI(ctx, openai.ChatModelConfig{
			BaseURL: ragCfg.LLM.BaseUrl,
			Model:   ragCfg.LLM.Model, // 使用的模型版本
			APIKey:  ragCfg.LLM.APIKey,
		})

		if err != nil {
			panic("open rag model error" + err.Error())
		}

		llmModel = m
	default:
		panic("rag config llm.type is null ")
	}

	if err != nil {
		panic("create import.RagModel error" + err.Error())
	}
	ragService := &RagService{}
	ragService.chatModel = llmModel
	return ragService
}

func (s *RagService) Query(ctx context.Context, qry *query.RagQueryRequest, streams ...func(txt string)) (string, error) {
	userPrompt := &schema.Message{
		Role:    "user",
		Content: qry.UserPrompt,
	}
	input := []*schema.Message{s.getSystemPrompt(), userPrompt}
	streamReader, err := s.chatModel.Stream(ctx, input)
	if err != nil {
		return "", fmt.Errorf("大模型问题分析时出错: %w", err)
	}

	sb, err := my_rag.Reader(streamReader, streams...)
	if err != nil {
		return "", fmt.Errorf("读取返回结果时出错: %w", err)
	}
	return sb.String(), nil
}

func (s *RagService) getSystemPrompt() *schema.Message {
	return &schema.Message{
		Role: "system",
		Content: "## 角色：" +
			"你是一个数据分析师。你的任务是根据给定的json格式数据分析标准列数据json中script属性怎么取值和mapKeys属性怎么取值。" +
			"## 核心规则：" +
			"1. 你的输出必须严格遵守JSON格式，不要包含Markdown标记。" +
			"2. 标准列中属性值可以只是个常量也可以是公式。" +
			"3. 标准列是JSON格式数据。" +
			"**标准列数据**" +
			"```json" +
			"{" +
			"    \"properties\": {" +
			"        \"date\": {\"key\": \"date\",\"name\": \"交易日期\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"iden\": {\"key\": \"iden\",\"name\": \"我方标识\",\"script\": \"\",\"allowNull\": true,\"isHide\": true,\"mapKeys\": null}," +
			"        \"name\": {\"key\": \"name\",\"name\": \"我方名称\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}, " +
			"        \"acct\": {\"key\": \"acct\",\"name\": \"我方账号\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"bankName\": {\"key\": \"bankName\",\"name\": \"我方开户行\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"acctType\": {\"key\": \"acctType\",\"name\": \"我方账号类型\",\"script\": \"\",\"allowNull\": true,\"isHide\": true,\"mapKeys\": null}," +
			"        \"category\": {\"key\": \"category\",\"name\": \"类别\",\"script\": \"\",\"allowNull\": true,\"isHide\": true,\"mapKeys\": null}," +
			"        \"oppIden\": {\"key\": \"oppIden\",\"name\": \"对方标识\",\"script\": \"\",\"allowNull\": true,\"isHide\": true,\"mapKeys\": null}," +
			"        \"oppName\": {\"key\": \"oppName\",\"name\": \"对方名称\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"oppAcct\": {\"key\": \"oppAcct\",\"name\": \"对方账号\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"oppBankName\": {\"key\": \"oppBankName\",\"name\": \"对方开户行\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"oppAcctType\": {\"key\": \"oppAcctType\",\"name\": \"对方账号类型\",\"script\": \"\",\"allowNull\": true,\"isHide\": true,\"mapKeys\": null}," +
			"        \"oppCategory\": {\"key\": \"oppCategory\",\"name\": \"对方类别\",\"script\": \"\",\"allowNull\": true,\"isHide\": true,\"mapKeys\": null}," +
			"        \"income\": {\"key\": \"income\",\"name\": \"收入金额(贷)\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"payout\": {\"key\": \"payout\",\"name\": \"支出金额(借)\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"amount\": {\"key\": \"amount\",\"name\": \"交易金额\",\"script\": \"\",\"allowNull\": false,\"isHide\": null,\"mapKeys\": null}," +
			"        \"balance\": {\"key\": \"balance\",\"name\": \"余额\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}," +
			"        \"cash\": {\"key\": \"cash\",\"name\": \"是否现金\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}," +
			"        \"ccy\": {\"key\": \"ccy\",\"name\": \"币种\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}," +
			"        \"serial\": {\"key\": \"serial\",\"name\": \"流水号\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}," +
			"        \"type\": {\"key\": \"type\",\"name\": \"交易类型\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}," +
			"        \"place\": {\"key\": \"place\",\"name\": \"交易地点\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}," +
			"        \"summary\": {\"key\": \"summary\",\"name\": \"摘要\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}," +
			"        \"notes\": {\"key\": \"notes\",\"name\": \"备注\",\"script\": \"\",\"allowNull\": true,\"isHide\": null,\"mapKeys\": null}" +
			"    }" +
			"}" +
			"```" +
			"4. 根据标准列数据填写最终数据。" +
			"5. mapKeys属性是标准列需要使用到的表格列，mapKeys属性值是字符串数组。" +
			"6. script属性是标准列需要使用到的公式，公式可以进行嵌套，公式名取title属性值，类型是字符串，公式参数字段名是变量。" +
			"7. 只处理mapKeys和script两个属性，其它属性原样输出。" +
			"8. 最后输出完整列json数据。" +
			"9. 我方名称属性值从户名或我方户名字段取值，如果数据中没有户名或我方户名字段则从数据json中表名属性值中分析取出值。" +
			"10. 我方账号属性值从我方账号或账号字段取值。" +
			"11. 我方开户行属性值从开户行或我方开户行字段取值，如果数据中没有开户行或我方开户行字段则从数据json中非结构化数据属性值中分析取出值。" +
			"**公式json数据**" +
			"```json" +
			"{" +
			"	\"properties\": {" +
			"		{ \"name\": \"文字替换\", \"parameters\": [\"原始文本字段\", \"需要替换的文本\", \"新文本\"], \"desc\": \"文字替换(【文本】,【替换文本】,【新文本】)\" }," +
			"		{ \"name\": \"取时间\", \"parameters\": [\"日期部分字段\", \"时间部分字段\"], \"desc\": \"取时间(【日期部分】,【时间部分...】)\" }," +
			"		{ \"name\": \"取浮点值\", \"parameters\": [\"需要转成浮点数的字段\"], \"desc\": \"取浮点值(【值】)\" }," +
			"		{ \"name\": \"取绝对值\", \"parameters\": [\"需要取绝对值的字段\"], \"desc\": \"取绝对值(【值】)\" }," +
			"		{ \"name\": \"是否有负号\", \"parameters\": [\"需要判断是否有负号的字段\"], \"desc\": \"是否有负号(【值】)\" }," +
			"		{ \"name\": \"取支出金额\", \"parameters\": [\"支出金额字段\"], \"desc\": \"取支出金额(【值】)\" }," +
			"		{ \"name\": \"取收入金额\", \"parameters\": [\"收入金额字段\"], \"desc\": \"取收入金额(【值】)\" }," +
			"		{ \"name\": \"取交易金额\", \"parameters\": [\"支出金额字段\", \"收入金额字段\"], \"desc\": \"取交易金额(【支出金额】,【收入金额】)\" }," +
			"		{ \"name\": \"取文本\", \"parameters\": [\"需要取文本的字段\", \"默认值\"], \"desc\": \"取文本(【文本】,【默认值】)\" }," +
			"		{ \"name\": \"根据标识取支出金额\", \"parameters\": [\"'借'\", \"支出标识字段\", \"支出金额字段\"], \"desc\": \"根据标识取支出金额(【标识文本】,【支出标识】,【金额】)\" }," +
			"		{ \"name\": \"根据标识取收入金额\", \"parameters\": [\"'贷'\", \"收入标识字段\", \"收入金额字段\"], \"desc\": \"根据标识取收入金额(【标识文本】,【收入标识】,【金额】)\" }," +
			"		{ \"name\": \"取数字文本\", \"parameters\": [\"需要取数字文本的字段\", \"默认值\", \"顺序号\"], \"desc\": \"取数字文本(【文本】,【默认值】,【顺序号】)\" }," +
			"		{ \"name\": \"取中间文本\", \"parameters\": [\"需要取文本的字段\", \"开始文本\", \"结尾文本\", \"替换文本\"], \"desc\": \"取中间文本(【文本】,【开始文本】,【结尾文本】,【替换文本...】)\"  }," +
			"		{ \"name\": \"取支付宝账号\", \"parameters\": [\"支付宝公司名称\", \"对方公司名称\", \"备注或摘要字段\", \"默认值\", \"替换文本\"], \"desc\": \"取支付宝账号(【支付宝公司名称】,【对方公司名称】,【备注或摘要】,【默认值】,【替换文本...】)\" }," +
			"		{ \"name\": \"取支持宝人名\", \"parameters\": [\"支付宝公司名称\", \"对方公司名称\", \"备注或摘要字段\", \"开始文本\", \"结尾文本\", \"默认值\", \"替换文本\"], \"desc\": \"取支持宝人名(【支付宝公司名称】,【对方公司名称】,【备注或摘要】,【开始文本】,【结尾文本】,【默认值】,【替换文本...】)\" }," +
			"		{ \"name\": \"取币种\", \"parameters\": [\"币种字段\"], \"desc\": \"取币种(【字段...】)\" }," +
			"		{ \"name\": \"是否现金交易\", \"parameters\": [\"现金交易关键词'现金,取现,卡取'\", \"提取现金交易关键词的字段\"], \"desc\": \"是否现金交易('现金,取现,卡取',【备注字段...】)\" }," +
			"		{ \"name\": \"取开户行\", \"parameters\": [\"开户银行字段\"], \"desc\": \"取开户行(【字段...】)\" }" +
			"	}" +
			"}" +
			"```",
	}
}
