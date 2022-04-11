package database

import (
	"sync"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var (
	redisOnce sync.Once
	redisDB   *redis.Client
)

func initRedis() {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := redisDB.Ping().Result()
	if err != nil {
		log.Error(err)
	}
	log.Info("redis init done")
}

func GetRedisInstance() *redis.Client {
	redisOnce.Do(func() {
		initRedis()
	})
	return redisDB
}
