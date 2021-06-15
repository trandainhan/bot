package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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
	flag.IntVar(&numWorker, "numWorker", 8, "Numer of worker with each worker control one order")
	flag.Float64Var(&quantityToGetPrice, "quantityToGetPrice", 8.0, "Quantity To Get Price")
	flag.Parse()

	// Setup chatID, chatErrorID
	var err error
	chatID, err = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		log.Panic("Missing ChatID")
	}
	chatErrorID, err = strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)
	if err != nil {
		log.Panic("Missing ChatErrorID")
	}

	// setup client
	ctx := context.Background()
	redisURL := os.Getenv("REDIS_URL")
	redisClient = rediswrapper.NewRedisClient(ctx, redisURL)
	teleClient = telegram.NewTeleBot(os.Getenv("TELE_BOT_TOKEN"))

	// get environment for login
	fiahubToken := login()
	fia = &fiahub.Fiahub{
		RedisClient: redisClient,
		Token:       fiahubToken,
	}

	// Init value in redis
	initValuesInRedis()

	// Cancel all order before starting
	log.Println("Cancel all fiahub orders before starting")
	fia.CancelAllOrder()
	time.Sleep(2 * time.Second)

	setCoinGiatotParams()

	// Set usdtvnd rate
	rate, _ := fiahub.GetUSDVNDRate()
	log.Printf("Set fiahub usdtvnd rate %v", rate)
	redisClient.Set("usdtvnd_rate", rate)

	// Set offet time
	binanceTimeDifference := binance.GetOffsetTimeUnix()

	bn = &binance.Binance{
		RedisClient:     redisClient,
		TimeDifferences: binanceTimeDifference,
	}

	// Calculate Per profit
	calculatePerProfit()
}

func initValuesInRedis() {
	log.Println("Init values in redis")
	redisClient.Set("per_cancel", 0.1/100)
	redisClient.Set("per_fee_binance", 0.075/100)
	redisClient.Set("per_profit_ask", 0.0)
	redisClient.Set("per_profit_bid", 0.0)
	for i := 1; i <= numWorker/2; i++ { // numWorker == 4
		// Set bot runable
		runnableAskKey := fmt.Sprintf("%s_ask%d_runable", coin, i)
		redisClient.Set(runnableAskKey, true)
		runnableBidKey := fmt.Sprintf("%s_bid%d_runable", coin, i)
		redisClient.Set(runnableBidKey, true)
	}
}

func setCoinGiatotParams() {
	params := fiahub.GetCoinGiaTotParams()
	validated := validateCoinGiaTotParams(params)
	if validated {
		renewCoinGiaTotParams(params)
	}
}

func login() string {
	email := os.Getenv("FIAHUB_EMAIL")
	password := os.Getenv("FIAHUB_PASSWORD")
	fiahubToken := fiahub.Login(email, password)
	log.Println("Successfully login in fiahub")
	return fiahubToken
}
