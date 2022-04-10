package model

type User struct {
	BaseModel
	Username  string `gorm:"uniqueIndex;size:255" form:"username" json:"username,omitempty"`
	Password  string `gorm:"size:255;notNull;default:''" form:"password" json:"password"`
	PublicKey []byte `form:"public_key" json:"public_key,omitempty"`
	Nickname  string `gorm:"size:255;notNull;default:''" json:"nickname,omitempty"`
	Email     string `gorm:"size:255;notNull;default:''" json:"email,omitempty"`
	Phone     string `gorm:"size:255;notNull;default:''" json:"phone,omitempty"`
	Avatar    string `gorm:"size:255;notNull;default:''" json:"avatar,omitempty"`
}

func init() {
	db.AutoMigrate(User{})
}
