package storage

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dbschema"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/env"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/logs"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/utils/maputils"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// Neo4jGraphStorage provides a Neo4j graph database implementation of storage interfaces.
// It handles database connections and operations for storing and retrieving graph entities
// and relationships.
type Neo4jGraphStorage struct {
	client neo4j.DriverWithContext
	dao    idao.Dao[map[string]any]
	logger logs.Logger
}

// NewNeo4jGraphStorage creates a new Neo4j client connection with the provided connection parameters.
// It returns an initialized Neo4J struct and any error encountered during connection setup.
// The returned Neo4J instance must be closed with Close() when no longer needed to free up resources.
func NewNeo4jGraphStorage(dbKey string, logger logs.Logger) *Neo4jGraphStorage {
	nodeCfg := &dao.DaoConfig{
		DBKey:              dbKey,
		GraphType:          idao.GraphType_Node,
		GraphLabels:        []string{"doc"},
		IsCancelModified:   true,
		IsCancelSoftDelete: true,
		DBSchema:           dbschema.NewDBSchema("graph", "graph"),
	}
	newDao := dao.NewDao[map[string]any](nodeCfg)
	dbItem := env.GetEnv().GetDB(dbKey)
	if dbItem == nil {
		panic(fmt.Errorf("dbKey is nil"))
	}

	return &Neo4jGraphStorage{
		client: dbItem.GetNeo4j(),
		dao:    newDao,
		logger: logger,
	}
}

func (n *Neo4jGraphStorage) CreateTenant(ctx context.Context, tenantId string) error {
	return nil
}

func (n *Neo4jGraphStorage) CreateCase(ctx context.Context, tenantId, caseId string) error {
	return nil
}

func (n *Neo4jGraphStorage) DeleteTenant(ctx context.Context, tenantId string) error {
	labels := fmt.Sprintf(":tenant_%s:doc", tenantId)
	_, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := fmt.Sprintf("MATCH (n%s) DETACH DELETE n", labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}
			return queryRes.Record(), nil
		})
	})
	return err
}

func (n *Neo4jGraphStorage) DeleteCase(ctx context.Context, tenantId, caseId string) error {
	labels := fmt.Sprintf(":tenant_%s:case_%s:doc", tenantId, caseId)
	_, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := fmt.Sprintf("MATCH (n%s) DETACH DELETE n", labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}
			return queryRes.Record(), nil
		})
	})
	return err
}

func (n *Neo4jGraphStorage) DeleteDoc(ctx context.Context, tenantId, caseId, docId string) error {
	labels := fmt.Sprintf(":tenant_%s:case_%s:doc_%s:doc", tenantId, caseId, docId)
	_, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := fmt.Sprintf("MATCH (n%s) DETACH DELETE n", labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}
			return queryRes.Record(), nil
		})
	})
	return err
}

func (n *Neo4jGraphStorage) LoadTenant(ctx context.Context, tenantId, caseId string) error {
	return nil
}

func (n *Neo4jGraphStorage) GraphSaveDoc(ctx context.Context, tenantId, caseId, docId string, entries []*GraphEntity, relationships []*GraphRelationship) error {
	err := n.GraphSaveDocEntities(ctx, tenantId, caseId, docId, entries)
	if err != nil {
		return err
	}
	return n.GraphSaveDocRelationships(ctx, tenantId, caseId, docId, relationships)
}

func (n *Neo4jGraphStorage) GraphSaveDocEntities(ctx context.Context, tenantId, caseId, docId string, entries []*GraphEntity) error {
	labelGroup := n.getGroupEntity(ctx, entries)
	// 并发处理不同标签组
	var wg sync.WaitGroup
	errChan := make(chan error, len(labelGroup))

	for label, list := range labelGroup {
		err := n.graphSaveDocEntities(ctx, tenantId, caseId, docId, label, list)
		if err != nil {
			errChan <- err
		}
	}

	wg.Wait()
	close(errChan)

	// 收集所有错误
	var errors []string
	for err := range errChan {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("batch processing errors:\n%s", strings.Join(errors, "\n"))
	}

	return nil
}

