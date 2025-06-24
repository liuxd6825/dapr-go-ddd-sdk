package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/entity"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/llm"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"sort"
	"strings"
	"time"
)

type DocEventType string

type DocEvent func(ctx context.Context, doc *entity.Document, eventType DocEventType, eventData interface{}) error

type DocOptions struct {
	OnStartInsert func(ctx context.Context, doc *entity.Document)
	OnDoneInsert  func(ctx context.Context, doc *entity.Document, err error)
	OnChunkCount  func(ctx context.Context, doc *entity.Document, chunkCount int)

	OnStartExtractEntities func(ctx context.Context, doc *entity.Document, source *Source)
	OnDoneExtractEntities  func(ctx context.Context, doc *entity.Document, source *Source, err error)
}

type DocHandle struct {
	config  Config
	storage Storage
	llm     llm.LLM
	logger  *logrus.Logger
	opts    *DocOptions
}

func NewDocHandle(config Config, storage Storage, llm llm.LLM, logger *logrus.Logger, options ...func(opts *DocOptions)) *DocHandle {
	opts := &DocOptions{}
	for _, option := range options {
		option(opts)
	}

	return &DocHandle{
		config:  config,
		storage: storage,
		llm:     llm,
		logger:  logger,
		opts:    opts,
	}
}

// SaveGraph processes a document and stores it in the provided storage.
// It chunks the document content, extracts entities and relationships using the provided
// document handler, and stores the results in the appropriate storage.
// It returns an error if any step in the process fails.
func (d *DocHandle) SaveGraph(ctx context.Context, doc *entity.Document) (err error) {
	if d.opts.OnStartInsert != nil {
		d.opts.OnStartInsert(ctx, doc)
	}

	if d.opts.OnDoneInsert != nil {
		defer func() {
			err = errors.GetRecoverError(err, recover())
			d.opts.OnDoneInsert(ctx, doc, err)
		}()
	}

	content := CleanContent(doc.Text)
	logger := d.logger.WithFields(logrus.Fields{"package": "my_rag.storage", "function": "Insert"})

	chunks, err := d.config.GetChunksDocument(doc.TenantId, doc.CaseId, doc.Id, content)
	if err != nil {
		return fmt.Errorf("failed to chunk string: %w", err)
	}

	// The chunks returned from the ChunksDocument doesn't have an ID, generate one here
	// based on the document ID and the order of the chunks. This ID would be used to retrieve
	// the chunk in the Query function.
	chunksWithID := make([]Source, len(chunks))
	for i, chunk := range chunks {
		id := chunk.GenID(doc.Id)
		chunksWithID[i] = Source{
			Id:         id,
			DocId:      doc.Id,
			TenantId:   doc.TenantId,
			CaseId:     doc.CaseId,
			Content:    chunk.Content,
			TokenSize:  chunk.TokenSize,
			OrderIndex: chunk.OrderIndex,
		}
	}
	logger.Info("Upserting sources", "count", len(chunks))

	if d.opts.OnChunkCount != nil {
		d.opts.OnChunkCount(ctx, doc, len(chunks))
	}

	if err := d.storage.KVUpsertSources(ctx, chunksWithID); err != nil {
		return fmt.Errorf("failed to upsert sources kv: %w", err)
	}

	llmConcurrencyCount := d.config.GetConcurrencyCount()
	if llmConcurrencyCount == 0 {
		llmConcurrencyCount = 1
	}

	err = d.ExtractEntities(ctx, doc, chunks)

	if err != nil {
		return fmt.Errorf("failed to extract entities: %w", err)
	}

	return nil
}

