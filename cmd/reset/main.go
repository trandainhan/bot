package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"

	"gitlab.com/fiahub/bot/internal/rediswrapper"
)

var (
	coin string
)

func init() {
	flag.StringVar(&coin, "coin", "ALICE", "Coin")
}

func main() {

	ctx := context.Background()
	redisURL := os.Getenv("REDIS_URL")
	redisDBNum, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	redisClient := rediswrapper.NewRedisClient(ctx, redisURL, redisDBNum)

	log.Println("Start reset redis key")

	iter := redisClient.Client.Scan(ctx, 0, coin+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if strings.Contains(key, "price") {
			continue
		}
		log.Printf("Delete key %s", key)
		redisClient.Del(iter.Val())
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	log.Println("=============")
	log.Println("Done")
}
