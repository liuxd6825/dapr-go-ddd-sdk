package service

import (
	"github.com/liuxd6825/dapr-go-ddd-sdk/app/domain/rag/dao"
)

type MessageService struct {
	*dao.MessageDao
}

func NewMessageService() *MessageService {
	return &MessageService{
		MessageDao: dao.NewMessageDao(DBKey),
	}
}
