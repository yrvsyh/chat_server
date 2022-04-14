package model

type User struct {
	BaseModel
	Username  string `gorm:"uniqueIndex;size:255" form:"username" json:"username"`
	Password  string `gorm:"size:255;notNull;default:''" form:"password"`
	PublicKey []byte `form:"public_key" json:"public_key"`
	Nickname  string `gorm:"size:255;notNull;default:''" json:"nickname"`
	Email     string `gorm:"size:255;notNull;default:''" json:"email"`
	Phone     string `gorm:"size:255;notNull;default:''" json:"phone"`
	Avatar    string `gorm:"size:255;notNull;default:''" json:"avatar"`
}
