package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
	"gitlab.com/fiahub/bot/internal/telegram"
)

func init() {
	flag.StringVar(&coin, "coin", "ALICE", "Coin")
	flag.IntVar(&decimalsToRound, "decimalsToRound", 2, "Decimal to round")
	flag.IntVar(&numWorker, "numWorker", 2, "Numer of worker with each worker control one order")
	flag.Float64Var(&quantityToGetPrice, "quantityToGetPrice", 20, "Quantity To Get Price")
	flag.Float64Var(&defaultQuantity, "defaultQuantity", 1, "Order Quantity")
	flag.Float64Var(&jumpPricePercentage, "jumpPricePercentage", 1, "Price Jump Percentage to cancel order")
	flag.Parse()

	// get currentExchange
	currentExchange = os.Getenv("EXCHANGE_CLIENT")

	// Setup chatID, chatProfitID, chatErrorID
	var err error
	chatID, err = strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	if err != nil {
		log.Panic("Missing ChatID ENV Variable")
	}
	chatProfitID, err = strconv.ParseInt(os.Getenv("CHAT_PROFIT_ID"), 10, 64)
	if err != nil {
		log.Panic("Missing ChatProfitID ENV Variable")
	}
	chatErrorID, err = strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)
	if err != nil {
		log.Panic("Missing ChatErrorID ENV Variable")
	}

	// setup client
	ctx := context.Background()
	redisURL := os.Getenv("REDIS_URL")
	redisDBNum, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	redisClient = rediswrapper.NewRedisClient(ctx, redisURL, redisDBNum)
	teleClient = telegram.NewTeleBot(os.Getenv("TELE_BOT_TOKEN"))

	// Init value in redis
	initValuesInRedis()

	// Init current Price
	updateCurrentAskPrice()
	updateCurrentBidPrice()

	// Set offet time
	binanceTimeDifference := binance.GetOffsetTimeUnix()
	bn := &binance.Binance{
		TimeDifferences: binanceTimeDifference,
	}
	ftxClient := ftx.FtxClient{}
	exchangeClient = &exchanges.ExchangeClient{
		Ftx: &ftxClient,
		Bn:  bn,
	}
}

func initValuesInRedis() {
	log.Println("Init values in redis")
	redisClient.Set(currentExchange+"_auto_mode", true, 0)
	redisClient.Set(coin+"_buy_worker_runable", true, 0)
	redisClient.Set(coin+"_sell_worker_runable", true, 0)
	upTrendKey := coin + "_up_trend_percentage"
	_, err := redisClient.GetFloat64(upTrendKey)
	if err != nil {
		redisClient.Set(coin+"_up_trend_percentage", "0", 0)
	}
}