func (n *Neo4jGraphStorage) graphSaveDocEntities(ctx context.Context, tenantId, caseId, docId, label string, entries []*GraphEntity) error {

	// 2. 准备批量数据（实际可从JSON/CSV加载）
	batchData := []map[string]any{}
	for _, entry := range entries {
		batchData = append(batchData, map[string]any{
			"id":          entry.Id,
			"name":        entry.Name,
			"case_id":     entry.CaseId,
			"doc_id":      entry.DocId,
			"tenant_id":   entry.TenantId,
			"description": entry.Descriptions,
			"source_ids":  entry.SourceIDs,
			"source_type": "doc",
			"type":        entry.Type,
		})
	}

	// 3. 构建APOC执行参数
	params := map[string]interface{}{
		"batch": batchData,
	}

	labels := fmt.Sprintf(":tenant_%s:case_%s:doc_%s:%s:doc", tenantId, caseId, docId, label)

	cypher := fmt.Sprintf(`
	UNWIND $batch AS row
	MERGE (n%s{name:row.name})
		ON CREATE
		  SET n = row, n.created_at = datetime()
		ON MATCH
		  SET n.description =  n.description + "\n" +  row.description
    `, labels)

	_, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to run write: %w", err)
			}
			err = result.Err()
			if err != nil {
				return nil, fmt.Errorf("failed to run write: %w", err)
			}
			return nil, nil
		})
	})
	return err

}

func (n *Neo4jGraphStorage) GraphSaveDocRelationships(ctx context.Context, tenantId, caseId, docId string, relationships []*GraphRelationship) error {
	// 2. 准备批量数据（实际可从JSON/CSV加载）
	batchData := []map[string]any{}
	for _, rel := range relationships {
		batchData = append(batchData, map[string]any{
			"id":          rel.Id,
			"target":      rel.Target,
			"source":      rel.Source,
			"case_id":     rel.CaseId,
			"tenant_id":   rel.TenantId,
			"doc_id":      rel.DocId,
			"description": rel.Descriptions,
			"source_ids":  rel.SourceIDs,
			"source_type": "doc",
			"keywords":    rel.Keywords,
		})
	}

	// 3. 构建APOC执行参数
	params := map[string]interface{}{
		"batch": batchData,
	}

	labels := fmt.Sprintf(":tenant_%s:case_%s:doc_%s:doc", tenantId, caseId, docId)

	cypher := fmt.Sprintf(`
	UNWIND $batch AS row
	MATCH (n%s{id:row.source}),(m%s{id:row.target})
	MERGE (n)-[r:DIRECTED{id:row.id}]->(m)
		ON CREATE
		  SET r = row, r.created_at = datetime()
		ON MATCH
		  SET r.description =  r.description + "\n" +  row.description
              
    `, labels, labels)

	_, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			result, err := tx.Run(ctx, cypher, params)
			if err != nil {
				return nil, fmt.Errorf("failed to run write: %w", err)
			}
			err = result.Err()
			if err != nil {
				return nil, fmt.Errorf("failed to run write: %w", err)
			}
			return nil, nil
		})
	})
	return err

}

func (n *Neo4jGraphStorage) getGroupEntity(ctx context.Context, entries []*GraphEntity) map[string][]*GraphEntity {
	res := make(map[string][]*GraphEntity)
	for _, entry := range entries {
		if list, ok := res[entry.Type]; !ok {
			var list []*GraphEntity
			list = append(list, entry)
			res[entry.Type] = list
		} else {
			res[entry.Type] = append(list, entry)
		}
	}
	return res
}

func (n *Neo4jGraphStorage) GraphQuery(ctx context.Context, query GraphQueryParam, opts Options) ([]string, error) {
	contents := []string{}
	nodes, rels, err := n.FindNodes(ctx, query, opts)
	if err != nil {
		return nil, err
	}
	for _, node := range nodes {
		contents = append(contents, node.Descriptions)
	}
	for _, edge := range rels {
		source, ok1 := nodes[edge.Source]
		target, ok2 := nodes[edge.Target]
		if ok1 && ok2 {
			contents = append(contents, fmt.Sprintf("%s与%s之间存在关系是:%s, %s", source.Name, target.Name, strings.Join(edge.Keywords, ","), edge.Descriptions))
		}
	}
	return contents, nil
}

