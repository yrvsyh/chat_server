package model

type Base struct {
	Id uint `gorm:"primarykey;notNull;autoIncrement" json:"id,omitempty"`
}
