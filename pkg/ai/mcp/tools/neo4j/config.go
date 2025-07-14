package neo4j

type Config struct {
	Kind string     `json:"kind" yaml:"kind"`
	Tool ToolConfig `json:"tool" yaml:"tool"`
	DB   DBConfig   `json:"db" yaml:"db"`
}

type ToolConfig struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Prompt      string `json:"prompt" yaml:"prompt"`
	Schema      string `json:"schema" yaml:"schema"`
}

type DBConfig struct {
	Uri      string
	Username string
	Password string
}
