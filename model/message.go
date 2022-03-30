package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID      int64  `gorm:"primaryKey" json:"id,omitempty"`
	Type    int32  `gorm:"index;notNull" json:"type,omitempty"`
	From    uint32 `gorm:"index;notNull;size:255" json:"from,omitempty"`
	To      uint32 `gorm:"index;notNull;size:255" json:"to,omitempty"`
	Content []byte `json:"content,omitempty"`

	CreatedAt time.Time      `json:"created_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func init() {
	db.AutoMigrate(Message{})
}
