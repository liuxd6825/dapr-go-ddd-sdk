package ddd_mongodb

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func Test_Search(t *testing.T) {
	ctx := context.Background()
	mongodb, coll := newCollection("test_users")
	repository := newRepository(mongodb, coll)
	objId := NewObjectID()
	id := string(objId)
	user := &User{
		Id:        objId,
		TenantId:  "001",
		UserName:  "UserName",
		UserCode:  "UserCode",
		Address:   "address",
		Email:     "lxd@163.com",
		Telephone: "17767788888",
	}

	_ = repository.Insert(ctx, user).OnSuccess(func(data *User) error {
		println(data)
		return nil
	}).OnError(func(err error) error {
		assert.Error(t, err)
		return err
	})

	search := ddd_repository.NewFindPagingQuery()
	search.SetTenantId("001")
	search.SetFilter(fmt.Sprintf("id=='%s'", id))

	_ = repository.FindPaging(ctx, search).OnSuccess(func(data *[]*User) error {
		println(data)
		return nil
	}).OnNotFond(func() error {
		err := ddd_errors.NewNotFondError()
		assert.Error(t, err)
		return err
	}).OnError(func(err error) error {
		assert.Error(t, err)
		return err
	}).GetError()

}

func TestMongoSession_UseTransaction(t *testing.T) {
	mongodb, coll := newCollection("test_users")
	repository := newRepository(mongodb, coll)
	err := ddd_repository.StartSession(context.Background(), NewSession(mongodb), func(ctx context.Context) error {
		for i := 0; i < 5; i++ {
			user := &User{
				Id:        NewObjectID(),
				TenantId:  "001",
				UserName:  "userName" + fmt.Sprint(i),
				UserCode:  "UserCode",
				Address:   "address",
				Email:     "lxd@163.com",
				Telephone: "17767788888",
			}

			err := repository.Insert(ctx, user).OnSuccess(func(data *User) error {
				println(data)
				return nil
			}).GetError()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func newRepository(mongodb *MongoDB, coll *mongo.Collection) *Repository[*User] {
	return NewRepository[*User](func() *User { return &User{} }, mongodb, coll)
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
}

type User struct {
	Id        ObjectId `json:"id" validate:"gt=0" bson:"_id"`
	TenantId  string   `json:"tenantId" validate:"gt=0" bson:"tenantId"`
	UserCode  string   `json:"userCode" validate:"gt=0" bson:"userCode"`
	UserName  string   `json:"userName" validate:"gt=0" bson:"userName"`
	Email     string   `json:"email" validate:"gt=0" bson:"email"`
	Telephone string   `json:"telephone" validate:"gt=0" bson:"telephone"`
	Address   string   `json:"address" validate:"gt=0" bson:"address"`
}

func (u *User) GetTenantId() string {
	return u.TenantId
}

func (u *User) GetId() string {
	return string(u.Id)
}
