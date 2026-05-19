package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/config"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/model"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/db/dao/idao"
	"github.com/liuxd6825/dapr-go-ddd-sdk/pkg/errors"
)

type MessageService struct {
	*dao.MessageDao
	chatSummaryService *ChatSummaryService
	compressConfig     *config.CompressConfig
	compressLock       sync.Map
}

func NewMessageService() *MessageService {
	return &MessageService{
		MessageDao: dao.NewMessageDao(DBKey),
	}
}

func (s *MessageService) SetChatSummaryService(chatSummaryService *ChatSummaryService) {
	s.chatSummaryService = chatSummaryService
}

func (s *MessageService) SetCompressConfig(cfg *config.CompressConfig) {
	s.compressConfig = cfg
}

func (s *MessageService) FindByChatId(ctx context.Context, chatId string) ([]*model.Message, error) {
	sort := "order_num"
	opts := idao.NewCallOptions()
	opts.SetSort(&sort)

	return s.Dao.FindByRSQL(ctx, fmt.Sprintf("chat_id=='%s'", chatId), opts)
}

func (s *MessageService) FindByChatIdAfterPosition(ctx context.Context, chatId string, position int) ([]*model.Message, error) {
	if position == 0 {
		return s.FindByChatId(ctx, chatId)
	}
	return s.Dao.FindByRSQL(ctx, "chat_id=='"+chatId+"' and order_num>"+string(rune(position))+" sort by order_num asc")
}

func (s *MessageService) GetUncompressedMessages(ctx context.Context, chatId string) ([]*model.Message, error) {
	messages, err := s.FindByChatId(ctx, chatId)
	if err != nil {
		return nil, err
	}

	if s.chatSummaryService == nil {
		return messages, nil
	}

	return s.chatSummaryService.GetUncompressedMessages(ctx, chatId, messages), nil
}

func (s *MessageService) Create(ctx context.Context, msg *model.Message) error {
	if s.chatSummaryService != nil {
		summary, err := s.chatSummaryService.GetByChatId(ctx, msg.ChatId)
		if err == nil && summary.IsLocked {
			return errors.New("chat_is_compressing")
		}
	}

	msg.OrderNum = s.getNextOrderNum(ctx, msg.ChatId)
	return s.MessageDao.Create(ctx, msg).GetError()
}

func (s *MessageService) getNextOrderNum(ctx context.Context, chatId string) float32 {
	messages, err := s.FindByChatId(ctx, chatId)
	if err != nil || len(messages) == 0 {
		return 1
	}

	maxOrder := float32(0)
	for _, msg := range messages {
		if msg.OrderNum > maxOrder {
			maxOrder = msg.OrderNum
		}
	}
	return maxOrder + 1
}

func (s *MessageService) GetMessagesBefore(ctx context.Context, chatId string, position int) ([]*model.Message, error) {
	if position == 0 {
		return []*model.Message{}, nil
	}
	return s.Dao.FindByRSQL(ctx, "chat_id=="+chatId+" and order_num<="+string(rune(position))+" sort by order_num asc")
}

func (s *MessageService) GetAllMessages(ctx context.Context, chatId string) ([]*model.Message, error) {
	return s.FindByChatId(ctx, chatId)
}

func (s *MessageService) GetLatestMessages(ctx context.Context, chatId string, limit int) ([]*model.Message, error) {
	messages, err := s.FindByChatId(ctx, chatId)
	if err != nil {
		return nil, err
	}

	if len(messages) <= limit {
		return messages, nil
	}

	return messages[len(messages)-limit:], nil
}

func (s *MessageService) CountTokens(ctx context.Context, chatId string) int {
	messages, err := s.FindByChatId(ctx, chatId)
	if err != nil {
		return 0
	}

	total := 0
	for _, msg := range messages {
		total += len(msg.Content)
	}
	return total
}

func (s *MessageService) CheckAndCompress(ctx context.Context, chatId string) error {
	if s.chatSummaryService == nil || s.compressConfig == nil {
		return nil
	}

	if _, loaded := s.compressLock.LoadOrStore(chatId, true); loaded {
		return nil
	}
	defer s.compressLock.Delete(chatId)

	summary, err := s.chatSummaryService.GetByChatId(ctx, chatId)
	if err != nil {
		return err
	}

	position := summary.CompressedPosition
	var messages []*model.Message

	if position == 0 {
		messages, err = s.FindByChatId(ctx, chatId)
	} else {
		messages, err = s.GetMessagesBefore(ctx, chatId, position)
	}

	if err != nil {
		return err
	}

	totalTokens := s.chatSummaryService.CountTokens(messages)
	if totalTokens < s.compressConfig.GetThreshold() {
		return nil
	}

	_, err = s.chatSummaryService.Compress(ctx, chatId, false, messages)
	return err
}

func (s *MessageService) IsLocked(ctx context.Context, chatId string) bool {
	if s.chatSummaryService == nil {
		return false
	}
	summary, err := s.chatSummaryService.GetByChatId(ctx, chatId)
	if err != nil {
		return false
	}
	return summary.IsLocked
}

func (s *MessageService) GetNextOrderNumWithLock(ctx context.Context, chatId string) (float32, error) {
	messages, err := s.FindByChatId(ctx, chatId)
	if err != nil {
		return 1, err
	}

	maxOrder := float32(0)
	for _, msg := range messages {
		if msg.OrderNum > maxOrder {
			maxOrder = msg.OrderNum
		}
	}

	return maxOrder + 1, nil
}

func (s *MessageService) CreateWithTimeout(ctx context.Context, msg *model.Message, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		done <- s.Create(ctx, msg)
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return errors.New("create_message_timeout")
	}
}

func (s *MessageService) DeleteByChatId(ctx context.Context, chatId string) error {
	return s.MessageDao.DeleteByRSQL(ctx, "chat_id=="+chatId).GetError()
}
