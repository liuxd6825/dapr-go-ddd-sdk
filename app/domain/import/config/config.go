package config

const DBKey = "$db"
const EventBus = "pubsub"
const ImportAppId = "import"
const Neo4jDBKey = "$neo4jDb"

func GetDBKey() string {
	return DBKey
}
