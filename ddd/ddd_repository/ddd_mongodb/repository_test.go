package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func Test_Search(t *testing.T) {
	ctx := context.Background()
	mongodb, coll := newCollection("users")

	entityBuilder := ddd_repository.NewEntityBuilder(func() interface{} {
		return &User{}
	}, func() interface{} {
		return make([]User, 0)
	})
	repository := NewRepository(entityBuilder, mongodb, coll)
	user := &User{
		Id:        "001",
		TenantId:  "001",
		UserName:  "UserName",
		UserCode:  "UserCode",
		Address:   "address",
		Email:     "lxd@163.com",
		Telephone: "17767788888",
	}

	err := repository.DoCreate(ctx, user).OnSuccess(func(data interface{}) error {
		println(data)
		return nil
	}).GetError()

	assert.Error(t, err)

	search := &ddd_repository.SearchQuery{
		TenantId: "001",
		Filter:   "id=='003'",
	}
	err = repository.DoSearch(ctx, search).OnSuccess(func(data interface{}) error {
		println(data)
		return nil
	}).OnNotFond(func() error {
		return ddd_errors.NewNotFondError()
	}).OnError(func(err error) error {
		return err
	}).GetError()

	assert.Error(t, err)
}

func newCollection(name string) (*MongoDB, *mongo.Collection) {
	mongodb := NewMongoDB()
	err := mongodb.Init(&Config{
		Host:         "192.168.64.4",
		DatabaseName: "example-query-service",
		UserName:     "dapr",
		Password:     "123456",
	})
	if err != nil {
		panic(err)
	}
	coll := mongodb.NewCollection(name)
	return mongodb, coll
}

type User struct {
	Id        string `json:"id" validate:"gt=0" bson:"id"`
	TenantId  string `json:"tenantId" validate:"gt=0" bson:"tenantId"`
	UserCode  string `json:"userCode" validate:"gt=0" bson:"userCode"`
	UserName  string `json:"userName" validate:"gt=0" bson:"userName"`
	Email     string `json:"email" validate:"gt=0" bson:"email"`
	Telephone string `json:"telephone" validate:"gt=0" bson:"telephone"`
	Address   string `json:"address" validate:"gt=0" bson:"address"`
}

func (u *User) GetTenantId() string {
	return u.TenantId
}

func (u *User) GetId() string {
	return u.Id
}
