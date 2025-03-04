package ddd_neo4j

import (
	"context"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/assert"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd"
	"github.com/liuxd6825/dapr-go-ddd-sdk/ddd/ddd_repository"
	"github.com/liuxd6825/dapr-go-ddd-sdk/errors"
	"github.com/liuxd6825/dapr-go-ddd-sdk/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/gp"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/jsonutils"
	"github.com/liuxd6825/dapr-go-ddd-sdk/utils/reflectutils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
	"strings"
	"time"
)

/*type Dao[T any] interface {
	Insert(ctx context.Context, entity T, opts ...ddd_repository.Options) (setResult *ddd_repository.SetResult[T])
	InsertMany(ctx context.Context, entities []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T]

	Update(ctx context.Context, entity T, opts ...ddd_repository.Options) *ddd_repository.SetResult[T]
	UpdateMany(ctx context.Context, list []T, opts ...ddd_repository.Options) *ddd_repository.SetManyResult[T]

	DeleteById(ctx context.Context, tenantId string, id string, opts ...ddd_repository.Options) error
	DeleteByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) error
	DeleteAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) error
	DeleteByFilter(ctx context.Context, tenantId string, filter string, opts ...ddd_repository.Options) error
	DeleteByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) error

	FindListByMap(ctx context.Context, tenantId string, filterMap map[string]interface{}, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
	FindByFilter(ctx context.Context, tenantId, filter string) *ddd_repository.FindListResult[T]
	FindById(ctx context.Context, tenantId, id string, opts ...ddd_repository.Options) (T, bool, error)
	FindByIds(ctx context.Context, tenantId string, ids []string, opts ...ddd_repository.Options) ([]T, bool, error)
	FindAll(ctx context.Context, tenantId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
	FindByGraphId(ctx context.Context, tenantId string, graphId string, opts ...ddd_repository.Options) *ddd_repository.FindListResult[T]
}*/

type Dao[T any] struct {
	driver    neo4j.DriverWithContext
	cypher    Cypher
	eb        ddd.EntityBuilder[T]
	tableName string
}

func (d *Dao[T]) NewEntity() T {
	return d.eb.NewEntity()
}

func (d *Dao[T]) NewEntityList() []T {
	return d.eb.NewEntityList()
}

func (d *Dao[T]) GetTenantId(entity T) string {
	return d.eb.GetTenantId(entity)
}

func (d *Dao[T]) SetTenantId(entity T, tenantId string) {
	d.eb.SetTenantId(entity, tenantId)
}

func (d *Dao[T]) GetId(entity T) string {
	return d.eb.GetId(entity)
}

func (d *Dao[T]) SetId(entity T, id string) {
	d.eb.SetId(entity, id)
}

type Options[T interface{}] struct {
	newOne  func() T
	newList func() []T
}

type ImportCsvCmd struct {
	TenantId         string                 `json:"tenantId" desc:"租户ID"`
	CaseId           string                 `json:"caseId" desc:""`
	ImportFile       string                 `json:"importFile"`
	Labels           []string               `json:"label"`
	Fields           []string               `json:"fields"`
	ImportType       ImportType             `json:"importType"`
	Data             ImportCsvCmdData       `json:"data"`
	SaveFileCallback ImportSaveFileCallback `json:"-"`
}

type ImportCsvCmdData interface {
	Data() any
	List() any
	Item(index int) any
	Append(item any)
	Length() int
}

type ImportType int

type ImportJsonCmd struct {
	TenantId   string     `json:"tenantId" desc:"租户ID"`
	CaseId     string     `json:"caseId" desc:""`
	Neo4jPath  string     `json:"neo4JPath"`
	ImportFile string     `json:"importFile"`
	Nodes      []Node     `json:"nodes"`
	Relations  []Relation `json:"relations"`
}

type importCsvCmdData struct {
	list []any
	data any
}

func (i *importCsvCmdData) List() any {
	return i.list
}

