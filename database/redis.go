package database

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping().Result()
	if err != nil {
		log.Error(err)
	}
	log.Info("redis init done")
}
