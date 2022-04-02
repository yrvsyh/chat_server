package model

type GroupMessage struct {
	Message
}

func init() {
	db.AutoMigrate(GroupMessage{})
}
