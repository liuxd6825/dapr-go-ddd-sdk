package hugedao

import (
	"context"
	"testing"
	"time"

	"github.com/apache/hugegraph-toolchain/hugegraph-client-go"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbevent"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/schema"
)

func getClient(t *testing.T) *hugegraph.CommonClient {
	t.Helper()
	if testClient == nil {
		t.Skip("hugegraph not available")
	}
	return testClient
}

func cleanAll(t *testing.T, client *hugegraph.CommonClient) {
	t.Helper()
	resp, err := client.Gremlin.Post(
		client.Gremlin.Post.WithGremlin("g.V().drop().iterate();g.E().drop().iterate()"),
	)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err != nil {
		t.Logf("cleanAll: %v", err)
	}
}

func newTestGraph() *schema.Graph {
	return &schema.Graph{
		IsEnable: true,
		Name:     "name",
		RelType:  "${relation_type}",
		RelStart: "source",
		RelEnd:   "target",
	}
}

func newDBSchema(name string) *dbschema.DBSchema {
	return &dbschema.DBSchema{
		Name:      name,
		TableName: name,
	}
}

// buildTestDBSchema 构造一个含 graph 元信息的 DBSchema,
// 让 model.NewMasterNode / NewMasterRelation 可正常工作
func buildTestDBSchema() *dbschema.DBSchema {
	jsonStr := `{
		"type": "object",
		"title": "master_node",
		"meta": {
			"graph": {
				"isEnable": true,
				"name": "name",
				"relType": "${relation_type}",
				"relStart": "${source}",
				"relEnd": "${target}",
				"labels": ["master"],
				"graphType": ["node", "rel"]
			}
		},
		"properties": {
			"id": {"type": "string"},
			"name": {"type": "string"},
			"case_id": {"type": "string"},
			"tenant_id": {"type": "string"},
			"source": {"type": "string"},
			"target": {"type": "string"},
			"relation_type": {"type": "string"},
			"description": {"type": "string"},
			"table": {"type": "string"},
			"type": {"type": "string"},
			"source_ids": {"type": "string"},
			"source_type": {"type": "string"},
			"keywords": {"type": "string"}
		}
	}`
	return dbschema.NewDBSchemaWithJsonSchemaText("master_node", jsonStr)
}

// assertCount 执行 Gremlin count 查询并验证返回值
// HugeGraph 1.7 存在事务间最终一致性问题, 写后立即读可能看不到; 失败时重试 3 次
func assertCount(t *testing.T, client *hugegraph.CommonClient, script string, expect int) {
	t.Helper()
	var got int
	var lastResp any
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			time.Sleep(200 * time.Millisecond)
		}
		resp, err := client.Gremlin.Post(client.Gremlin.Post.WithGremlin(script))
		if err != nil {
			t.Fatalf("assertCount script=%q err=%v", script, err)
		}
		if resp == nil || resp.Data == nil {
			if expect == 0 {
				return
			}
			lastResp = nil
			continue
		}
		if resp.Data.Exception != "" || resp.Data.Message != "" {
			t.Fatalf("assertCount script=%q gremlin error: %s cause: %s",
				script, resp.Data.Message, resp.Data.Cause)
		}
		list, ok := resp.Data.Result.Data.([]any)
		if !ok {
			t.Fatalf("assertCount script=%q result.data type=%T value=%v", script, resp.Data.Result.Data, resp.Data.Result.Data)
		}
		if len(list) == 0 {
			if expect == 0 {
				return
			}
			lastResp = nil
			continue
		}
		v, ok := list[0].(float64)
		if !ok {
			t.Fatalf("assertCount script=%q value type=%T value=%v", script, list[0], list[0])
		}
		got = int(v)
		lastResp = v
		if got == expect {
			return
		}
	}
	t.Fatalf("assertCount script=%q got %v (expected %d) after retries, lastResp=%v", script, lastResp, expect, lastResp)
}

