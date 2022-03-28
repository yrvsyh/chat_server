package model

import "time"

type (
	User struct {
		Name      string    `gorm:"primaryKey" form:"username" json:"name,omitempty"`
		UserAuth  *UserAuth `json:"user_auth,omitempty"`
		PublicKey []byte    `form:"public_key" json:"public_key,omitempty"`
		Nickname  string    `json:"nickname,omitempty"`
		Email     string    `json:"email,omitempty"`
		Phone     string    `json:"phone,omitempty"`
		Avatar    string    `json:"avatar,omitempty"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`

		Friends         []*UserFriends `json:"friends,omitempty"`
		Groups          []*UserGroups  `json:"groups,omitempty"`
		SendMessages    []*Message     `gorm:"foreignKey:From" json:"send_messages,omitempty"`
		ReceiveMessages []*Message     `gorm:"foreignKey:To" json:"receive_messages,omitempty"`
	}

	UserAuth struct {
		UserName string `form:"username"`
		Password string `form:"password"`
	}

	UserFriends struct {
		UserName   string `gorm:"primaryKey" json:"user_name,omitempty"`
		FriendName string `gorm:"primaryKey" json:"friend_name,omitempty"`
		Friend     *User  `json:"friend,omitempty" json:"friend,omitempty"`
		Remark     string `json:"remark,omitempty"`
		Accept     bool   `json:"accept,omitempty"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserGroups struct {
		UserName string `gorm:"primaryKey" json:"user_name,omitempty"`
		GroupID  uint   `gorm:"primaryKey" json:"group_id,omitempty"`
		Group    *Group `json:"group,omitempty"`
		Remark   string `json:"remark,omitempty"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
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
