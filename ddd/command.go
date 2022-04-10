package ddd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"strings"
)

type AggregateId interface {
	RootId() string
	ItemIds() *[]string
	ItemCount() int
	ItemEmpty() bool
}

type aggregateId struct {
	rootId  string
	itemIds *[]string
}

func NewAggregateId(rootId string, itemIds ...string) AggregateId {
	return &aggregateId{
		rootId:  rootId,
		itemIds: &itemIds,
	}
}

func NewAggregateIds(rootId string, itemIds *[]string) AggregateId {
	return &aggregateId{
		rootId:  rootId,
		itemIds: itemIds,
	}
}

func (a *aggregateId) RootId() string {
	return a.rootId
}

func (a *aggregateId) ItemIds() *[]string {
	return a.itemIds
}

func (a *aggregateId) ItemCount() int {
	if a.itemIds == nil {
		return 0
	}
	return len(*a.itemIds)
}

func (a *aggregateId) ItemEmpty() bool {
	if a.ItemCount() == 0 {
		return true
	}
	return false
}

type BaseCommand struct {
	CommandId   string `json:"commandId"  validate:"gt=0"`
	IsValidOnly bool   `json:"isValidOnly"`
}

func (d *BaseCommand) GetCommandId() string {
	return d.CommandId
}

func (d *BaseCommand) Validate() error {
	verifyError := ddd_errors.NewVerifyError()
	d.ValidateError(verifyError)
	return verifyError
}

func (d *BaseCommand) ValidateError(verifyError *ddd_errors.VerifyError) {
	if strings.Index(d.CommandId, " ") > -1 {
		verifyError.AppendField("commandId", "不能包含“空格”")
	}
	if len(d.CommandId) == 0 {
		verifyError.AppendField("commandId", "不能为空")
	}
}

func (d *BaseCommand) GetIsValidOnly() bool {
	return d.IsValidOnly
}

type Command interface {
	NewDomainEvent() DomainEvent
	GetAggregateId() AggregateId
	GetTenantId() string
	GetCommandId() string
	Verify
}

type CreateCommand interface {
	Command
	GetIsValidOnly() bool
}

type UpdateCommand interface {
	Command
	GetIsValidOnly() bool
	GetUpdateMask() string
}
