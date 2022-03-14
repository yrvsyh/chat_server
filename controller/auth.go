package controller

import (
	"chat_server/config"
	"chat_server/middleware"
	"chat_server/model"
	"chat_server/service"
	"chat_server/utils"
	"crypto/sha256"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Register(c *gin.Context) {
	user := &model.User{}
	c.Bind(user)

	if dbUser := service.GetUserById(user.Username); dbUser != nil {
		utils.Error(c, -1, "用户已存在")
		c.Abort()
		return
	}

	user.Password = hashPasswd(user.Password)

	if err := service.InsertUser(user); err != nil {
		utils.Error(c, -1, "注册失败")
		c.Abort()
		return
	}

	utils.Success(c, "注册成功")
}

func Login(c *gin.Context) {
	tokenString := middleware.GetToken(c)
	if middleware.VerifyToken(tokenString) {
		// utils.Error(c, utils.ERR_ALREADY_LOGIN)
		return
	}

	user := &model.User{}
	c.Bind(user)

	dbUser := service.GetUserById(user.Username)
	if dbUser == nil {
		utils.Error(c, -1, "用户不存在")
		c.Abort()
		return
	}

	if !verifyPassword(user.Password, dbUser.Password) {
		utils.Error(c, -1, "密码错误")
		c.Abort()
		return
	}

	// tokenString, err := middleware.GenToken(user.Name, user.Role)
	tokenString, err := middleware.GenToken(user.Username)
	if err != nil {
		utils.Error(c, -1, "Token生成失败")
		c.Abort()
		return
	}

	utils.SuccessWithData(c, "登陆成功", gin.H{"token": tokenString})
}

func Refresh(c *gin.Context) {
	tokenString := middleware.GetToken(c)
	newTokenString, err := middleware.RefreshToken(tokenString)
	if err != nil {
		utils.Error(c, -1, "Token刷新失败")
		return
	}
	utils.SuccessWithData(c, "token刷新成功", gin.H{"token": newTokenString})
}

func Logout(c *gin.Context) {
	tokenString := middleware.GetToken(c)
	if err := middleware.DelToken(tokenString); err != nil {
		utils.Error(c, -1, "Token删除失败")
		return
	}
	utils.Success(c, "注销成功")
}

func verifyPassword(formPasswd string, dbPasswd string) bool {
	hashPasswd := hashPasswd(formPasswd)
	log.WithFields(log.Fields{
		"password":   formPasswd,
		"dbPasswd":   dbPasswd,
		"hashPasswd": hashPasswd,
	}).Info("PASS")
	return dbPasswd == hashPasswd
}

func hashPasswd(password string) string {
	sha256sum := sha256.Sum256([]byte(password + config.PasswordSalt))
	return hex.EncodeToString(sha256sum[:])
}
