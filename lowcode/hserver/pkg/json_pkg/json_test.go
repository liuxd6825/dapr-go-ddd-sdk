package json_pkg

import "testing"

func TestJsonPkg_Format(t *testing.T) {
	jsonPkg, err := NewJsonPkg()
	if err != nil {
		t.Fatal(err)
	}
	// 原始 JSON 字符串
	rawJSON := `{"name":"John","age":30,"address":{"city":"New York","zip":"10001"},"hobbies":["reading","coding"]}`
	json := jsonPkg.Format([]byte(rawJSON))
	t.Log(string(json))
}