func (d *DocHandle) ExtractEntities(
	ctx context.Context,
	doc *entity.Document,
	sources []Source,
) error {
	// Sort sources by order index to maintain document flow
	orderedSources := make([]Source, len(sources))
	copy(orderedSources, sources)
	sort.Slice(orderedSources, func(i, j int) bool {
		return orderedSources[i].OrderIndex < orderedSources[j].OrderIndex
	})

	d.logger.Info("Extracting entities ", " count ", len(orderedSources))
	extractPromptData := d.config.GetEntityExtractionPromptData()

	eg := new(errgroup.Group)
	// Semaphore to limit concurrent LLM calls
	sem := make(chan struct{}, d.config.GetConcurrencyCount())
	for _, source := range orderedSources {
		eg.Go(func() error {
			// Acquire semaphore before making LLM call
			sem <- struct{}{}
			var err error
			logger := d.logger
			defer func() {
				if d.opts.OnDoneExtractEntities != nil {
					err = errors.GetRecoverError(err, recover())
					d.opts.OnDoneExtractEntities(ctx, doc, &source, err)
				}
				<-sem
			}()

			orderIndex := source.OrderIndex
			logger.Info("call LLM extract orderIndex:", orderIndex)
			// Extract entities and relationships for this source chunk
			entities, relationships, err := d.llmExtractEntities(ctx, doc, source.Content)
			if err != nil {
				return fmt.Errorf("failed to extract entities with LLM: %w", err)
			}

			logger.Info("Done call LLM orderIndex:", orderIndex, " entities:", len(entities), " relationships:", len(relationships))

			if err != nil {
				return fmt.Errorf("failed to merge entities with LLM: %w", err)
			}

			newEntities := []*GraphEntity{}
			// Process each entity group by name
			for name, unmergedEntities := range entities {
				newEntity, err := d.mergeGraphEntities(ctx, name, doc, &source, extractPromptData.Language, unmergedEntities)
				if err == nil {
					newEntities = append(newEntities, newEntity)
				}
			}

			newRelationships := []*GraphRelationship{}
			// Process each relationship group by source-target pair
			for key, unmergedRelationships := range relationships {
				if rel, err := d.mergeGraphRelationships(ctx, doc, key, source.GenID(doc.Id), extractPromptData.Language,
					unmergedRelationships); err != nil {
					return fmt.Errorf("failed to process graph relationship: %w", err)
				} else if rel != nil {
					newRelationships = append(newRelationships, rel)
				}
			}

			if err := d.storage.GraphSaveDoc(ctx, doc.TenantId, doc.CaseId, doc.Id, newEntities, newRelationships); err != nil {
				return err
			}

			logger.Info("Processed source ", "orderIndex:", orderIndex)

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

// LlmExtractEntities 导入实体与关系
func (d *DocHandle) llmExtractEntities(
	ctx context.Context,
	doc *entity.Document,
	content string,
) (map[string][]*GraphEntity, map[string][]*GraphRelationship, error) {
	logger := d.logger

	promptData := d.config.GetEntityExtractionPromptData()
	promptData.Input = content

	extractPrompt, err := PromptTemplate("extract-entities", extractEntitiesPrompt, promptData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate extract entities prompt: %w", err)
	}

	gleanPrompt, err := PromptTemplate("glean-entities", gleanEntitiesPrompt, promptData)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate glean entities prompt: %w", err)
	}

	logger.Debug("Use LLM to extract entities from source",
		"extractPrompt", extractPrompt, "gleanPrompt", gleanPrompt, "source", content)

	var results LLMResult

	retry := 0
	backoffDuration := d.config.GetBackoffDuration()
	maxRetries := d.config.GetMaxRetries()
	maxGleanCount := d.config.GetGleanCount()

	for {
		// If this is not a first retry, add backoff delay.
		if retry > 0 {
			time.Sleep(backoffDuration)
		}
		// LLM sometimes returns incorrect format, retry up to maxRetries() times.
		if retry >= maxRetries {
			return nil, nil, fmt.Errorf("failed to extract entities after %d retries", maxRetries)
		}

		logger.Debug("Use LLM to extract entities from source", "extractPrompt", extractPrompt)

		// Initial extraction conversation
		histories := []string{extractPrompt}
		sourceResult, err := d.llm.Generate(ctx, newMessages(histories))
		if err != nil {
			nErr := fmt.Errorf("failed to call LLM: %w", err)
			retry++
			logger.Warn("Retry extract", "retry", retry, "error", nErr)
			continue
		}

		// Parse initial extraction results
		var sourceParsed LLMResult
		jsonData := GetJsonString(sourceResult)
		err = json.Unmarshal([]byte(jsonData), &sourceParsed)
		if err != nil {
			nErr := fmt.Errorf("failed to parse llm result: %w", err)
			retry++
			logger.Warn("Retry parse result ", "retry ", retry, "error ", nErr)
			continue
		}
		results.Entities = append(results.Entities, sourceParsed.Entities...)
		results.Relationships = append(results.Relationships, sourceParsed.Relationships...)

		initGraphData(doc, results.Entities, results.Relationships)

		histories = append(histories, sourceResult.Content)

		// "Gleaning" process: attempt to extract additional entities that might have been missed
		gleanCount := 0
		for {
			if gleanCount >= maxGleanCount {
				break
			}

			logger.Debug("Use LLM to glean entities from source", "gleanPrompt", gleanPrompt)
			histories = append(histories, gleanPrompt)
			gleanResult, err := d.llm.Generate(ctx, newMessages(histories))
			if err != nil {
				nErr := fmt.Errorf("failed to call LLM on glean: %w", err)
				retry++
				logger.Warn("Retry glean", "retry", retry, "error", nErr)
				continue
			}

			histories = append(histories, gleanResult.Content)

			var gleanParsed LLMResult
			jsonData := GetJsonString(sourceResult)
			err = json.Unmarshal([]byte(jsonData), &gleanParsed)
			if err != nil {
				nErr := fmt.Errorf("failed to parse llm result: %w", err)
				retry++
				logger.Warn("Retry parse result", "retry", retry, "error", nErr)
				continue
			}

			results.Entities = append(results.Entities, gleanParsed.Entities...)
			results.Relationships = append(results.Relationships, gleanParsed.Relationships...)
			initGraphData(doc, results.Entities, results.Relationships)
			gleanCount++

			// Ask LLM if we should continue gleaning more entities
			decideMessages := make([]string, 0)
			decideMessages = append(decideMessages, histories...)
			decideMessages = append(decideMessages, gleanDecideContinuePrompt)

			decideResult, err := d.llm.Generate(ctx, newMessages(decideMessages))
			if err != nil {
				nErr := fmt.Errorf("failed to call LLM on decide: %w", err)
				retry++
				logger.Warn("Retry decide", "retry", retry, "error", nErr)
				continue
			}

			decideResultContent := strings.ToLower(strings.TrimSpace(strings.Trim(strings.Trim(decideResult.Content, `"`), `'`)))

			logger.Debug("Decide result from LLM", "decideResult", decideResultContent)

			// Only continue gleaning if the LLM explicitly says "yes"
			if decideResultContent != "yes" {
				break
			}
		}

		// Organize entities by name and relationships by source-target pair
		entities, relationships := DedupeLLMResult(ctx, results.Entities, results.Relationships, promptData.EntityTypes)

		if logger.IsLevelEnabled(logrus.DebugLevel) {
			for _, entity := range results.Entities {
				logger.Debug("实体", entity.Name)
			}

			for _, rel := range results.Relationships {
				logger.Debug(fmt.Errorf("关系：%s - %s", rel.Target, rel.Source))
			}
		}

		return entities, relationships, nil
	}
}

func (d *DocHandle) mergeGraphEntities(
	ctx context.Context,
	name string,
	doc *entity.Document,
	source *Source,
	language string,
	unmergedEntities []*GraphEntity,
) (*GraphEntity, error) {
	llm := d.llm
	logger := d.logger
	summariesMaxToken := d.config.GetMaxSummariesTokenLength()

	// Collect data from existing entity (if found) to merge with new data
	existingTypes := make([]string, 0)
	existingSourceIDs := make([]string, 0)
	existingDescriptions := make([]string, 0)
	// Merge data from new entities
	for _, entity := range unmergedEntities {
		existingTypes = append(existingTypes, entity.Type)
		existingDescriptions = AppendIfUnique(existingDescriptions, entity.Descriptions)
	}
	existingSourceIDs = AppendIfUnique(existingSourceIDs, source.Id)

	// Choose the most frequent entity type from all type mentions
	entityType := MostFrequentItem(existingTypes)
	sourceIDs := strings.Join(existingSourceIDs, GraphFieldSeparator)

	// Summarize descriptions if they exceed token limit
	description, err := DescriptionsSummary(name, language, summariesMaxToken, existingDescriptions, llm)
	if err != nil {
		return nil, fmt.Errorf("failed to summarize descriptions: %w", err)
	}

	ent := &GraphEntity{
		Id:           name,
		TenantId:     doc.TenantId,
		CaseId:       doc.CaseId,
		DocId:        doc.Id,
		Name:         name,
		Type:         entityType,
		Descriptions: description,
		SourceIDs:    sourceIDs,
		CreatedAt:    time.Now(),
	}

	logger.Debug("Upserting graph entity", "entity", ent)

	/*
		// Update both graph and vector storage for entity
		if err := storage.GraphUpsertEntity(ctx, ent, opts); err != nil {
			return nil, fmt.Errorf("failed to upsert graph entity in graph storage: %w", err)
		}
		vectorUpsertEntity := &VectorUpsertEntity{
			Name:     ent.Name,
			TenantId: doc.TenantId,
			CaseId:   doc.CaseId,
			FileName: doc.FileName,
			DocId:    doc.Id,
			Content:  []string{ent.Name + ":" + ent.Descriptions},
		}
		if err := storage.VectorUpsertEntity(ctx, vectorUpsertEntity); err != nil {
			return nil, fmt.Errorf("failed to upsert entity in vector storage: %w", err)
		}
	*/
	return ent, nil
}

func (d *DocHandle) mergeGraphRelationships(
	ctx context.Context,
	doc *entity.Document,
	key, sourceID, language string,
	relationships []*GraphRelationship,
) (*GraphRelationship, error) {
	summariesMaxToken := d.config.GetMaxSummariesTokenLength()
	storage := d.storage
	logger := d.logger

	// Track existing relationship properties to merge with new data
	existingWeight := 0.0
	existingDescriptions := make([]string, 0)
	existingKeywords := make([]string, 0)
	existingSourceIDs := make([]string, 0)

	// Parse composite key format "SOURCE-TARGET" into separate entity names
	arrKey := strings.Split(key, "-")
	sourceEntity := arrKey[0]
	targetEntity := arrKey[1]
	opts := newOptionsWithDoc(doc)
	// Retrieve existing relationship data from storage if it exists
	existingRelationship, err := storage.GraphRelationship(ctx, sourceEntity, targetEntity, opts)
	if err != nil {
		if !errors.Is(err, ErrRelationshipNotFound) {
			return nil, fmt.Errorf("failed to get relationship: %w", err)
		}
		// If relationship not found, continue with empty existing data
	} else if existingRelationship != nil {
		// Accumulate existing weight (weights are additive for relationship strength)
		existingWeight += existingRelationship.Weight

		// Extract and parse data from existing relationship
		arrDescriptions := strings.Split(existingRelationship.Descriptions, GraphFieldSeparator)
		existingDescriptions = append(existingDescriptions, arrDescriptions...)

		existingKeywords = append(existingKeywords, existingRelationship.Keywords...)

		arrSourceIDs := strings.Split(existingRelationship.SourceIDs, GraphFieldSeparator)
		existingSourceIDs = append(existingSourceIDs, arrSourceIDs...)
	}

	// Merge new relationship data with existing data
	for _, relationship := range relationships {
		existingWeight += relationship.Weight
		existingDescriptions = AppendIfUnique(existingDescriptions, relationship.Descriptions)
		for _, keyword := range relationship.Keywords {
			existingKeywords = AppendIfUnique(existingKeywords, keyword)
		}
	}
	existingSourceIDs = AppendIfUnique(existingSourceIDs, sourceID)

	// Summarize all descriptions if they exceed token limit
	description, err := d.descriptionsSummary(key, language, summariesMaxToken, existingDescriptions)
	if err != nil {
		return nil, fmt.Errorf("failed to summarize descriptions: %w", err)
	}
	sourceIDs := strings.Join(existingSourceIDs, GraphFieldSeparator)

	// Create source entity if it doesn't exist
	// This ensures relationship integrity by avoiding dangling references
	_, err = storage.GraphEntity(ctx, sourceEntity, opts)
	if err != nil {
		if !errors.Is(err, ErrEntityNotFound) {
			return nil, fmt.Errorf("failed to get source entity with name %s: %w", sourceEntity, err)
		}
		logger.Debug("Entity not found, upserting", "entity", sourceEntity)

		// Create a minimal placeholder entity with UNKNOWN type
		if err := storage.GraphUpsertEntity(ctx, &GraphEntity{
			Name:         sourceEntity,
			Type:         "UNKNOWN",
			Descriptions: description,
			SourceIDs:    sourceID,
			CreatedAt:    time.Now(),
		}, opts); err != nil {
			return nil, fmt.Errorf("failed to upsert source node with name %s: %w", sourceEntity, err)
		}
	}

	// Create target entity if it doesn't exist
	// Similar to source entity creation for relationship integrity
	_, err = storage.GraphEntity(ctx, targetEntity, opts)
	if err != nil {
		if !errors.Is(err, ErrEntityNotFound) {
			return nil, fmt.Errorf("failed to get target entity with name %s: %w", targetEntity, err)
		}
		logger.Debug("Entity not found, upserting", "entity", targetEntity)
		if err := storage.GraphUpsertEntity(ctx, &GraphEntity{
			Name:         targetEntity,
			Type:         "UNKNOWN",
			Descriptions: description,
			SourceIDs:    sourceID,
			CreatedAt:    time.Now(),
		}, opts); err != nil {
			return nil, fmt.Errorf("failed to upsert target node with name %s: %w", targetEntity, err)
		}
	}

	// Create a combined content string for vector storage
	// This enables semantic search over relationships
	keywords := strings.Join(existingKeywords, GraphFieldSeparator)
	content := keywords + sourceEntity + targetEntity + description
	id := fmt.Sprintf("doc_%s_%s_%s", opts.DocId, sourceEntity, targetEntity)
	// Create final relationship with merged data
	rel := &GraphRelationship{
		Id:           id,
		TenantId:     opts.TenantId,
		CaseId:       opts.CaseId,
		DocId:        opts.DocId,
		Source:       sourceEntity,
		Target:       targetEntity,
		Weight:       existingWeight,
		Descriptions: content,
		Keywords:     existingKeywords,
		SourceIDs:    sourceIDs,
		CreatedAt:    time.Now(),
	}

	/*
		// Update both graph and vector storage for the relationship
		if err := storage.GraphUpsertRelationship(ctx, rel, opts); err != nil {
			return nil,  fmt.Errorf("failed to upsert graph relationship: %w", err)
		}
	*/
	/*
		vRel := &VectorUpsertRelationship{
			TenantId: opts.TenantId,
			CaseId:   opts.CaseId,
			DocId:    opts.DocId,
			Source:   sourceEntity,
			Target:   targetEntity,
			Content:  []string{content},
			FileName: doc.FileName,
		}

			if err := storage.VectorUpsertRelationship(ctx, vRel); err != nil {
				return fmt.Errorf("failed to upsert relationship vector: %w", err)
			}
	*/
	return rel, nil
}

func (d *DocHandle) descriptionsSummary(name, language string, maxToken int, descriptions []string) (string, error) {
	// Join all descriptions with separator
	joinedDescriptions := strings.Join(descriptions, GraphFieldSeparator)

	// Check if the joined descriptions exceed the token limit
	tokens, err := EncodeStringByTiktoken(joinedDescriptions)
	if err != nil {
		return "", fmt.Errorf("failed to encode string: %w", err)
	}

	// If descriptions are under token limit, no need to summarize
	if len(tokens) < maxToken {
		return joinedDescriptions, nil
	}

	// Format descriptions for LLM summarization
	descString := strings.Join(descriptions, ", ")
	descString = "[" + descString + "]"

	// Generate summary prompt and get LLM to create a condensed description
	summarizePrompt, err := PromptTemplate("summarize-descriptions", summarizeDescriptionsPrompt,
		summarizeDescriptionsPromptData{
			EntityName:   name,
			Descriptions: descString,
			Language:     language,
		})
	if err != nil {
		return "", fmt.Errorf("failed to generate summarize descriptions prompt: %w", err)
	}
	ctx := context.Background()
	msg, err := d.llm.Generate(ctx, newMessages([]string{summarizePrompt}))
	if err != nil {
		return "", fmt.Errorf("failed to generate summarize descriptions prompt: %w", err)
	}
	return msg.Content, nil
}
