package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/llm"
	ragcontext "github.com/liuxd6825/dapr-go-ddd-sdk/pkg/ai/rag/my_rag/context"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type ChatSummaryService struct {
	*dao.ChatSummaryDao
	config    *config.CompressConfig
	compressor *ragcontext.Compressor
}

func NewChatSummaryService(cfg *config.CompressConfig) *ChatSummaryService {
	return &ChatSummaryService{
		ChatSummaryDao: dao.NewChatSummaryDao(DBKey),
		config:          cfg,
	}
}

func (s *ChatSummaryService) SetLLM(llm llm.LLM) {
	summaryCfg := ragcontext.SummaryConfig{
		Threshold: s.config.GetThreshold(),
		Model:     "",
	}
	s.compressor = ragcontext.NewCompressor(llm, summaryCfg)
}

func (s *ChatSummaryService) GetByChatId(ctx context.Context, chatId string) (*model.ChatSummary, error) {
	result, err := s.ChatSummaryDao.FindByRSQL(ctx, "chat_id=="+chatId)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("chat_summary_not_found")
	}
	return result[0], nil
}

func (s *ChatSummaryService) CreateIfNotExists(ctx context.Context, chatId string) error {
	existing, _ := s.GetByChatId(ctx, chatId)
	if existing != nil {
		return nil
	}

	summary := &model.ChatSummary{
		ChatID:             chatId,
		Summary:            "",
		CompressedPosition: 0,
		MessageCount:       0,
		OriginalTokens:     0,
		CompressedTokens:   0,
		IsLocked:           false,
	}

	return s.ChatSummaryDao.Create(ctx, summary).GetError()
}

type CompressResult struct {
	Success            bool
	Skipped           bool
	Reason             string
	Summary            string
	CompressedPosition int
	OriginalTokens      int
	CompressedTokens    int
}

