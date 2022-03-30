package model

import (
	"chat_server/database"
	"time"
)

type BaseModel struct {
	ID        uint32    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

var db = database.DB
