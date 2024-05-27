package model

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/domain/{{.Class.FN}}/command"
	"github.com/liuxu6825/duxm-graph-service/pkg/cmd-service/domain/{{.Class.FN}}/factory"
)

{{- range $method := .Class.Methods }}
// PaymentNodeCreateCommand
// @Description:
// @receiver a
// @param ctx
// @param cmd
// @param metadata
// @return any
// @return error
func (a *{{.Class.Name}}Aggregate) {{$method.Name}}(ctx context.Context, cmd *command.PaymentNodeCreateCommand, metadata ddd.Metadata)  (any, error) {
	e, err := factory.Event.NewPaymentNodeCreateEvent(ctx, cmd, metadata)
	if err != nil {
		return nil, err
	}
	return ddd.ApplyEvent(ctx, a, e, ddd.NewApplyEventOptions(metadata))
}
{{- end}}
