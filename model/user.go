package model

type User struct {
	BaseModel
	Username  string         `gorm:"uniqueIndex;size:255" form:"username" json:"username"`
	Password  string         `gorm:"size:255;notNull;default:''" form:"password"`
	PublicKey string         `form:"public_key" json:"public_key"`
	Nickname  string         `gorm:"index;size:255;notNull;default:''" json:"nickname"`
	Email     string         `gorm:"size:255;notNull;default:''" json:"email"`
	Phone     string         `gorm:"size:255;notNull;default:''" json:"phone"`
	Avatar    string         `gorm:"size:255;notNull;default:'default.png'" json:"avatar"`
	LastMsgID *int64         `json:"last_msg_id"`
	LastMsg   *FriendMessage `json:"last_msg,omitempty"`
}
