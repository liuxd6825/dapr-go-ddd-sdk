package config

import (
	"time"
)

type ElasticConfig struct {
	Addresses    []string    `json:"addresses" yaml:"addresses" title:"ES节点地址"`
	Username     string      `json:"username" yaml:"username" title:"用户名"`
	Password     string      `json:"password" yaml:"password" title:"密码"`
	APIKey       string      `json:"apiKey" yaml:"apiKey" title:"API Key"`
	ServiceToken string      `json:"serviceToken" yaml:"serviceToken" title:"Service Token"`
	CloudID      string      `json:"cloudId" yaml:"cloudId" title:"Elastic Cloud ID"`
	CACert       string      `json:"caCert" yaml:"caCert" title:"CA证书路径"`
	ClientCert   string      `json:"clientCert" yaml:"clientCert" title:"客户端证书路径"`
	ClientKey    string      `json:"clientKey" yaml:"clientKey" title:"客户端密钥路径"`
	Retry        RetryConfig `json:"retry" yaml:"retry" title:"重试配置"`
	Debug        bool        `json:"debug" yaml:"debug" title:"调试模式"`
}

type RetryConfig struct {
	MaxRetries     int           `json:"maxRetries" yaml:"maxRetries" title:"最大重试次数"`
	RetryOnTimeout bool          `json:"retryOnTimeout" yaml:"retryOnTimeout" title:"超时重试"`
	RetryBackoff   time.Duration `json:"retryBackoff" yaml:"retryBackoff" title:"重试间隔"`
}
