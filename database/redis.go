package database

import (
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var RDB *redis.Client

func init() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := RDB.Ping().Result()
	if err != nil {
		log.Error(err)
	}
}