func (n *Neo4jGraphStorage) getNames(values []string) string {
	count := len(values)
	sb := strings.Builder{}
	for i, value := range values {
		sb.WriteString("\"" + value + "\"")
		if i < count-1 {
			sb.WriteString(",")
		}
	}
	return sb.String()
}

// FindNodes
/*
	MATCH p=(n)-[*..5]-(m) 	WHERE n.name IN ['名称1', '名称2'] RETURN p
*/
func (n *Neo4jGraphStorage) FindNodes(ctx context.Context, query GraphQueryParam, opts Options) (nodes map[string]*GraphEntity, rels map[string]*GraphRelationship, err error) {
	namesStr := n.getNames(query.Keys)
	limit := query.Limit
	if limit <= 0 {
		limit = 1000
	}
	maxDeep := query.MaxDeep
	if maxDeep <= 0 {
		maxDeep = 5
	}
	labels := n.getLabels(opts)

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(` MATCH p=(n%s)-[*..%d]-(m) `, labels, query.MaxDeep))
	sb.WriteString(fmt.Sprintf(` WHERE ANY(word IN [%s] WHERE n.description CONTAINS word) `, namesStr))
	sb.WriteString(fmt.Sprintf(` OR n.name IN [%s] RETURN p LIMIT 1000`, namesStr))
	n.logger.Info(sb.String())

	session := n.client.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	nodes = make(map[string]*GraphEntity)
	rels = make(map[string]*GraphRelationship)

	err = n.executeQuery(ctx, session, sb.String(), nil, func(result neo4j.ResultWithContext) error {
		for result.Next(ctx) {
			record := result.Record()
			for _, key := range record.Keys {
				if v, ok := record.Get(key); ok {
					switch v.(type) {
					case dbtype.Path:
						if path, ok := v.(dbtype.Path); ok {
							for _, n := range path.Nodes {
								node := newGraphEntity(n)
								nodes[node.Id] = node
							}
							for _, r := range path.Relationships {
								rel := newGraphRelationship(r)
								rels[rel.Id] = rel
							}
						}
						break
					}
				}
			}
		}
		return result.Err()
	})
	return
}

func newGraphEntity(node neo4j.Node) *GraphEntity {
	id := maputils.GetStringErr(node.Props, "id", "")
	name := maputils.GetStringErr(node.Props, "name", "")
	typeName := maputils.GetStringErr(node.Props, "type", "")
	caseId := maputils.GetStringErr(node.Props, "case_id", "")
	docId := maputils.GetStringErr(node.Props, "doc_id", "")
	descriptions := maputils.GetStringErr(node.Props, "description", "")
	sourceIDs := maputils.GetStringErr(node.Props, "source_ids", "")
	return &GraphEntity{
		Id:           id,
		Name:         name,
		Type:         typeName,
		CaseId:       caseId,
		DocId:        docId,
		Descriptions: descriptions,
		SourceIDs:    sourceIDs,
	}
}

func newGraphRelationship(rel neo4j.Relationship) *GraphRelationship {
	id := maputils.GetStringErr(rel.Props, "id", "")
	source := maputils.GetStringErr(rel.Props, "source", "")
	target := maputils.GetStringErr(rel.Props, "target", "")
	caseId := maputils.GetStringErr(rel.Props, "case_id", "")
	docId := maputils.GetStringErr(rel.Props, "doc_id", "")
	descriptions := maputils.GetStringErr(rel.Props, "description", "")
	sourceIDs := maputils.GetStringErr(rel.Props, "source_ids", "")
	if id == "uUMkiqE2VKvkLNik2b96-10-uUMkiqE2VKvkLNik2b96-11" {
		println(id)
	}
	keywords := maputils.GetStringsErr(rel.Props, "keywords", nil)
	return &GraphRelationship{
		Id:           id,
		Source:       source,
		Target:       target,
		CaseId:       caseId,
		DocId:        docId,
		Descriptions: descriptions,
		SourceIDs:    sourceIDs,
		Keywords:     keywords,
	}
}

