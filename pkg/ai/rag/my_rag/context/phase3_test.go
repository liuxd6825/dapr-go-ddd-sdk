package context

import (
	"testing"
)

func TestCompressor_ShouldCompress(t *testing.T) {
	cfg := SummaryConfig{
		Threshold: 1000,
		Model:     "test-model",
	}
	compressor := NewCompressor(nil, cfg)

	if !compressor.ShouldCompress(1000) {
		t.Error("expected ShouldCompress(1000) to be true")
	}

	if !compressor.ShouldCompress(2000) {
		t.Error("expected ShouldCompress(2000) to be true")
	}

	if compressor.ShouldCompress(500) {
		t.Error("expected ShouldCompress(500) to be false")
	}
}

func TestCompressor_GetConfig(t *testing.T) {
	cfg := SummaryConfig{
		Threshold: 1000,
		Model:     "test-model",
	}
	compressor := NewCompressor(nil, cfg)

	result := compressor.GetConfig()
	if result.Threshold != 1000 {
		t.Errorf("expected threshold 1000, got %d", result.Threshold)
	}
	if result.Model != "test-model" {
		t.Errorf("expected model 'test-model', got %s", result.Model)
	}
}

func TestQueryRewriter_RewriteWithEmptyHistory(t *testing.T) {
	rewriter := NewQueryRewriter(nil)
	query := "GraphRAG是什么？"

	result, err := rewriter.Rewrite(nil, query, []*Content{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != query {
		t.Errorf("expected '%s', got '%s'", query, result)
	}
}

func TestQueryRewriter_RewriteStandalone(t *testing.T) {
	rewriter := NewQueryRewriter(nil)
	if rewriter.llm == nil {
		t.Skip("skipping test: LLM is nil")
	}
	query := "它和传统RAG有什么区别？"

	result, err := rewriter.RewriteStandalone(nil, query)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != query {
		t.Errorf("expected original query (no LLM), got '%s'", result)
	}
}

func TestReranker_NewReranker(t *testing.T) {
	reranker := NewReranker("bge-reranker")
	if reranker.GetModel() != "bge-reranker" {
		t.Errorf("expected model 'bge-reranker', got '%s'", reranker.GetModel())
	}
}

func TestReranker_ReorderSingleChunk(t *testing.T) {
	reranker := NewReranker("test")
	chunks := []string{"只有一段内容"}

	result, err := reranker.Reorder(nil, "查询", chunks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(result))
	}
}

func TestReranker_ReorderMultipleChunks(t *testing.T) {
	reranker := NewReranker("test")
	chunks := []string{
		"GraphRAG是混合检索技术",
		"传统RAG只依赖向量检索",
		"天气今天很好",
	}

	result, err := reranker.Reorder(nil, "GraphRAG和传统RAG的区别", chunks)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(result))
	}
}

func TestReranker_SetModel(t *testing.T) {
	reranker := NewReranker("original")
	reranker.SetModel("new-model")
	if reranker.GetModel() != "new-model" {
		t.Errorf("expected 'new-model', got '%s'", reranker.GetModel())
	}
}

func TestWordCount(t *testing.T) {
	words := wordCount("GraphRAG 是 混合 检索 技术")
	if words["GraphRAG"] != 1 {
		t.Errorf("expected GraphRAG count 1, got %d", words["GraphRAG"])
	}
	if words["是"] != 1 {
		t.Errorf("expected 是 count 1, got %d", words["是"])
	}

	words = wordCount("今天天气很好")
	if words["今"] != 1 || words["天"] != 2 {
		t.Error("Chinese character counting failed")
	}
}

func TestIsChinese(t *testing.T) {
	if !isChinese('中') {
		t.Error("expected '中' to be Chinese")
	}
	if isChinese('A') {
		t.Error("expected 'A' to not be Chinese")
	}
	if isChinese('1') {
		t.Error("expected '1' to not be Chinese")
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	if !isAlphaNumeric('a') {
		t.Error("expected 'a' to be alphanumeric")
	}
	if !isAlphaNumeric('Z') {
		t.Error("expected 'Z' to be alphanumeric")
	}
	if !isAlphaNumeric('0') {
		t.Error("expected '0' to be alphanumeric")
	}
	if isAlphaNumeric('中') {
		t.Error("expected '中' to not be alphanumeric")
	}
	if isAlphaNumeric(' ') {
		t.Error("expected ' ' to not be alphanumeric")
	}
}

func TestCalculateJaccardSimilarity(t *testing.T) {
	words1 := map[string]int{"a": 1, "b": 1, "c": 1}
	words2 := map[string]int{"b": 1, "c": 1, "d": 1}

	score := calculateJaccardSimilarity(words1, words2)
	if score <= 0 {
		t.Errorf("expected positive score, got %f", score)
	}

	score = calculateJaccardSimilarity(map[string]int{}, words2)
	if score != 0 {
		t.Errorf("expected 0 for empty words1, got %f", score)
	}

	score = calculateJaccardSimilarity(words1, map[string]int{})
	if score != 0 {
		t.Errorf("expected 0 for empty words2, got %f", score)
	}
}

func TestExtractJSONFromResponse(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`{"key": "value"}`, `{"key": "value"}`},
		{`这里有一些文字{"key": "value"}后面也有`, `{"key": "value"}`},
		{`[1, 2, 3]`, `[1, 2, 3]`},
		{`无JSON在这里`, ``},
		{`{"nested": {"key": "value"}}`, `{"nested": {"key": "value"}}`},
	}

	for _, tt := range tests {
		result := extractJSONFromResponse(tt.input)
		if result != tt.expected {
			t.Errorf("extractJSONFromResponse(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}