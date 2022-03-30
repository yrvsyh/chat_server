package model

import (
	"time"

	"gorm.io/gorm"
)

type (
	Message struct {
		ID      int64  `gorm:"primaryKey" json:"id,omitempty"`
		Type    int32  `gorm:"index;notNull" json:"type,omitempty"`
		From    string `gorm:"index;notNull;size:255" json:"from,omitempty"`
		To      string `gorm:"index;notNull;size:255" json:"to,omitempty"`
		Content []byte `json:"content,omitempty"`

		CreatedAt time.Time      `json:"created_at,omitempty"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	}

	UserMessage struct {
		Message
	}

	GroupMessage struct {
		Message
	}
)

func init() {
	db.AutoMigrate(UserMessage{}, GroupMessage{})
}
