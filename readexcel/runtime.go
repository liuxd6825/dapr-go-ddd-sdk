package readexcel

import (
	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/script"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"math"
	"regexp"
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

	if err := vm.Set("根据标识取支出金额", payoutByTag); err != nil {
		return nil, err
	}

	if err := vm.Set("income", income); err != nil {
		return nil, err
	}
	if err := vm.Set("取收入金额", income); err != nil {
		return nil, err
	}

	if err := vm.Set("根据标识取收入金额", incomeByTag); err != nil {
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

	if err := vm.Set("regexpNum", regexpNum); err != nil {
		return nil, err
	}
	if err := vm.Set("提取数字文本", regexpNum); err != nil {
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

func toDateTime(value ...string) (resText string) {
	defer func() {
		if e := recover(); e != nil {
			if err, ok := e.(error); ok {
				resText = err.Error()
			}
		}
	}()
	var err error
	res := time.Time{}
	count := len(value)
	if count == 1 {
		res, err = timeutils.AnyToTime(value[0], time.Time{})
	} else if count > 1 {
		res, err = timeutils.ToDateTime(value[0], value[1])
	}
	if err != nil {
		return err.Error()
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
		v = strings.ReplaceAll(v, "\t", "")
		v = strings.ReplaceAll(v, "\r", "")
		v = strings.ReplaceAll(v, "\n", "")
		v = strings.ReplaceAll(v, " ", "")
		v = strings.ReplaceAll(v, ",", "")
		v = strings.ReplaceAll(v, "，", "")
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

// isMinus
// @Description: 是否有负号
// @param val
// @return bool
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

// payout
// @Description: 取得支出金额
// @param val
// @return float64
func payout(val any) float64 {
	v := toFloat(val)
	return 0 - abs(v)
}

func payoutByTag(tagValue, tagName string, money any) float64 {
	if tagValue == tagName {
		return payout(money)
	}
	return 0
}

// income
// @Description: 取得收入金额
// @param val
// @return float64
func income(val any) float64 {
	v := toFloat(val)
	return abs(v)
}

func incomeByTag(tagValue, tagName string, money any) float64 {
	if tagValue == tagName {
		return income(money)
	}
	return 0
}

// amount
// @Description: 取交易金额
// @param val
// @return float64
func amount(v1 any, v2 any) float64 {
	a1 := abs(v1)
	a2 := abs(v2)
	if a1 != 0 {
		return a1
	}
	return a2
}

func regexpNum(val string, def string, idx int) string {
	re := regexp.MustCompile("[0-9]+")
	numbers := re.FindAllString(val, -1)
	lg := len(numbers)
	if lg == 1 || idx < 0 {
		return numbers[0]
	}
	if lg >= idx {
		return numbers[idx-1]
	}
	return def
}
