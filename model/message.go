package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	From    string
	To      string
	Type    uint8
	Content []byte
}

func init() {
	db.AutoMigrate(Message{})
}
