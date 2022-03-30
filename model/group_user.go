package model

import "time"

type GroupUser struct {
	GroupID       uint32 `gorm:"primaryKey" json:"group_id,omitempty"`
	Group         *Group
	UserID        uint32 `gorm:"primaryKey" json:"user_id,omitempty"`
	User          *User
	LastMessageID uint64 `json:"last_message_id,omitempty"`
	LastMessage   *Message
	Remark        string `gorm:"size:255" json:"remark,omitempty"`

	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func init() {
	db.AutoMigrate(GroupUser{})
}
