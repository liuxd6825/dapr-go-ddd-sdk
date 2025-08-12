package sql2mongo

import (
	"fmt"
	"regexp"
	"strings"
)

// SQLToMongoDBAggregate 将 SQL 查询转换为 MongoDB 聚合管道
func SQLToMongoDBAggregate(sql string) ([]map[string]interface{}, error) {
	// 去除多余的空格
	sql = strings.TrimSpace(sql)

	// 创建管道数组
	var pipeline []map[string]interface{}

	// 解析 SELECT 部分
	selectRe := regexp.MustCompile(`SELECT\s+(.*?)\s+FROM\s+(\w+)`)
	selectMatches := selectRe.FindStringSubmatch(sql)
	if len(selectMatches) < 3 {
		return nil, fmt.Errorf("无法解析 SELECT 子句")
	}
	selectFields := selectMatches[1]
	//tableName := selectMatches[2]

	// 解析 WHERE 部分
	whereRe := regexp.MustCompile(`WHERE\s+(.*)`)
	whereMatches := whereRe.FindStringSubmatch(sql)
	var whereClause string
	if len(whereMatches) > 1 {
		whereClause = whereMatches[1]
	}

	// 解析 GROUP BY 部分
	groupByRe := regexp.MustCompile(`GROUP BY\s+(.*?)\s+`)
	groupByMatches := groupByRe.FindStringSubmatch(sql)
	var groupByClause string
	if len(groupByMatches) > 1 {
		groupByClause = groupByMatches[1]
	}

	// 解析 ORDER BY 部分
	orderByRe := regexp.MustCompile(`ORDER BY\s+(.*?)\s+`)
	orderByMatches := orderByRe.FindStringSubmatch(sql)
	var orderByClause string
	if len(orderByMatches) > 1 {
		orderByClause = orderByMatches[1]
	}

	// 构建聚合管道
	// 第一步：$match - 根据 WHERE 条件过滤
	if whereClause != "" {
		matchStage := make(map[string]interface{})
		val, err := parseWhereClause(whereClause)
		if err != nil {
			panic(err)
		}
		matchStage["$match"] = val
		pipeline = append(pipeline, matchStage)
	}

	// 第二步：$group - 根据 GROUP BY 聚合数据
	if groupByClause != "" {
		groupStage := make(map[string]interface{})
		groupStage["$group"] = parseGroupByClause(groupByClause, selectFields)
		pipeline = append(pipeline, groupStage)
	}

	// 第三步：$sort - 根据 ORDER BY 排序
	if orderByClause != "" {
		sortStage := make(map[string]interface{})
		sortStage["$sort"] = parseOrderByClause(orderByClause)
		pipeline = append(pipeline, sortStage)
	}

	// 第四步：$project - 选择字段
	projectStage := make(map[string]interface{})
	projectStage["$project"] = parseSelectFields(selectFields)
	pipeline = append(pipeline, projectStage)

	return pipeline, nil
}

// 解析 WHERE 条件
func parseWhereClause2(where string) map[string]interface{} {
	// 示例：解析简单的条件，处理 AND、OR、=、!=、>、< 等
	where = strings.ReplaceAll(where, " AND ", " && ")
	where = strings.ReplaceAll(where, " OR ", " || ")
	where = strings.ReplaceAll(where, "=", ":")
	where = strings.ReplaceAll(where, ">", ":gt")
	where = strings.ReplaceAll(where, "<", ":lt")
	// 实际上需要更复杂的处理，进行具体的操作符替换等
	return map[string]interface{}{
		"condition": where,
	}
}

// 解析 SELECT 字段
func parseSelectFields(selectFields string) map[string]interface{} {
	fields := strings.Split(selectFields, ",")
	project := make(map[string]interface{})
	for _, field := range fields {
		// 处理 STRFTIME 转换为 MongoDB 的 $dateToString
		if strings.Contains(field, "STRFTIME") {
			dateRe := regexp.MustCompile(`STRFTIME\('(.+?)', (.*?)\) as (.+)`)
			dateMatches := dateRe.FindStringSubmatch(field)
			if len(dateMatches) > 0 {
				dateFormat := dateMatches[1]
				dateField := dateMatches[2]
				alias := dateMatches[3]
				project[alias] = map[string]interface{}{
					"$dateToString": map[string]interface{}{
						"format": dateFormat,
						"date":   "$" + dateField,
					},
				}
			}
		} else {
			project[strings.TrimSpace(field)] = 1
		}
	}
	return project
}

// 解析 GROUP BY
func parseGroupByClause(groupBy string, selectFields string) map[string]interface{} {
	groupByFields := strings.Split(groupBy, ",")
	group := make(map[string]interface{})
	for _, field := range groupByFields {
		group[strings.TrimSpace(field)] = "$" + strings.TrimSpace(field)
	}
	// 对于聚合函数的处理
	if strings.Contains(selectFields, "SUM") {
		group["total_amount"] = map[string]interface{}{"$sum": "$amount"}
	}
	return group
}

// 解析 ORDER BY
func parseOrderByClause(orderBy string) map[string]any {
	// 简单解析 ORDER BY，实际情况可能更复杂
	sortFields := strings.Split(orderBy, ",")
	sort := make(map[string]any)
	for _, field := range sortFields {
		sort[strings.TrimSpace(field)] = 1 // 默认升序
	}
	return sort
}
