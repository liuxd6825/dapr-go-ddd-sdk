package light_rag

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
)

type Model string

const (
	ModelLocal  Model = "local"
	ModelGlobal Model = "global"
	ModelHybrid Model = "hybrid"
	ModelNaive  Model = "naive"
	ModelMix    Model = "mix"
	ModelBypass Model = "bypass"
)

type ConversationHistoryItem struct {
	Role    string `json:"role"`    // user|assistant
	Content string `json:"content"` // message
}

type QueryRequest struct {
	Query                    string                     `json:"query"`
	Mode                     Model                      `json:"mode"`
	OnlyNeedContext          *bool                      `json:"only_need_context"`
	OnlyNeedPrompt           *bool                      `json:"only_need_prompt"`
	ResponseType             *string                    `json:"response_type"`
	TopK                     *int                       `json:"top_k"`
	MaxTokenForTextUnit      *int                       `json:"max_token_for_text_unit"`
	MaxTokenForGlobalContext *int                       `json:"max_token_for_global_context"`
	MaxTokenForLocalContext  *int                       `json:"max_token_for_local_context"`
	ConversationHistory      []*ConversationHistoryItem `json:"conversation_history"`
	HistoryTurns             *int                       `json:"history_turns"`
	UserPrompt               string                     `json:"user_prompt"`
}

type QueryResponse struct {
	Response string `json:"response"`
}

type TextRequest struct {
	Text string `json:"text"`
}

type TextsRequest []string

type QueryService struct {
	cli *Client
}

func NewQueryService(cli *Client) *QueryService {
	return &QueryService{cli: cli}
}

func (s *QueryService) Query(ctx context.Context, query *QueryRequest) (*QueryResponse, error) {
	resp, err := s.cli.Post(ctx, "/query", query)
	if err != nil {
		panic(err)
	}
	var result *QueryResponse
	if err := s.cli.getJsonData(ctx, resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *QueryService) QueryStream(ctx context.Context, query *QueryRequest, rec func(resp *QueryResponse)) error {
	resp, err := s.cli.Post(ctx, "/query/stream", query)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var data *QueryResponse
		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			fmt.Printf("解码错误: %v\n", err)
			continue
		}
		if rec != nil {
			rec(data)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func NewQueryRequest(query string, model Model) *QueryRequest {
	return &QueryRequest{
		Mode:  model,
		Query: query,
	}
}
