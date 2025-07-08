package sql2mongo

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
你是一个golang程序员，写一个方法将sql转换成mongo的Aggregate() Pipeline。
### 要求
1. 支持 select name,sex form users
2. 支持 select count(sex) form users
3. 支持 select * from users where name='张三' or name like '张三%' and sex = 1 支持and、 or、 is null、is not null, != ,>, <等所有比较符
4. 支持  select sum(sex) form users
5. 支持  SELECT SUM(age), AVG(money) FROM mcp_item WHERE name = '张三' AND date BETWEEN '2026-05-01' AND '2026-05-31'
*/

// Pipeline 将简单SQL转换为MongoDB Aggregate Pipeline
func Pipeline(sql string) ([]map[string]interface{}, error) {
	sql = strings.TrimSpace(sql)
	sqlLower := strings.ToLower(sql)
	pipeline := []map[string]interface{}{}

	// 1. 解析 select 字段
	selectRe := regexp.MustCompile(`(?i)^select\s+(.+?)\s+from\s+(\w+)`)
	selectMatches := selectRe.FindStringSubmatch(sqlLower)
	if len(selectMatches) < 3 {
		return nil, fmt.Errorf("无法解析SELECT语句")
	}
	selectFields := strings.TrimSpace(selectMatches[1])
	tableName := selectMatches[2]
	if tableName == "" {
		return nil, errors.New("tableName is empty")
	}

	// 2. 解析 where 条件
	whereClause := ""
	whereIdx := strings.Index(strings.ToLower(sql), " where ")
	if whereIdx != -1 {
		whereClause = sql[whereIdx+7:]
	}

	// 2.1 处理 where 子句，支持 and/or、like、=、is null、is not null、!=、>、<、>=、<=、between
	if whereClause != "" {
		matchExpr, err := parseWhereClause(whereClause)
		if err != nil {
			return nil, err
		}
		if len(matchExpr) > 0 {
			pipeline = append(pipeline, map[string]interface{}{"$match": matchExpr})
		}
	}

	// 2.2 解析 limit/offset/skip
	limit := -1
	skip := -1
	// 支持 limit n offset m 或 limit m,n 或 skip n
	limitRe := regexp.MustCompile(`(?i)limit\\s+(\\d+)(?:\\s*,\\s*(\\d+))?`)
	offsetRe := regexp.MustCompile(`(?i)offset\\s+(\\d+)`)
	skipRe := regexp.MustCompile(`(?i)skip\\s+(\\d+)`)

	if m := limitRe.FindStringSubmatch(sqlLower); len(m) > 0 {
		if len(m) > 2 && m[2] != "" {
			skip, _ = strconv.Atoi(m[1])
			limit, _ = strconv.Atoi(m[2])
		} else {
			limit, _ = strconv.Atoi(m[1])
		}
	}
	if m := offsetRe.FindStringSubmatch(sqlLower); len(m) > 1 {
		skip, _ = strconv.Atoi(m[1])
	}
	if m := skipRe.FindStringSubmatch(sqlLower); len(m) > 1 {
		skip, _ = strconv.Atoi(m[1])
	}

	// select sum(age), avg(money), count(money) ...
	if groupFields := parseGroupFields(selectFields); len(groupFields) > 0 {
		group := map[string]interface{}{"_id": nil}
		for _, gf := range groupFields {
			switch gf.fn {
			case "sum":
				group["sum_"+gf.field] = map[string]interface{}{"$sum": "$" + gf.field}
			case "avg":
				group["avg_"+gf.field] = map[string]interface{}{"$avg": "$" + gf.field}
			case "count":
				group["count_"+gf.field] = map[string]interface{}{"$sum": map[string]interface{}{"$cond": []interface{}{map[string]interface{}{"$ne": []interface{}{"$" + gf.field, nil}}, 1, 0}}}
			}
		}
		pipeline = append(pipeline, map[string]interface{}{"$group": group})
	}

	// select * ...
	if selectFields == "*" {
		if skip > -1 {
			pipeline = append(pipeline, map[string]interface{}{"$skip": skip})
		}
		if limit > -1 {
			pipeline = append(pipeline, map[string]interface{}{"$limit": limit})
		}
		return pipeline, nil
	}

	// 如果是聚合函数则不生成 $project
	if isAllAggFields(selectFields) {
		return pipeline, nil
	}

	// select name,sex ...
	fieldsArr := strings.Split(selectFields, ",")
	project := map[string]interface{}{}
	for _, f := range fieldsArr {
		f = strings.TrimSpace(f)
		project[f] = 1
	}
	pipeline = append(pipeline, map[string]interface{}{"$project": project})
	if skip > -1 {
		pipeline = append(pipeline, map[string]interface{}{"$skip": skip})
	}
	if limit > -1 {
		pipeline = append(pipeline, map[string]interface{}{"$limit": limit})
	}
	return pipeline, nil
}

