package ddd_mongodb

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/stringutils"
	"go.mongodb.org/mongo-driver/mongo"
)

type Decoder[T any] interface {
	Single(ctx context.Context, result *mongo.SingleResult, data T) error
	List(ctx context.Context, cursor *mongo.Cursor, data *[]T) error
}

type mapDecoder[T any] struct {
}

func NewMapDecoder[T any]() Decoder[T] {
	return &mapDecoder[T]{}
}

func (d *mapDecoder[T]) Single(ctx context.Context, result *mongo.SingleResult, data T) error {
	err := result.Decode(data)
	if err != nil {
		return err
	}
	d.rename(data)
	return nil
}

func (d *mapDecoder[T]) List(ctx context.Context, cursor *mongo.Cursor, list *[]T) error {
	if err := cursor.All(ctx, list); err != nil {
		return err
	}
	/*
		for _, item := range *list {
			d.rename(item)
		}*/
	return nil
}

func (d *mapDecoder[T]) rename(val T) {
	data := d.as(val)
	for k, v := range data {
		name := stringutils.MongoFieldAsJsonName(k)
		data[name] = v
		if name != k {
			delete(data, k)
		}
	}
}

func (d *mapDecoder[T]) as(val any) map[string]any {
	if m, ok := val.(ddd.MapEntity); ok {
		return m
	} else if m, ok := val.(map[string]any); ok {
		return m
	}
	return nil
}

type structDecoder[T any] struct {
}

func NewStructDecoder[T any]() Decoder[T] {
	return &structDecoder[T]{}
}

func (d *structDecoder[T]) Single(ctx context.Context, result *mongo.SingleResult, data T) error {
	return result.Decode(&data)
}

func (d *structDecoder[T]) List(ctx context.Context, cursor *mongo.Cursor, list *[]T) error {
	return cursor.All(ctx, list)
}
