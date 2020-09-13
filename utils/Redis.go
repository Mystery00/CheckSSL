package utils

import (
	"time"

	"github.com/go-redis/redis/v8"

	"CheckSSL/config"
)

var DomainKey = "check:ssl:domains"
var ExpireTime = 10 * time.Hour

func RedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.EnvConfig.RedisHost,
		Password: config.EnvConfig.RedisPassword,
		DB:       config.EnvConfig.RedisDatabase,
	})
}
