package hugedao

import (
	"context"
	"fmt"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/store/graph"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/hugegraph"
)

type MasterGraphDao struct {
	store *hugegraph.Store
}

func NewMasterGraphDao() *MasterGraphDao {
	store, err := GetStore()
	if err != nil {
		panic(err)
	}
	return &MasterGraphDao{
		store: store,
	}
}

func (d *MasterGraphDao) FindByCaseId(ctx context.Context, caseId string) *graph.GraphView {
	gremlin := fmt.Sprintf("g.V().hasLabel('human').has('case_id','%s').outE('human_interact').inV().path()", caseId)
	res, err := d.store.Query(ctx, gremlin)
	if err != nil {
		panic(err)
	}
	return res.GetGraphView()
}

func (d *MasterGraphDao) FindById(ctx context.Context, caseId string, id string) *graph.GraphView {
	cypher := fmt.Sprintf("g.V('%s').has('case_id','%s')", id, caseId)
	res, err := d.store.Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.GetGraphView()
}

func (d *MasterGraphDao) FindByName(ctx context.Context, caseId string, name string) *graph.GraphView {
	cypher := fmt.Sprintf("g.V().has('case_id','%s').has('name','%s')", caseId, name)
	res, err := d.store.Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.GetGraphView()
}

func (d *MasterGraphDao) FindByContainName(ctx context.Context, caseId string, name string) *graph.GraphView {
	cypher := fmt.Sprintf("g.V().has('case_id','%s').has('name', Text.contains('%s'))", caseId, name)
	res, err := d.store.Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.GetGraphView()
}

func (d *MasterGraphDao) FindByStartWithName(ctx context.Context, caseId string, name string) *graph.GraphView {
	// g.V().hasLabel('human').has('name').filter{ it.get().value('name').startsWith('任') }
	cypher := fmt.Sprintf("g.V().has('case_id','%s').has('name', Text.contains('%s')).filter{ it.get().value('name').startsWith('%s')}", caseId, name, name)
	res, err := d.store.Query(ctx, cypher, nil)
	if err != nil {
		panic(err)
	}
	return res.GetGraphView()
}
