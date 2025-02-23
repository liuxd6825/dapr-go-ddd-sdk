package ddd_mongodb

import (
	"errors"
	"fmt"
	"github.com/dapr/components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql/rsql_mongo"
	"go.mongodb.org/mongo-driver/bson"
	"strconv"
	"strings"
	"time"
)

type QueryGroup struct {
	TenantId  string
	Filter    string
	GroupCols []*ddd_repository.GroupCol
	ValueCols []*ddd_repository.ValueCol
	GroupKeys []any
	Sort      string
	Query     ddd_repository.FindPagingQuery
}

func NewQueryGroup(qry ddd_repository.FindPagingQuery) *QueryGroup {
	var err error
	if qry == nil {
		panic(errors.New("query is nil"))
	}
	f1 := qry.GetFilter()
	f2 := qry.GetMustFilter()
	f3 := ""
	mustWhere, ok := qry.(ddd_repository.FindPagingQueryMustWhere)
	if ok {
		f3, err = mustWhere.GetMustWhere()
		if err != nil {
			panic(err)
		}
	}
	filter := getRsqlAnds(f1, f2, f3)
	baseGroup := &QueryGroup{
		Query:     qry,
		TenantId:  qry.GetTenantId(),
		Filter:    filter,
		GroupCols: qry.GetGroupCols(),
		GroupKeys: qry.GetGroupKeys(),
		ValueCols: qry.GetValueCols(),
		Sort:      qry.GetSort(),
	}
	return baseGroup
}

// IsPaging
// @Description:
// @receiver b
// @return bool
func (b *QueryGroup) IsPaging() bool {
	if !b.IsGroup() {
		return true
	}

	if !b.IsExpand() {
		return true
	}
	return false
}

// IsGroup
// @Description:
// @receiver b
// @return bool
func (b *QueryGroup) IsGroup() bool {
	if b.GroupCols == nil || len(b.GroupCols) == 0 {
		return false
	}
	return true
}

// IsExpand
// @Description: 分组是否展开
// @receiver b
// @return bool
func (b *QueryGroup) IsExpand() bool {
	if b.GroupKeys == nil || len(b.GroupKeys) == 0 {
		return false
	}
	return true
}

// IsLeaf 是树型查询的子数据
func (b *QueryGroup) IsLeaf() bool {
	if b.IsGroup() && b.IsExpand() && len(b.GroupCols) == len(b.GroupKeys) {
		return true
	}
	return false
}

// GetGroup
// @Description:
// @receiver b
// @return bson.D
// @return error
func (b *QueryGroup) GetGroup() bson.D {
	if b.GroupCols == nil || len(b.GroupCols) == 0 {
		return nil
	}

	gSubMap := make(map[string]any)
	groupIndex := 0
	if b.GroupKeys != nil && len(b.GroupKeys) > 0 && len(b.GroupKeys) < len(b.GroupCols) {
		groupIndex = len(b.GroupKeys)
	}

	ids := make([]any, 0)
	for i := 0; i <= groupIndex; i++ {
		col := b.GroupCols[i]
		var newId interface{} = map[string]interface{}{"$toString": "$" + utils.SnakeString(col.Field)}
		if col.DataType.IsDateTime() || col.DataType.IsDate() {
			newId = map[string]any{"$dateToString": map[string]any{"date": "$" + utils.SnakeString(col.Field)}}
		}
		if i == 0 {
			ids = append(ids, newId)
		} else {
			ids = append(ids, "_")
			ids = append(ids, newId)
		}
	}

	gSubMap["_id"] = map[string]any{"$concat": ids}
	field := utils.SnakeString(b.GroupCols[groupIndex].Field)
	gSubMap[field] = map[string]any{"$max": "$" + field}

	if b.ValueCols != nil && len(b.ValueCols) > 0 {
		for _, col := range b.ValueCols {
			gSubMap[utils.SnakeString(col.Field)] = map[string]any{"$" + col.AggFunc.Name(): "$" + utils.SnakeString(col.Field)}
		}
	}

	group := bson.D{{
		"$group", gSubMap,
	}}

	return group
}

func (b *QueryGroup) GetTotalGroup() bson.D {
	projectMap := make(map[string]interface{})
	projectMap["_id"] = "null"
	pushMap := make(map[string]interface{})
	pushMap["_id"] = "$_id"

	groupIndex := 0
	if b.GroupKeys != nil && len(b.GroupKeys) > 0 && len(b.GroupKeys) < len(b.GroupCols) {
		groupIndex = len(b.GroupKeys)
	}
	if b.GroupCols != nil && len(b.GroupCols) > 0 {
		pushMap[utils.SnakeString(b.GroupCols[groupIndex].Field)] = "$" + utils.SnakeString(b.GroupCols[groupIndex].Field)
	}
	if b.ValueCols != nil && len(b.ValueCols) > 0 {
		for _, col := range b.ValueCols {
			pushMap[utils.SnakeString(col.Field)] = "$" + utils.SnakeString(col.Field)
		}
	}
	projectMap["data"] = map[string]interface{}{
		"$push": pushMap,
	}
	projectMap["total_rows"] = map[string]interface{}{"$sum": 1}
	return bson.D{{"$group", projectMap}}
}

func (q *QueryGroup) GetPageNum() int64 {
	return q.Query.GetPageNum()
}

