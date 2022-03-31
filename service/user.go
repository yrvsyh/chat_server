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
func (u UserService) verifyPassword(formPassword string, dbPassword string) bool {
	hashPassword := u.hashPassword(formPassword)

	log.WithFields(log.Fields{
		"password":     formPassword,
		"dbPassword":   dbPassword,
		"hashPassword": hashPassword,
	}).Info("PASSWORD CHECK")

	return dbPassword == hashPassword
}

func (UserService) verifyPublicKey(publicKey []byte, dbKey []byte) bool {
	return string(publicKey) == string(dbKey)
}

func (u UserService) Register(username string, password string, publicKey []byte) error {
	password = u.hashPassword(password)
	user := &model.User{
		Username:  username,
		Password:  password,
		PublicKey: publicKey,
	}
	err := db.Create(user).Error
	return errors.Wrap(err, "create data error")
}

func (u UserService) Login(username string, password string, publicKey []byte) (uint32, string, error) {
	user, err := u.GetUserByUsername(username)

	if err != nil {
		return user.ID, user.Username, errors.Wrap(err, "no such user")
	}

	if !u.verifyPassword(password, user.Password) {
		return user.ID, user.Username, errors.Wrap(err, "password error")
	}

	// if !u.verifyPublicKey(publicKey, user.PublicKey) {
	// 	return user.ID, user.Username, errors.Wrap(err, "public key invalid")
	// }

	return user.ID, user.Username, nil
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
	err := db.Where("user_id = ? and accept = ?", id, true).Find(&friends).Error
	return friends, err
}

func (UserService) GetUserFriendsDetailByID(id uint32) ([]model.UserFriend, error) {
	friends := []model.UserFriend{}
	err := db.Preload("Friend").Where("user_id = ? and accept = ?", id, true).Find(&friends).Error
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

func (UserService) GetUserFriendIDSet(id uint32) (map[uint32]bool, error) {
	ret := make(map[uint32]bool)

	var userFriends []model.UserFriend
	err := db.Select("friend_id").Where("user_id = ?", id).Find(&userFriends).Error
	for _, userFriend := range userFriends {
		ret[userFriend.FriendID] = true
	}

	return ret, err
}

func (UserService) GetUserGroupIDSet(id uint32) (map[uint32]bool, error) {
	ret := make(map[uint32]bool)

	var groupUsers []model.GroupUser
	err := db.Select("group_id").Where("user_id = ?", id).Find(&groupUsers).Error
	for _, groupUser := range groupUsers {
		ret[groupUser.GroupID] = true
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

// func (s UserService) CheckUserExistByName(name string) bool {
// 	_, err := s.GetUserByName(name)
// 	return err == nil
// }

func (UserService) AddUserFriend(id uint32, friendID uint32) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 双向好友, 默认为未同意
		userFriend1 := &model.UserFriend{UserID: id, FriendID: friendID, Accept: false}
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
	return db.Transaction(func(tx *gorm.DB) error {
		userFriend1 := &model.UserFriend{UserID: id, FriendID: friendID}
		userFriend2 := &model.UserFriend{UserID: friendID, FriendID: id}

		if err := tx.Model(userFriend1).Update("accept", true).Error; err != nil {
			return err
		}
		if err := tx.Model(userFriend2).Update("accept", true).Error; err != nil {
			return err
		}

		return nil
	})
}
