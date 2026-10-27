package envs

import (
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/yamlutils"
)

func copyObject(src, dst any) error {
	return yamlutils.Copy(src, dst)
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
