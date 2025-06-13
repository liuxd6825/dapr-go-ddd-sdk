package my_rag

import (
	"os"
	"strconv"
)

type Content struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// QueryParam Configuration parameters for query execution in LightRAG.
type QueryParam struct {
	TenantId string `json:"tenantId"`
	CaseId   string `json:"caseId"`
	Query    string `json:"query"`
	MaxDeep  int    `json:"maxDeep"`
	// Mode specifies the retrieval mode:
	// - "local": Focuses on context-dependent information.
	// - "global": Utilizes global knowledge.
	// - "hybrid": Combines local and global retrieval methods.
	// - "naive": Performs a basic search without advanced techniques.
	// - "mix": Integrates knowledge graph and vector retrieval.
	Mode string `default:"global" json:"mode"`

	// OnlyNeedContext if true, only returns the retrieved context without generating a response.
	OnlyNeedContext bool `default:"false" json:"onlyNeedContext"`

	// OnlyNeedPrompt if true, only returns the generated prompt without producing a response.
	OnlyNeedPrompt bool `default:"false" json:"onlyNeedPrompt"`

	// ResponseType defines the response format. Examples: 'Multiple Paragraphs', 'Single Paragraph', 'Bullet Points'.
	ResponseType string `default:"Multiple Paragraphs" json:"responseType"`

	// Stream if true, enables streaming output for real-time responses.
	Stream bool `default:"true" json:"stream"`

	// TopK number of top items to retrieve. Represents entities in 'local' mode and relationships in 'global' mode.
	TopK int `default:"60" json:"topK"`

	// MaxTokenForTextUnit maximum number of tokens allowed for each retrieved text chunk.
	MaxTokenForTextUnit int `default:"4000" json:"maxTokenForTextUnit"`

	// MaxTokenForGlobalContext maximum number of tokens allocated for relationship descriptions in global retrieval.
	MaxTokenForGlobalContext int `default:"4000" json:"maxTokenForGlobalContext"`

	// MaxTokenForLocalContext maximum number of tokens allocated for entity descriptions in local retrieval.
	MaxTokenForLocalContext int `default:"4000" json:"maxTokenForLocalContext"`

	// HLKeywords list of high-level keywords to prioritize in retrieval.
	HLKeywords []string `json:"hlKeywords"`

	// LLKeywords list of low-level keywords to refine retrieval focus.
	LLKeywords []string `json:"llKeywords"`

	// ConversationHistory stores past conversation history to maintain context.
	// Format: [{"role": "user/assistant", "content": "message"}].
	ConversationHistory []*Content `json:"conversationHistory"`

	// HistoryTurns number of complete conversation turns (user-assistant pairs) to consider in the response context.
	HistoryTurns int `default:"3" json:"historyTurns"`

	// IDs list of ids to filter the results.
	IDs []string `json:"ids"`

	// ModelFunc optional override for the LLM model function to use for this specific query.
	// If provided, this will be used instead of the global model function.
	// This allows using different models for different query modes.
	//ModelFunc func(...interface{}) interface{}

	// UserPrompt user-provided prompt for the query.
	// If provided, this will be use instead of the default value from prompt template.
	UserPrompt string `json:"userPrompt"`
}

// NewQueryParam creates a new QueryParam with default values
func NewQueryParam() *QueryParam {
	topK, _ := strconv.Atoi(getEnv("TOP_K", "60"))
	maxTokenTextChunk, _ := strconv.Atoi(getEnv("MAX_TOKEN_TEXT_CHUNK", "4000"))
	maxTokenRelationDesc, _ := strconv.Atoi(getEnv("MAX_TOKEN_RELATION_DESC", "4000"))
	maxTokenEntityDesc, _ := strconv.Atoi(getEnv("MAX_TOKEN_ENTITY_DESC", "4000"))

	return &QueryParam{
		Mode:                     "global",
		OnlyNeedContext:          false,
		OnlyNeedPrompt:           false,
		ResponseType:             "Multiple Paragraphs",
		Stream:                   false,
		TopK:                     topK,
		MaxTokenForTextUnit:      maxTokenTextChunk,
		MaxTokenForGlobalContext: maxTokenRelationDesc,
		MaxTokenForLocalContext:  maxTokenEntityDesc,
		HLKeywords:               []string{},
		LLKeywords:               []string{},
		ConversationHistory:      []*Content{},
		HistoryTurns:             3,
		IDs:                      nil,
		//ModelFunc:                nil,
		UserPrompt: "",
	}
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
