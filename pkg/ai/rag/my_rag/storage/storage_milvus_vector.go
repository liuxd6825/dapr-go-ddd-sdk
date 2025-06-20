package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/milvus-io/milvus/client/v2/index"
	"strconv"
	"time"

	"github.com/milvus-io/milvus/client/v2/entity"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

type MilvusVector struct {
	client        *milvusclient.Client
	cfg           MilvusConfig
	loadTenants   map[string]bool
	embedderStore EmbedderStore
}

type MilvusConfig struct {
	Addr           string `json:"addr"`
	CollectionName string `json:"collectionName"`
	Dim            int    `json:"dim"`
	DBName         string `json:"dbName"`
	TopK           int    `json:"topK"`
}

const (
	milvusEntitiesCollectionName      = "entities"
	milvusRelationshipsCollectionName = "relationships"
	cosineThreshold                   = 0.2
)

func NewMilvusVector(embedderStore EmbedderStore, cfg MilvusConfig) *MilvusVector {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Connect to Milvus
	client, err := milvusclient.New(ctx, &milvusclient.ClientConfig{
		Address: cfg.Addr,
		DBName:  cfg.DBName,
	})
	if err != nil {
		panic(fmt.Errorf("failed to connect to Milvus: %w", err))
	}

	if cfg.Dim <= 0 {
		panic(fmt.Errorf("dim should be greater than zero"))
	}

	return &MilvusVector{client: client, cfg: cfg, loadTenants: map[string]bool{}, embedderStore: embedderStore}
}

func (m *MilvusVector) DeleteTenant(ctx context.Context, tenantId string) error {
	entName := m.getEntityCollName(tenantId)
	relName := m.getRelCollName(tenantId)
	docName := m.getDocCollName(tenantId)
	return m.deleteCollection(ctx, entName, relName, docName)
}

func (m *MilvusVector) deleteCollection(ctx context.Context, name ...string) error {
	for _, itemName := range name {
		opts := milvusclient.NewDropCollectionOption(itemName)
		if err := m.client.DropCollection(ctx, opts); err != nil {
			return err
		}
	}
	return nil
}

func (m *MilvusVector) DeleteCase(ctx context.Context, tenantId, caseId string) error {
	collName := m.getEntityCollName(tenantId)
	opts := milvusclient.NewDropPartitionOption(collName, "")
	return m.client.DropPartition(ctx, opts)
}

func (m *MilvusVector) DeleteDoc(ctx context.Context, tenantId, caseId, docId string) error {
	collName := m.getEntityCollName(tenantId)
	delOption := milvusclient.NewDeleteOption(collName)
	delOption.WithExpr(fmt.Sprintf(`case_id=="%s" and doc_id=="%s"`, caseId, docId))
	_, err := m.client.Delete(ctx, delOption)
	return err
}

func (m *MilvusVector) CreateTenant(ctx context.Context, tenantId string) error {
	if err := m.createEntitiesCollection(ctx, tenantId); err != nil {
		return err
	}
	if err := m.createRelationshipsCollection(ctx, tenantId); err != nil {
		return err
	}
	if err := m.createDocCollection(ctx, tenantId); err != nil {
		return err
	}
	return nil
}

func (m *MilvusVector) LoadTenant(ctx context.Context, tenantId string) (err error) {
	if _, ok := m.loadTenants[tenantId]; ok {
		return nil
	}

	m.loadTenants[tenantId] = true
	entName := m.getEntityCollName(tenantId)
	relName := m.getRelCollName(tenantId)
	docName := m.getDocCollName(tenantId)

	if err := m.LoadCollection(ctx, entName, relName, docName); err != nil {
		m.loadTenants[tenantId] = false
		return err
	}
	return nil
}

func (m *MilvusVector) LoadCollection(ctx context.Context, collName ...string) error {
	for _, itemName := range collName {
		hasColl, err := m.client.HasCollection(ctx, milvusclient.NewHasCollectionOption(itemName))
		if err != nil {
			return err
		}
		if !hasColl {
			return errors.New(fmt.Sprintf("collection %s not exist", itemName))
		}
		loadOpts := milvusclient.NewLoadCollectionOption(itemName)
		if _, err := m.client.LoadCollection(ctx, loadOpts); err != nil {
			return err
		}
	}
	return nil
}

func (m *MilvusVector) CreateCase(ctx context.Context, tenantId, caseId string) error {
	return nil
}

func (m *MilvusVector) VectorFlush(ctx context.Context, opts Options) error {
	entCollName := m.getEntityCollName(opts.TenantId)
	entOpt := milvusclient.NewFlushOption(entCollName)
	if _, err := m.client.Flush(ctx, entOpt); err != nil {
		return err
	}
	return nil
}

