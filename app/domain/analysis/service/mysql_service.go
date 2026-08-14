package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/analysis/query"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"gorm.io/gorm"
)

type MysqlService struct {
	env      *env.Env
	rootPath string
	db       *gorm.DB
}

func NewMysqlService(env *env.Env, rootPath string) *MysqlService {
	dbKey := env.GetDBKeyValue("doris")
	dbItem := env.GetDB(dbKey)
	if dbItem == nil {
		panic("mysqlDB is nil")
	}
	db, ok := dbItem.GetDB().(*gorm.DB)
	if !ok {
		panic("mysqlDB is nil")
	}
	return &MysqlService{
		env:      env,
		rootPath: rootPath,
		db:       db,
	}
}

func (s *MysqlService) Aggregate(ctx context.Context, qry *query.MysqlQuery) (any, error) {
	qry.Init()

	match := qry.GetMatch()
	whereClause := qry.BuildWhereClause(match)
	groupClause := qry.BuildGroupClause()
	selectClause := qry.BuildSelectClause()
	orderClause := qry.BuildOrderClause()

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s GROUP BY %s ORDER BY %s",
		selectClause, qry.Collection, whereClause, groupClause, orderClause)

	rows, err := s.db.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	results := make([]map[string]any, 0)

	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any)
		for i, col := range columns {
			row[col] = convertValue(values[i])
		}
		results = append(results, row)
	}

	if qry.Options.Sort != nil {
		results = convertTimeValues(results, qry.Options.Sort.TimeKey)
		results, _ = CompleteAndSortMysqlData(results, qry.Options.Sort.TimeKey, qry.Options.Sort.TimeType)
	}

	return results, nil
}

func (s *MysqlService) Summary(ctx context.Context, qry *query.SummaryQuery) (any, error) {
	sql := fmt.Sprintf("SELECT "+
		"case_id,name,opp_name, min_date,max_date,total_amount,total_count,total_acct FROM "+
		"("+
		" SELECT case_id,name,opp_name,MIN(DATE) AS min_date,MAX(DATE) AS max_date,ROUND(SUM(amount),2) AS total_amount,COUNT(id) AS total_count, COUNT(DISTINCT opp_acct) as total_acct "+
		" FROM master_record "+
		" WHERE 1=1 and case_id='%s'"+
		" %s "+
		" GROUP BY case_id, name, opp_name"+
		") t "+
		" WHERE 1=1 %s "+
		" order by total_amount desc", qry.CaseId, getDateCondition(qry), getCondition(qry))

	rows, err := s.db.Raw(sql).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	results := make([]map[string]any, 0)

	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any)
		for i, col := range columns {
			row[col] = convertValue(values[i])
		}
		results = append(results, row)
	}

	return results, nil
}

func convertTimeValues(data []map[string]any, timeKey string) []map[string]any {
	for _, item := range data {
		if v, ok := item[timeKey]; ok {
			switch val := v.(type) {
			case []uint8:
				if t, err := time.Parse("2006-01-02", string(val)); err == nil {
					item[timeKey] = t
				}
			case string:
				if t, err := time.Parse("2006-01-02", val); err == nil {
					item[timeKey] = t
				}
			case time.Time:
				item[timeKey] = val
			}
		}
	}
	return data
}

func convertValue(v any) any {
	switch val := v.(type) {
	case []uint8:
		str := string(val)
		if f, err := strconv.ParseFloat(str, 64); err == nil {
			return f
		}
		return str
	}
	return v
}

func CompleteAndSortMysqlData(data []map[string]any, timeKey string, timeType model.SortTimeType) ([]map[string]any, error) {
	if len(data) < 2 {
		return data, nil
	}

	completedData := make([]map[string]any, 0, len(data)*2)

	for i := 0; i < len(data)-1; i++ {
		completedData = append(completedData, data[i])

		currentTime := data[i][timeKey].(time.Time)
		nextTime := data[i+1][timeKey].(time.Time)

		tempTime := currentTime
		for {
			switch timeType {
			case model.Year:
				tempTime = tempTime.AddDate(1, 0, 0)
			case model.Month:
				tempTime = tempTime.AddDate(0, 1, 0)
			case model.Day:
				tempTime = tempTime.AddDate(0, 0, 1)
			default:
				return completedData, nil
			}

			if tempTime.Before(nextTime) {
				missingEntry := map[string]any{
					timeKey: tempTime,
				}
				completedData = append(completedData, missingEntry)
			} else {
				break
			}
		}
	}

	completedData = append(completedData, data[len(data)-1])
	return completedData, nil
}

func getDateCondition(qry *query.SummaryQuery) string {
	condition := ""
	if qry.DateMin != "" {
		condition += fmt.Sprintf(" and date >= '%s'", qry.DateMin)
	}
	if qry.DateMax != "" {
		condition += fmt.Sprintf(" and date <= '%s 23:59:59'", qry.DateMax)
	}
	return condition
}

func getCondition(qry *query.SummaryQuery) string {
	condition := ""
	if qry.Name != "" {
		condition += fmt.Sprintf(" and name='%s'", qry.Name)
	}
	if qry.OppName != "" {
		oppName := fmt.Sprintf("%v%v%v", "'", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(qry.OppName, " ", ""), "，", ","), ",", "','"), "'")
		condition += fmt.Sprintf(" and opp_name in (%s)", oppName)
	}
	if qry.TotalAmountMin != nil {
		condition += fmt.Sprintf(" and total_amount >= %v", qry.TotalAmountMin)
	}
	if qry.TotalAmountMax != nil {
		condition += fmt.Sprintf(" and total_amount <= %v", qry.TotalAmountMax)
	}
	if qry.TotalCountMin != nil {
		condition += fmt.Sprintf(" and total_count >= %v", qry.TotalCountMin)
	}
	if qry.TotalCountMax != nil {
		condition += fmt.Sprintf(" and total_count <= %v", qry.TotalCountMax)
	}
	return condition
}
