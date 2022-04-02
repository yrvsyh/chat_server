package model

type UserMessage struct {
	Message
}

func init() {
	db.AutoMigrate(UserMessage{})
}