func newMasterNodeDao(t *testing.T, client *hugegraph.CommonClient) dao.IMasterNodeDao {
	t.Helper()
	gm := newTestGraph()
	if err := EnsureSchema(context.Background(), client, gm); err != nil {
		t.Fatalf("EnsureSchema: %v", err)
	}
	registerTestEnv(client)
	return NewMasterNodeDao([]string{"master"}, newDBSchema("master_node"), gm)
}

func newBusRelationDao(t *testing.T, client *hugegraph.CommonClient) *BusRelationDao {
	t.Helper()
	gm := newTestGraph()
	nodeDao := newMasterNodeDao(t, client)
	return NewBusRelationDao(newDBSchema("master_relation"), nodeDao, gm)
}

func TestMasterNodeDao_CreateMain(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	dao := newMasterNodeDao(t, client)
	node := &model.MasterNode{
		Id:          "n1",
		Name:        "alpha",
		CaseId:      "c1",
		TenantId:    "t1",
		SourceIds:   "n1",
		SourceType:  "master",
		Type:        "entity",
		Description: "test node",
		Table:       "master_node",
	}
	dao.CreateMain(context.Background(), node)
}

func TestMasterNodeDao_DeleteMain(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	dao := newMasterNodeDao(t, client)
	sch := buildTestDBSchema()

	node := &model.MasterNode{Id: "n2", Name: "beta", CaseId: "c1", TenantId: "t1"}
	dao.CreateMain(context.Background(), node)
	assertCount(t, client, "g.V().hasLabel('master').has('id','n2').count()", 1)

	rec := &dbevent.CDCRecord{
		DBSchema: sch,
		OpType:   dbevent.OpTypeDelete,
		Before: map[string]any{
			"id": "n2", "name": "beta", "case_id": "c1", "tenant_id": "t1",
		},
	}
	dao.DeleteMain(context.Background(), rec, nil)

	assertCount(t, client, "g.V().hasLabel('master').has('id','n2').count()", 0)
}

func TestMasterNodeDao_ClearAll(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	dao := newMasterNodeDao(t, client)
	dao.CreateMain(context.Background(), &model.MasterNode{Id: "n3", Name: "gamma", CaseId: "c1", TenantId: "t1"})
	dao.ClearAll(context.Background())
}

func TestMasterNodeDao_CreateRelNode(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	dao := newMasterNodeDao(t, client)
	sch := buildTestDBSchema()

	// 准备 src / dst 节点
	dao.CreateMain(context.Background(), &model.MasterNode{Id: "s1", Name: "src", CaseId: "c1", TenantId: "t1"})
	dao.CreateMain(context.Background(), &model.MasterNode{Id: "d1", Name: "dst", CaseId: "c1", TenantId: "t1"})

	rel := model.NewMasterRelation(map[string]any{
		"id":            "r1",
		"source":        "s1",
		"target":        "d1",
		"case_id":       "c1",
		"tenant_id":     "t1",
		"relation_type": "depends_on",
		"description":   "test rel",
		"table":         "master_relation",
	}, sch)

	relNode := &model.MasterNode{
		Id:         "rn1",
		Name:       "rn1",
		CaseId:     "c1",
		TenantId:   "t1",
		SourceType: "master",
		SourceIds:  "r1",
	}

	dao.CreateRelNode(context.Background(), rel, relNode)

	// 验证: relNode 已创建为 master 顶点
	assertCount(t, client, "g.V().hasLabel('master').has('id','rn1').count()", 1)
	// master_rel 边已创建
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','r1').count()", 1)
	// 边的 rel_type 属性
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','r1').has('rel_type','depends_on').count()", 1)
	// 边起点是 s1
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','r1').outV().has('id','s1').count()", 1)
	// 边终点是 rn1
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','r1').inV().has('id','rn1').count()", 1)
}

