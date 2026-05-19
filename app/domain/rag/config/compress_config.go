package config

const (
	DefaultCompressThreshold = 50000
	DefaultLockTimeout       = 30
	DefaultMaxRetries        = 3
)

type CompressConfig struct {
	Threshold     int  `json:"threshold"`
	LockTimeout   int  `json:"lockTimeout"`
	MaxRetries    int  `json:"maxRetries"`
	UseCheapModel bool `json:"useCheapModel"`
}

func DefaultCompressConfig() *CompressConfig {
	return &CompressConfig{
		Threshold:   DefaultCompressThreshold,
		LockTimeout: DefaultLockTimeout,
		MaxRetries:  DefaultMaxRetries,
	}
}

func (c *CompressConfig) GetThreshold() int {
	if c.Threshold <= 0 {
		return DefaultCompressThreshold
	}
	return c.Threshold
}

func (c *CompressConfig) GetLockTimeout() int {
	if c.LockTimeout <= 0 {
		return DefaultLockTimeout
	}
	return c.LockTimeout
}

func (c *CompressConfig) GetMaxRetries() int {
	if c.MaxRetries <= 0 {
		return DefaultMaxRetries
	}
	return c.MaxRetries
}