package database

import (
	"sync"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mysqlOnce sync.Once
	mysqlDB   *gorm.DB
)

func initMysql() {
	var err error
	dsn := "yzy:yuan@tcp(127.0.0.1:3306)/chat_server?charset=utf8mb4&parseTime=True&loc=Local"
	mysqlDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	if mysqlDB.Error != nil {
		logrus.Fatal(mysqlDB.Error)
	}
	logrus.Info("mysql init done")
}

func GetMysqlInstance() *gorm.DB {
	mysqlOnce.Do(func() {
		initMysql()
	})
	return mysqlDB
}
