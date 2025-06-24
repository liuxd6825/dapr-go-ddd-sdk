package storage

import "time"

// Config provides an interface for processing documents and interacting with language models.
type Config interface {
	// GetChunksDocument splits a document's content into smaller, manageable chunks.
	// It returns a slice of Source objects representing the document chunks,
	// without assigning IDs (IDs will be generated in the Insert function).
	GetChunksDocument(tenantId, caseId, docId, content string) ([]Source, error)
	// GetEntityExtractionPromptData returns the data needed to generate prompts for extracting
	// entities and relationships from text content.
	// The implementation doesn't need to fill the Input field, as it will be filled in the
	// Insert function.
	GetEntityExtractionPromptData() EntityExtractionPromptData
	// GetMaxRetries determines the maximum number of retries allowed for the Chat function.
	// This is especially used when extracting entities and relationships from text content,
	// due to the incorrect format that sometimes LLM returns.
	GetMaxRetries() int
	// GetConcurrencyCount determines the number of concurrent requests to the LLM.
	GetConcurrencyCount() int
	// GetBackoffDuration determines the backoff duration between retries.
	GetBackoffDuration() time.Duration
	// GetGleanCount returns the maximum number of additional extraction attempts
	// to perform after the initial entity extraction to find entities that might
	// have been missed.
	GetGleanCount() int
	// GetMaxSummariesTokenLength returns the maximum token length allowed for entity
	// and relationship descriptions before they need to be summarized by the LLM.
	GetMaxSummariesTokenLength() int
	// GetBatchSize returns
	GetBatchSize() int
}
