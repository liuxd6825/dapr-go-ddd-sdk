package config

const DBKey = "$db"
const EventBus = "pubsub"
const ImportAppId = "import"
const MasterAppId = "master"
const Neo4jDBKey = "$neo4jDb"
const RedisKey = "$redis"

func GetDBKey() string {
	return DBKey
}

func GetRedisKey() string {
	return RedisKey
}
