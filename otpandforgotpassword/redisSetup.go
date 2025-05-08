package otpandforgotpassword

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error to load the env file")
	}
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	add := host + ":" + port
	RedisClient = redis.NewClient(&redis.Options{
		Addr: add,
		DB:   0,
	})

	_, err = RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Redis client initialized successfully")
}
