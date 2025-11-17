package service

import (
	"fmt"
	"testing"
	"time"
)

func Test_CompleteAndSortData(t *testing.T) {
	// 您的原始数据，格式为 []map[string]any
	input := []map[string]any{
		{"amount": 11520.0, "month": 1, "oppCount": 2, "recordCount": 4, "time": toTime("2015-01-01 08:00:00")},
		{"amount": 11520.0, "month": 1, "oppCount": 2, "recordCount": 4, "time": toTime("2015-02-01 08:00:00")},
		{"amount": 2764233671.04, "month": 3, "oppCount": 2, "recordCount": 20, "time": toTime("2015-03-01 08:00:00")},
		{"amount": 100.0, "month": 3, "oppCount": 1, "recordCount": 1, "time": toTime("2015-06-01 08:00:00")}, // 新增一条按天测试的数据
	}

	fmt.Println("------ 按 'yearMonth' 补全 ------")
	completedByMonth, err := CompleteAndSortData(input, "time", "month")
	if err != nil {
		fmt.Println("处理数据时发生错误:", err)
		return
	}
	for _, d := range completedByMonth {
		fmt.Printf("%v\n", d)
	}

	fmt.Println("\n------ 按 'day' 补全 ------")
	// 只截取 2015-03 的数据来演示按天补全
	dayInput := []map[string]any{
		{"amount": 2764233671.04, "month": 3, "oppCount": 2, "recordCount": 20, "time": toTime("2015-03-01 08:00:00")},
		{"amount": 100.0, "month": 3, "oppCount": 1, "recordCount": 1, "time": toTime("2015-03-03 10:00:00")},
	}
	completedByDay, err := CompleteAndSortData(dayInput, "time", "day")
	if err != nil {
		fmt.Println("处理数据时发生错误:", err)
		return
	}
	for _, d := range completedByDay {
		fmt.Printf("%v\n", d)
	}
}

func toTime(val string) time.Time {
	res, _ := time.Parse(time.DateTime, val)
	return res
}