func (n *Neo4jGraphStorage) executeQuery(ctx context.Context, session neo4j.SessionWithContext, cypher string, params map[string]any, data func(withContext neo4j.ResultWithContext) error) error {
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		result, err := tx.Run(ctx, cypher, params)
		if err != nil {
			return nil, err
		}
		return result, data(result)
	})
	return err
}

func (n *Neo4jGraphStorage) getNode(ctx context.Context, res any, getNode func(node *neo4j.Node), keys ...string) {
	for _, key := range keys {
		if result, ok := res.(neo4j.ResultWithContext); ok {
			if result.Next(ctx) {
				record := result.Record()
				if record == nil {
					continue
				}
				if data, ok := record.Get(key); ok {
					fmt.Println(data)
				}
			}
		}
	}
}

// GraphEntity retrieves a graph entity by name from the Neo4j database.
// It returns the found entity or an error if the entity doesn't exist or if the query fails.
func (n *Neo4jGraphStorage) GraphEntity(ctx context.Context, name string, opts Options) (*GraphEntity, error) {
	labels := n.getDocLabels(opts)
	res, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := fmt.Sprintf("MATCH (n%s {name:$name}) RETURN n", labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{
				"name": name,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}
			return queryRes.Record(), nil
		})
	})
	if err != nil {
		return nil, err
	}
	record, ok := res.(*db.Record)
	if !ok {
		return nil, fmt.Errorf("invalid result type, got %T, want *db.Record", res)
	}
	if record == nil {
		return nil, nil
	}
	nNode, ok := record.Get("n")
	if !ok {
		return nil, nil
	}
	node, ok := nNode.(dbtype.Node)
	if !ok {
		return nil, fmt.Errorf("invalid n type, got %T, want dbtype.Node", n)
	}

	return graphEntityFromNode(node), nil
}

func (n *Neo4jGraphStorage) getLabels(opts Options) string {
	switch opts.NodeLabel {
	case NodeLabel_Master:
		return fmt.Sprintf(":tenant_%s:case_%s:master", opts.TenantId, opts.CaseId)
	case NodeLabel_Draw:
		return fmt.Sprintf(":tenant_%s:case_%s:draw", opts.TenantId, opts.CaseId)
	case NodeLabel_Doc, "":
		return fmt.Sprintf(":tenant_%s:case_%s:doc", opts.TenantId, opts.CaseId)
	case NodeLabel_All:
		return fmt.Sprintf(":tenant_%s:case_%s", opts.TenantId, opts.CaseId)
	default:
		panic(fmt.Sprintf("invalid label type, got %s, want %s,%s,%s,%s", opts.NodeLabel, NodeLabel_Master, NodeLabel_Draw, NodeLabel_Doc, NodeLabel_All))
	}
}

func (n *Neo4jGraphStorage) getDocLabels(opts Options) string {
	return fmt.Sprintf(":tenant_%s:case_%s:doc_%s:doc", opts.TenantId, opts.CaseId, opts.DocId)
}

// GraphRelationship retrieves a relationship between two entities from the Neo4j database.
// It returns the found relationship or an error if the relationship doesn't exist or if the query fails.
func (n *Neo4jGraphStorage) GraphRelationship(ctx context.Context, sourceEntity, targetEntity string, opts Options) (*GraphRelationship, error) {
	res, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			labels := n.getLabels(opts)
			query := fmt.Sprintf(`
MATCH (start%s{name: $source_entity_id})-[r]-(end%s{name: $target_entity_id})
RETURN properties(r) as edge_properties
`, labels, labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{
				"source_entity_id": sourceEntity,
				"target_entity_id": targetEntity,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}

			return queryRes.Record(), nil
		})
	})
	if err != nil {
		return nil, err
	}
	record, ok := res.(*db.Record)
	if !ok {
		return nil, fmt.Errorf("invalid result type, got %T, want *db.Record", res)
	}
	if record == nil {
		return nil, nil
	}
	edgeProps, ok := record.Get("edge_properties")
	if !ok {
		return nil, fmt.Errorf("expected edge_properties key is not found")
	}
	props, ok := edgeProps.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid edge_properties type, got %T, want map[string]any", edgeProps)
	}

	return graphRelationshipFromEdge(sourceEntity, targetEntity, props), nil
}

