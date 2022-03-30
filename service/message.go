package service

import (
	"chat_server/message"
	"chat_server/model"
)

type MessageService struct{}

func (MessageService) SaveUserMessage(msg *message.Message) error {
	userMsg := &model.UserMessage{
		Message: model.Message{
			ID:      msg.Id,
			Type:    int32(msg.Type),
			From:    msg.From,
			To:      msg.To,
			Content: msg.Content,
		},
	}
	return db.Create(userMsg).Error
}

func (MessageService) SaveGroupMessage(msg *message.Message) error {
	groupMsg := &model.GroupMessage{
		Message: model.Message{
			ID:      msg.Id,
			Type:    int32(msg.Type),
			From:    msg.From,
			To:      msg.To,
			Content: msg.Content,
		},
	}
	return db.Create(groupMsg).Error
}
