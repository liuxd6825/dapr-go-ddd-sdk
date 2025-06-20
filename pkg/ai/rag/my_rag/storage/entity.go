package storage

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"text/template"
	"time"
)

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
}

// Source represents a document chunk with metadata.
// It contains the text content, size information, and position data.
type Source struct {
	Id         string `json:"id"`
	TenantId   string `json:"tenantId"`
	CaseId     string `json:"caseId"`
	DocId      string `json:"docId"`
	Content    string
	TokenSize  int
	OrderIndex int
}

func (s Source) GenID(docID string) string {
	return fmt.Sprintf("%s-chunk-%d", docID, s.OrderIndex)
}

// GraphEntity represents an entity in the knowledge graph.
// It contains information about the entity's name, type,
// descriptions, sources, and creation timestamp.
type GraphEntity struct {
	Id       string `json:"id"`
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	DocId    string `json:"docId"`

	Name         string `json:"name"`
	Type         string `json:"type"`
	Descriptions string `json:"description"`
	SourceIDs    string
	CreatedAt    time.Time
}

// GraphRelationship represents a relationship between two entities in the knowledge graph.
// It contains information about the source and target entities,
// relationship weight, descriptions, keywords, sources, and creation timestamp.
type GraphRelationship struct {
	Id       string `json:"id"`
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	DocId    string `json:"docId"`

	Source       string   `json:"source"`
	Target       string   `json:"target"`
	Weight       float64  `json:"strength"`
	Descriptions string   `json:"description"`
	Keywords     []string `json:"keywords"`
	SourceIDs    string
	CreatedAt    time.Time
}

var (
	// ErrEntityNotFound is returned when an entity is not found in the storage.
	ErrEntityNotFound = errors.New("entity not found")
	// ErrRelationshipNotFound is returned when a relationship is not found in the storage.
	ErrRelationshipNotFound = errors.New("relationship not found")
)

func CleanContent(content string) string {
	// Removes spaces and null characters.
	str := strings.TrimSpace(content)
	return strings.ReplaceAll(str, "\x00", "")
}

func PromptTemplate(name, templ string, data any) (string, error) {
	buf := strings.Builder{}
	tmpl := template.New(name).Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
	tmpl = template.Must(tmpl.Parse(templ))
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func AppendIfUnique(slice []string, item string) []string {
	if slices.Contains(slice, item) {
		return slice
	}
	return append(slice, item)
}

func MostFrequentItem(list []string) string {
	// Create a map to store counts
	counts := make(map[string]int)

	// Count occurrences of each string
	for _, item := range list {
		counts[item]++
	}

	// Find the item with highest count
	maxCount := 0
	var mostFreqItem string

	for item, count := range counts {
		if count > maxCount {
			maxCount = count
			mostFreqItem = item
		}
	}

	return mostFreqItem
}

func ThreeBacktick(caption string) string {
	return "```" + caption
}
