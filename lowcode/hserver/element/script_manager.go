package element

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/types"
	"github.com/sirupsen/logrus"
)

type ScriptManager interface {
	AddScript(config *ScriptConfig, logger logrus.FieldLogger, pkg *types.CMap[any]) error
	RunScript(funcName string, runValues *RunValues, checkHave bool, opts ...RunOptions) (res any, err error)
}