// GraphUpsertEntity creates or updates an entity in the Neo4j graph database.
// It returns an error if the database operation fails.
func (n *Neo4jGraphStorage) GraphUpsertEntity(ctx context.Context, entity *GraphEntity, opts Options) error {
	_, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			labels := n.getDocLabels(opts)
			return tx.Run(
				ctx,
				fmt.Sprintf(`
MERGE (n%s {name: $properties.name})
SET n += $properties
SET n:%s`, labels, "`"+entity.Type+"`"),
				map[string]any{
					"properties": map[string]any{
						"name":        entity.Name,
						"type":        entity.Type,
						"description": entity.Descriptions,
						"source_ids":  entity.SourceIDs,
						"created_at":  entity.CreatedAt.Format(time.RFC3339),
					},
				},
			)
		})
	})

	return err
}

// GraphUpsertRelationship creates or updates a relationship between two entities in the Neo4j graph database.
// It returns an error if the database operation fails.
func (n *Neo4jGraphStorage) GraphUpsertRelationship(ctx context.Context, relationship *GraphRelationship, opts Options) error {
	_, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			keywords := strings.Join(relationship.Keywords, GraphFieldSeparator)
			labels := n.getLabels(opts)
			return tx.Run(
				ctx,
				fmt.Sprintf(`
MATCH (source%s {name: $source})
WITH source
MATCH (target%s {name: $target})
MERGE (source)-[r:DIRECTED]-(target)
SET r += $properties
`, labels, labels),
				map[string]any{
					"source": relationship.Source,
					"target": relationship.Target,
					"properties": map[string]any{
						"weight":      relationship.Weight,
						"description": relationship.Descriptions,
						"keywords":    keywords,
						"source_ids":  relationship.SourceIDs,
						"created_at":  relationship.CreatedAt.Format(time.RFC3339),
					},
				},
			)
		})
	})

	return err
}

// GraphEntities retrieves multiple graph entities by their names from the Neo4j database.
// It returns a map of entity names to GraphEntity objects, or an error if the query fails.
func (n *Neo4jGraphStorage) GraphEntities(ctx context.Context, names []string, opts Options) (map[string]*GraphEntity, error) {
	if len(names) == 0 {
		return nil, nil
	}
	labels := n.getDocLabels(opts)
	res, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := fmt.Sprintf(`
MATCH (n%s) 
WHERE n.name IN $entityIDs 
RETURN n, n.name as entity_id`, labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{
				"entityIDs": names,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}

			result := make(map[string]dbtype.Node)
			for record, err := range queryRes.Records(ctx) {
				if err != nil {
					return nil, fmt.Errorf("failed to get result: %w", err)
				}

				node, ok := record.Get("n")
				if !ok {
					continue
				}

				entityID, ok := record.Get("name")
				if !ok {
					continue
				}

				entityIDStr, ok := entityID.(string)
				if !ok {
					continue
				}

				dbNode, ok := node.(dbtype.Node)
				if !ok {
					continue
				}

				result[entityIDStr] = dbNode
			}

			return result, nil
		})
	})
	if err != nil {
		return nil, err
	}

	nodeMap, ok := res.(map[string]dbtype.Node)
	if !ok {
		return nil, fmt.Errorf("invalid result type, got %T, want map[string]dbtype.Node", res)
	}

	entities := make(map[string]*GraphEntity)
	for name, node := range nodeMap {
		entities[name] = graphEntityFromNode(node)
	}

	return entities, nil
}

