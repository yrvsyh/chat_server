package service

import (
	"chat_server/message"
	"chat_server/model"
	"errors"

	"github.com/sirupsen/logrus"
)

type MessageService struct{}

func (MessageService) LoadMessage(dbMsg *model.Message) *message.Message {
	msg := &message.Message{
		Id:      dbMsg.ID,
		Type:    message.Type(dbMsg.Type),
		State:   message.State(dbMsg.State),
		From:    dbMsg.From,
		To:      dbMsg.To,
		Content: dbMsg.Content,
	}
	return msg
}

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
		userMsg := &model.FriendMessage{Message: baseMsg}
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
		userMsg := &model.FriendMessage{Message: model.Message{ID: msg.Id}}
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
		// FIXME 分开存储每个好友的离线消息状态
		// userFriend := &model.UserFriend{}
		// userFriend.UserID =
		// userFriend.FriendID = msg.From
		// return db.Model(userFriend).Update("last_msg_id", msg.Id).Error
		user := &model.User{}
		user.ID = msg.To
		return db.Model(user).Update("last_msg_id", msg.Id).Error
	case message.MESSAGE_TYPE_GROUP:
		groupUser := &model.GroupUser{}
		groupUser.GroupID = msg.To
		groupUser.UserID = msg.From
		return db.Model(groupUser).Update("last_msg_id", msg.Id).Error
	default:
		return errors.New("msg type error")
	}
}

func (MessageService) GetFriendOfflineMessages(id uint32) ([]model.FriendMessage, error) {
	var messages []model.FriendMessage
	// FIXME 分开存储每个好友的离线消息状态
	user := &model.User{}
	if err := db.First(user, id).Error; err != nil {
		return messages, err
	}
	lastMsgID := user.LastMsgID
	// WHY 需要加上反引号
	err := db.Where("`to` = ?", id).Where("`id` > ?", lastMsgID).Find(&messages).Error
	return messages, err
}

func (MessageService) GetGroupOfflineMessages(id uint32) ([]model.GroupMessage, error) {
	var messages []model.GroupMessage

	var groupUsers []model.GroupUser
	if err := db.Where("user_id = ?", id).Find(&groupUsers).Error; err != nil {
		return messages, err
	}
	var retErr error
	for _, groupUser := range groupUsers {
		lastMsgID := groupUser.LastMsgID
		var groupMessages []model.GroupMessage
		if err := db.Where("`to` = ?", groupUser.GroupID).Where("`id` > ?", lastMsgID).Where("`from` != ?", id).Find(&groupMessages).Error; err != nil {
			logrus.Error(err)
			retErr = err
		}
		messages = append(messages, groupMessages...)
	}
	return messages, retErr
}
