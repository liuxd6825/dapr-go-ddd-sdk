package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapper_Search(t *testing.T) {
	ctx := context.Background()
	mongodb, coll := newCollection("test_mapper_users")
	mapper := NewDao[*MapperUser](mongodb, coll)
	objId := NewObjectID()
	// id := string(objId)
	user := &MapperUser{
		Id:        objId,
		TenantId:  "001",
		UserName:  "UserName",
		UserCode:  "UserCode",
		Address:   "address",
		Email:     "lxd@163.com",
		Telephone: "17767788888",
	}

	_ = mapper.Insert(ctx, user).OnSuccess(func(data *MapperUser) error {
		println(data)
		return nil
	}).OnError(func(err error) error {
		assert.Error(t, err)
		return err
	})

	search := ddd_repository.NewFindPagingQuery()
	search.SetTenantId("001")
	// search.SetFilter(fmt.Sprintf("id=='%s'", id))

	_ = mapper.FindPaging(ctx, search).OnSuccess(func(data []*MapperUser) error {
		println(data)
		return nil
	}).OnNotFond(func() error {
		err := errors.NewNotFondError()
		assert.Error(t, err)
		return err
	}).OnError(func(err error) error {
		assert.Error(t, err)
		return err
	}).GetError()

}

/*func newRepository(mongodb *MongoDB, coll *mongo.Collection) *Dao[*User] {
	return NewDao[*User](func() *User { return &User{} }, mongodb, coll)
}

func newCollection(name string) (*MongoDB, *mongo.Collection) {
	config := &Config{
		Host:         "192.168.64.8:27019",
		DatabaseName: "query-example",
		UserName:     "query-example",
		Password:     "123456",
	}
	mongodb, err := NewMongoDB(config)
	if err != nil {
		panic(err)
	}
	_ = mongodb.CreateCollection(name)
	coll := mongodb.GetCollection(name)
	return mongodb, coll
}*/

type MapperUser struct {
	Id        ObjectId `json:"id" validate:"gt=0" bson:"_id"`
	TenantId  string   `json:"tenantId" validate:"gt=0" bson:"tenant_id"`
	UserCode  string   `json:"userCode" validate:"gt=0" bson:"user_code"`
	UserName  string   `json:"userName" validate:"gt=0" bson:"user_name"`
	Email     string   `json:"email" validate:"gt=0" bson:"email"`
	Telephone string   `json:"telephone" validate:"gt=0" bson:"telephone"`
	Address   string   `json:"address" validate:"gt=0" bson:"address"`
}

func (u *MapperUser) GetTenantId() string {
	return u.TenantId
}

func (u *MapperUser) GetId() string {
	return string(u.Id)
}
