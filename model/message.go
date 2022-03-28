package model

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID      int64  `gorm:"primaryKey" json:"id,omitempty"`
	Type    int32  `gorm:"index;notNull" json:"type,omitempty"`
	From    string `gorm:"index;notNull" json:"from,omitempty"`
	To      string `gorm:"index;notNull" json:"to,omitempty"`
	Content []byte `json:"content,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func init() {
	db.AutoMigrate(Message{})
}
