package reflectutils

import (
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"strconv"
	"strings"
)

func ParseDescMap(data any) (map[string]string, error) {
	res := map[string]string{}
	obj := NewRefObj(data)
	fields := obj.FieldsAll()
	for _, f := range fields {
		desc, err := f.Tag("desc")
		if err != nil {
			return nil, err
		}
		if desc == "" {
			desc = f.Name()
		}
		val, err := f.Get()
		if err != nil {
			return nil, err
		}
		res[desc] = stringutils.AnyToString(val)
	}
	return res, nil
}
func ParseDesc(data any) (string, error) {
	sb, err := ParseDescOptions(data, "desc", "ulog")
	if err != nil {
		return "", err
	}
	return sb.String(), nil
}

func ParseDescOptions(data any, nameTag string, ulogTag string) (*strings.Builder, error) {
	res := &strings.Builder{}
	obj := NewRefObj(data)
	fields := obj.FieldsAll()

	for _, f := range fields {
		userLog, err := f.Tag(ulogTag)
		if err != nil {
			return nil, err
		}
		if userLog == "-" {
			continue
		}

		desc, err := f.Tag(nameTag)
		if err != nil {
			return nil, err
		}
		if desc == "" {
			desc = f.Name()
		}
		val, err := f.Get()
		if err != nil {
			return nil, err
		}

		res.WriteString(fmt.Sprintf("[%s]=`%s`; ", desc, stringutils.AnyToString(val)))
	}

	return res, nil
}

// ParseTag parses a golang struct tag into a map.
func ParseTag(tag string) (map[string]string, error) {
	res := map[string]string{}

	// This code is copied/modified from: reflect/type.go:
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := string(tag[:i])
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := string(tag[:i+1])
		tag = tag[i+1:]

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil, fmt.Errorf("Cannot unquote tag %s in %s: %s", name, tag, err.Error())
		}
		res[name] = value
	}

	return res, nil
}
