package tokenutils

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/appctx"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/timeutils"
)
import "time"

type EditEntity interface {
	SetCreatedTime(value *time.Time)
	SetCreatorId(value string)
	SetCreatorName(value string)
	SetDeletedTime(value *time.Time)
	SetDeleterId(value string)
	SetDeleterName(value string)
	SetUpdatedTime(value *time.Time)
	SetUpdaterId(value string)
	SetUpdaterName(value string)
}

type TokenData struct {
	UserId   string
	UserName string
	TenantId string
}

func GetUser(ctx context.Context, getFunc func(userId, userName string), errFunc func(err error)) (appctx.AuthUser, error) {
	user, isFound := appctx.GetAuthUser(ctx)
	if !isFound {
		return nil, errors.ErrNotFoundLoginUser
	}
	if getFunc != nil {
		getFunc(user.GetId(), user.GetName())
	}
	if errFunc != nil {
		//errFunc()
	}
	return user, nil
}

func SetCreateUser(ctx context.Context, entity EditEntity) error {
	user, isFound := appctx.GetAuthUser(ctx)
	if !isFound {
		return errors.ErrNotFoundLoginUser
	}
	userId := user.GetId()
	userName := user.GetName()

	t := timeutils.Now()
	entity.SetCreatorId(userId)
	entity.SetCreatorName(userName)
	entity.SetCreatedTime(&t)

	entity.SetUpdaterId(userId)
	entity.SetUpdaterName(userName)
	entity.SetUpdatedTime(&t)
	return nil
}

func SetUpdateUser(ctx context.Context, entity EditEntity) error {
	user, isFound := appctx.GetAuthUser(ctx)
	if !isFound {
		return errors.ErrNotFoundLoginUser
	}

	userId := user.GetId()
	userName := user.GetName()

	t := timeutils.Now()
	entity.SetUpdaterId(userId)
	entity.SetUpdaterName(userName)
	entity.SetUpdatedTime(&t)
	return nil
}

func SetDeleteUser(ctx context.Context, entity EditEntity) error {
	user, isFound := appctx.GetAuthUser(ctx)
	if !isFound {
		return errors.ErrNotFoundLoginUser
	}

	userId := user.GetId()
	userName := user.GetName()

	t := timeutils.Now()
	entity.SetDeleterId(userId)
	entity.SetDeleterName(userName)
	entity.SetDeletedTime(&t)
	return nil
}
