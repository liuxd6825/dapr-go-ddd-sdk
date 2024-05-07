package db

import (
	"context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/schema"
)

type Record map[string]any

type Model struct {
	schema *schema.Schema
	db     *Db
}

func NewModel(schema *schema.Schema, db *Db) *Model {
	return &Model{schema: schema, db: db}
}

func (m *Model) Update(ctx context.Context, field string, val any) {
	println("Update")
}

func (m *Model) Updates(ctx context.Context, data any) {
	println("Updates")
}

func (m *Model) First(ctx context.Context, id ...string) Record {
	println("First")
	return nil
}

func (m *Model) Find(ctx context.Context, rSQL string) []Record {
	println("Find")
	return nil
}

func (m *Model) Create(ctx context.Context, record Record) error {
	println("Create")
	return nil
}

func (m *Model) Delete(ctx context.Context, val ...string) error {
	println("Delete")
	return nil
}
