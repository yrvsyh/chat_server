package model

import (
	"time"
)

type BaseModel struct {
	ID        uint32    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
