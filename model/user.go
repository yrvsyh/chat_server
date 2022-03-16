package model

import "time"

type (
	User struct {
		Name      string `gorm:"primaryKey" form:"username"`
		UserAuth  *UserAuth
		PublicKey []byte `form:"public_key"`
		Nickname  string
		Email     string
		Phone     string
		Avatar    string

		CreatedAt time.Time
		UpdatedAt time.Time

		Friends         []*UserFriends
		Groups          []*UserGroups
		SendMessages    []*Message `gorm:"foreignKey:From"`
		ReceiveMessages []*Message `gorm:"foreignKey:To"`
	}

	UserAuth struct {
		UserName string `form:"username"`
		Password string `form:"password"`
	}

	UserFriends struct {
		UserName   string `gorm:"primaryKey"`
		FriendName string `gorm:"primaryKey"`
		Friend     *User  `json:"friend,omitempty"`
		Remark     string
		Accept     bool

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	UserGroups struct {
		UserName string `gorm:"primaryKey"`
		GroupID  uint   `gorm:"primaryKey"`
		Group    *Group `json:"group,omitempty"`
		Remark   string

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
