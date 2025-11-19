package readexcel

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/dop251/goja"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/import/pkg/readexcel/script"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/timeutils"
)

const LocalTimeLayoutLine = "2006-01-02 15:04:05"

type Dictionary interface {
	GetName() string
	GetKeywords() string
}

type Runtime struct {
	*goja.Runtime
	options *RuntimeOptions
}

type RuntimeOptions struct {
	CurrencyMap map[string]Dictionary // 币种
	BankMap     map[string]Dictionary // 开户行
}

func newRuntime(options RuntimeOptions) (*Runtime, error) {
	vm, err := script.NewRuntime()
	if err != nil {
		return nil, err
	}
	runtime := &Runtime{Runtime: vm, options: &options}

	setValue(vm, replace, "文字替换", "replace")
	setValue(vm, toDateTime, "取时间", "toDateTime")
	setValue(vm, abs, "取绝对值", "abs")
	setValue(vm, isMinus, "是否有负号", "isMinus")
	setValue(vm, payout, "取支出金额", "payout")
	setValue(vm, payoutByTag, "根据标识取支出金额", "payoutByTag")
	setValue(vm, income, "取收入金额", "income")
	setValue(vm, incomeByTag, "根据标识取收入金额", "incomeByTag")
	setValue(vm, amount, "取交易金额", "amount")
	setValue(vm, toFloat, "取浮点值", "toFloat")
	setValue(vm, toString, "取文本", "toString")
	setValue(vm, regexpNum, "取数字文本", "regexpNum")
	setValue(vm, match, "取中间文本", "match")
	setValue(vm, GetZhiFuBaoOppAcc, "取支付宝账号", "getZhiFuBaoOppAcc")
	setValue(vm, GetZhiFuBaoOppName, "取支付宝人名", "getZhiFuBaoOppName")
	setValue(vm, Contains, "包含", "contains")
	setValue(vm, NotContains, "不包含", "notContains")
	setValue(vm, IsCash, "是否现金交易", "isCash")
	setValue(vm, runtime.GetCurrencyType, "取币种", "getCurrencyType")
	setValue(vm, runtime.GetBankName, "取开户行", "getBankName")

	return runtime, nil
}
func (r *Runtime) GetRuntime() *goja.Runtime {
	return r.Runtime
}

func (r *Runtime) GetBankName(strList ...string) string {
	return GetDictionary(r.options.BankMap, "", strList...)
}

func (r *Runtime) GetCurrencyType(strList ...string) string {
	return GetDictionary(r.options.CurrencyMap, "人民币", strList...)
}

func GetDictionary(distMap map[string]Dictionary, defaultVal string, strList ...string) string {
	for _, str := range strList {
		if _, ok := distMap[str]; ok {
			return str
		}
	}
	for _, str := range strList {
		for _, cur := range distMap {
			keywords := strings.Split(cur.GetKeywords(), ",")
			for _, key := range keywords {
				if key != "" && strings.Contains(str, key) {
					return cur.GetName()
				}
			}
		}
	}
	if strList == nil || len(strList) == 0 {
		return defaultVal
	}
	return strList[0]
}

func setValue(vm *goja.Runtime, value any, names ...string) {
	for _, name := range names {
		if name != "" {
			_ = vm.Set(name, value)
		}
	}
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

func toDateTime(value ...string) (res string) {
	var err error
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()

	timeStr := strings.Join(value, " ")
	res, err = timeutils.FormatStr(LocalTimeLayoutLine, timeStr)
	if err != nil {
		return ""
	}
	return res
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
func payout(val any) *float64 {
	if isNull(val) {
		return nil
	}
	v := toFloat(val)
	v = 0 - abs(v)
	return getMoney(v)
}

func getMoney(val float64) *float64 {
	if val == 0 {
		return nil
	}
	return &val
}

func isNull(val any) bool {
	if val == nil {
		return true
	}
	if v, ok := val.(string); ok {
		if v == "" {
			return true
		}
	}
	return false
}

func payoutByTag(tagValue, tagName string, money any) *float64 {
	if tagValue == tagName {
		return payout(money)
	}
	return nil
}

// income
// @Description: 取得收入金额
// @param val
// @return float64
func income(val any) *float64 {
	if isNull(val) {
		return nil
	}
	v := toFloat(val)
	v = abs(v)
	return getMoney(v)
}

func incomeByTag(tagValue, tagName string, money any) *float64 {
	if tagValue == tagName {
		return income(money)
	}
	return nil
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

func Contains(txt string, sub ...string) bool {
	for _, s := range sub {
		if strings.Contains(txt, s) {
			return true
		}
	}
	return false
}

func NotContains(txt string, sub ...string) bool {
	for _, s := range sub {
		if strings.Contains(txt, s) {
			return false
		}
	}
	return true
}

func IsCash(keyText string, values ...string) string {
	text := strings.Join(values, ",")
	keys := strings.Split(keyText, ",")
	for _, key := range keys {
		if strings.Contains(text, key) {
			return "是"
		}
	}
	return "否"
}
