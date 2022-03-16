package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	var err error
	//err = os.Remove("database/sqlite.db")
	DB, err = gorm.Open(sqlite.Open("database/sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	if DB.Error != nil {
		log.Fatal(DB.Error)
	}
	log.Info("database init done")
}
