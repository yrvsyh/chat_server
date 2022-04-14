package model

import "time"

type UserFriend struct {
	UserID        uint32       `gorm:"primaryKey" json:"user_id"`
	User          *User        `json:"user,omitempty"`
	FriendID      uint32       `gorm:"primaryKey" json:"friend_id"`
	Friend        *User        `json:"friend,omitempty"`
	LastMessageID uint64       `json:"last_message_id"`
	LastMessage   *UserMessage `json:"last_message,omitempty"`
	Remark        string       `gorm:"size:255;notNull;default:''" json:"remark"`
	Accept        bool         `gorm:"notNull;default:false" json:"accept"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
