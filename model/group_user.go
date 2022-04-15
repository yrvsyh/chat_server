package model

import "time"

type GroupUser struct {
	GroupID   uint32        `gorm:"primaryKey" json:"group_id"`
	Group     *Group        `json:"group,omitempty"`
	UserID    uint32        `gorm:"primaryKey" json:"user_id"`
	User      *User         `json:"user,omitempty"`
	LastMsgID uint64        `json:"last_msg_id"`
	LastMsg   *GroupMessage `json:"last_msg,omitempty"`
	Remark    string        `gorm:"size:255;notNull;default:''" json:"remark"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