// parseWhereClause 解析SQL where子句为mongo表达式
func parseWhereClause(where string) (map[string]interface{}, error) {
	where = strings.TrimSpace(where)
	// 递归解析and/or
	return parseOr(where)
}

// parseOr 解析 or
func parseOr(expr string) (map[string]interface{}, error) {
	parts := splitLogical(expr, "or")
	if len(parts) == 1 {
		return parseAnd(parts[0])
	}
	var orArr []interface{}
	for _, part := range parts {
		a, err := parseAnd(part)
		if err != nil {
			return nil, err
		}
		orArr = append(orArr, a)
	}
	return map[string]interface{}{"$or": orArr}, nil
}

// parseAnd 解析 and
func parseAnd(expr string) (map[string]interface{}, error) {
	parts := splitLogical(expr, "and")
	if len(parts) == 1 {
		return parseCondition(parts[0])
	}
	var andArr []interface{}
	for _, part := range parts {
		a, err := parseCondition(part)
		if err != nil {
			return nil, err
		}
		andArr = append(andArr, a)
	}
	return map[string]interface{}{"$and": andArr}, nil
}

// splitLogical 按逻辑运算符分割，忽略括号和BETWEEN ... AND ...中的运算符
func splitLogical(expr, op string) []string {
	var res []string
	level := 0
	start := 0
	opLower := strings.ToLower(op)
	lowerExpr := strings.ToLower(expr)
	i := 0
	for i < len(expr)-len(op)+1 {
		if expr[i] == '(' {
			level++
			i++
		} else if expr[i] == ')' {
			level--
			i++
		} else if level == 0 && strings.HasPrefix(lowerExpr[i:], "between") {
			// 跳过 between ... and ...
			j := i + len("between")
			for j < len(expr) {
				if strings.ToLower(expr[j:j+4]) == " and" {
					j += 4
					break
				}
				j++
			}
			i = j
		} else if level == 0 && strings.HasPrefix(lowerExpr[i:], opLower) &&
			(i == 0 || isSQLSpace(expr[i-1])) &&
			(i+len(op) == len(expr) || isSQLSpace(expr[i+len(op)])) {
			res = append(res, strings.TrimSpace(expr[start:i]))
			i += len(op)
			start = i
		} else {
			i++
		}
	}
	res = append(res, strings.TrimSpace(expr[start:]))
	return res
}

func isSQLSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n'
}

