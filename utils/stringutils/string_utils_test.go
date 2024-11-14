package stringutils

import (
	"fmt"
	"testing"
)

func Test_ReplacePlaceholders(t *testing.T) {
	// 定义字符串模板和替换值
	template := "{name} is {age} years old and lives in {city}."
	values := map[string]any{
		"name": "John",
		"age":  "30",
		"city": "New York",
	}

	// 替换占位符
	result := ReplacePlaceholders(template, values)
	fmt.Println(result) // 输出: "John is 30 years old and lives in New York."
	if result != "John is 30 years old and lives in New York." {
		t.Fail()
	}
}
