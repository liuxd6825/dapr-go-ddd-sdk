package strings_pkg

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
	"strconv"
	"strings"
	"time"
)

type StringsPkg struct {
}

func New() *StringsPkg {
	return NewStringsPkg()
}

func NewStringsPkg() *StringsPkg {
	return &StringsPkg{}
}

func (s *StringsPkg) ToByte(str string) []byte {
	return []byte(str)
}

func (s *StringsPkg) ToString(data any) string {
	if b, ok := data.([]byte); ok {
		return string(b)
	} else if s, ok := data.(string); ok {
		return s
	} else if m, ok := data.(map[string]interface{}); ok {
		objStr, err := jsonutils.Marshal(m)
		if err != nil {
			panic(err)
		}
		return objStr
	} else if i, ok := data.(int); ok {
		return strconv.Itoa(i)
	} else if f, ok := data.(float64); ok {
		return strconv.FormatFloat(f, 'f', -1, 64)
	} else if t, ok := data.(*time.Time); ok {
		return timeutils.ToTimeString(t)
	}
	return fmt.Sprintf("%v", data)
}

func (s *StringsPkg) Fmt(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (s *StringsPkg) Join(elems ...string) string {
	return strings.Join(elems, " ")
}

func (s *StringsPkg) ToLower(str string) string {
	return strings.ToLower(str)
}

func (s *StringsPkg) ToTitle(str string) string {
	return strings.ToTitle(str)
}

func (s *StringsPkg) ToUpper(str string) string {
	return strings.ToUpper(str)
}

func (s *StringsPkg) Index(str string, substr string) int {
	return strings.Index(str, substr)
}

func (s *StringsPkg) Split(str string, substr string) []string {
	return strings.Split(str, substr)
}

func (s *StringsPkg) LastIndex(str string, substr string) int {
	return strings.LastIndex(str, substr)
}

func (s *StringsPkg) TrimLeft(str string, cutset string) string {
	return strings.TrimLeft(str, cutset)
}

func (s *StringsPkg) TrimRight(str string, cutset string) string {
	return strings.TrimRight(str, cutset)
}

func (s *StringsPkg) TrimLeftAll(str string, cutset string) string {
	return strings.TrimLeft(str, cutset)
}

func (s *StringsPkg) TrimRightAll(str string, cutset string) string {
	return strings.TrimRight(str, cutset)
}

func (s *StringsPkg) Trim(str string, cutset string) string {
	return strings.Trim(str, cutset)
}

func (s *StringsPkg) Count(str string, sub string) int {
	return strings.Count(str, sub)
}

func (s *StringsPkg) Fields(str string) []string {
	return strings.Fields(str)
}

func (s *StringsPkg) HasPrefix(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

func (s *StringsPkg) HasSuffix(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

func (s *StringsPkg) Contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func (s *StringsPkg) HasPrefixFold(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

func (s *StringsPkg) HasSuffixFold(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

func (s *StringsPkg) HasPrefixFoldFold(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

func (s *StringsPkg) HasSuffixFoldFold(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}

func (s *StringsPkg) HasPrefixFoldFoldFold(str, prefix string) bool {
	return strings.HasPrefix(str, prefix)
}

func (s *StringsPkg) Len(str string) int {
	return len(str)
}
