package service

import (
	"chat_server/message"
	"chat_server/model"
)

type MessageService struct{}

func (MessageService) SaveMessage(message *message.Message) error {
	msg := &model.Message{
		ID:      message.Id,
		Type:    int32(message.Type),
		From:    message.From,
		To:      message.To,
		Content: message.Content,
	}
	return db.Create(msg).Error
}
