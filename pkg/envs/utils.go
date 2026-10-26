package envs

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/yamlutils"
)

func copyObject(src, dst any) error {
	return yamlutils.Copy(src, dst)
}
