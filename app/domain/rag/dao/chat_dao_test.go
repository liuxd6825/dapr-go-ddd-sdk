package dao

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/randomutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
	"testing"
)

func Test_ChatDao_Insert(t *testing.T) {
	ctx := xtest.NewContext()
	e := xtest.NewEnvConfig_Mongo("db", "test")
	env.SetEnv(e)
	dao := NewChatDao("db")
	chat := &model.Chat{
		Title: randomutils.NameCN(),
		Base: model.Base{
			Id:       randomutils.NewId(),
			TenantId: "test",
			CaseId:   "1001",
		},
	}
	res := dao.Create(ctx, chat)
	assert.Equal(res.RowsAffected, int64(1))

	chat.Remark = "remark"
	res = dao.Update(ctx, chat)
	assert.Equal(res.RowsAffected, int64(1))
}
