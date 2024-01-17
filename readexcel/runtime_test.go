package readexcel

import (
	"fmt"
	"testing"
)

func TestField_toDateTime(t *testing.T) {
	date := "20110812"
	time := "11.56.23"
	res := toDateTime(date, time)

	date = "20221001 11:09:22"
	res = toDateTime(date)

	date = "20221001T11:09:22"
	res = toDateTime(date)

	println(res)
}

func Test(t *testing.T) {

}

func Test_regexpNum(t *testing.T) {
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
		fmt.Printf("%s => %s\n", item, regexpNum(item, item, 2))
	}

}

func Test_Submatch(t *testing.T) {
	list := []string{
		"交易对手卡号: 696070521, 入账: 田海博支付宝转账",
		"交易对手卡号: 696070521, 王磊支付宝转账",
		"交易对手卡号: 696070521, 朱秀芳支付宝转账",
		"交易对手卡号: 696070521, 王光军支付宝转账",
		"交易对手卡号: 696070521, 王光军支付宝转账",
		"交易对手卡号: 696070521, 谭雪支付宝转账",
	}
	for _, item := range list {
		fmt.Printf("%s => %s\n", item, match(item, ",", "支付宝转账", "入账: "))
	}
}

func Test_GetZhiFuBao(t *testing.T) {
	list := []zhiFuBaoAccByRemarks{
		{"696070521", "支付宝（中国）网络技术有限公司客户备付金", "交易对手卡号:696070521,云霞支付宝转账	"},
		{"696070521	", "支付宝（中国）网络技术有限公司客户备付金	", "交易对手卡号:696070521,郑华支付宝转账	"},
		{"696070521	", "支付宝（中国）网络技术有限公司客户备付金	", "交易对手卡号:696070521,丁璐支付宝转账	"},
		{"33001616783059000667	", "支付宝（中国）网络技术有限公司客户备付金	", "交易对手卡号:33001616783059000667,宋卫东支付宝转账-支付宝（中国）网络技术有限公司客户备付金	"},
		{"33001616783059000667	", "支付宝（中国）网络技术有限公司客户备付金	", "交易对手卡号:33001616783059000667,杨玉军支付宝转账-支付宝（中国）网络技术有限公司客户备付金	"},
		{"33001616783059000667	", "支付宝（中国）网络技术有限公司客户备付金	", "交易对手卡号:33001616783059000667,宋卫东支付宝转账-支付宝（中国）网络技术有限公司客户备付金	"},
	}

	for _, item := range list {
		fmt.Printf("%s => %s, %s\n", item.Remarks,
			GetZhiFuBaoOppAcc("支付宝（中国）", item.Company, item.Remarks, "", 0),
			GetZhiFuBaoOppName("支付宝（中国）", item.Company, item.Remarks, ",", "支付宝", "", 0),
		)
	}
}

type zhiFuBaoAccByRemarks struct {
	Account string
	Company string
	Remarks string
}