func (m *MilvusVector) VectorInsertDoc(ctx context.Context, data InsertDocData, opts Options) error {
	docCollName := m.getDocCollName(opts.TenantId)
	id := m.getId(opts.TenantId, opts.CaseId, data.DocId)
	id = fmt.Sprintf("%s_%d", id, data.ChunkId)
	// 向量化文本块
	vector, err := m.embedderStore.EmbedTexts(ctx, data.Content)
	if err != nil {
		return err
	}
	insertOpts := milvusclient.NewColumnBasedInsertOption(docCollName).
		WithVarcharColumn("id", []string{id}).
		WithVarcharColumn("content", data.Content).
		WithVarcharColumn("case_id", []string{opts.CaseId}).
		WithVarcharColumn("doc_id", []string{data.DocId}).
		WithVarcharColumn("file_name", []string{data.FileName}).
		WithFloatVectorColumn("vector", m.cfg.Dim, vector)
	_, err = m.client.Insert(ctx, insertOpts)
	if err != nil {
		println(err)
	}
	return err
}

func (m *MilvusVector) VectorSearchDoc(ctx context.Context, vector []float32, topK int, opts Options) ([]string, error) {
	if err := m.LoadTenant(ctx, opts.TenantId); err != nil {
		return nil, err
	}
	//sp, _ := entity.NewIndexHNSWSearchParam(200) // 搜索参数
	entCollName := m.getEntityCollName(opts.TenantId)

	// 设置过滤条件
	// 使用 expression 来构造查询条件
	expr := fmt.Sprintf("case_id='%s'", opts.CaseId)

	searchOpts := milvusclient.NewSearchOption(entCollName, topK, []entity.Vector{entity.FloatVector(vector)}).WithFilter(expr)
	results, err := m.client.Search(ctx, searchOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to search Milvus: %w", err)
	}

	texts := make([]string, 0, topK)
	for _, result := range results {
		for _, field := range result.Fields {
			if field.Name() == "context" {
				for i := 0; i < field.Len(); i++ {
					text, _ := field.GetAsString(i)
					texts = append(texts, text)
				}
			}
		}
	}

	return texts, nil
}

func (m *MilvusVector) VectorQueryEntity(ctx context.Context, keywords string, opts Options) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	vectorList, err := m.embedderStore.EmbedTexts(ctx, []string{keywords})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding for query: %w", err)
	}
	vectors := []entity.Vector{}
	for _, t := range vectorList {
		vectors = append(vectors, entity.FloatVector(t))
	}

	annParam := index.NewCustomAnnParam()
	annParam.WithRadius(cosineThreshold)
	opt := milvusclient.
		NewSearchOption(milvusEntitiesCollectionName, m.cfg.TopK, vectors).
		WithOutputFields("name").
		WithAnnParam(annParam)
	searchResult, err := m.client.Search(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to query entities: %w", err)
	}

	results := make([]string, 0, m.cfg.TopK)
	for _, result := range searchResult {
		for i := 0; i < result.ResultCount; i++ {
			entityName, err := result.GetColumn("name").Get(i)
			if err != nil {
				return nil, fmt.Errorf("failed to get entity name from result: %w", err)
			}
			entityNameStr, ok := entityName.(string)
			if !ok {
				return nil, fmt.Errorf("entity name not string")
			}
			// Milvus returns strings with surrounding quotes, remove them
			cleanStr, err := strconv.Unquote(entityNameStr)
			if err != nil {
				if !errors.Is(err, strconv.ErrSyntax) {
					return nil, fmt.Errorf("failed to unquote entity name: %w", err)
				}
				// ErrSyntax means the string is not surrounded by quotes, so we can use it as is
				cleanStr = entityNameStr
			}
			results = append(results, cleanStr)
		}
	}

	return results, nil
}

