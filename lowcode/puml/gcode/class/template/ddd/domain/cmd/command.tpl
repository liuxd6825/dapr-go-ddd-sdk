 package command

 import (
 	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
 	"github.com/liuxu6825/{{.ServiceName}}-service/pkg/cmd-service/domain/{{.AggreateName}}/field"
 )

 //
 // {{.Name}
 // @Description: 添加节点命令
 //
 type {{.Name}} struct {
 	CommandId   string                        `json:"commandId" validate:"required"` // 命令ID
 	IsValidOnly bool                          `json:"isValidOnly" validate:"-"`      // 是否仅验证，不执行
 	Data        field.{{.Name}}Fields `json:"data" validate:"-"`             // 业务数据
 }

 //
 // GetAggregateId
 // @Description: 获取聚合根Id
 //
 func (c *{{.Name}}) GetAggregateId() ddd.AggregateId {
 	return ddd.NewAggregateId(c.Data.Id)
 }

 //
 // GetCommandId
 // @Description: 获取命令Id
 //
 func (c *{{.Name}}) GetCommandId() string {
 	return c.CommandId
 }

 //
 // GetTenantId
 // @Description: 获取租户Id
 //
 func (c *{{.Name}}) GetTenantId() string {
 	return c.Data.TenantId
 }

 //
 // GetIsValidOnly
 // @Description: 是否只验证不执行。
 //
 func (c *{{.Name}}) GetIsValidOnly() bool {
 	return c.IsValidOnly
 }

 //
 // Validate
 // @Description: 命令数据验证
 //
 func (c *{{.Name}}) Validate() error {
 	ve := ddd.ValidateCommand(c, nil)
 	/* 添加其他数据检查
 	   if c.Items.Identity == "" {
 	       ve.AppendField("data.id", "不能为空")
 	   }
 	*/
 	return ve.GetError()
 }
