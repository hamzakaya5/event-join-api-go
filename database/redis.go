package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "mypassword", // no password set
		DB:       0,            // use default DB
	})

	fmt.Print(*RedisClient)
	err := RedisClient.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("Connected to Redis successfully!")

}
