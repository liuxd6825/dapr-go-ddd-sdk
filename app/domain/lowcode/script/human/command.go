package human

import "github.com/liuxd6825/dapr-go-ddd-sdk/app/pkg/xcommon/xbase"

type CreateCommand struct {
	xbase.Command[any]
}