func (i *importCsvCmdData) Data() any {
	return i.data
}

func (i *importCsvCmdData) Item(index int) any {
	return i.list[index]
}

func (i *importCsvCmdData) Append(item any) {
	i.list = append(i.list, item)
}

func (i *importCsvCmdData) Length() int {
	return len(i.list)
}

func NewImportCsvCmdData(data any) ImportCsvCmdData {
	return &importCsvCmdData{data: data}
}

const (
	ImportTypeNode = iota
	ImportTypeRelation
)

func NewOptions[T interface{}](opts ...*Options[T]) *Options[T] {
	n := &Options[T]{}
	for _, o := range opts {
		if o.newList != nil {
			n.newList = o.newList
		}
		if o.newOne != nil {
			n.newOne = o.newOne
		}
	}
	return n
}

func NewNodeDao[T any](driver neo4j.DriverWithContext, tableName string, eb ddd.EntityBuilder[T], opts ...*Options[T]) ddd_repository.Dao[T] {
	cypher := NewNodeCypher()
	return NewDao(driver, tableName, cypher, eb, opts...)
}

func NewRelationDao[T any](driver neo4j.DriverWithContext, tableName string, eb ddd.EntityBuilder[T], opts ...*Options[T]) ddd_repository.Dao[T] {
	cypher := NewRelationCypher("")
	return NewDao(driver, tableName, cypher, eb, opts...)
}

func NewDao[T any](driver neo4j.DriverWithContext, tableName string, cypher Cypher, eb ddd.EntityBuilder[T], opts ...*Options[T]) ddd_repository.Dao[T] {
	return newDao(driver, tableName, cypher, eb, opts...)
}

func newDao[T any](driver neo4j.DriverWithContext, tableName string, cypher Cypher, eb ddd.EntityBuilder[T], opts ...*Options[T]) *Dao[T] {
	dao := &Dao[T]{
		driver:    driver,
		cypher:    cypher,
		eb:        eb,
		tableName: tableName,
	}
	return dao
}

func (d *Dao[T]) init(driver neo4j.DriverWithContext, cypher Cypher, eb ddd.EntityBuilder[T], opts ...*Options[T]) {
	d.driver = driver
	d.cypher = cypher
	d.eb = eb
}

func (d *Dao[T]) query(ctx context.Context, query string, data map[string]any) (any, error) {
	result, err := neo4j.ExecuteQuery(ctx, d.driver, query, data, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase("neo4j"))
	return result, err
}

func (d *Dao[T]) doSession(ctx context.Context, fun func(tx neo4j.ManagedTransaction) (*Neo4jResult, error), opts ...*SessionOptions) (result *Neo4jResult, err error) {
	if fun == nil {
		return nil, errors.New("doSession(ctx, fun) fun is nil")
	}
	if sc, ok := GetSessionContext(ctx); ok {
		tx := sc.GetTransaction()
		_, err := fun(tx)
		return nil, err
	}

	opt := NewSessionOptions(opts...)
	opt.setDefault()

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: *opt.AccessMode})
	defer func() {
		_ = session.Close(ctx)
		if e1 := errors.GetError(recover()); e1 != nil {
			err = e1
		}
	}()
	/*
		ex, err := session.BeginTransaction(ctx, func(config *neo4j.TransactionConfig) {
			config.Timeout = 50 * time.Second
		})
		if err != nil {
			return nil, err
		}
	*/

	var res any
	if *opt.AccessMode == neo4j.AccessModeRead {
		res, err = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			return fun(tx)
		})
	} else if *opt.AccessMode == neo4j.AccessModeWrite {
		res, err = session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			return fun(tx)
		})
	}
	if err != nil {
		return nil, err
	}
	if result, ok := res.(*Neo4jResult); ok {
		return result, nil
	}
	return nil, err
}

