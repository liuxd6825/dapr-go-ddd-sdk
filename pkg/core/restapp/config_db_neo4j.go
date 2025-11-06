package restapp

type Neo4jConfig struct {
	DbKey        string
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	Database     string `yaml:"dbname"`
	UserName     string `yaml:"user"`
	Password     string `yaml:"pwd"`
	EventPublish bool   `yaml:"eventPublish" json:"eventPublish"` // 是否发送领域事件
}
