package ddd_mongodb

import (
	"errors"
	"fmt"
	"github.com/liuxd6825/dapr-components-contrib/liuxd/common/utils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/rsql"
	"go.mongodb.org/mongo-driver/bson"
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
}

func NewQueryGroup(qry ddd_repository.FindPagingQuery) (*QueryGroup, error) {
	var err error
	f1 := qry.GetFilter()
	f2 := qry.GetMustFilter()
	f3 := ""
	mustWhere, ok := qry.(ddd_repository.FindPagingQueryMustWhere)
	if ok {
		f3, err = mustWhere.GetMustWhere()
		if err != nil {
			return nil, err
		}
	}
	filter := getRsqlAnds(f1, f2, f3)
	baseGroup := &QueryGroup{
		TenantId:  qry.GetTenantId(),
		Filter:    filter,
		GroupCols: qry.GetGroupCols(),
		GroupKeys: qry.GetGroupKeys(),
		ValueCols: qry.GetValueCols(),
		Sort:      qry.GetSort(),
	}
	return baseGroup, nil
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
func (b *QueryGroup) GetGroup() (bson.D, error) {
	if b.GroupCols == nil || len(b.GroupCols) == 0 {
		return nil, nil
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

	return group, nil
}

func (b *QueryGroup) GetTotalGroup() (bson.D, error) {
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
	return bson.D{{"$group", projectMap}}, nil
}

// GetFilter
// @Description: 不分组分页条件，即原始网格数据
// @receiver b
// @return map[string]interface{}
// @return error
func (b *QueryGroup) GetFilter() (map[string]interface{}, error) {
	if b.Filter == "" {
		return nil, nil
	}
	p := NewMongoProcess()
	if err := rsql.ParseProcess(b.Filter, p); err != nil {
		return nil, err
	}
	return p.GetFilter(b.TenantId)
}

// GetGroupExpandFilter
// @Description: 分组不分页，即分组全部展开时过滤器
// @receiver b
// @return map[string]interface{}
// @return error
func (b *QueryGroup) GetGroupExpandFilter() (map[string]interface{}, error) {
	mMatch, err := b.GetFilter()
	if err != nil {
		return nil, err
	}

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
			} else {
				val = append(val, map[string]interface{}{utils.SnakeString(f.Field): b.GroupKeys[i]})
			}
		}
		mMatch["$and"] = val
	}
	return mMatch, nil
}

// GetFilterSort
// @Description:
// @receiver b
// @return bson.D
// @return error
func (b *QueryGroup) GetFilterSort() (bson.D, error) {
	if len(b.Sort) == 0 {
		return bson.D{}, nil
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
			return nil, oerr
		}
		item := bson.E{Key: utils.SnakeString(name), Value: orderVal}
		res = append(res, item)
	}
	return res, nil
}

func (b *QueryGroup) GetBsonFilterSort() (bson.D, error) {
	_sort := bson.D{}
	flag := false
	if len(b.Sort) > 0 {
		list := strings.Split(b.Sort, ",")
		for _, s := range list {
			if flag {
				break
			}
			for _, _rowGroupCol := range b.GroupCols {
				if strings.Contains(s, _rowGroupCol.Field) {
					flag = true
					break
				}
			}
		}
	}
	if (len(b.Sort) == 0 || !flag) && b.IsGroup() {
		for _, _rowGroupCol := range b.GroupCols {
			_sort = append(_sort, bson.E{Key: utils.SnakeString(_rowGroupCol.Field), Value: 1})
		}
	}
	_sort1, _ := b.GetFilterSort()
	_sort = append(_sort, _sort1...)
	return bson.D{{"$sort", _sort}}, nil
}

func toDate(v interface{}) time.Time {
	if v == nil {
		return time.Time{}
	}
	_v := strings.Trim(v.(string), " ")
	if _v == "" {
		return time.Time{}
	}

	timeLayout := "2006-01-02T15:04:05+08:00" //转化所需模板
	loc, _ := time.LoadLocation("Local")      //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, _v, loc)
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