func (d *Dao[T]) Write(ctx context.Context, cypher string) (*Neo4jResult, error) {
	return d.doSession(ctx, func(tx neo4j.ManagedTransaction) (*Neo4jResult, error) {
		result, err := tx.Run(ctx, cypher, nil)
		if err != nil {
			return nil, err
		}
		return NewNeo4jResult(ctx, result), err
	})
}

func (d *Dao[T]) Query(ctx context.Context, cypher string, params map[string]interface{}) (*Neo4jResult, error) {
	var resultData *Neo4jResult
	_, err := d.doSession(ctx, func(tx neo4j.ManagedTransaction) (*Neo4jResult, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			log.Println("wirte to DB with error:", err)
			return nil, err
		}
		resultData = NewNeo4jResult(ctx, result)
		return nil, err
	})
	return resultData, err
}

func (d *Dao[T]) SumEntity(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]T, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) SumMap(ctx context.Context, qry ddd_repository.FindPagingQuery, opts ...ddd_repository.Options) ([]map[string]any, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) Sum(ctx context.Context, qry ddd_repository.FindPagingQuery, resData any, opts ...ddd_repository.Options) (any, bool, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) SumByRSQL(ctx context.Context, tenantId, rSql string, valueCols []*ddd_repository.ValueCol, opts ...ddd_repository.Options) map[string]any {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) CountByMap(ctx context.Context, tenantId string, filterData any, opts ...ddd_repository.Options) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) CountByRSQL(ctx context.Context, tenantId string, rsql string, opts ...ddd_repository.Options) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) StartTx(ctx context.Context, fun ddd_repository.TxFunc, options ...*ddd_repository.SessionOptions) error {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) SetMetadata(metadata map[string]any) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) GetMetadata() map[string]any {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) AddMetadata(key string, val any) {
	//TODO implement me
	panic("implement me")
}

func (d *Dao[T]) DoFilter(ctx context.Context, tenantId string, fun func() (*ddd_repository.FindPagingResult[T], bool, error), opts ...ddd_repository.Options) *ddd_repository.FindPagingResult[T] {
	data, _, err := fun()
	if err != nil {
		return ddd_repository.NewFindPagingResultWithError[T](err)
	}
	return data
}

func (d *Dao[T]) DoList(ctx context.Context, tenantId string, fun func() (*ddd_repository.FindListResult[T], bool, error), opts ...ddd_repository.Options) *ddd_repository.FindListResult[T] {
	data, _, err := fun()
	if err != nil {
		return ddd_repository.NewFindListResultError[T](err)
	}
	return data
}

func (d *Dao[T]) newSetManyResult(result *Neo4jResult, err error) *ddd_repository.SetManyResult[T] {
	if err != nil {
		return ddd_repository.NewSetManyResultError[T](err)
	}
	var data []T
	if err := result.GetList("n", &data); err != nil {
		ddd_repository.NewSetResultError[T](err)
	}
	return ddd_repository.NewSetManyResult[T](data, err)
}

func (d *Dao[T]) doSet(ctx context.Context, tenantId string, cypher string, params map[string]interface{}, opts ...ddd_repository.Options) (*Neo4jResult, error) {
	if err := assert.NotEmpty(tenantId, assert.NewOptions("tenantId is empty")); err != nil {
		return nil, err
	}
	return d.Run(ctx, cypher, params, true, opts...)
}

func (d *Dao[T]) Run(ctx context.Context, cypher string, params map[string]any, isWriteMode bool, opts ...ddd_repository.Options) (*Neo4jResult, error) {
	sOptionsBuilder := NewSessionOptionsBuilder().SetAccessMode(neo4j.AccessModeRead)
	if isWriteMode {
		sOptionsBuilder.SetAccessMode(neo4j.AccessModeWrite)
	}

	res, err := d.doSession(ctx, func(tx neo4j.ManagedTransaction) (*Neo4jResult, error) {
		r, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return NewNeo4jResult(ctx, r), nil
	}, sOptionsBuilder.Build())

	return res, err
}