func (m *MilvusVector) VectorQueryRelationship(ctx context.Context, keywords string, opts Options) ([][2]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	vectorList, err := m.embedderStore.EmbedTexts(ctx, []string{keywords})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding for query: %w", err)
	}

	vectors := []entity.Vector{}
	for _, t := range vectorList {
		vectors = append(vectors, entity.FloatVector(t))
	}

	annParam := index.NewCustomAnnParam()
	annParam.WithRadius(cosineThreshold)
	opt := milvusclient.
		NewSearchOption(milvusRelationshipsCollectionName, m.cfg.TopK, vectors).
		WithOutputFields("source", "target").
		WithAnnParam(annParam)
	searchResult, err := m.client.Search(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to query relationships: %w", err)
	}

	results := make([][2]string, 0, m.cfg.TopK)
	for _, result := range searchResult {
		for i := 0; i < result.ResultCount; i++ {
			sourceEntity, err := result.GetColumn("source_entity").Get(i)
			if err != nil {
				return nil, fmt.Errorf("failed to get source entity from result: %w", err)
			}
			sourceEntityStr, ok := sourceEntity.(string)
			if !ok {
				return nil, fmt.Errorf("source entity not string")
			}
			// Milvus returns strings with surrounding quotes, remove them
			sourceCleanStr, err := strconv.Unquote(sourceEntityStr)
			if err != nil {
				if !errors.Is(err, strconv.ErrSyntax) {
					return nil, fmt.Errorf("failed to unquote source entity: %w", err)
				}
				// ErrSyntax means the string is not surrounded by quotes, so we can use it as is
				sourceCleanStr = sourceEntityStr
			}

			targetEntity, err := result.GetColumn("target_entity").Get(i)
			if err != nil {
				return nil, fmt.Errorf("failed to get target entity from result: %w", err)
			}
			targetEntityStr, ok := targetEntity.(string)
			if !ok {
				return nil, fmt.Errorf("target entity not string")
			}
			// Milvus returns strings with surrounding quotes, remove them
			targetCleanStr, err := strconv.Unquote(targetEntityStr)
			if err != nil {
				if !errors.Is(err, strconv.ErrSyntax) {
					return nil, fmt.Errorf("failed to unquote target entity: %w", err)
				}
				// ErrSyntax means the string is not surrounded by quotes, so we can use it as is
				targetCleanStr = targetEntityStr
			}

			results = append(results, [2]string{sourceCleanStr, targetCleanStr})
		}
	}

	return results, nil
}

func (m *MilvusVector) getId(tenantId, caseId, name string) string {
	return fmt.Sprintf("%s_%s_%s", tenantId, caseId, name)
}

func (m *MilvusVector) VectorUpsertEntity(ctx context.Context, entity *VectorUpsertEntity) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	vector, err := m.embedderStore.EmbedTexts(ctx, entity.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding for upsert: %w", err)
	}
	id := m.getId(entity.TenantId, entity.CaseId, entity.Name)
	collName := m.getEntityCollName(entity.TenantId)
	upsertOpts := milvusclient.NewColumnBasedInsertOption(collName).
		WithVarcharColumn("id", []string{id}).
		WithVarcharColumn("name", []string{entity.Name}).
		WithVarcharColumn("content", entity.Content).
		WithVarcharColumn("case_id", []string{entity.CaseId}).
		WithVarcharColumn("doc_id", []string{entity.DocId}).
		WithVarcharColumn("file_name", []string{entity.FileName}).
		WithFloatVectorColumn("vector", m.cfg.Dim, vector)
	_, err = m.client.Upsert(ctx, upsertOpts)
	if err != nil {
		return fmt.Errorf("failed to upsert entity: %w", err)
	}
	return nil
}

func (m *MilvusVector) VectorUpsertRelationship(ctx context.Context, rel *VectorUpsertRelationship) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	vector, err := m.embedderStore.EmbedTexts(ctx, rel.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding for upsert: %w", err)
	}
	collName := m.getRelCollName(rel.TenantId)
	id := fmt.Sprintf("%s-%s-%s", rel.CaseId, rel.Source, rel.Target)
	upsertOpts := milvusclient.NewColumnBasedInsertOption(collName).
		WithVarcharColumn("id", []string{id}).
		WithVarcharColumn("source", []string{rel.Source}).
		WithVarcharColumn("target", []string{rel.Target}).
		WithVarcharColumn("case_id", []string{rel.CaseId}).
		WithVarcharColumn("doc_id", []string{rel.DocId}).
		WithVarcharColumn("content", rel.Content).
		WithVarcharColumn("file_name", []string{rel.FileName}).
		WithFloatVectorColumn("vector", m.cfg.Dim, vector)

	_, err = m.client.Upsert(ctx, upsertOpts)
	if err != nil {
		return fmt.Errorf("failed to upsert relationship: %w", err)
	}
	return nil
}

func (m *MilvusVector) getDocCollName(tenantId string) string {
	return fmt.Sprintf("%s_%s_doc", m.cfg.CollectionName, tenantId)
}

func (m *MilvusVector) getEntityCollName(tenantId string) string {
	return fmt.Sprintf("%s_%s_ent", m.cfg.CollectionName, tenantId)
}