// parseCondition 支持 > < >= <= != = is null is not null like between
func parseCondition(cond string) (map[string]interface{}, error) {
	cond = strings.Trim(cond, "() \t\n")
	// between（兼容无引号和有引号的写法）
	if m := regexp.MustCompile(`(?i)^(\w+)\s+between\s+'?([^']+)'?\s+and\s+'?([^']+)'?$`).FindStringSubmatch(cond); len(m) == 4 {
		field, v1, v2 := m[1], m[2], m[3]
		return map[string]interface{}{field: map[string]interface{}{"$gte": tryParseValue(v1), "$lte": tryParseValue(v2)}}, nil
	}
	// is null / is not null
	if m := regexp.MustCompile(`(?i)^(\w+)\s+is\s+not\s+null$`).FindStringSubmatch(cond); len(m) == 2 {
		return map[string]interface{}{m[1]: map[string]interface{}{"$ne": nil}}, nil
	}
	if m := regexp.MustCompile(`(?i)^(\w+)\s+is\s+null$`).FindStringSubmatch(cond); len(m) == 2 {
		return map[string]interface{}{m[1]: nil}, nil
	}
	// like
	if m := regexp.MustCompile(`(?i)^(\w+)\s+like\s+'(.+)'$`).FindStringSubmatch(cond); len(m) == 3 {
		field, pattern := m[1], m[2]
		regex := "^" + regexp.QuoteMeta(strings.TrimRight(pattern, "%")) + ".*"
		return map[string]interface{}{field: map[string]interface{}{"$regex": regex}}, nil
	}
	// !=
	if m := regexp.MustCompile(`(?i)^(\w+)\s*!=\s*'?(.*?)'?(\s|$)`).FindStringSubmatch(cond); len(m) >= 3 {
		field, val := m[1], m[2]
		return map[string]interface{}{field: map[string]interface{}{"$ne": tryParseValue(val)}}, nil
	}
	// >=
	if m := regexp.MustCompile(`(?i)^(\w+)\s*>=\s*'?(.*?)'?(\s|$)`).FindStringSubmatch(cond); len(m) >= 3 {
		field, val := m[1], m[2]
		return map[string]interface{}{field: map[string]interface{}{"$gte": tryParseValue(val)}}, nil
	}
	// <=
	if m := regexp.MustCompile(`(?i)^(\w+)\s*<=\s*'?(.*?)'?(\s|$)`).FindStringSubmatch(cond); len(m) >= 3 {
		field, val := m[1], m[2]
		return map[string]interface{}{field: map[string]interface{}{"$lte": tryParseValue(val)}}, nil
	}
	// >
	if m := regexp.MustCompile(`(?i)^(\w+)\s*>\s*'?(.*?)'?(\s|$)`).FindStringSubmatch(cond); len(m) >= 3 {
		field, val := m[1], m[2]
		return map[string]interface{}{field: map[string]interface{}{"$gt": tryParseValue(val)}}, nil
	}
	// <
	if m := regexp.MustCompile(`(?i)^(\w+)\s*<\s*'?(.*?)'?(\s|$)`).FindStringSubmatch(cond); len(m) >= 3 {
		field, val := m[1], m[2]
		return map[string]interface{}{field: map[string]interface{}{"$lt": tryParseValue(val)}}, nil
	}
	// =
	if m := regexp.MustCompile(`(?i)^(\w+)\s*=\s*'?(.*?)'?(\s|$)`).FindStringSubmatch(cond); len(m) >= 3 {
		field, val := m[1], m[2]
		return map[string]interface{}{field: tryParseValue(val)}, nil
	}
	return nil, fmt.Errorf("无法解析条件: %s", cond)
}

// 尝试把字符串转为int、float或time.Time，否则返回原字符串
func tryParseValue(val string) interface{} {
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	// 尝试解析日期（仅支持常见格式）
	dateLayouts := []string{"2006-01-02", "2006-01-02 15:04:05", time.RFC3339, "2006/01/02"}
	for _, layout := range dateLayouts {
		if t, err := time.Parse(layout, val); err == nil {
			return t
		}
	}
	return val
}

// 解析 sum 字段列表
func parseSumFields(selectFields string) []string {
	fields := []string{}
	selectFields = strings.ReplaceAll(selectFields, " ", "")
	if !strings.Contains(selectFields, "sum(") {
		return nil
	}
	parts := strings.Split(selectFields, ",")
	for _, part := range parts {
		if m := regexp.MustCompile(`(?i)^sum\((\w+)\)$`).FindStringSubmatch(part); len(m) == 2 {
			fields = append(fields, m[1])
		}
	}
	return fields
}

// 解析 sum/avg/count 字段列表
// 返回 [{fn:"sum", field:"age"}, ...]
type groupField struct{ fn, field string }

func parseGroupFields(selectFields string) []groupField {
	fields := []groupField{}
	selectFields = strings.ReplaceAll(selectFields, " ", "")
	parts := strings.Split(selectFields, ",")
	for _, part := range parts {
		if m := regexp.MustCompile(`(?i)^(sum|avg|count)\((\w+)\)$`).FindStringSubmatch(part); len(m) == 3 {
			fields = append(fields, groupField{fn: strings.ToLower(m[1]), field: m[2]})
		}
	}
	return fields
}

// 判断是否全是聚合函数
func isAllAggFields(selectFields string) bool {
	selectFields = strings.ReplaceAll(selectFields, " ", "")
	if selectFields == "" {
		return false
	}
	parts := strings.Split(selectFields, ",")
	for _, part := range parts {
		if !regexp.MustCompile(`^(sum|avg|count)\(\w+\)$`).MatchString(strings.ToLower(part)) {
			return false
		}
	}
	return true
}
