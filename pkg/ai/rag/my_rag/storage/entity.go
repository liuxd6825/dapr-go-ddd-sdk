package storage

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"text/template"
	"time"
)

// Source represents a document chunk with metadata.
// It contains the text content, size information, and position data.
type Source struct {
	Id         string `json:"id"`
	TenantId   string `json:"tenantId"`
	CaseId     string `json:"caseId"`
	DocId      string `json:"docId"`
	Content    string `json:"content"`
	TokenSize  int    `json:"tokenSize"`
	OrderIndex int    `json:"orderIndex"`
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

	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Descriptions string    `json:"description"`
	SourceIDs    string    `json:"sourceIDs"`
	SourceUrl    string    `json:"sourceUrl"`
	SourceName   string    `json:"sourceName"`
	CreatedAt    time.Time `json:"createdAt"`
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
