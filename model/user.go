package model

import "time"

type (
	User struct {
		BaseModel
		Name      string `gorm:"uniquIndex;size:255" form:"username" json:"username,omitempty"`
		Password  string `gorm:"size:255" form:"password"`
		PublicKey []byte `form:"public_key" json:"public_key,omitempty"`
		Nickname  string `gorm:"size:255" json:"nickname,omitempty"`
		Email     string `gorm:"size:255" json:"email,omitempty"`
		Phone     string `gorm:"size:255" json:"phone,omitempty"`
		Avatar    string `gorm:"size:255" json:"avatar,omitempty"`

		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`

		// Friends         []*UserFriends `json:"friends,omitempty"`
		// Groups          []*UserGroups  `json:"groups,omitempty"`
		// SendMessages    []*Message     `gorm:"foreignKey:From" json:"send_messages,omitempty"`
		// ReceiveMessages []*Message     `gorm:"foreignKey:To" json:"receive_messages,omitempty"`
	}

	// UserAuth struct {
	// 	UserName string `form:"username"`
	// 	Password string `gorm:"size:255" form:"password"`
	// 	// PublicKeyHash string `gorm:"size:255" form:"public_key_hash"`
	// }

	UserFriend struct {
		UserName   string `gorm:"primaryKey" json:"user_name,omitempty"`
		FriendName string `gorm:"primaryKey;" json:"friend_name,omitempty"`
		Friend     *User  `json:"friend,omitempty"`
		Remark     string `gorm:"size:255" json:"remark,omitempty"`
		Accept     bool   `gorm:"default:false" json:"accept,omitempty"`

		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserGroup struct {
		UserName      string   `gorm:"primaryKey" json:"user_name,omitempty"`
		GroupID       uint     `gorm:"primaryKey" json:"group_id,omitempty"`
		Group         *Group   `json:"group,omitempty"`
		LastMessageID uint64   `json:"last_message_id,omitempty"`
		LastMessage   *Message `json:"last_message,omitempty"`
		Remark        string   `gorm:"size:255" json:"remark,omitempty"`

		CreatedAt time.Time `json:"created_at,omitempty"`
		UpdatedAt time.Time `json:"updated_at,omitempty"`
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
