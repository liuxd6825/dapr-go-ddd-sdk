package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
)

const DBKey string = "$ragDb"

type ChatService struct {
	*dao.ChatDao
}

func NewChatService() *ChatService {
	return &ChatService{
		ChatDao: dao.NewChatDao(DBKey),
	}
}
