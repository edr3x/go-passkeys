package redisclient

import (
	"cmp"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

func RedisConnect() (client *redis.Client) {
	addr := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.Atoi(cmp.Or(os.Getenv("REDIS_DB"), "0"))

	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping().Err(); err != nil {
		log.Println("Error connecting to redis: ", err)
	} else {
		log.Println("Redis Connection Successful")
	}

	return
}
