package store_huge

import (
	"context"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/apache/hugegraph-toolchain/hugegraph-client-go/api/v1/model"
)

type VertexDao[T IVertex] struct {
	*BaseDao[T]
}

func NewVertexDao[T IVertex](client *hugegraph.CommonClient, config *Config[T]) *VertexDao[T] {
	return &VertexDao[T]{
		BaseDao: NewBaseDao[T](client, config),
	}
}

func (d *VertexDao[T]) GetModeProps(ctx context.Context, data T) error {
	v, err1 := d.getModel(ctx, data)
	if err1 != nil {
		return err1
	}
	_, err2 := d.client.Vertex.Create(d.client.Vertex.Create.WithVertex(v))
	if err2 != nil {
		return err2
	}
	return nil
}

func (d *VertexDao[T]) toMap(m T) (map[string]any, error) {
	return nil, nil
}

func (d *VertexDao[T]) Create(ctx context.Context, data T) error {
	v, err1 := d.getModel(ctx, data)
	if err1 != nil {
		return err1
	}
	_, err2 := d.client.Vertex.Create(d.client.Vertex.Create.WithVertex(v))
	if err2 != nil {
		return err2
	}
	return nil
}

func (d *VertexDao[T]) BatchCreate(ctx context.Context, data []T) error {
	list := []model.Vertex[any]{}
	for _, item := range data {
		v, err1 := d.getModel(ctx, item)
		if err1 != nil {
			return err1
		}
		list = append(list, v)
	}

	_, err := d.client.Vertex.BatchCreate(d.client.Vertex.BatchCreate.WithVertices(list))
	if err != nil {
		return err
	}
	return nil
}

func (d *VertexDao[T]) DeleteByID(ctx context.Context, data T) error {
	v, err1 := d.getModel(ctx, data)
	if err1 != nil {
		return err1
	}

	_, err := d.client.Vertex.DeleteByID(d.client.Vertex.Create.WithVertex(v))
	if err != nil {
		return err
	}
	return nil
}

func (d *VertexDao[T]) GetByID(ctx context.Context, data T) error {
	v, err1 := d.getModel(ctx, data)
	if err1 != nil {
		return err1
	}

	_, err := d.client.Vertex.GetByID(d.client.Vertex.Create.WithVertex(v))
	if err != nil {
		return err
	}
	return nil
}

func (d *VertexDao[T]) DeleteProperties(ctx context.Context, data T) error {
	v, err1 := d.getModel(ctx, data)
	if err1 != nil {
		return err1
	}

	_, err := d.client.Vertex.DeleteProperties(d.client.Vertex.Create.WithVertex(v))
	if err != nil {
		return err
	}
	return nil
}

func (d *VertexDao[T]) getModel(ctx context.Context, data T) (model.Vertex[any], error) {
	v := model.Vertex[any]{}
	v.SetID(data.GetId())
	v.SetTyp("")
	v.SetLabel(d.GetLabel())
	v.SetProperties(data)
	return v, nil
}
