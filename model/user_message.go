package model

type UserMessage struct {
	Message
	// SenderID   uint32 `gorm:"primaryKey"`
	// Sender     *User
	// ReceiverID uint32 `gorm:"primaryKey"`
	// Receiver   *User
	// LocalMsgID uint64 `gorm:"index;autoIncrement"`
	// MessageID  uint64
	// Message    *Message
}

func init() {
	db.AutoMigrate(UserMessage{})
}
