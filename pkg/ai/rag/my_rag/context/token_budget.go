package context

import (
	"github.com/pkoukk/tiktoken-go"
)

type TokenBudget struct {
	encoding  *tiktoken.Tiktoken
	maxTokens int
	reserved  int
	available int
}

func NewTokenBudget(maxTokens, reserved int) (*TokenBudget, error) {
	encoding, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		return nil, err
	}
	return &TokenBudget{
		encoding:  encoding,
		maxTokens: maxTokens,
		reserved:  reserved,
		available: maxTokens - reserved,
	}, nil
}

func (tb *TokenBudget) CountTokens(text string) int {
	return len(tb.encoding.Encode(text, nil, nil))
}

func (tb *TokenBudget) Allocate(ratio float64) int {
	return int(float64(tb.available) * ratio)
}

func (tb *TokenBudget) Remaining() int {
	return tb.available
}

func (tb *TokenBudget) Consume(tokens int) {
	tb.available -= tokens
	if tb.available < 0 {
		tb.available = 0
	}
}

func (tb *TokenBudget) MaxTokens() int {
	return tb.maxTokens
}

func (tb *TokenBudget) Reserved() int {
	return tb.reserved
}

func (tb *TokenBudget) Available() int {
	return tb.available
}

func (tb *TokenBudget) Reset() {
	tb.available = tb.maxTokens - tb.reserved
}