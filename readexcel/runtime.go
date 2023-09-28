package readexcel

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/script"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"math"
	"strconv"
	"strings"
	"time"
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
	if err := vm.Set("文字替换", replace); err != nil {
		return nil, err
	}

	if err := vm.Set("toDateTime", toDateTime); err != nil {
		return nil, err
	}
	if err := vm.Set("取时间", toDateTime); err != nil {
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
	if err := vm.Set("取支出金额", payout); err != nil {
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
	if err := vm.Set("取浮点值", toFloat); err != nil {
		return nil, err
	}

	if err := vm.Set("toString", toString); err != nil {
		return nil, err
	}
	if err := vm.Set("取文本", toString); err != nil {
		return nil, err
	}

	return vm, nil
}

func toString(s string, def string) string {
	if len(s) == 0 {
		return def
	}
	return s
}

func replace(s string, old string, new string) string {
	return strings.ReplaceAll(s, old, new)
}

func toDateTime(value ...string) string {
	res := time.Time{}
	count := len(value)
	if count == 1 {
		if v, err := timeutils.AnyToTime(value[0], time.Time{}); err == nil {
			res = v
		}
	} else if count > 1 {
		res = timeutils.ToDateTime(value[0], value[1])
	}
	return res.Format(LocalTimeLayoutLine)
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
//  @Description: 取得支出金额
//  @param val
//  @return float64
//
func payout(val any) float64 {
	v := toFloat(val)
	return 0 - abs(v)
}

//
//  income
//  @Description: 取得收入金额
//  @param val
//  @return float64
//
func income(val any) float64 {
	v := toFloat(val)
	return abs(v)
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
