package model

type User struct {
	BaseModel
	Username  string `gorm:"uniqueIndex;size:255" form:"username" json:"username,omitempty"`
	Password  string `gorm:"size:255;notNull" form:"password"`
	PublicKey []byte `form:"public_key" json:"public_key,omitempty"`
	Nickname  string `gorm:"size:255" json:"nickname,omitempty"`
	Email     string `gorm:"size:255" json:"email,omitempty"`
	Phone     string `gorm:"size:255" json:"phone,omitempty"`
	Avatar    string `gorm:"size:255" json:"avatar,omitempty"`
}

func init() {
	db.AutoMigrate(User{})
}
