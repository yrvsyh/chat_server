package service

import (
	"chat_server/message"
	"chat_server/model"
	"errors"
)

type MessageService struct{}

func (MessageService) SaveMessage(msg *message.Message) error {
	baseMsg := model.Message{
		ID:      msg.Id,
		Type:    int32(msg.Type),
		State:   int32(message.State_SERVER_RECV),
		From:    msg.From,
		To:      msg.To,
		Content: msg.Content,
	}

	t := message.CheckMessageType(msg)
	switch t {
	case message.MESSAGE_TYPE_USER:
		userMsg := &model.UserMessage{Message: baseMsg}
		return db.Create(userMsg).Error
	case message.MESSAGE_TYPE_GROUP:
		groupMsg := &model.GroupMessage{Message: baseMsg}
		return db.Create(groupMsg).Error
	default:
		return errors.New("msg type error")
	}
}

func (MessageService) UpdateMessageState(msg *message.Message, state message.State) error {
	t := message.CheckMessageType(msg)
	switch t {
	case message.MESSAGE_TYPE_USER:
		userMsg := &model.UserMessage{Message: model.Message{ID: msg.Id}}
		return db.Model(userMsg).Update("state", state).Error
	case message.MESSAGE_TYPE_GROUP:
		groupMsg := &model.GroupMessage{Message: model.Message{ID: msg.Id}}
		return db.Model(groupMsg).Update("state", state).Error
	default:
		return errors.New("msg type error")
	}
}