func (m *MilvusVector) getRelCollName(tenantId string) string {
	return fmt.Sprintf("%s_%s_rel", m.cfg.CollectionName, tenantId)
}

// Close closes the connection to Milvus.
func (m *MilvusVector) Close(ctx context.Context) error {
	if m.client != nil {
		return m.client.Close(ctx)
	}
	return nil
}

func (m *MilvusVector) createEntitiesCollection(ctx context.Context, tenantId string) error {
	collName := m.getEntityCollName(tenantId)
	// 定义集合 schema
	schema := entity.NewSchema().
		WithName(collName).
		WithDescription("RAG Knowledge Base TenantId " + tenantId).
		WithAutoID(false).
		WithField(entity.NewField().
			WithName("id").
			WithDataType(entity.FieldTypeVarChar).
			WithIsPrimaryKey(true).
			WithIsAutoID(false).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("vector").
			WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(m.cfg.Dim))).
		WithField(entity.NewField().
			WithName("name").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("content").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(65535)).
		WithField(entity.NewField().
			WithName("doc_id").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("case_id").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("file_name").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256))

	createOpts := milvusclient.NewCreateCollectionOption(collName, schema)
	err := m.client.CreateCollection(ctx, createOpts)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	m.createIndexVector(ctx, collName, "vector")
	m.createIndexVarChar(ctx, collName, "case_id")
	m.createIndexVarChar(ctx, collName, "doc_id")
	return err
}

func (m *MilvusVector) createDocCollection(ctx context.Context, tenantId string) error {
	collName := m.getDocCollName(tenantId)
	// 定义集合 schema
	schema := entity.NewSchema().
		WithName(collName).
		WithDescription("RAG Knowledge Base TenantId " + tenantId).
		WithAutoID(false).
		WithField(entity.NewField().
			WithName("id").
			WithDataType(entity.FieldTypeVarChar).
			WithIsPrimaryKey(true).
			WithIsAutoID(false).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("vector").
			WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(m.cfg.Dim))).
		WithField(entity.NewField().
			WithName("content").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(65535)).
		WithField(entity.NewField().
			WithName("doc_id").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("case_id").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("file_name").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256))

	createOpts := milvusclient.NewCreateCollectionOption(collName, schema)
	err := m.client.CreateCollection(ctx, createOpts)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	m.createIndexVector(ctx, collName, "vector")
	m.createIndexVarChar(ctx, collName, "case_id")
	m.createIndexVarChar(ctx, collName, "doc_id")
	return err
}

func (m *MilvusVector) createRelationshipsCollection(ctx context.Context, tenantId string) error {
	relCollName := m.getRelCollName(tenantId)
	// 定义集合 schema
	schema := entity.NewSchema().
		WithName(relCollName).
		WithDescription("RAG Knowledge Base TenantId " + tenantId).
		WithAutoID(false).
		WithField(entity.NewField().
			WithName("id").
			WithDataType(entity.FieldTypeVarChar).
			WithIsPrimaryKey(true).
			WithIsAutoID(false).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("source").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("target").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("vector").
			WithDataType(entity.FieldTypeFloatVector).
			WithDim(int64(m.cfg.Dim))).
		WithField(entity.NewField().
			WithName("content").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(65535)).
		WithField(entity.NewField().
			WithName("doc_id").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("case_id").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256)).
		WithField(entity.NewField().
			WithName("file_name").
			WithDataType(entity.FieldTypeVarChar).
			WithMaxLength(256))

	createOpts := milvusclient.NewCreateCollectionOption(relCollName, schema)
	err := m.client.CreateCollection(ctx, createOpts)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	m.createIndexVector(ctx, relCollName, "vector")
	m.createIndexVarChar(ctx, relCollName, "case_id")
	m.createIndexVarChar(ctx, relCollName, "doc_id")
	return err
}

func (m *MilvusVector) createIndexVector(ctx context.Context, collName string, fieldName string) error {
	idx := index.NewHNSWIndex(entity.L2, 8, 200)
	idxOption := milvusclient.NewCreateIndexOption(collName, fieldName, idx)
	_, err := m.client.CreateIndex(ctx, idxOption)
	return err
}

func (m *MilvusVector) createIndexVarChar(ctx context.Context, collName string, fieldName string) error {
	idx := index.NewInvertedIndex()
	idxOption := milvusclient.NewCreateIndexOption(collName, fieldName, idx)
	_, err := m.client.CreateIndex(ctx, idxOption)
	return err
}

// truncateString 截断字符串到指定长度
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