// GraphRelationships retrieves multiple relationships between entity pairs from the Neo4j database.
// It returns a map where the key is "sourceEntity-targetEntity" and the value is the GraphRelationship.
func (n *Neo4jGraphStorage) GraphRelationships(ctx context.Context, pairs [][2]string, opts Options) (map[string]*GraphRelationship, error) {
	if len(pairs) == 0 {
		return map[string]*GraphRelationship{}, nil
	}

	// Prepare parameters for the query
	sources := make([]string, len(pairs))
	targets := make([]string, len(pairs))
	for i, pair := range pairs {
		sources[i] = pair[0]
		targets[i] = pair[1]
	}

	res, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			labels := n.getDocLabels(opts)
			query := fmt.Sprintf(`
UNWIND $pairs AS pair
MATCH (start%s {name: pair[0]})-[r]-(end%s{name: pair[1]})
RETURN pair[0] as source, pair[1] as target, properties(r) as edge_properties
			`, labels, labels)

			// Convert pairs to a format suitable for the query
			pairsParam := make([][]string, len(pairs))
			for i, pair := range pairs {
				pairsParam[i] = []string{pair[0], pair[1]}
			}

			queryRes, err := tx.Run(ctx, query, map[string]any{
				"pairs": pairsParam,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}

			result := make(map[string]map[string]any)
			for record, err := range queryRes.Records(ctx) {
				if err != nil {
					return nil, fmt.Errorf("failed to get result: %w", err)
				}

				source, sourceOK := record.Get("source")
				target, targetOK := record.Get("target")
				edgeProps, propsOK := record.Get("edge_properties")

				if !sourceOK || !targetOK || !propsOK {
					continue
				}

				sourceStr, sourceOK := source.(string)
				targetStr, targetOK := target.(string)
				props, propsOK := edgeProps.(map[string]any)

				if !sourceOK || !targetOK || !propsOK {
					continue
				}

				key := fmt.Sprintf("%s-%s", sourceStr, targetStr)
				result[key] = props
			}

			return result, nil
		})
	})
	if err != nil {
		return nil, err
	}

	propsMap, ok := res.(map[string]map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid result type, got %T, want map[string]map[string]any", res)
	}

	relationships := make(map[string]*GraphRelationship)
	for key, props := range propsMap {
		parts := strings.Split(key, "-")
		if len(parts) != 2 {
			continue
		}

		rel := graphRelationshipFromEdge(parts[0], parts[1], props)
		relationships[key] = rel
	}

	return relationships, nil
}

// GraphCountEntitiesRelationships counts the number of relationships for multiple entities.
// It returns a map of entity names to their relationship counts.
func (n *Neo4jGraphStorage) GraphCountEntitiesRelationships(ctx context.Context, names []string, opts Options) (map[string]int, error) {
	if len(names) == 0 {
		return map[string]int{}, nil
	}

	res, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			labels := n.getDocLabels(opts)
			query := fmt.Sprintf(`
MATCH (n%s)
WHERE n.name IN $entity_ids
OPTIONAL MATCH (n)-[r]-()
RETURN n.name AS entity_id, COUNT(r) AS degree
            `, labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{
				"entity_ids": names,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}

			result := make(map[string]int64)
			for record, err := range queryRes.Records(ctx) {
				if err != nil {
					return nil, fmt.Errorf("failed to get result: %w", err)
				}

				entityID, idOK := record.Get("entity_id")
				degree, degreeOK := record.Get("degree")

				if !idOK || !degreeOK {
					continue
				}

				entityIDStr, idOK := entityID.(string)
				degreeCnt, degreeOK := degree.(int64)

				if !idOK || !degreeOK {
					continue
				}

				result[entityIDStr] = degreeCnt
			}

			return result, nil
		})
	})
	if err != nil {
		return nil, err
	}

	countMap, ok := res.(map[string]int64)
	if !ok {
		return nil, fmt.Errorf("invalid result type, got %T, want map[string]int64", res)
	}

	// Convert int64 to int
	counts := make(map[string]int)
	for name, count := range countMap {
		counts[name] = int(count)
	}

	return counts, nil
}