// TestMasterNodeDao_Create 测试主子表双向关系
// 场景:
//  1. 人员主数据新增"张三", 人员公司子表新增"大豆科技" (作为张三的子公司)
//  2. 公司主数据新增"大豆科技", 公司人员子表新增"张三" (作为大豆科技的公司人)
//
// 验证:
//   - 两个顶点都存在 (CreateMain 按 id 幂等, 重复调用不会创建重复顶点)
//   - 双向边都存在且 rel_type 正确 (owns_company 与 has_person)
//   - 边的源/目标方向正确
//   - 顶点的 type / node_type 属性正确 (区分 person/company 业务实体)
func TestMasterNodeDao_Create(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	dao := newMasterNodeDao(t, client)
	sch := buildTestDBSchema()
	ctx := context.Background()

	// === 步骤 1: 人员主数据新增"张三" ===
	zhangsan := &model.MasterNode{
		Id:       "p_zhangsan",
		Name:     "张三",
		Type:     "person",
		CaseId:   "c1",
		TenantId: "t1",
	}
	dao.CreateMain(ctx, zhangsan)

	// === 步骤 2: 人员公司子表新增"大豆科技"(作为张三的子公司) ===
	dadou := &model.MasterNode{
		Id:       "c_dadou",
		Name:     "大豆科技",
		Type:     "company",
		CaseId:   "c1",
		TenantId: "t1",
	}
	rel1 := model.NewMasterRelation(map[string]any{
		"id":            "e_zs_to_ds",
		"source":        "p_zhangsan",
		"target":        "c_dadou",
		"case_id":       "c1",
		"tenant_id":     "t1",
		"relation_type": "owns_company",
		"description":   "张三→大豆科技",
	}, sch)
	dao.CreateRelNode(ctx, rel1, dadou)

	// === 步骤 3: 公司主数据新增"大豆科技" (顶点已存在, 期望幂等) ===
	dao.CreateMain(ctx, dadou)

	// === 步骤 4: 公司人员子表新增"张三" (顶点已存在, 期望幂等) ===
	rel2 := model.NewMasterRelation(map[string]any{
		"id":            "e_ds_to_zs",
		"source":        "c_dadou",
		"target":        "p_zhangsan",
		"case_id":       "c1",
		"tenant_id":     "t1",
		"relation_type": "has_person",
		"description":   "大豆科技→张三",
	}, sch)
	dao.CreateRelNode(ctx, rel2, zhangsan)

	// === 验证 ===

	// 1. 两个顶点都存在 (按 id 幂等, 各只一个)
	assertCount(t, client, "g.V().hasLabel('master').has('id','p_zhangsan').count()", 1)
	assertCount(t, client, "g.V().hasLabel('master').has('id','c_dadou').count()", 1)

	// 2. 顶点 type 属性区分业务实体
	assertCount(t, client, "g.V().hasLabel('master').has('id','p_zhangsan').has('type','person').count()", 1)
	assertCount(t, client, "g.V().hasLabel('master').has('id','c_dadou').has('type','company').count()", 1)

	// 3. 双向边存在且 rel_type 正确
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','e_zs_to_ds').has('rel_type','owns_company').count()", 1)
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','e_ds_to_zs').has('rel_type','has_person').count()", 1)

	// 4. 边的方向正确 (源→目标)
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','e_zs_to_ds').outV().has('id','p_zhangsan').count()", 1)
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','e_zs_to_ds').inV().has('id','c_dadou').count()", 1)
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','e_ds_to_zs').outV().has('id','c_dadou').count()", 1)
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','e_ds_to_zs').inV().has('id','p_zhangsan').count()", 1)
}

