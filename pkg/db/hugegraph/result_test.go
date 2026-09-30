package hugegraph

import (
	"encoding/json"
	"testing"
)

func Test_Result(t *testing.T) {
	jsonText := `[
      {
        "labels": [
          ["11111"],
          ["22222"]
        ],
        "objects": [
          {
            "id": "2:c",
            "label": "same",
            "type": "vertex",
            "properties": {
              "name": "c"
            }
          },
          {
            "id": "1:a_c",
            "label": "master",
            "type": "vertex",
            "properties": {
              "id": "a_c",
              "name": "c",
              "case_id": "1001",
              "tenant_id": "test",
              "source_ids": "a_c",
              "source_type": "master",
              "description": "公司名称:c; ",
              "type": "case_1001,company,master",
              "table": "company_company",
              "node_type": "case_1001,company,master"
            }
          }
        ]
      },
      {
        "labels": [
          [],
          []
        ],
        "objects": [
          {
            "id": "2:c",
            "label": "same",
            "type": "vertex",
            "properties": {
              "name": "c"
            }
          },
          {
            "id": "1:b_c",
            "label": "master",
            "type": "vertex",
            "properties": {
              "id": "b_c",
              "name": "c",
              "case_id": "1001",
              "tenant_id": "test",
              "source_ids": "b_c",
              "source_type": "master",
              "description": "公司名称:c; ",
              "type": "case_1001,company,master",
              "table": "company_company",
              "node_type": "case_1001,company,master"
            }
          }
        ]
      },
	{"id": "v1", "label": "human", "type": "vertex", "properties": {"name": "任萌"}},
	{"id": "e1", "label": "owns", "type": "edge", "outV": "v1", "inV": "account_123"},
	{"month": "202608", "amount": 500000.0},
	"标量字符串测试",
	12580,
	["列表元素1", "列表元素2"]
]`
	var data any
	err := json.Unmarshal([]byte(jsonText), &data)
	if err != nil {
		t.Error(err)
	}
	result := NewResult(data)
	t.Log("paths:", result.Paths)
	t.Log("nodes:", result.Nodes)
	t.Log("edges:", result.Edges)
	t.Log("records:", result.Records)
}
