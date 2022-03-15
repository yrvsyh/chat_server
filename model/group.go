package model

import (
	"chat_server/database"
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name      string
	PublicKey []byte

	Users []*User `gorm:"many2many:user_groups;"`
}

func init() {
	database.DB.AutoMigrate(Group{})
}
