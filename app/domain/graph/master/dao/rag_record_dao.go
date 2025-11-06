package dao

import (
	"context"
	"fmt"
	"sync"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/graph/master/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/stringutils"
)

type RagRecordDao struct {
	*Base[*model.Record]
	labels []string
}

var _recordDao *RagRecordDao
var _recordOnce sync.Once

func NewRagRecordDao(dbKey string) *RagRecordDao {
	_recordOnce.Do(func() {
		_recordDao = newRecordDao(dbKey, "record")
	})
	return _recordDao
}

func newRecordDao(dbKey string, labels ...string) *RagRecordDao {
	dbSch := dbschema.NewDBSchemaWithStruct("record", &model.Record{}, "record")
	nodeCfg := &dao.DaoConfig{
		DBKey:              dbKey,
		GraphType:          idao.GraphType_Node,
		GraphLabels:        labels,
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
	}
	newDao := dao.NewDao[*model.Record](nodeCfg)
	return &RagRecordDao{
		labels: labels,
		Base:   &Base[*model.Record]{DBSchema: dbSch, Dao: newDao},
	}
}

func (d *RagRecordDao) Create(ctx context.Context, record *model.Record) error {
	storeDao := d.GetStore()
	props, dataMap, err := storeDao.Cypher.GetCreateProperties(ctx, record)
	if err != nil {
		panic(err)
	}
	c := storeDao.Cypher

	labels := c.GetLabels(ctx, record)

	fb := stringutils.NewFmtBuilder()
	fb.String("labels", labels)
	fb.String("props", props)

	fb.Varchar("name", record.Name)
	fb.Varchar("acct", record.Acct)
	fb.Varchar("oppName", record.OppName)
	fb.Varchar("oppAcct", record.OppAcct)
	fb.Varchar("id", record.Id)
	fb.Varchar("source", record.Acct)
	fb.Varchar("target", record.OppAcct)
	fb.Varchar("tenantId", record.TenantId)
	fb.Varchar("caseId", record.CaseId)
	fb.Varchar("keywords", "")
	fb.Varchar("description", "")
	fb.Varchar("sourceType", record.FileId)
	fb.Varchar("table", "record")

	fmtStr := `
	MERGE (n1$<labels>:human{id:$<name>})
	MERGE (a1$<labels>:account{id:$<acct>})   
	MERGE (n1)-[:owner]->(a1)  

	MERGE (n2$<labels>:human{id:$<oppName>})
	MERGE (a2$<labels>:account{id:$<oppAcct>})   
	MERGE (n2)-[:owner]->(a2)  
`
	if err = d.write(ctx, fmtStr, fb, dataMap); err != nil {
		return err
	}

	fmtStr = `
	MATCH 
		(n$<labels>:account{id:$<acct>}),
		(m$<labels>:account{id:$<oppAcct>})   
	CREATE 
		(n)-[r:record{$<props>}]->(m)
	SET 
		r.id=$<id>,r.keywords=$<keywords>,r.case_id=$<caseId>,r.tenant_id=$<tenantId>, 
		r.description=$<description>,r.source=$<source>,r.target=$<target>, 
		r.source_ids=$<id>,r.source_type=$<sourceType>,r.table=$<table>
`
	err = d.write(ctx, fmtStr, fb, dataMap)
	return err
}

func (d *RagRecordDao) GetDescription(ctx context.Context, node *model.Record) string {
	date := node.Date.Format("2006-01-02 15:04:05")
	return fmt.Sprintf("%s于%s通过账号%s向%s的%s账号汇入%s%d", date, node.Name, node.Acct, node.OppName, node.OppAcct, node.Ccy, node.Amount)
}

func (d *RagRecordDao) GetLabels(ctx context.Context, record *model.Record) string {
	return fmt.Sprintf(":tenant_%s:case_%s:file_%s", record.TenantId, record.CaseId, record.FileId)
}

func (d *RagRecordDao) Update(ctx context.Context, record *model.Record) error {
	storeDao := d.GetStore()
	props, dataMap, err := storeDao.Cypher.GetUpdateProperties(ctx, record, "r")
	if err != nil {
		return err
	}
	c := storeDao.Cypher

	labels := c.GetLabels(ctx, record)

	fb := stringutils.NewFmtBuilder()
	fb.String("labels", labels)
	fb.String("props", props)

	fb.Varchar("name", record.Name)
	fb.Varchar("acct", record.Acct)
	fb.Varchar("oppName", record.OppName)
	fb.Varchar("oppAcct", record.OppAcct)
	fb.Varchar("id", record.Id)
	fb.Varchar("source", record.Acct)
	fb.Varchar("target", record.OppAcct)
	fb.Varchar("tenantId", record.TenantId)
	fb.Varchar("caseId", record.CaseId)
	fb.Varchar("keywords", "")
	fb.Varchar("description", "")
	fb.Varchar("sourceType", record.FileId)
	fb.Varchar("table", "record")

	fmtStr := `
	MERGE (n1$<labels>:human{id:$<name>})
	MERGE (a1$<labels>:account{id:$<acct>})   
	MERGE (n1)-[:owner]->(a1)  

	MERGE (n2$<labels>:human{id:$<oppName>})
	MERGE (a2$<labels>:account{id:$<oppAcct>})   
	MERGE (n2)-[:owner]->(a2)  
`
	if err = d.write(ctx, fmtStr, fb, dataMap); err != nil {
		return err
	}

	fmtStr = `
	MATCH 
		(n$<labels>:account{id:$<acct>}),
		(m$<labels>:account{id:$<oppAcct>})   
	MERGE 
		(n)-[r:record{id:$<id>}]->(m)
	SET 
		r.keywords=$<keywords>,r.case_id=$<caseId>,
		r.description=$<description>,r.source=$<source>,r.target=$<target>, 
		r.source_ids=$<id>,r.source_type=$<sourceType>,r.table=$<table>,
		$<props>
	`
	err = d.write(ctx, fmtStr, fb, dataMap)
	return err
}

func (d *RagRecordDao) Delete(ctx context.Context, record *model.Record) (err error) {
	storeDao := d.GetStore()
	c := storeDao.Cypher
	labels := c.GetLabels(ctx, record)

	fb := stringutils.NewFmtBuilder()
	fb.String("labels", labels)
	fb.Varchar("id", record.Id)
	fb.Varchar("acct", record.Acct)
	fb.Varchar("oppAcct", record.OppAcct)
	fmtStr := `
	MATCH (a$<labels>:account{id:$<acct>})-[r:record{id:$<id>}]->(a1$<labels>:account{id:$<oppAcct>})	
	DELETE r
`
	if err = d.write(ctx, fmtStr, fb, nil); err != nil {
		return err
	}
	return nil
}
