package mongo

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
### 人员表
CREATE TABLE mcp_item (

	    id varchar(255) NOT NULL PRIMARY KEY, // 主键
	    case_id varchar(255) NOT NULL, // 案装ID
	    tenant_id varchar(255) NOT NULL, // 租户ID
		name varchar(255) NOT NULL, // 姓名
	    age int NOT NULL,       // 年龄
	    money float NOT NULL,  // 金额
	    date datetime NOT NULL, // 交易日期

);
*/

const SQLSchema = `
### 交易记录表，银行流水表
CREATE TABLE import_record (
  id VARCHAR(36) NOT NULL COMMENT '主键',
  tenant_id VARCHAR(36) COMMENT '租户ID',
  case_id VARCHAR(36) COMMENT '案件ID',
  is_deleted TINYINT(1) DEFAULT 0 COMMENT '是否删除',
  remark VARCHAR(500) COMMENT '备注',
  
  row_num BIGINT COMMENT '行号',
  task_id VARCHAR(36) COMMENT '任务id',
  doc_id VARCHAR(36) COMMENT '文档id',
  file_id VARCHAR(36) COMMENT '文件id',
  
  iden VARCHAR(100) COMMENT '我方标识',
  name VARCHAR(100) COMMENT '我方名称', ### 名称、姓名、人员
  acct VARCHAR(100) COMMENT '我方账号',
  acct_type VARCHAR(50) COMMENT '我方账号类型',
  category VARCHAR(50) COMMENT '我方类别',
  bank_name VARCHAR(100) COMMENT '我方开户银行',
  balance DECIMAL(20,6) COMMENT '我方余额账户',
  
  opp_iden VARCHAR(100) COMMENT '对方标识',
  opp_name VARCHAR(100) COMMENT '对方名称',
  opp_acct VARCHAR(100) COMMENT '对方账号',
  opp_acct_type VARCHAR(50) COMMENT '对方账号类型',
  opp_category VARCHAR(50) COMMENT '对方类别',
  opp_bank_name VARCHAR(100) COMMENT '对方开户银行',
  
  serial VARCHAR(100) COMMENT '流水号',
  payout DECIMAL(20,6) COMMENT '支出金额',
  income DECIMAL(20,6) COMMENT '收入金额',
  amount DECIMAL(20,6) COMMENT '交易金额',

  date DATETIME COMMENT '交易时间',
  year INT COMMENT '年', // 交易时间中的年份，统计字段
  month INT COMMENT '月', // 交易时间中的月份，统计字段
  day INT COMMENT '日',  // 交易时间中的月的日期，统计字段

  type VARCHAR(50) COMMENT '交易类型',
  ccy VARCHAR(10) COMMENT '交易币种',
  place VARCHAR(100) COMMENT '地点',
  summary VARCHAR(500) COMMENT '摘要',
  notes VARCHAR(500) COMMENT '备注'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='交易记录表，银行流水表';

`

type McpRecord struct {
	Id       string     `gorm:"id;primaryKey"`
	TenantId string     `gorm:"tenant_id"`
	Name     string     `gorm:"name"`
	Age      int        `gorm:"age"`
	Money    float64    `gorm:"money"`
	Date     *time.Time `gorm:"data"`
}

func Test_MongoQueryTool_Handler(t *testing.T) {
	ctx := xtest.NewContext()
	logger := logrus.New()

	cfg := &Config{
		Tool: ToolConfig{
			Name:        "import_record",
			Description: "获取，查询，统计：银行流水和交易数据",
			Schema:      SQLSchema,
		},
		DB: DBConfig{
			Host:           "127.0.0.1:27017,127.0.0.1:27018,127.0.0.1:27019",
			UserName:       "super_admin",
			Password:       "123456",
			ReplicaSet:     "rs0",
			AuthSource:     "admin",
			DBName:         "test",
			CollName:       "import_record",
			ConnectTimeout: 2000,
		},
		LLM: LLMConfig{
			Type:    LLMType_Ollama,
			BaseURL: "http://localhost:11434",
			Model:   "qwen2.5-coder:14b",
			ApiKey:  "",
		},
	}

	tool, err := NewMongoQueryTool(ctx, cfg, logger)
	assert.NoError(t, err)

	params := map[string]any{
		"tenantId": "test",
		"caseId":   "1001",
		"userId":   "test",
		"query":    "查询日期在2012年我方名称是张宇的银行流水，并以年和月字段分组汇总",
	}
	request := mcp.CallToolRequest{
		Request: mcp.Request{
			Method: "query",
			Params: mcp.RequestParams{
				Meta: &mcp.Meta{},
			},
		},
		Params: mcp.CallToolParams{
			Name:      "query",
			Arguments: params,
		},
	}

	resp, err := tool.Handler(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range resp.Content {
		if val, ok := c.(mcp.TextContent); ok {
			t.Log(val.Text)
		}
	}
}
