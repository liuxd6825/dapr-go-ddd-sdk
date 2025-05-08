package json_pkg

import "testing"

func TestJsonPkg_Format(t *testing.T) {
	jsonPkg := NewJsonPkg()
	// 原始 JSON 字符串
	rawJSON := `{"name":"John","age":30,"address":{"city":"New York","zip":"10001"},"hobbies":["reading","coding"]}`
	json := jsonPkg.Format([]byte(rawJSON))
	t.Log(string(json))
}

func TestJsonPkg_NewMap(t *testing.T) {
	jsonPkg := NewJsonPkg()
	// 原始 JSON 字符串
	rawJSON := `{"name":"John","age":30,"address":{"city":"New York","zip":"10001"},"hobbies":["reading","coding"]}`
	json := jsonPkg.NewMap(rawJSON)
	t.Log(json)
}
