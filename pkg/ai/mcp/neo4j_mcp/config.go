package neo4j_mcp

type Config struct {
	Name  string
	Addr  string
	Neo4j Neo4jConfig
}

type Neo4jConfig struct {
	Uri      string
	Username string
	Password string
}
