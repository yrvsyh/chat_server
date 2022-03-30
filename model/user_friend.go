package model

import "time"

type UserFriend struct {
	UserID        uint32 `gorm:"primaryKey" json:"user_id,omitempty"`
	User          *User
	FriendID      uint32 `gorm:"primaryKey" json:"friend_id,omitempty"`
	Friend        *User
	LastMessageID uint64 `json:"last_message_id,omitempty"`
	LastMessage   *Message
	Remark        string `gorm:"size:255" json:"remark,omitempty"`
	Accept        bool   `gorm:"default:false" json:"accept,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func init() {
	db.AutoMigrate(UserFriend{})
}
