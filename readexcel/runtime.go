package readexcel

import (
	"fmt"
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
	if err := vm.Set("文字替换", replace); err != nil {
		return nil, err
	}
	if err := vm.Set("取时间", toDateTime); err != nil {
		return nil, err
	}
	if err := vm.Set("取绝对值", abs); err != nil {
		return nil, err
	}
	if err := vm.Set("是否有负号", isMinus); err != nil {
		return nil, err
	}
	if err := vm.Set("取支出金额", payout); err != nil {
		return nil, err
	}
	if err := vm.Set("根据标识取支出金额", payoutByTag); err != nil {
		return nil, err
	}
	if err := vm.Set("取收入金额", income); err != nil {
		return nil, err
	}
	if err := vm.Set("根据标识取收入金额", incomeByTag); err != nil {
		return nil, err
	}
	if err := vm.Set("取交易金额", amount); err != nil {
		return nil, err
	}
	if err := vm.Set("取浮点值", toFloat); err != nil {
		return nil, err
	}
	if err := vm.Set("取文本", toString); err != nil {
		return nil, err
	}
	if err := vm.Set("取数字文本", regexpNum); err != nil {
		return nil, err
	}
	if err := vm.Set("取中间文本", match); err != nil {
		return nil, err
	}
	if err := vm.Set("取支付宝账号", GetZhiFuBaoOppAcc); err != nil {
		return nil, err
	}
	if err := vm.Set("取支付宝人名", GetZhiFuBaoOppName); err != nil {
		return nil, err
	}

	return vm, nil
}

func toString(list ...string) string {
	for _, s := range list {
		if s != "" {
			return s
		}
	}
	return ""
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

// match
//
//	@Description: 提取中间文本
//	@param str  要提取的文本
//	@param begin 以..开始
//	@param end 以..结束
//	@param replace 替换的字符串数组
//	@return string
func match(str, begin, end string, replace ...string) string {
	compileRegex := regexp.MustCompile(fmt.Sprintf("%s(.*?)%s", begin, end)) // 正则表达式的分组，以括号()表示，每一对括号就是我们匹配到的一个文本，可以把他们提取出来。
	matchArr := compileRegex.FindStringSubmatch(str)                         // FindStringSubmatch 方法是提取出匹配的字符串，然后通过[]string返回。我们可以看到，第1个匹配到的是这个字符串本身，从第2个开始，才是我们想要的字符串。
	if len(matchArr) > 0 {
		s := matchArr[0]
		s, _ = strings.CutPrefix(s, begin)
		s, _ = strings.CutSuffix(s, end)
		for _, r := range replace {
			s = strings.ReplaceAll(s, r, "")
		}
		return s
	}
	return ""
}

// GetZhiFuBaoOppAcc
//
//	@Description: 从摘要中取支付宝转账账号
//	@param companyName  支付宝公司的名称
//	@param companyValue 对手名称
//	@param remarks 摘要
//	@param defStr 默认值
//	@param idx 对
//	@return string
func GetZhiFuBaoOppAcc(companyName, companyValue string, remarks string, defStr string, idx int) string {
	if strings.Contains(companyValue, companyName) {
		return regexpNum(remarks, defStr, idx)
	}
	return defStr
}

// GetZhiFuBaoOppName
//
//	@Description: 从摘要中取支付宝转账人姓名
//	@param companyName  支付宝公司的名称
//	@param companyValue 对手名称
//	@param remarks 摘要
//	@param defStr 默认值
//	@param idx 对
//	@return string
func GetZhiFuBaoOppName(companyName, companyValue string, remarks string, begin, end string, defStr string, idx int, replace ...string) string {
	if strings.Contains(companyValue, companyName) {
		return match(remarks, begin, end, replace...)
	}
	return defStr
}
