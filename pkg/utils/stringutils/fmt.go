package stringutils

import (
	"fmt"
	"strings"
	"time"
)

type FmtType int

const (
	FmtTypeVarchar FmtType = iota
	FmtTypeInteger
	FmtTypeTime
	FmtTypeDate
	FmtTypeJson
	FmtTypeString
)

type FmtVal struct {
	Key   string
	Value any
	Type  FmtType
}

type FmtBuilder struct {
	data map[string]FmtVal
}

func NewFmtBuilder() *FmtBuilder {
	return &FmtBuilder{
		data: make(map[string]FmtVal),
	}
}

func (f *FmtBuilder) Add(k string, v any, t FmtType) *FmtBuilder {
	f.data[k] = FmtVal{
		Key:   k,
		Value: v,
		Type:  t,
	}
	return f
}

func (f *FmtBuilder) Varchar(k string, v any) *FmtBuilder {
	return f.Add(k, v, FmtTypeVarchar)
}

func (f *FmtBuilder) Int(k string, v any) *FmtBuilder {
	return f.Add(k, v, FmtTypeInteger)
}

func (f *FmtBuilder) Time(k string, v *time.Time) *FmtBuilder {
	return f.Add(k, v, FmtTypeTime)
}

func (f *FmtBuilder) Date(k string, v *time.Time) *FmtBuilder {
	return f.Add(k, v, FmtTypeDate)
}

func (f *FmtBuilder) Json(k string, v any) *FmtBuilder {
	return f.Add(k, v, FmtTypeJson)
}

func (f *FmtBuilder) String(k string, v any) *FmtBuilder {
	return f.Add(k, v, FmtTypeString)
}

func (f *FmtBuilder) StringsJoin(k string, vals []string, sep string) *FmtBuilder {
	list := strings.Join(vals, sep)
	return f.String(k, list)
}

func (f *FmtBuilder) Build() map[string]FmtVal {
	return f.data
}

func (f *FmtBuilder) Format(format string) string {
	return Format(format, f.data)
}

func Format(format string, vals map[string]FmtVal) string {
	for k, v := range vals {
		var str = ""
		switch v.Type {
		case FmtTypeInteger:
			str = fmt.Sprintf("%d", v.Value)
		case FmtTypeTime:
			str = fmt.Sprintf("%d", v.Value)
		case FmtTypeDate:
			str = fmt.Sprintf("\"%s\"", v.Value)
		case FmtTypeString:
			str = fmt.Sprintf("%s", v.Value)
		case FmtTypeVarchar:
			str = fmt.Sprintf("\"%s\"", v.Value)
		case FmtTypeJson:
			str = fmt.Sprintf("%v", v.Value)
		default:
			panic("unhandled default case")
		}
		format = strings.Replace(format, "$<"+k+">", str, -1)
	}
	return format
}
