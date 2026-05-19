package dao

import (
	"testing"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/randomutils"
	xtest2 "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/xtest"
)

func Test_ChatDao_Insert(t *testing.T) {
	ctx := xtest2.NewContext()
	e := xtest2.NewEnvConfigMongo(&env.Mongo{
		DbKey:  "db",
		DbName: "test",
	})
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
