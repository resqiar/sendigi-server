package configs

import (
	"os"

	"github.com/gofiber/storage/redis/v3"
)

var RedisStore *redis.Storage

func InitRedis() {
	URL := os.Getenv("RDSN")
	RedisStore = redis.New(redis.Config{
		URL: URL,
	})
}
