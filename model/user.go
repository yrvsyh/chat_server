package model

import "time"

// Model定义
type (
	User struct {
		BaseModel
		Username  string    `gorm:"uniqueIndex" form:"username"`
		PublicKey []byte    `form:"public_key"`
		UserAuth  *UserAuth `gorm:"foreignKey:Username;references:Username"`

		Email  string
		Phone  string
		Avatar string

		Friends         []*UserFriends `gorm:"foreignKey:UserID"`
		Groups          []*UserGroups  `gorm:"foreignKey:UserID"`
		SendMessages    []*Message     `gorm:"foreignKey:From;references:Username"`
		ReceiveMessages []*Message     `gorm:"foreignKey:To;references:Username"`
	}

	UserAuth struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}

	UserFriends struct {
		UserID   uint  `gorm:"primaryKey"`
		FriendID uint  `gorm:"primaryKey"`
		Friend   *User `json:"friend,omitempty"`
		Remark   string

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	UserGroups struct {
		UserID  uint `gorm:"primaryKey"`
		GroupID uint `gorm:"primaryKey"`
		Group   *Group
		Remark  string

		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

//type (
//	UserDTO struct {
//		Username  string `form:"username" binding:"required"`
//		PublicKey string
//
//		Email  string
//		Phone  string
//		Avatar string
//	}
//
//	UserAuthDTO struct {
//		UserDTO
//		Password string `form:"password" binding:"required"`
//	}
//)

func init() {
	db.AutoMigrate(User{}, UserAuth{}, UserFriends{}, UserGroups{})
}

//func (user *User) ToDTO() *UserDTO {
//	return &UserDTO{
//		Username:  user.Username,
//		PublicKey: string(user.PublicKey),
//		Email:     user.Email,
//		Phone:     user.Phone,
//		Avatar:    user.Avatar,
//	}
//}
//
//func (user *UserAuthDTO) ToUserAuth() *UserAuth {
//	return &UserAuth{
//		User: User{
//			Username:  user.Username,
//			PublicKey: []byte(user.PublicKey),
//			Email:     user.Email,
//			Phone:     user.Phone,
//			Avatar:    user.Phone,
//		},
//		Password: user.Password,
//	}
//}
