package store_huge

import (
	"encoding/json"
	"fmt"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/gremlin"
)

// HugeResult 封装 Gremlin 返回的数据,提供常用取值/转换方法
type HugeResult struct {
	resp *gremlin.PostResponseData
}

func NewHugeResult(respData *gremlin.PostResponseData) *HugeResult {
	return &HugeResult{resp: respData}
}

func (r *HugeResult) Raw() any {
	return r.resp.Result.Data
}

// Records 返回 []map[string]any 形式的查询结果
func (r *HugeResult) Records() []map[string]any {
	return nil
}

// First 取第一条记录
func (r *HugeResult) First() map[string]any {
	list := r.Records()
	if len(list) == 0 {
		return nil
	}
	return list[0]
}

// OneValue 取单个标量值(如 count、id 等)
func (r *HugeResult) OneValue() (any, error) {
	switch v := r.resp.Result.Data.(type) {
	case nil:
		return nil, nil
	case []any:
		if len(v) == 0 {
			return nil, nil
		}
		if len(v) > 1 {
			return nil, fmt.Errorf("multiple values, expect single")
		}
		return v[0], nil
	default:
		return v, nil
	}
}

func toMapSlice(list []any) []map[string]any {
	out := make([]map[string]any, 0, len(list))
	for _, item := range list {
		switch m := item.(type) {
		case map[string]any:
			out = append(out, m)
		default:
			b, err := json.Marshal(item)
			if err != nil {
				continue
			}
			var mm map[string]any
			if err := json.Unmarshal(b, &mm); err == nil {
				out = append(out, mm)
			}
		}
	}
	return out
}
