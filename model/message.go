package model

import (
	"chat_server/database"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	From    uint
	To      uint
	Content []byte
}

func init() {
	database.DB.AutoMigrate(Message{})
}
