package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/xtest"
)

func init() {
	xtest.InitEnv_MongoRemoteTest(xtest.NewMongoOptions().SetDBName("master"))
}
func newTask() *model.SuTask {
	task := &model.SuTask{}
	task.Id = "001"
	account := "6235822099004087593"
	accounts := []*model.SuTaskAccount{}
	accounts = append(accounts, &model.SuTaskAccount{
		Account: account,
	})
}
