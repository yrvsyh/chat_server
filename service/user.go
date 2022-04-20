package service

import (
	"chat_server/config"
	"chat_server/model"
	"crypto/sha256"
	"encoding/hex"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserService struct{}

func (UserService) hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password + config.PasswordSalt))

	return hex.EncodeToString(hash.Sum(nil))
}
func (UserService) verifyPassword(formPassword string, dbPassword string) bool {
	hashPassword := userService.hashPassword(formPassword)

	log.WithFields(log.Fields{
		"password":     formPassword,
		"dbPassword":   dbPassword,
		"hashPassword": hashPassword,
	}).Info("PASSWORD CHECK")

	return dbPassword == hashPassword
}

func (UserService) Register(username string, password string, publicKey string) error {
	password = userService.hashPassword(password)
	user := &model.User{
		Username:  username,
		Password:  password,
		PublicKey: publicKey,
	}
	err := db.Create(user).Error
	return errors.Wrap(err, "create data error")
}

func (UserService) Login(username string, password string) (uint32, string, string, error) {
	user, err := userService.GetUserByUsername(username)

	if err != nil {
		return 0, "", "", errors.Wrap(err, "no such user")
	}

	if !userService.verifyPassword(password, user.Password) {
		return 0, "", "", errors.Wrap(err, "password error")
	}

	return user.ID, user.Username, user.PublicKey, nil
}

func (UserService) SearchUserByName(name string) ([]model.User, error) {
	var users []model.User
	err := db.Select("id", "username", "nickname", "avatar").Where("username LIKE ? or nickname LIKE ?", "%"+name+"%", "%"+name+"%").Find(&users).Error
	return users, err
}

func (UserService) GetUserByID(id uint32) (*model.User, error) {
	user := &model.User{}
	err := db.First(user, id).Error
	return user, err
}

func (UserService) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := db.First(user, "username = ?", username).Error
	return user, err
}

func (UserService) UpdateUser(user *model.User) error {
	return db.Model(user).Updates(user).Error
}

func (UserService) UpdateUserAvatarByID(id uint32, path string) error {
	user := &model.User{}
	user.ID = id
	return db.Model(user).Update("avatar", path).Error
}

func (UserService) GetUserFriendsByID(id uint32) ([]model.UserFriend, error) {
	friends := []model.UserFriend{}
	err := db.Joins("JOIN user_friends as t on user_friends.user_id = t.friend_id and user_friends.friend_id = t.user_id and user_friends.accept = true and t.accept = true").Where("user_friends.user_id = ?", id).Find(&friends).Error
	return friends, err
}

func (UserService) GetUserFriendsDetailByID(id uint32) ([]model.UserFriend, error) {
	friends := []model.UserFriend{}
	err := db.Preload("Friend").Joins("JOIN user_friends as t on user_friends.user_id = t.friend_id and user_friends.friend_id = t.user_id and user_friends.accept = true and t.accept = true").Where("user_friends.user_id = ?", id).Find(&friends).Error
	return friends, err
}

// func (UserService) GetUserFriendNameList(name string) ([]string, error) {
// 	var ret []string
// 	var userFriends []model.UserFriends
// 	err := db.Model(&model.UserFriends{}).Where("user_name = ?", name).Find(&userFriends).Error
// 	for _, userFriend := range userFriends {
// 		ret = append(ret, userFriend.FriendName)
// 	}
// 	return ret, err
// }

func (UserService) GetUserFriendIDSet(id uint32) (map[uint32]struct{}, error) {
	ret := make(map[uint32]struct{})

	var userFriends []model.UserFriend
	err := db.Select("friend_id").Where("user_id = ?", id).Find(&userFriends).Error
	for _, userFriend := range userFriends {
		ret[userFriend.FriendID] = struct{}{}
	}

	return ret, err
}

func (UserService) GetUserGroupIDSet(id uint32) (map[uint32]struct{}, error) {
	ret := make(map[uint32]struct{})

	var groupUsers []model.GroupUser
	err := db.Select("group_id").Where("user_id = ?", id).Find(&groupUsers).Error
	for _, groupUser := range groupUsers {
		ret[groupUser.GroupID] = struct{}{}
	}

	return ret, err
}

func (UserService) GetUserFriendDetailByFriendID(id uint32, friendID uint32) (*model.UserFriend, error) {
	userFriend := &model.UserFriend{}
	err := db.Model(userFriend).First(userFriend, "user_id= ? and friend_id = ?", id, friendID).Error
	return userFriend, err
}

func (UserService) UpdateUserFriend(userFriend *model.UserFriend) error {
	return db.Save(userFriend).Error
}

func (UserService) AddUserFriend(id uint32, friendID uint32) error {
	return db.Transaction(func(tx *gorm.DB) error {
		userFriend1 := &model.UserFriend{UserID: id, FriendID: friendID, Accept: true}
		userFriend2 := &model.UserFriend{UserID: friendID, FriendID: id, Accept: false}

		if err := tx.Create(userFriend1).Error; err != nil {
			return err
		}
		if err := tx.Create(userFriend2).Error; err != nil {
			return err
		}

		return nil
	})
}

func (UserService) AcceptUserFriend(id uint32, friendID uint32) error {
	peerInfo := &model.UserFriend{UserID: friendID, FriendID: id, Accept: true}
	if err := db.Where(peerInfo).First(peerInfo).Error; err != nil {
		// 对方未发起请求
		log.Error(err)
		return errors.New("对方未发起请求")
	}

	userFriend := &model.UserFriend{UserID: id, FriendID: friendID}
	return db.Model(userFriend).Update("accept", true).Error
}

func (UserService) DeleteUserFriend(id uint32, friendID uint32) error {
	return db.Transaction(func(tx *gorm.DB) error {
		userFriend := &model.UserFriend{UserID: id, FriendID: friendID}
		return tx.Delete(userFriend).Error
	})
}
