package mongoutils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"time"
)

func GetMQL(p mongo.Pipeline) string {
	sb := strings.Builder{}
	sb.WriteString("[\n")
	for i, stage := range p {
		sb.WriteString(getAny(stage))
		if i != len(p)-1 {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("]\n")
	return sb.String()
}

func PrintPipeline(pipeline mongo.Pipeline) {
	mql := GetMQL(pipeline)
	print(mql)
}

func getMap(m map[string]any) string {
	sb := strings.Builder{}
	count := len(m)
	i := 0
	sb.WriteString("\n{")
	for key, val := range m {
		switch val.(type) {
		case map[string]any:
			mVal, _ := val.(map[string]any)
			sb.WriteString(fmt.Sprintf(`"%s"":`, key))
			sb.WriteString(getMap(mVal))
			break
		case primitive.M:
			mVal, _ := val.(primitive.M)
			sb.WriteString(fmt.Sprintf(`"%s":`, key))
			sb.WriteString(getMap(mVal))
			break
		case []any:
			array, _ := val.([]any)
			sb.WriteString(fmt.Sprintf(`"%s":%s`, key, getArray(array)))
			break
		case primitive.A:
			array, _ := val.(primitive.A)
			sb.WriteString(fmt.Sprintf(`"%s":%s"`, key, getArray(array)))
			break
		case primitive.Regex:
			regex, _ := val.(primitive.Regex)
			str := fmt.Sprintf(`"%s":{"$regex":"%s", "$options":"%s"}`, key, regex.Pattern, regex.Options)
			sb.WriteString(str)
		default:
			sb.WriteString(fmt.Sprintf(`"%s":"%v"`, key, val))
			break
		}
		if i < count-1 {
			sb.WriteString(",")
			i++
		}
	}
	sb.WriteString("}")
	return sb.String()
}

func getArray(array []any) string {
	sb := strings.Builder{}
	sb.WriteString("[")
	for i, v := range array {
		sb.WriteString(getAny(v))
		if i < len(array)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]")
	return sb.String()
}

func getAny(val any) string {
	switch val.(type) {
	case time.Time:
		return fmt.Sprintf(`"%s"`, val)
	case *time.Time:
		return fmt.Sprintf(`"%s"`, val)
	case string:
		return fmt.Sprintf(`"%s"`, val)
	case bson.D:
		m, _ := val.(bson.D)
		return getBsonD(m)
	case []any:
		return getArray(val.([]any))
	case map[string]any:
		m, _ := val.(map[string]any)
		return getMap(m)
	case primitive.M:
		m, _ := val.(primitive.M)
		return getMap(m)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func getBsonD(d bson.D) string {
	sb := strings.Builder{}
	sb.WriteString("{")
	for i, elem := range d {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(`"%s":%v`, elem.Key, getAny(elem.Value)))
	}
	sb.WriteString("}")
	return sb.String()
}
