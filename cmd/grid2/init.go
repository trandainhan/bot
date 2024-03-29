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
	flag.StringVar(&fiat, "fiat", "USDT", "Fiat")
	flag.IntVar(&decimalsToRound, "decimalsToRound", 2, "Decimal to round")
	flag.IntVar(&orderQuantityToRound, "orderQuantityToRound", 1, "Order quantity decimals to round")
	flag.IntVar(&numWorker, "numWorker", 3, "Numer of worker with each worker control maximum two order")
	flag.Float64Var(&quantityToGetPrice, "quantityToGetPrice", 20, "Quantity To Get Price")
	flag.Float64Var(&orderQuantity, "orderQuantity", 1, "Order Quantity")
	flag.Float64Var(&maximumOrderUsdt, "maximumOrderUsdt", 100, "Order Quantity")
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
	chatRunableID, err = strconv.ParseInt(os.Getenv("CHAT_RUNNABLE_ID"), 10, 64)
	if err != nil {
		log.Panic("Missing chatRunableID ENV Variable")
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

	// Validate fund when starting bot
	validateFund()
}

func initValuesInRedis() {
	log.Println("Init values in redis")
	redisClient.Set(currentExchange+"_auto_mode", true, 0)
	redisClient.Set(currentExchange+coin+"_worker_runable", true, 0)

	var err error

	buyKey := coin + "_open_buy_order"
	_, err = redisClient.GetFloat64(buyKey)
	if err != nil {
		redisClient.Set(buyKey, 0, 0)
	}

	sellKey := coin + "_open_sell_order"
	_, err = redisClient.GetFloat64(sellKey)
	if err != nil {
		redisClient.Set(sellKey, 0, 0)
	}

	upTrendKey := coin + "_up_trend_percentage"
	_, err = redisClient.GetFloat64(upTrendKey)
	if err != nil {
		redisClient.Set(coin+"_up_trend_percentage", "0", 0)
	}

	upTrendFactor := coin + "_up_trend_percentage_factor"
	_, errFactor := redisClient.GetFloat64(upTrendFactor)
	if errFactor != nil {
		redisClient.Set(coin+"_up_trend_percentage_factor", 1.5, 0)
	}
}
