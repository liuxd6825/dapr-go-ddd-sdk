package storage

import (
	"context"

	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
)

type InsertDocData struct {
	ChunkId  int
	TenantId string
	CaseId   string
	FileName string
	DocId    string
	Content  []string
}

type VectorUpsertEntity struct {
	Name     string
	TenantId string
	CaseId   string
	FileName string
	DocId    string
	Content  []string
}

type VectorUpsertRelationship struct {
	TenantId string
	CaseId   string
	Source   string
	Target   string
	Content  []string
	DocId    string
	FileName string
}

type NodeLabel string

const (
	NodeLabel_All    NodeLabel = "all"
	NodeLabel_Master NodeLabel = "master"
	NodeLabel_Draw   NodeLabel = "draw"
	NodeLabel_Doc    NodeLabel = "doc"
)

type Options struct {
	TenantId  string
	CaseId    string
	DocId     string
	NodeLabel NodeLabel
}

type GraphQueryParam struct {
	Keys    []string
	MaxDeep int
	Limit   int
}

type BaseStorage interface {
	CreateTenant(ctx context.Context, tenantId string) error
	CreateCase(ctx context.Context, tenantId, caseId string) error
	DeleteTenant(ctx context.Context, tenantId string) error
	DeleteCase(ctx context.Context, tenantId, caseId string) error
	DeleteDoc(ctx context.Context, tenantId, caseId, docId string) error
}

// GraphStorage defines the interface for graph database operations.
// It provides methods to query and manipulate entities and relationships
// in a knowledge graph.
type GraphStorage interface {
	BaseStorage
	// GraphEntity retrieves a single entity by name from the graph storage.
	// Returns ErrEntityNotFound if the entity doesn't exist.
	GraphEntity(ctx context.Context, name string, opts Options) (*GraphEntity, error)

	// GraphRelationship retrieves a relationship between sourceEntity and targetEntity.
	// Returns ErrRelationshipNotFound if the relationship doesn't exist.
	GraphRelationship(ctx context.Context, sourceEntity, targetEntity string, opts Options) (*GraphRelationship, error)

	// GraphUpsertEntity creates a new entity or updates an existing entity in the graph storage.
	// If the entity already exists, it should merge the new data with existing data.
	GraphUpsertEntity(ctx context.Context, entity *GraphEntity, opts Options) error

	// GraphUpsertRelationship creates a new relationship or updates an existing relationship
	// between two entities in the graph storage.
	// If the relationship already exists, it should merge the new data with existing data.
	GraphUpsertRelationship(ctx context.Context, relationship *GraphRelationship, opts Options) error

	// GraphEntities batch retrieves multiple entities by their names.
	// Returns a map with entity names as keys and entity objects as values.
	// If an entity doesn't exist, it should be omitted from the result map.
	GraphEntities(ctx context.Context, names []string, opts Options) (map[string]*GraphEntity, error)

	// GraphRelationships batch retrieves multiple relationships by their source-target pairs.
	// Returns a map with composite keys (formatted as "source-target") as keys and
	// relationship objects as values.
	// If a relationship doesn't exist, it should be omitted from the result map.
	GraphRelationships(ctx context.Context, pairs [][2]string, opts Options) (map[string]*GraphRelationship, error)

	// GraphCountEntitiesRelationships counts the number of relationships each entity has.
	// Returns a map with entity names as keys and relationship counts as values.
	// This is used to determine entity importance during queries.
	GraphCountEntitiesRelationships(ctx context.Context, names []string, opts Options) (map[string]int, error)

	// GraphRelatedEntities finds entities directly connected to the specified entities.
	// Returns a map with entity names as keys and slices of directly connected entities as values.
	// Used to expand the context during queries.
	GraphRelatedEntities(ctx context.Context, names []string, opts Options) (map[string][]*GraphEntity, error)

	GraphSaveDoc(ctx context.Context, doc *entity.Document, entries []*GraphEntity, rels []*GraphRelationship) error

	GraphQuery(ctx context.Context, param GraphQueryParam, opts Options) ([]string, error)
}

// VectorStorage defines the interface for vector database operations.
// It provides methods to query and store entities and relationships
// in a vector space for semantic search capabilities.
type VectorStorage interface {
	BaseStorage
	// VectorSearchDoc performs a semantic search for entities based on the provided keywords.
	// Returns a slice of entity names that semantically match the keywords.
	// The results should be ordered by relevance.
	VectorSearchDoc(ctx context.Context, vector []float32, topK int, opts Options) ([]string, error)
	// VectorQueryEntity performs a semantic search for entities based on the provided keywords.
	// Returns a slice of entity names that semantically match the keywords.
	// The results should be ordered by relevance.
	VectorQueryEntity(ctx context.Context, keywords string, opts Options) ([]string, error)
	// VectorQueryRelationship performs a semantic search for relationships based on the provided keywords.
	// Returns a slice of source-target entity name pairs that semantically match the keywords.
	// The results should be ordered by relevance.
	VectorQueryRelationship(ctx context.Context, keywords string, opts Options) ([][2]string, error)

	// VectorUpsertEntity creates or updates the vector representation of an entity.
	// The content parameter should contain the text used for semantic matching.
	// This typically includes the entity name and description.
	VectorUpsertEntity(ctx context.Context, entity *VectorUpsertEntity) error

	// VectorUpsertRelationship creates or updates the vector representation of a relationship.
	// The content parameter should contain the text used for semantic matching.
	// This typically includes keywords, descriptions, and entity names.
	VectorUpsertRelationship(ctx context.Context, entity *VectorUpsertRelationship) error

	// VectorInsertDoc
	VectorInsertDoc(ctx context.Context, data InsertDocData, opts Options) error

	VectorFlush(ctx context.Context, opts Options) error
	LoadTenant(ctx context.Context, tenantId string) error
}

