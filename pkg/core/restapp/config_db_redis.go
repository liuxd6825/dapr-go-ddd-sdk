package restapp

import (
	"time"
)

type RedisConfig struct {
	DbKey           string
	Host            string  `yaml:"host"`
	Database        int     `yaml:"db"`
	Password        string  `yaml:"pwd"`
	PoolSize        *int    `yaml:"poolSize"`
	ReadTimeout     *string `yaml:"readTimeout"`
	WriteTimeout    *string `yaml:"writeTimeout"`
	DialTimeout     *string `yaml:"dialTimeout"`
	PoolTimeout     *string `yaml:"poolTimeout"`
	MaxRetries      *int    `yaml:"maxRetries"`
	MinRetryBackoff *string `yaml:"minRetryBackoff"`
	MaxRetryBackoff *string `yaml:"maxRetryBackoff"`
}

func parseInt(name string, val *int, defaultVal int) int {
	if val == nil {
		return defaultVal
	}
	return *val
}

func parseDuration(name string, val *string, defaultDur time.Duration) time.Duration {
	if val == nil || *val == "" {
		return defaultDur
	}
	dur, err := time.ParseDuration(*val)
	if err != nil {
		panic(name + err.Error())
	}
	return dur
}