func (q *QueryGroup) GetPageSize() int64 {
	return q.Query.GetPageSize()
}

// GetFilter
// @Description: 不分组分页条件，即原始网格数据
// @receiver b
// @return map[string]interface{}
// @return error
func (b *QueryGroup) GetFilter() *rsql_mongo.Filter {
	if b.Filter == "" {
		return rsql_mongo.NewMongoFilter()
	}

	p := rsql_mongo.NewProcess(b.TenantId)
	if err := rsql.ParseProcess(b.Filter, p); err != nil {
		panic(err)

	}
	filter := p.GetFilter().(*rsql_mongo.Filter)
	return filter
}

// GetGroupExpandFilter
// @Description: 分组不分页，即分组全部展开时过滤器
// @receiver b
// @return map[string]interface{}
// @return error
func (b *QueryGroup) GetGroupExpandFilter() *rsql_mongo.Filter {
	filter := b.GetFilter()
	mMatch := filter.Match

	if mMatch == nil {
		mMatch = make(map[string]interface{})
	}

	if b.GroupKeys != nil && len(b.GroupKeys) > 0 {
		subMap, ok := mMatch["$and"]
		if !ok {
			subMap = make([]interface{}, 0)
		}
		val, _ := subMap.([]interface{})
		for i := 0; i < len(b.GroupKeys); i++ {
			f := b.GroupCols[i]
			if f.DataType.IsDate() || f.DataType.IsDateTime() {
				val = append(val, map[string]interface{}{utils.SnakeString(f.Field): toDate(b.GroupKeys[i])})
			} else if f.DataType.IsFloat() || f.DataType.IsInt() || f.DataType.IsMoney() || f.DataType.IsYear() || f.DataType.IsMonth() || f.DataType.IsDay() {
				val = append(val, map[string]interface{}{utils.SnakeString(f.Field): toNumber(b.GroupKeys[i])})
			} else {
				val = append(val, map[string]interface{}{utils.SnakeString(f.Field): b.GroupKeys[i]})
			}
		}
		mMatch["$and"] = val
	}
	return filter
}

func toNumber(v interface{}) *float64 {
	if v == nil {
		return nil
	}
	val := strings.Trim(v.(string), " ")
	if val == "" {
		return nil
	}

	num, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return nil
	}
	return &num
}

// GetFilterSort
// @Description:
// @receiver b
// @return bson.D
// @return error
func (b *QueryGroup) GetFilterSort() bson.D {
	if len(b.Sort) == 0 {
		return bson.D{}
	}
	// 输入
	// name:desc,id:asc
	// 输出
	/*	sort := bson.D{
		bson.E{"update_time", -1},
		bson.E{"goods_id", -1},
	}*/
	res := bson.D{}
	list := strings.Split(b.Sort, ",")
	for _, s := range list {
		sortItem := strings.Split(s, ":")
		name := sortItem[0]
		name = strings.Trim(name, " ")
		if name == "id" {
			name = "_id"
		}
		order := "asc"
		if len(sortItem) > 1 {
			order = sortItem[1]
			order = strings.ToLower(order)
			order = strings.Trim(order, " ")
		}

		// 其中 1 为升序排列，而-1是用于降序排列.
		orderVal := 1
		var oerr error
		switch order {
		case "asc":
			orderVal = 1
		case "desc":
			orderVal = -1
		default:
			oerr = errors.New("order " + order + " is error")
		}
		if oerr != nil {
			print(oerr)
		}
		item := bson.E{Key: utils.SnakeString(name), Value: orderVal}
		res = append(res, item)
	}
	return res
}

func (b *QueryGroup) GetBsonFilterSort() bson.D {
	sort := bson.D{}
	flag := false
	if len(b.Sort) > 0 {
		list := strings.Split(b.Sort, ",")
		for _, s := range list {
			if flag {
				break
			}
			for _, rowGroupCol := range b.GroupCols {
				if strings.Contains(s, rowGroupCol.Field) {
					flag = true
					break
				}
			}
		}
	}
	if (len(b.Sort) == 0 || !flag) && b.IsGroup() {
		for _, rowGroupCol := range b.GroupCols {
			sort = append(sort, bson.E{Key: utils.SnakeString(rowGroupCol.Field), Value: 1})
		}
	}
	sort1 := b.GetFilterSort()
	if len(sort1) > 0 {
		sort = append(sort, sort1...)
	}
	if len(sort) == 0 {
		return nil
	}
	return bson.D{{"$sort", sort}}
}

func toDate(v interface{}) time.Time {
	if v == nil {
		return time.Time{}
	}
	val := strings.Trim(v.(string), " ")
	if val == "" {
		return time.Time{}
	}

	timeLayout := "2006-01-02T15:04:05+08:00" //转化所需模板
	loc, _ := time.LoadLocation("Local")      //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, val, loc)
	return theTime
}

func getRsqlAnds(s ...string) string {
	res := ""
	for _, item := range s {
		res = getRsqlAnd(res, item)
	}
	return res
}

func getRsqlAnd(s1 string, s2 string) string {
	b1 := len(s1) > 0
	b2 := len(s2) > 0
	if b1 && b2 {
		return fmt.Sprintf("(%s) and (%s)", s1, s2)
	} else if b1 {
		return s1
	}
	return s2
}