// KeyValueStorage defines the interface for key-value storage operations.
// It provides methods to access and store source documents.
type KeyValueStorage interface {
	BaseStorage
	// KVSource retrieves a source document chunk by its ID.
	// Returns an error if the source doesn't exist or can't be retrieved.
	KVSource(ctx context.Context, id string) (Source, error)
	// KVUpsertSources creates or updates multiple source document chunks at once.
	// Each source should be stored with its ID as the key.
	// This is called during document processing to store chunked documents.
	KVUpsertSources(ctx context.Context, sources []Source) error
}

type EmbedderStore interface {
	EmbedTexts(ctx context.Context, strList []string) ([][]float32, error)
}

// Storage is a composite interface that combines GraphStorage,
// VectorStorage, and KeyValueStorage interfaces to provide
// comprehensive data storage capabilities.
type Storage interface {
	BaseStorage
	GraphStorage
	VectorStorage
	KeyValueStorage
	EmbedderStore
}
type storage struct {
	GraphStorage
	VectorStorage
	KeyValueStorage
	EmbedderStore
}

type NewStorageOptions struct {
	Graph    GraphStorage
	Vector   VectorStorage
	KeyValue KeyValueStorage
	Embedder EmbedderStore
}

func NewStorage(graph GraphStorage, vector VectorStorage, keyValue KeyValueStorage, embedder EmbedderStore) Storage {
	if graph == nil {
		panic("graph cannot be nil")
	}
	if vector == nil {
		panic("vector cannot be nil")
	}
	if embedder == nil {
		panic("embedder cannot be nil")
	}
	if keyValue == nil {
		panic("keyValue cannot be nil")
	}
	return &storage{
		GraphStorage:    graph,
		VectorStorage:   vector,
		KeyValueStorage: keyValue,
		EmbedderStore:   embedder,
	}
}

func (s *storage) CreateTenant(ctx context.Context, tenantId string) error {
	if err := s.GraphStorage.CreateTenant(ctx, tenantId); err != nil {
		return err
	}
	if err := s.KeyValueStorage.CreateTenant(ctx, tenantId); err != nil {
		return err
	}
	if err := s.VectorStorage.CreateTenant(ctx, tenantId); err != nil {
		return err
	}
	return nil
}

func (s *storage) CreateCase(ctx context.Context, tenantId string, caseId string) error {
	if err := s.GraphStorage.CreateCase(ctx, tenantId, caseId); err != nil {
		return err
	}
	if err := s.KeyValueStorage.CreateCase(ctx, tenantId, caseId); err != nil {
		return err
	}
	if err := s.VectorStorage.CreateCase(ctx, tenantId, caseId); err != nil {
		return err
	}
	return nil
}

func (s *storage) DeleteTenant(ctx context.Context, tenantId string) error {
	if err := s.GraphStorage.DeleteTenant(ctx, tenantId); err != nil {
		return err
	}
	if err := s.KeyValueStorage.DeleteTenant(ctx, tenantId); err != nil {
		return err
	}
	if err := s.VectorStorage.DeleteTenant(ctx, tenantId); err != nil {
		return err
	}
	return nil
}

func (s *storage) DeleteCase(ctx context.Context, tenantId, caseId string) error {
	if err := s.GraphStorage.DeleteCase(ctx, tenantId, caseId); err != nil {
		return err
	}
	if err := s.KeyValueStorage.DeleteCase(ctx, tenantId, caseId); err != nil {
		return err
	}
	if err := s.VectorStorage.DeleteCase(ctx, tenantId, caseId); err != nil {
		return err
	}
	return nil
}

func (s *storage) DeleteDoc(ctx context.Context, tenantId, caseId, docId string) error {
	if err := s.GraphStorage.DeleteDoc(ctx, tenantId, caseId, docId); err != nil {
		return err
	}
	if err := s.KeyValueStorage.DeleteDoc(ctx, tenantId, caseId, docId); err != nil {
		return err
	}
	if err := s.VectorStorage.DeleteDoc(ctx, tenantId, caseId, docId); err != nil {
		return err
	}
	return nil
}
