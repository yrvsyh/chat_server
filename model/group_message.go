package model

type GroupMessage struct {
	Message
	// GroupID    uint32 `gorm:"primaryKey"`
	// Group      *Group
	// SenderID   uint32 `gorm:"primaryKey"`
	// Sender     *User
	// LocalMsgID uint64 `gorm:"index;autoIncrement"`
	// MessageID  uint64
	// Message    *Message
}

func init() {
	db.AutoMigrate(GroupMessage{})
}
