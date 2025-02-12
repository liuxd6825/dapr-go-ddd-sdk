package ddd_sql

import "github.com/liuxd6825/dapr-go-ddd-sdk/rsql"

func GetSQL(rSql string) (string, error) {
	res, err := rsql.SqlParseProcess(rSql)
	return res, err
}
