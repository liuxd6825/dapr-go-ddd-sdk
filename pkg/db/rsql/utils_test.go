package rsql

import (
	"fmt"
	"testing"
)

func Test_GetFieldValues(t *testing.T) {
	// 示例1：结构体切片
	type Person struct {
		Name string
		Age  int
	}
	people := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}
	names, err := GetFieldValues(people, "Name")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Names:", names) // 输出: Names: ['Alice' 'Bob']
	}

	// 示例2：map 切片
	peopleMaps := []map[string]interface{}{
		{"Name": "Charlie", "Age": 35},
		{"Name": "David", "Age": 40},
	}
	names, err = GetFieldValues(peopleMaps, "Name")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Names:", names) // 输出: Names: ['Charlie' 'David']
	}

	// 示例3：混合类型
	mixedData := []map[string]interface{}{
		{"Name": "Eve", "Age": 28, "IsStudent": true},
		{"Name": "Frank", "Age": 45, "IsStudent": false},
	}
	names, err = GetFieldValues(mixedData, "Name")
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Names:", names) // 输出: Names: ['Eve' 'Frank']
	}
}