func (d *Dao[T]) CreateIndex(ctx context.Context, index, label, property string) (err error) {
	cypher := fmt.Sprintf("CREATE INDEX %s IF NOT EXISTS FOR (n:%s) ON (n.%s) ", index, label, property)
	idxSession := d.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})

	_, err = idxSession.Run(ctx, cypher, nil)
	defer func() {
		_ = idxSession.Close(ctx)
	}()
	return err
}

type ImportJsonRelation struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Label      string                  `json:"label"`
	Properties any                     `json:"properties"`
	Start      ImportJsonRelationStart `json:"start"`
	End        ImportJsonRelationEnd   `json:"end"`
}
type ImportJsonRelationStart struct {
	Id         string   `json:"id"`
	Labels     []string `json:"labels"`
	Properties any      `json:"properties"`
}
type ImportJsonRelationEnd struct {
	Id         string   `json:"id"`
	Labels     []string `json:"labels"`
	Properties any      `json:"properties"`
}

type ImportJsonNode struct {
	Id         string   `json:"id"`
	Type       string   `json:"type"`
	Labels     []string `json:"labels"`
	Properties any      `json:"properties"`
}
type Null struct {
}

func (d *Dao[T]) ImportJson(ctx context.Context, cmd ImportJsonCmd, opts ...ddd_repository.Options) (err error) {
	defer func() {
		err = errors.GetRecoverError(err, recover())
	}()
	fileName := cmd.Neo4jPath + "/import/" + cmd.ImportFile

	var jsonFile *os.File
	jsonFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer func() {
		_ = jsonFile.Close()
	}()

	//labelTenant := fmt.Sprintf("tenant_%s", cmd.TenantId)
	labelCase := fmt.Sprintf("case_%s", cmd.CaseId)
	lables := []string{
		//labelTenant,
		labelCase,
	}
	ids := map[string]int{}
	for _, item := range cmd.Nodes {
		id := item.GetId()
		//nodeLables := append(lables, "human")
		if _, ok := ids[id]; ok {
			continue
		}
		ids[id] = 0
		props := map[string]any{"id": item.GetId()}
		node := ImportJsonNode{
			Id:         item.GetId(),
			Type:       "node",
			Labels:     lables,
			Properties: props,
		}
		if item, err := jsonutils.Marshal(node); err != nil {
			return err
		} else {
			jsonFile.WriteString(item)
			jsonFile.WriteString("\r\n")
		}
	}

	for _, item := range cmd.Relations {
		rel := ImportJsonRelation{
			Id:         item.GetId(),
			Type:       "relationship",
			Label:      item.GetRelType(),
			Properties: item.GetProperties(),
			Start: ImportJsonRelationStart{
				Id:         item.GetStartId(),
				Labels:     lables,
				Properties: Null{},
			},
			End: ImportJsonRelationEnd{
				Id:         item.GetEndId(),
				Labels:     lables,
				Properties: Null{},
			},
		}
		if item, err := jsonutils.Marshal(rel); err != nil {
			return err
		} else {
			jsonFile.WriteString(item)
			jsonFile.WriteString("\r\n")
		}
	}

	cypher := fmt.Sprintf(`CALL apoc.import.json("file:///%s",{cleanup:false, importIdName:"id"} )`, cmd.ImportFile)

	fmt.Println("***********")
	logs.Debug(ctx, "", logs.Fields{"cypher": cypher})
	fmt.Println("***********")

	session := d.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	_, err = session.Run(ctx, cypher, nil)
	if err != nil {
		fmt.Println(err)

	} else {
		//summary, _ := result.Consume(ctx)
		//fmt.Println("Query updated the database?", summary.Counters().ContainsUpdates())
	}

	return err
}

func getLabels(labels ...string) string {
	var s string
	for _, l := range labels {
		if len(l) > 0 {
			s = fmt.Sprintf("%v :`%v`", s, l)
		}
	}
	return strings.ToLower(s)
}
