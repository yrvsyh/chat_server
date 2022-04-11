package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var TestDB *gorm.DB

func InitSqlite() {
	var err error
	//err = os.Remove("database/sqlite.db")
	TestDB, err = gorm.Open(sqlite.Open("database/sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	if TestDB.Error != nil {
		log.Fatal(TestDB.Error)
	}
	log.Info("sqlite init done")
}