func (s *ChatSummaryService) Compress(ctx context.Context, chatId string, force bool, messages []*model.Message) (*CompressResult, error) {
	summary, err := s.GetByChatId(ctx, chatId)
	if err != nil {
		return nil, err
	}

	if !s.AcquireLock(ctx, chatId) {
		return &CompressResult{
			Success: false,
			Reason:  "compress_in_progress",
		}, errors.New("compress_in_progress")
	}
	defer s.ReleaseLock(ctx, chatId)

	totalTokens := s.CountTokens(messages)

	if !force && totalTokens < s.config.GetThreshold() {
		return &CompressResult{
			Success: false,
			Skipped: true,
			Reason:  "threshold_not_reached",
		}, nil
	}

	var newSummary string
	if s.compressor != nil {
		contextMessages := s.convertToContextMessages(messages)
		newSummary, err = s.compressor.GenerateSummary(ctx, summary.Summary, contextMessages)
		if err != nil {
			newSummary = summary.Summary
		}
	} else {
		newSummary = s.SimpleSummarize(messages)
	}

	maxOrderNum := s.GetMaxOrderNum(messages)

	summary.Summary = newSummary
	summary.CompressedPosition = maxOrderNum
	summary.MessageCount = len(messages)
	summary.OriginalTokens = totalTokens
	summary.CompressedTokens = s.CountMessageTokens(newSummary)
	summary.UpdatedTime = ptrTime(time.Now())

	opts := idao.NewCallOptions().SetUpdateFields([]string{
		"summary", "compressed_position", "message_count",
		"original_tokens", "compressed_tokens", "updated_time",
	})
	err = s.ChatSummaryDao.Update(ctx, summary, opts).GetError()
	if err != nil {
		return nil, err
	}

	return &CompressResult{
		Success:            true,
		Summary:            newSummary,
		CompressedPosition: maxOrderNum,
		OriginalTokens:     totalTokens,
		CompressedTokens:   s.CountMessageTokens(newSummary),
	}, nil
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func (s *ChatSummaryService) AcquireLock(ctx context.Context, chatId string) bool {
	summary, err := s.GetByChatId(ctx, chatId)
	if err != nil {
		return false
	}

	if summary.IsLocked {
		return false
	}

	summary.IsLocked = true
	summary.UpdatedTime = ptrTime(time.Now())
	opts := idao.NewCallOptions().SetUpdateFields([]string{"is_locked", "updated_time"})
	err = s.ChatSummaryDao.Update(ctx, summary, opts).GetError()
	return err == nil
}

func (s *ChatSummaryService) ReleaseLock(ctx context.Context, chatId string) {
	summary, err := s.GetByChatId(ctx, chatId)
	if err != nil {
		return
	}

	if !summary.IsLocked {
		return
	}

	summary.IsLocked = false
	summary.UpdatedTime = ptrTime(time.Now())
	opts := idao.NewCallOptions().SetUpdateFields([]string{"is_locked", "updated_time"})
	s.ChatSummaryDao.Update(ctx, summary, opts)
}

func (s *ChatSummaryService) ForceUnlock(ctx context.Context, chatId string) error {
	summary, err := s.GetByChatId(ctx, chatId)
	if err != nil {
		return err
	}

	summary.IsLocked = false
	summary.UpdatedTime = ptrTime(time.Now())
	opts := idao.NewCallOptions().SetUpdateFields([]string{"is_locked", "updated_time"})
	return s.ChatSummaryDao.Update(ctx, summary, opts).GetError()
}

func (s *ChatSummaryService) CountMessageTokens(text string) int {
	return len(text)
}

func (s *ChatSummaryService) CountTokens(messages []*model.Message) int {
	if messages == nil || len(messages) == 0 {
		return 0
	}
	total := 0
	for _, msg := range messages {
		total += len(msg.Content)
	}
	return total
}

func (s *ChatSummaryService) GetMaxOrderNum(messages []*model.Message) int {
	if messages == nil || len(messages) == 0 {
		return 0
	}
	max := float32(0)
	for _, msg := range messages {
		if msg.OrderNum > max {
			max = msg.OrderNum
		}
	}
	return int(max)
}

func (s *ChatSummaryService) convertToContextMessages(messages []*model.Message) []*ragcontext.Content {
	if messages == nil {
		return nil
	}
	result := make([]*ragcontext.Content, 0, len(messages))
	for _, msg := range messages {
		result = append(result, &ragcontext.Content{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}
	return result
}

func (s *ChatSummaryService) SimpleSummarize(messages []*model.Message) string {
	if messages == nil || len(messages) == 0 {
		return ""
	}
	
	var sb strings.Builder
	sb.WriteString("对话摘要：")
	
	count := len(messages)
	if count > 10 {
		sb.WriteString("(共")
		sb.WriteString(strconv.Itoa(count))
		sb.WriteString("条消息) ")
	}
	
	for i, msg := range messages {
		if i >= 5 {
			sb.WriteString("... (更多消息)")
			break
		}
		if msg.Role == "user" {
			content := msg.Content
			if len(content) > 50 {
				content = content[:50] + "..."
			}
			sb.WriteString("用户: ")
			sb.WriteString(content)
			sb.WriteString("; ")
		}
	}
	
	return sb.String()
}

func (s *ChatSummaryService) GetUncompressedMessages(ctx context.Context, chatId string, allMessages []*model.Message) []*model.Message {
	summary, err := s.GetByChatId(ctx, chatId)
	if err != nil {
		return allMessages
	}

	if summary.CompressedPosition == 0 {
		return allMessages
	}

	var result []*model.Message
	for _, msg := range allMessages {
		if int(msg.OrderNum) > summary.CompressedPosition {
			result = append(result, msg)
		}
	}
	return result
}

func (s *ChatSummaryService) ShouldCompress(ctx context.Context, totalTokens int) bool {
	return totalTokens >= s.config.GetThreshold()
}

func (s *ChatSummaryService) GetConfig() *config.CompressConfig {
	return s.config
}