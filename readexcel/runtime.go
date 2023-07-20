package readexcel

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/script"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"math"
	"strconv"
	"strings"
)

const LocalTimeLayoutLine = "2006-01-02 15:04:05"

func newRuntime() (*goja.Runtime, error) {
	vm, err := script.NewRuntime()
	if err != nil {
		return nil, err
	}
	if err := vm.Set("replace", replace); err != nil {
		return nil, err
	}
	if err := vm.Set("替换", replace); err != nil {
		return nil, err
	}

	if err := vm.Set("toDateTime", toDateTime); err != nil {
		return nil, err
	}
	if err := vm.Set("转时间", toDateTime); err != nil {
		return nil, err
	}

	if err := vm.Set("abs", abs); err != nil {
		return nil, err
	}
	if err := vm.Set("取绝对值", abs); err != nil {
		return nil, err
	}

	if err := vm.Set("isMinus", isMinus); err != nil {
		return nil, err
	}
	if err := vm.Set("是否有负号", isMinus); err != nil {
		return nil, err
	}

	if err := vm.Set("payout", payout); err != nil {
		return nil, err
	}
	if err := vm.Set("取支付金额", payout); err != nil {
		return nil, err
	}

	if err := vm.Set("income", income); err != nil {
		return nil, err
	}
	if err := vm.Set("取收入金额", income); err != nil {
		return nil, err
	}

	if err := vm.Set("amount", amount); err != nil {
		return nil, err
	}
	if err := vm.Set("取交易金额", amount); err != nil {
		return nil, err
	}

	if err := vm.Set("toFloat", toFloat); err != nil {
		return nil, err
	}
	if err := vm.Set("转浮点值", income); err != nil {
		return nil, err
	}

	return vm, nil
}

func replace(s string, old string, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func toDateTime(dateStr, timeStr string) string {
	v := timeutils.ToDateTime(dateStr, timeStr)
	return v.Format(LocalTimeLayoutLine)
}

func toFloat(val any) float64 {
	if v, ok := val.(float64); ok {
		return v
	} else if v, ok := val.(*float64); ok {
		return *v
	} else if v, ok := val.(int64); ok {
		return float64(v)
	} else if v, ok := val.(*int64); ok {
		return float64(*v)
	} else if v, ok := val.(int); ok {
		return float64(v)
	} else if v, ok := val.(*int); ok {
		return float64(*v)
	} else if v, ok := val.(string); ok {
		f, err := strconv.ParseFloat(v, 8)
		if err != nil {
			return 0
		}
		return f
	}
	return 0
}

func abs(val any) float64 {
	return math.Abs(toFloat(val))
}

//
//  isMinus
//  @Description: 是否有负号
//  @param val
//  @return bool
//
func isMinus(val any) bool {
	if v, ok := val.(string); ok {
		return strings.HasPrefix("-", v)
	} else if v, ok := val.(float64); ok {
		return v > 0
	} else if v, ok := val.(*float64); ok {
		return *v > 0
	} else if v, ok := val.(int64); ok {
		return v > 0
	} else if v, ok := val.(*int64); ok {
		return *v > 0
	} else if v, ok := val.(int); ok {
		return v > 0
	} else if v, ok := val.(*int); ok {
		return *v > 0
	}
	return false
}

//
//  payout
//  @Description: 取得支付金额
//  @param val
//  @return float64
//
func payout(val any) float64 {
	v := toFloat(val)
	if v > 0 {
		return 0
	}
	return abs(v)
}

//
//  income
//  @Description: 取得收入金额
//  @param val
//  @return float64
//
func income(val any) float64 {
	v := toFloat(val)
	if v < 0 {
		return 0
	}
	return v
}

//
//  amount
//  @Description: 取交易金额
//  @param val
//  @return float64
//
func amount(v1 any, v2 any) float64 {
	a1 := abs(v1)
	a2 := abs(v2)
	if a1 != 0 {
		return a1
	}
	return a2
}
