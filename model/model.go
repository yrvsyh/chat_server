package model

import (
	"chat_server/database"
	"time"
)

type BaseModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var db = database.DB
