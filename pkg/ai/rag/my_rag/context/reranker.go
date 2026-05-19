package context

import (
	"context"
	"sort"
)

type Reranker struct {
	rerankModel string
}

func NewReranker(rerankModel string) *Reranker {
	return &Reranker{
		rerankModel: rerankModel,
	}
}

func (r *Reranker) Reorder(ctx context.Context, query string, chunks []string) ([]string, error) {
	if len(chunks) <= 1 {
		return chunks, nil
	}

	simpleScores := r.calculateSimpleRelevance(query, chunks)

	sorted := make([]struct {
		chunk string
		score float64
	}, len(chunks))

	for i, chunk := range chunks {
		sorted[i] = struct {
			chunk string
			score float64
		}{chunk, simpleScores[i]}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})

	result := make([]string, 0, len(chunks))
	for i := len(sorted) - 1; i >= 0; i-- {
		result = append(result, sorted[i].chunk)
	}

	return result, nil
}

func (r *Reranker) calculateSimpleRelevance(query string, chunks []string) []float64 {
	queryWords := wordCount(query)
	scores := make([]float64, len(chunks))

	for i, chunk := range chunks {
		chunkWords := wordCount(chunk)
		score := calculateJaccardSimilarity(queryWords, chunkWords)
		scores[i] = score
	}

	return scores
}

func wordCount(text string) map[string]int {
	words := make(map[string]int)
	var currentWord []rune

	for _, char := range text {
		if isChinese(char) {
			if len(currentWord) > 0 {
				word := string(currentWord)
				words[word]++
				currentWord = nil
			}
			words[string(char)]++
		} else if isAlphaNumeric(char) {
			currentWord = append(currentWord, char)
		} else {
			if len(currentWord) > 0 {
				word := string(currentWord)
				words[word]++
				currentWord = nil
			}
		}
	}

	if len(currentWord) > 0 {
		word := string(currentWord)
		words[word]++
	}

	return words
}

func isChinese(char rune) bool {
	return char >= 0x4e00 && char <= 0x9fff
}

func isAlphaNumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func calculateJaccardSimilarity(words1, words2 map[string]int) float64 {
	if len(words1) == 0 || len(words2) == 0 {
		return 0
	}

	set1 := make(map[string]bool)
	for word := range words1 {
		set1[word] = true
	}

	intersection := 0
	union := len(set1)

	for word := range words2 {
		union++
		if set1[word] {
			intersection++
		} else {
			set1[word] = true
		}
	}

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

func (r *Reranker) GetModel() string {
	return r.rerankModel
}

func (r *Reranker) SetModel(model string) {
	r.rerankModel = model
}