package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
)

type ChatService struct {
	*dao.ChatDao
}

func NewChatService(dbKey string) *ChatService {
	return &ChatService{
		ChatDao: dao.NewChatDao(dbKey),
	}
}
