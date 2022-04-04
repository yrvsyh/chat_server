package model

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID      int64  `gorm:"primaryKey" json:"id"`
	Type    int32  `gorm:"index;notNull" json:"type"`
	State   int32  `gorm:"index;notNull;default:0"`
	From    uint32 `gorm:"index;notNull;size:255" json:"from"`
	To      uint32 `gorm:"index;notNull;size:255" json:"to"`
	Content []byte `json:"content,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
