package idao

type DbType string

const (
	DbType_Postgres DbType = "postgres"
	DbType_MySQL    DbType = "mysql"
	DbType_Sqlite   DbType = "sqlite"
	DbType_Neo4j    DbType = "neo4j"
	DbType_MongoDB  DbType = "mongodb"
	DbType_MsSQL    DbType = "mssql"
	DbType_Oracle   DbType = "oracle"
)

func (t DbType) String() string {
	return string(t)
}
