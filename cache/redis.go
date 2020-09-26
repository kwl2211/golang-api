package cache

import (
	"log"

	"github.com/go-redis/redis/v7"
)

var client *redis.Client

//InitRedisConfig function for cache
func InitRedisConfig(redisClient string) *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     redisClient,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		log.Println("redis error")
	}
	log.Println(pong)
	return client
}
