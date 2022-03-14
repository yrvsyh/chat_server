package model

import "chat_server/database"

type (
	User struct {
		Username string `gorm:"primarykey;notNull;default:''" json:"username,omitempty" form:"username"`
		Password string `gorm:"notNull;default:''" json:"password,omitempty" form:"password"`
		Email    string `gorm:"notNull;default:''" json:"email,omitempty" form:"email"`
		Avatar   string `gorm:"notNull;default:''" json:"avatar,omitempty"`
		// Role     int    `gorm:"notNull;default:1" json:"role,omitempty" form:"role"`
	}
)

// const (
// 	USER_ROLE_ADMIN  = 0
// 	USER_ROLE_NORMAL = 1
// )

func init() {
	database.DB.AutoMigrate(&User{})
}

// func (*UserModel) TableName() string { return "user" }

// func (user *UserModel) ToDTO() *dto.User {
// 	return &dto.User{
// 		Name:  user.Name,
// 		Email: user.Email,
// 		Role:  user.Role,
// 	}
// }

func GetUserById(username string) *User {
	user := &User{Username: username}
	if database.DB.First(user).Error != nil {
		return nil
	}
	return user
}

func InsertUser(user *User) error {
	return database.DB.Create(user).Error
}
