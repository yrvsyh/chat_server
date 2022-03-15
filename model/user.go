package model

import (
	"chat_server/database"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"uniqueIndex"`
	Password  string
	PublicKey []byte

	Email  string
	Phone  string
	Avatar string

	Friends         []*User    `gorm:"many2many:user_friends"`
	Groups          []*Group   `gorm:"many2many:user_groups;"`
	SendMessages    []*Message `gorm:"foreignKey:From"`
	ReceiveMessages []*Message `gorm:"foreignKey:To"`
}

type UserFriends struct {
	UserID   uint `gorm:"primaryKey"`
	FriendID uint `gorm:"primaryKey"`
	Remark   string
}

type UserGroups struct {
	UserID  uint `gorm:"primaryKey"`
	GroupID uint `gorm:"primaryKey"`
	Remark  string
}

func init() {
	database.DB.AutoMigrate(User{}, UserFriends{}, UserGroups{})
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
