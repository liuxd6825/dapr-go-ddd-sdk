package config

type ToolConfig struct {
	Name        string             `json:"name"`
	CollName    string             `json:"coll_name"`
	DBName      string             `json:"db_name"`
	Description string             `json:"description"`
	SQLPrompt   string             `json:"sql_prompt"`
	Options     []ToolOptionConfig `json:"options"`
}

type ToolOptionConfig struct {
	Name        string `json:"name"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}