// GraphRelatedEntities retrieves all entities related to multiple input entities.
// It returns a map of entity names to slices of related GraphEntity objects.
func (n *Neo4jGraphStorage) GraphRelatedEntities(ctx context.Context, names []string, opts Options) (map[string][]*GraphEntity, error) {
	if len(names) == 0 {
		return map[string][]*GraphEntity{}, nil
	}

	res, err := n.session(func(ctx context.Context, sess neo4j.SessionWithContext) (any, error) {
		return sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			labels := n.getDocLabels(opts)
			query := fmt.Sprintf(`
MATCH (n%s)
WHERE n.name IN $entity_ids
OPTIONAL MATCH (n)-[r]-(connected%s)
WHERE connected.entity_id IS NOT NULL
RETURN n.name as source_id, collect(connected) as connected_nodes
            `, labels, labels)
			queryRes, err := tx.Run(ctx, query, map[string]any{
				"entity_ids": names,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to run query: %w", err)
			}

			result := make(map[string][]dbtype.Node)
			for record, err := range queryRes.Records(ctx) {
				if err != nil {
					return nil, fmt.Errorf("failed to get result: %w", err)
				}

				sourceID, sourceOK := record.Get("source_id")
				connectedNodes, connectedOK := record.Get("connected_nodes")

				if !sourceOK || !connectedOK {
					continue
				}

				sourceIDStr, sourceOK := sourceID.(string)
				nodes, connectedOK := connectedNodes.([]any)

				if !sourceOK || !connectedOK {
					continue
				}

				nodeList := make([]dbtype.Node, 0, len(nodes))
				for _, node := range nodes {
					if dbNode, ok := node.(dbtype.Node); ok {
						nodeList = append(nodeList, dbNode)
					}
				}

				result[sourceIDStr] = nodeList
			}

			return result, nil
		})
	})
	if err != nil {
		return nil, err
	}

	nodesMap, ok := res.(map[string][]dbtype.Node)
	if !ok {
		return nil, fmt.Errorf("invalid result type, got %T, want map[string][]dbtype.Node", res)
	}

	relatedEntities := make(map[string][]*GraphEntity, len(nodesMap))
	for name, nodes := range nodesMap {
		entities := make([]*GraphEntity, 0, len(nodes))
		for _, node := range nodes {
			entities = append(entities, graphEntityFromNode(node))
		}
		relatedEntities[name] = entities
	}

	return relatedEntities, nil
}

// Close terminates the connection to the Neo4j database.
// It returns any error encountered during the closing operation.
func (n *Neo4jGraphStorage) Close(ctx context.Context) error {
	return n.client.Close(ctx)
}

func (n *Neo4jGraphStorage) session(sessFunc func(context.Context, neo4j.SessionWithContext) (any, error)) (any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	sess := n.client.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second*30)
		defer closeCancel()
		_ = sess.Close(closeCtx)
	}()

	trxCtx, trxCancel := context.WithTimeout(context.Background(), time.Second*30)
	defer trxCancel()

	return sessFunc(trxCtx, sess)
}

func graphEntityFromNode(node dbtype.Node) *GraphEntity {
	name, ok := node.Props["name"].(string)
	if !ok {
		name = ""
	}
	typ, ok := node.Props["type"].(string)
	if !ok {
		typ = ""
	}
	desc, ok := node.Props["description"].(string)
	if !ok {
		desc = ""
	}
	sourceIDs, ok := node.Props["source_ids"].(string)
	if !ok {
		sourceIDs = ""
	}
	createdAtStr, ok := node.Props["created_at"].(string)
	if !ok {
		createdAtStr = time.Now().Format(time.RFC3339)
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		createdAt = time.Now()
	}

	return &GraphEntity{
		Name:         name,
		Type:         typ,
		Descriptions: desc,
		SourceIDs:    sourceIDs,
		CreatedAt:    createdAt,
	}
}

func graphRelationshipFromEdge(source, target string, props map[string]any) *GraphRelationship {
	weight, ok := props["weight"].(float64)
	if !ok {
		weight = 1.0
	}
	description, ok := props["description"].(string)
	if !ok {
		description = ""
	}
	keywords, ok := props["keywords"].(string)
	if !ok {
		keywords = ""
	}
	arrKeywords := strings.Split(keywords, GraphFieldSeparator)
	sourceIDs, ok := props["source_ids"].(string)
	if !ok {
		sourceIDs = ""
	}
	createdAtStr, ok := props["created_at"].(string)
	if !ok {
		createdAtStr = time.Now().Format(time.RFC3339)
	}
	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		createdAt = time.Now()
	}

	return &GraphRelationship{
		Source:       source,
		Target:       target,
		Weight:       weight,
		Descriptions: description,
		Keywords:     arrKeywords,
		SourceIDs:    sourceIDs,
		CreatedAt:    createdAt,
	}
}
