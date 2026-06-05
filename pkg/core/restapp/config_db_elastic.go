package restapp

type ElasticConfig struct {
	DbKey        string   `json:"dbKey"`
	Name         string   `yaml:"name" json:"name"`
	Addresses    []string `yaml:"addresses" json:"addresses"`
	Username     string   `yaml:"username" json:"username"`
	Password     string   `yaml:"password" json:"password"`
	APIKey       string   `yaml:"apiKey" json:"apiKey"`
	ServiceToken string   `yaml:"serviceToken" json:"serviceToken"`
	CloudID      string   `yaml:"cloudId" json:"cloudId"`
	CACert       string   `yaml:"caCert" json:"caCert"`
	ClientCert   string   `yaml:"clientCert" json:"clientCert"`
	ClientKey    string   `yaml:"clientKey" json:"clientKey"`
	MaxRetries   int      `yaml:"maxRetries" json:"maxRetries"`
	Debug        bool     `yaml:"debug" json:"debug"`
}
