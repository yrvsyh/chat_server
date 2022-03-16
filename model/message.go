package model

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	ID      int64  `gorm:"primaryKey"`
	Type    int32  `gorm:"index"`
	From    string `gorm:"index"`
	To      string `gorm:"index"`
	Content []byte

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func init() {
	db.AutoMigrate(Message{})
}
