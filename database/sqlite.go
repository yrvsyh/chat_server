package database

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	sqliteOnce sync.Once
	sqliteDB   *gorm.DB
)

func initSqlite() {
	var err error
	//err = os.Remove("database/sqlite.db")
	sqliteDB, err = gorm.Open(sqlite.Open("database/sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}
	if sqliteDB.Error != nil {
		log.Fatal(sqliteDB.Error)
	}
	log.Info("sqlite init done")
}

func getSqliteInstance() *gorm.DB {
	sqliteOnce.Do(func() {
		initSqlite()
	})
	return sqliteDB
}