func TestMasterNodeDao_UpdateRelNode(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	dao := newMasterNodeDao(t, client)
	sch := buildTestDBSchema()

	// 准备 src / dst / relNode
	dao.CreateMain(context.Background(), &model.MasterNode{Id: "s1", Name: "src", CaseId: "c1", TenantId: "t1"})
	dao.CreateMain(context.Background(), &model.MasterNode{Id: "d1", Name: "dst", CaseId: "c1", TenantId: "t1"})

	rel := model.NewMasterRelation(map[string]any{
		"id":            "rn1",
		"source":        "s1",
		"target":        "d1",
		"case_id":       "c1",
		"tenant_id":     "t1",
		"relation_type": "depends_on",
		"description":   "old",
		"table":         "master_relation",
	}, sch)
	relNode := &model.MasterNode{Id: "rn1", Name: "rn1", CaseId: "c1", TenantId: "t1"}
	dao.CreateRelNode(context.Background(), rel, relNode)

	// 构造 CDC 记录: rename + relType 变更 + 业务字段
	before := map[string]any{
		"id":            "rn1",
		"name":          "rn1",
		"case_id":       "c1",
		"tenant_id":     "t1",
		"source":        "s1",
		"target":        "d1",
		"relation_type": "depends_on",
		"description":   "old",
		"table":         "master_relation",
		"source_ids":    "rn1",
		"source_type":   "master",
	}
	after := map[string]any{
		"id":            "rn1",
		"name":          "rn1-renamed",
		"case_id":       "c1",
		"tenant_id":     "t1",
		"source":        "s1",
		"target":        "d1",
		"relation_type": "links_to",
		"description":   "new",
		"table":         "master_relation",
		"source_ids":    "rn1",
		"source_type":   "master",
	}
	record := &dbevent.CDCRecord{
		DBSchema: sch,
		OpType:   dbevent.OpTypeUpdate,
		Before:   before,
		After:    after,
	}
	if err := dao.UpdateRelNode(context.Background(), record); err != nil {
		t.Fatalf("UpdateRelNode: %v", err)
	}

	// 验证: rename 生效
	assertCount(t, client, "g.V().hasLabel('master').has('id','rn1').has('name','rn1-renamed').count()", 1)
	// relType 变更: 边的 rel_type 属性更新
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','rn1').has('rel_type','links_to').count()", 1)
	// 旧 rel_type 的边已不存在(同 id 下 rel_type 应唯一)
	assertCount(t, client, "g.E().hasLabel('master_rel').has('id','rn1').has('rel_type','depends_on').count()", 0)
}

func TestMasterNodeDao_DeleteRelNode(t *testing.T) {
	client := getClient(t)
	cleanAll(t, client)
	dao := newMasterNodeDao(t, client)
	sch := buildTestDBSchema()

	dao.CreateMain(context.Background(), &model.MasterNode{Id: "s1", Name: "src", CaseId: "c1", TenantId: "t1"})
	dao.CreateMain(context.Background(), &model.MasterNode{Id: "d1", Name: "dst", CaseId: "c1", TenantId: "t1"})

	rel := model.NewMasterRelation(map[string]any{
		"id":            "r1",
		"source":        "s1",
		"target":        "d1",
		"case_id":       "c1",
		"tenant_id":     "t1",
		"relation_type": "depends_on",
		"table":         "master_relation",
	}, sch)
	relNode := &model.MasterNode{Id: "rn1", Name: "rn1", CaseId: "c1", TenantId: "t1"}
	dao.CreateRelNode(context.Background(), rel, relNode)

	// 删除前存在
	assertCount(t, client, "g.V().hasLabel('master').has('id','rn1').count()", 1)

	rec := &dbevent.CDCRecord{
		DBSchema: sch,
		OpType:   dbevent.OpTypeDelete,
		Before: map[string]any{
			"id": "rn1", "name": "rn1", "case_id": "c1", "tenant_id": "t1",
		},
	}
	dao.DeleteRelNode(context.Background(), rec)

	// 删除后不存在
	assertCount(t, client, "g.V().hasLabel('master').has('id','rn1').count()", 0)
}
