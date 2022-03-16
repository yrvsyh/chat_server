package service

import (
	"chat_server/message"
	"chat_server/model"
)

func SaveMessage(message *message.Message) error {
	msg := &model.Message{
		ID:      message.Id,
		Type:    message.Type,
		From:    message.From,
		To:      message.To,
		Content: message.Content,
	}
	return db.Create(msg).Error
}
