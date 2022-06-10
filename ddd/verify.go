package ddd

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"strings"
)

type Verify interface {
	Validate() error
}

//
//  ValidateCreateCommand
//  @Description: 验证创建接口
//  @param data
//  @param verifyError 值可以为nil
//  @return *ddd_errors.VerifyError
//
func ValidateCreateCommand(data CreateCommand, verifyError *ddd_errors.VerifyError) *ddd_errors.VerifyError {
	return validateCommand(data, verifyError)
}

//
//  ValidateUpdateCommand
//  @Description: 验证更新接口
//  @param data
//  @param verifyError 值可以为nil
//  @return *ddd_errors.VerifyError
//
func ValidateUpdateCommand(data UpdateCommand, verifyError *ddd_errors.VerifyError) *ddd_errors.VerifyError {
	return validateCommand(data, verifyError)
}

//
//  ValidateDeleteCommand
//  @Description: 验证删除接口
//  @param data
//  @param verifyError 值可以为nil
//  @return *ddd_errors.VerifyError
//
func ValidateDeleteCommand(data DeleteCommand, verifyError *ddd_errors.VerifyError) *ddd_errors.VerifyError {
	return validateCommand(data, verifyError)
}

func validateCommand(data Command, verifyError *ddd_errors.VerifyError) *ddd_errors.VerifyError {
	v := verifyError
	if v == nil {
		v = ddd_errors.NewVerifyError()
	}
	if tenantId, ok := data.(GetTenant); ok {
		validateId("tenantId", tenantId.GetTenantId(), v)
	}
	if commandId, ok := data.(GetCommandId); ok {
		validateId("commandId", commandId.GetCommandId(), v)
	}
	if aggId, ok := data.(GetAggregateId); ok {
		validateId("aggregateId", aggId.GetAggregateId().RootId(), v)
	}
	return v
}

func validateId(fieldName, idValue string, verifyError *ddd_errors.VerifyError) {
	if len(idValue) == 0 {
		verifyError.AppendField(fieldName, "不能为空")
	} else {
		if strings.Index(idValue, " ") > -1 {
			verifyError.AppendField(fieldName, "不能包含“空格”")
		}
	}
}
