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
	auth := &model.UserAuth{}
	err := c.Bind(auth)
	if err != nil {
		utils.Error(c, -1, "信息不完善")
		c.Abort()
		return
	}

	_, err = service.GetUserAuthByName(auth.UserName)
	if err == nil {
		utils.Error(c, -1, "用户已存在")
		c.Abort()
		return
	}

	auth.Password = hashPassword(auth.Password)

	if err := service.RegisterUser(auth); err != nil {
		utils.Error(c, -1, "注册失败")
		c.Abort()
		return
	}

	utils.Success(c, "注册成功")
}

func Login(c *gin.Context) {
	auth := &model.UserAuth{}
	err := c.Bind(auth)

	if err != nil {
		utils.Error(c, -1, "信息不完善")
		c.Abort()
		return
	}

	log.Info(auth)

	dbUser, err := service.GetUserAuthByName(auth.UserName)
	if err != nil {
		utils.Error(c, -1, "用户不存在")
		c.Abort()
		return
	}

	if !verifyPassword(auth.Password, dbUser.UserAuth.Password) {
		utils.Error(c, -1, "密码错误")
		c.Abort()
		return
	}

	session := middleware.GetAuthSession(c)
	session.Values["name"] = dbUser.Name

	err = session.Save(c.Request, c.Writer)
	if err != nil {
		utils.Error(c, -1, "服务器错误")
		c.Abort()
		return
	}

	utils.Success(c, "登陆成功")
}

func GetLoginUserName(c *gin.Context) string {
	session := middleware.GetAuthSession(c)
	return session.Values[config.SessionUserKey].(string)
}

func Logout(c *gin.Context) {
	session := middleware.GetAuthSession(c)
	delete(session.Values, config.SessionUserKey)
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		utils.Error(c, -1, "服务器错误")
		c.Abort()
		return
	}

	utils.Success(c, "注销成功")
}

func verifyPassword(formPassword string, dbPassword string) bool {
	hashPassword := hashPassword(formPassword)
	log.WithFields(log.Fields{
		"password":     formPassword,
		"dbPassword":   dbPassword,
		"hashPassword": hashPassword,
	}).Info("PASSWORD CHECK")
	return dbPassword == hashPassword
}

func hashPassword(password string) string {
	sha256sum := sha256.Sum256([]byte(password + config.PasswordSalt))
	return hex.EncodeToString(sha256sum[:])
}
