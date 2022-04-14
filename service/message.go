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

// func (MessageService) UpdateUserLastMsgID(userID uint32, friendID uint32, msgID int64) error {
// 	userFriend := &model.UserFriend{}
// 	userFriend.UserID = userID
// 	userFriend.FriendID = friendID
// 	return db.Model(userFriend).Update("last_msg_id", msgID).Error
// }

// func (MessageService) UpdateGroupLastMsgID(groupID uint32, userID uint32, msgID int64) error {
// 	groupUser := &model.GroupUser{}
// 	groupUser.GroupID = groupID
// 	groupUser.UserID = userID
// 	return db.Model(groupUser).Update("last_msg_id", msgID).Error
// }

func (MessageService) UpdateLastMsgID(msg *message.Message) error {
	t := message.CheckMessageType(msg)
	switch t {
	case message.MESSAGE_TYPE_USER:
		userFriend := &model.UserFriend{}
		userFriend.UserID = msg.To
		userFriend.FriendID = msg.From
		return db.Model(userFriend).Update("last_msg_id", msg.Id).Error
	case message.MESSAGE_TYPE_GROUP:
		groupUser := &model.GroupUser{}
		groupUser.GroupID = msg.From
		groupUser.UserID = msg.To
		return db.Model(groupUser).Update("last_msg_id", msg.Id).Error
	default:
		return errors.New("msg type error")
	}
}
