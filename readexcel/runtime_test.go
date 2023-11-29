package readexcel

import (
	"fmt"
	"testing"
)

func TestField_toDateTime(t *testing.T) {
	date := "20110812"
	time := "11.56.23"
	res := toDateTime(date, time)
	println(res)
}

func Test_regexpNumIndex(t *testing.T) {
	list := []string{
		"交易对手卡号:6212273100001667859,王森利息66666666666666666666666",
		"交易对手卡号:6214993760011176,李勇利息",
		"透支回补（本金） 钱生钱C",
		"交易对手卡号:6214993760011176,罗丽利息",
		"交易对手卡号:6214993760011176,丁万书还款",
		"交易对手卡号:11253000000873399,fk",
		"交易对手卡号:6226630901542281,网银转账",
		"交易对手卡号:6222023100085529787,网银转账",
		"交易对手卡号:6214993760011176,丁万书还款",
		"交易对手卡号:6214993760011176,丁万书还款",
	}

	for _, item := range list {
		fmt.Printf("%s = %s\n", item, regexpNum(item, item, 2))
	}

}
