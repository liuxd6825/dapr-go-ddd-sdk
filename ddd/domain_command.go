package ddd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"strings"
)

type DomainCommand interface {
	NewDomainEvent() DomainEvent
	GetAggregateId() string
	GetTenantId() string
	GetCommandId() string
	Verify
}

type BaseDomainCommand struct {
	CommandId   string `json:"commandId"  validate:"gt=0"`
	IsValidOnly bool   `json:"isValidOnly"`
}

func (d *BaseDomainCommand) GetCommandId() string {
	return d.CommandId
}

func (d *BaseDomainCommand) Validate() error {
	verifyError := ddd_errors.NewVerifyError()
	d.ValidateError(verifyError)
	return verifyError
}

func (d *BaseDomainCommand) ValidateError(verifyError *ddd_errors.VerifyError) {
	if strings.Index(d.CommandId, " ") > -1 {
		verifyError.AppendField("commandId", "不能包含“空格”")
	}
	if len(d.CommandId) == 0 {
		verifyError.AppendField("commandId", "不能为空")
	}
}

func (d *BaseDomainCommand) GetIsValidOnly() bool {
	return d.IsValidOnly
}

type CreateDomainCommand interface {
	DomainCommand
	GetIsValidOnly() bool
}

type UpdateDomainCommand interface {
	DomainCommand
	GetIsValidOnly() bool
	GetUpdateMask() string
}
