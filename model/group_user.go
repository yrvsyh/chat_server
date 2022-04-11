package model

import "time"

type GroupUser struct {
	GroupID       uint32        `gorm:"primaryKey" json:"group_id"`
	Group         *Group        `json:"group,omitempty"`
	UserID        uint32        `gorm:"primaryKey" json:"user_id"`
	User          *User         `json:"user,omitempty"`
	LastMessageID uint64        `json:"last_message_id,omitempty"`
	LastMessage   *GroupMessage `json:"last_message,omitempty"`
	Remark        string        `gorm:"size:255;notNull;default:''" json:"remark"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
