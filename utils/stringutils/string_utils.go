package stringutils

import (
	"errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/inflection"
	"strconv"
	"strings"
)

//
// Int64ToString
// @Description:
// @param v
// @return string
//
func Int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}

func IsEmptyStr(v string) bool {
	if len(v) == 0 {
		return true
	}
	return false
}

func ValidEmptyStr(v string, msg string) error {
	if IsEmptyStr(v) {
		return errors.New(msg)
	}
	return nil
}

func AsFieldName(s string) string {
	res := strings.Replace(s, " ", "", -1)
	res = SnakeString(res)
	res = strings.Replace(res, "._", ".", -1)
	if strings.HasSuffix(res, "_") {
		res = res[1:]
	}
	return res
}

func Relpace(s string, old string, new string) string {
	return strings.Replace(s, "._", ".", -1)
}

//
// FirstUpper
// @Description: 字符串首字母大写
// @param s
// @return string
//
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	v := strings.ToUpper(s[:1]) + s[1:]
	if v == "" {
		println(v)
	}
	return v
}

//
// FirstLower
// @Description: 字符串首字母小写
// @param s
// @return string
//
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

//
// ToUpper
// @Description: 大写
// @param s
// @return string
//
func ToUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s)
}

//
// ToLower
// @Description: 小写
// @param s
// @return string
//
func ToLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s)
}

//
// SnakeString
// @Description: 驼峰转蛇形
// @param s 要转换的字符串
// @return string
//
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	res := strings.ToLower(string(data[:]))
	if strings.HasPrefix(res, "_") {
		return res[1:]
	}
	return res
}

//
// EqualFold
// @Description: 可以检查两个字符串是否相等,同时忽略大小写
// @param s
// @param t
// @return bool
func EqualFold(s, t string) bool {
	return strings.EqualFold(s, t)
}

//
// MidlineString
// @Description: 驼峰转中线
// @param s 要转换的字符串
// @return string
//
func MidlineString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '-')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	res := strings.ToLower(string(data[:]))
	if strings.HasPrefix(res, "-") {
		return res[1:]
	}
	return res
}

//
// CamelString 蛇形转驼峰
// @Description:
// @param s 要转换的字符串
// @return string
//
func CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

//
// Plural
// @Description: 将单词的单数形式转换为复数形式
// @param str 单数
// @return string 复数
//
func Plural(str string) string {
	return inflection.Plural(MidlineString(str))
}

//
// Singular
// @Description: 复数转单数
// @param str 复数
// @return string 单数
//
func Singular(str string) string {
	return inflection.Singular(MidlineString(str))
}

func PStrList(s ...string) *[]string {
	var res []string
	for _, item := range s {
		res = append(res, item)
	}
	return &res
}
