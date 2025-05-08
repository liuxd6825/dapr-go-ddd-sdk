package db_pkg

type DB interface {
	Get(name string)
	NewModel(dbName, tableName string) *Model
}
