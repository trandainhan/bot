package main

import (
	"context"
	"flag"
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

	// Cancel all order before starting
	fiahub.CancelAllOrder(fiahubToken)
	time.Sleep(2 * time.Second)

	setCoinGiatotParams()

	// Init value in redis
	initValuesInRedis()

	// Calculate Per profit
	calculatePerProfit()

	rate, _ := fiahub.GetUSDVNDRate()
	redisClient.Set("usdtvnd_rate", rate)

	offset := binance.GetOffsetTimeUnix()
	redisClient.Set("local_binance_time_difference", offset)
}

func initValuesInRedis() {
	redisClient.Set("runable", true)
	redisClient.Set("per_fee_binance", 0.075/100)
	redisClient.Set("per_profit_ask", 0.0)
	redisClient.Set("per_profit_bid", 0.0)
}

func setCoinGiatotParams() {
	// Coin gia tot params()
	params := fiahub.GetCoinGiaTotParams()
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
}
