package element

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/lowcode/hserver/common"
)

type RunValues struct {
	Server     Server
	Service    Service
	WebContext WebContext
	Request    Request
	Alias      common.Alias
	WorkPath   string
	Self       any
}
