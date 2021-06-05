package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
	"gitlab.com/fiahub/bot/internal/telegram"
)

func init() {
	flag.StringVar(&coin, "coin", "btc", "Coin")
	flag.Float64Var(&minPrice, "minPrice", 0, "Min Price")
	flag.Float64Var(&maxPrice, "maxPrice", 1000000, "Min Price")
	flag.Int64Var(&defaultSleepSeconds, "defaultSleepSeconds", 18, "Sleep in second then restart")
	flag.IntVar(&decimalsToRound, "decimalsToRound", 3, "Decimal to round")
	flag.Float64Var(&quantityToGetPrice, "quantityToGetPrice", 8.0, "Quantity To Get Price")
	flag.Parse()

	// setup client
	ctx := context.Background()
	redisURL := os.Getenv("REDIS_URL")
	redisClient = rediswrapper.NewRedisClient(ctx, redisURL)
	teleClient = telegram.NewTeleBot(os.Getenv("TELE_BOT_TOKEN"))

	// get environment for login
	login()

	// Init value in redis
	initValuesInRedis()

	// Cancel all order before starting
	fia = fiahub.Fiahub{
		RedisClient: redisClient,
	}
	fia.CancelAllOrder(fiahubToken)
	time.Sleep(2 * time.Second)

	setCoinGiatotParams()

	// Set usdtvnd rate
	rate, _ := fiahub.GetUSDVNDRate()
	log.Printf("Set fiahub usdtvnd rate %v", rate)
	redisClient.Set("usdtvnd_rate", rate)

	// Set offet time
	offset := binance.GetOffsetTimeUnix()
	redisClient.Set("local_binance_time_difference", offset)

	// Calculate Per profit
	calculatePerProfit()
}

func initValuesInRedis() {
	log.Println("Init values in redis")
	redisClient.Set("runable", true)
	redisClient.Set("per_cancel", 0.1)
	redisClient.Set("per_fee_binance", 0.075/100)
	redisClient.Set("per_profit_ask", 0.0)
	redisClient.Set("per_profit_bid", 0.0)
}

func setCoinGiatotParams() {
	params := fiahub.GetCoinGiaTotParams()
	log.Printf("setCoinGiatotParams %v", params)
	validated := validateCoinGiaTotParams(params)
	if validated {
		renewCoinGiaTotParams(params)
	}
}

func login() {
	email := os.Getenv("FIAHUB_EMAIL")
	password := os.Getenv("FIAHUB_PASSWORD")
	fiahubToken = fiahub.Login(email, password)
	redisClient.Set("fiahub_token", fiahubToken)
	log.Println("Successfully login in fiahub")
}
