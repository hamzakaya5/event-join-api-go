package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"), // "localhost:6379",
		Password: os.Getenv("REDIS_PASSWORD"),                             // no password set
		DB:       0,                                                       // use default DB
	})

	fmt.Print(*RedisClient)
	err := RedisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully!")

}
