package utils

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"regexp"
	"strings"
)

func GetNameAndAlis(line string) (name string, alias string) {
	if strings.HasSuffix(line, "{") {
		line = strings.Trim(line[:len(line)-2], " ")
	}
	list := strings.Split(line, " ")
	count := len(list)
	if count == 4 {
		return list[1], list[2]
	}
	if count == 2 {
		return list[1], ""
	} else if count >= 5 {
		return list[1], list[3]
	}
	return "", ""
}

func SymbolStr(line string, symbol string) (symbolSub string, residue string) {
	idx := strings.Index(line, symbol)
	if idx == -1 {
		return "", ""
	}
	count := len(symbol)
	symbolSub = line[:idx]
	if idx+count < len(line) {
		residue = line[idx+count:]
	}
	return
}

func ParseNamespaceName(line string) (name string, alias string) {
	if strings.HasSuffix(line, "{") {
		line = strings.Trim(line[:len(line)-2], " ")
	}
	list := strings.Split(line, " ")
	count := len(list)
	if count == 4 {
		return list[1], list[2]
	}
	if count == 2 {
		return list[1], ""
	} else if count >= 5 {
		return list[1], list[3]
	}
	return "", ""
}

func ParseEnumName(line string) (name string, alias string, namespaceName string) {
	name, alias, namespaceName, _ = ParseName(line, "enum")
	return
}

func ParseName(line string, tagName string) (name string, alias string, namespaceName string, list []string) {
	if strings.HasSuffix(line, "{") {
		line = strings.Trim(line[:len(line)-2], " ")
	}
	list = strings.Split(line, " ")
	count := len(list)
	if count > 0 && list[0] != tagName {
		return
	}

	// 解析 class com.parse.Parse  {
	str := list[1]
	idx := strings.LastIndex(str, ".")
	if idx > -1 {
		namespaceName = str[:idx]
		name = str[idx+1:]
	} else {
		name = str
		namespaceName = ""
	}

	// 解决 class com.parse.Parse as parse
	for i := 0; i < count; i++ {
		if list[i] == "as" && i < count-1 {
			alias = list[i+1]
			break
		}
	}
	return
}

func Mkdir(rootPath string, namespace string) (string, error) {
	path := rootPath + "/" + strings.ReplaceAll(namespace, ".", "/")
	return path, os.MkdirAll(path, os.ModePerm)
}

func StringSplit(str string, sep string) []string {
	return strings.Split(strings.Trim(str, " "), sep)
}

func StringPrefixSuffix(str string, prefix string, suffix string) bool {
	s1 := strings.HasPrefix(str, prefix)
	s2 := strings.HasSuffix(str, suffix)
	if s1 && s2 {
		return true
	}
	return false
}
func StringTrim(str string, cutsets ...string) string {
	if len(cutsets) == 0 {
		return strings.Trim(str, " ")
	}
	s := str
	for _, c := range cutsets {
		s = strings.Trim(s, c)
	}
	return s
}

func StringClear(str string, dels ...string) string {
	s := str
	for _, c := range dels {
		s = strings.ReplaceAll(s, c, "")
	}
	return s
}

func StringContains(str string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			return true
		}
	}
	return false
}

func StringBetween(str string, startTag string, endTag string) (isHas bool, res string, startIdx int, endIdx int, other string) {
	text := str
	s := strings.Index(text, startTag)
	if s == -1 {
		return false, "", -1, -1, ""
	}
	text = text[s+len(startTag):]
	e := strings.Index(text, endTag)
	if e == -1 || e == 0 {
		return false, "", -1, -1, ""
	}
	val := text[0:e]
	if e == 0 || strings.Trim(val, " ") == "" {
		return false, "", -1, -1, ""
	}

	val = strings.Trim(val, " ")
	return true, val, s + 1, s + e + 1, text[e:]
}

func StringBetweenAll(str string, startTag string, endTag string) (isHas bool, res string, other string) {
	other = str
	for {
		has, val, _, end, _ := StringBetween(other, startTag, endTag)
		if has {
			isHas = true
			res = res + " " + val
			other = other[end+1:]
		} else {
			break
		}
	}
	return
}

func GetHtmlText(body string) (string, error) {
	reader := strings.NewReader(body)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatalf("Error parsing HTML from response body: %v", err)
	}

	var text string
	doc.Find("body").Each(func(_ int, s *goquery.Selection) {
		text += s.Text() + "\n"
	})

	return text, nil
}

func IsHtml(str string) bool {
	// HTML标记的正则表达式模式
	pattern := `<[^>]+>`
	matched, _ := regexp.MatchString(pattern, str)
	return matched
}

func FirstUpperName(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func FirstLowerName(str string) string {
	if len(str) == 0 {
		return ""
	}
	return strings.ToLower(str[:1]) + str[1:]
}

// SnakeString
// @Description: 驼峰转蛇形
// @param s 要转换的字符串
// @return string
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

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

var tagUnmarshalRegexp = regexp.MustCompile("([a-zA-Z]\\w*)")

// TagUnmarshal
//
//	@Description: 增加json中key的双引号
func TagUnmarshal(str string, data any) (err error) {
	regStr := tagUnmarshalRegexp.ReplaceAllString(str, `"$1"`)
	regStr = fmt.Sprintf("{%s}", regStr)
	if err = json.Unmarshal([]byte(regStr), data); err != nil {
		return err
	}
	return err
}

var removeTagsRegexp = regexp.MustCompile("@(.*?)((.*?))$")

// RemoveTags
//
//	@Description: 删除@tagName()文字
func RemoveTags(str string) string {
	text := removeTagsRegexp.ReplaceAllString(str, "")
	return strings.Trim(text, "")
}

var isEndNoteRegexp = regexp.MustCompile("end(.*?)note")

func StringIsEndNote(str string) bool {
	matched := isEndNoteRegexp.MatchString(str)
	return matched
}
