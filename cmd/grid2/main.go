package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
	"gitlab.com/fiahub/bot/internal/telegram"
)

var (
	coin                  string
	decimalsToRound       int
	orderQuantityToRound  int
	quantityToGetPrice    float64
	numWorker             int
	orderQuantity         float64
	maximumOrderUsdt      float64
	jumpPricePercentage   float64
	upTrendPercentage     float64
	chatID                int64
	chatRunableID         int64
	chatProfitID          int64
	chatErrorID           int64
	binanceTimeDifference int64
	redisClient           *rediswrapper.MyRedis
	teleClient            *telegram.TeleBot
	exchangeClient        *exchanges.ExchangeClient
	currentExchange       string
	currentAskPrice       float64
	currentBidPrice       float64
)

func main() {
	// run init.go

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		log.Println("Recieve OS signal, CancelAllOrder and stop bot")
		_, err := exchangeClient.CancelAllOrder(coin)
		if err != nil {
			log.Printf("Err CancelAllOrder: %s", err.Error())
		} else {
			resetOpenBuySellOrders()
		}
		log.Println("==================")
		log.Println("Finish trading bot")
		os.Exit(0)
	}()

	log.Println("=================")
	log.Println("Start trading bot")

	results := make(chan bool, numWorker)

	go trade_worker("worker1", coin, results)

	go func() {
		for {
			time.Sleep(time.Duration(1) * time.Minute)
			go updateCurrentAskPrice()
			go updateCurrentBidPrice()
		}
	}()

	go func() {
		for {
			period, err := strconv.Atoi(os.Getenv("VALIDATE_FUND_PERIOD"))
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Duration(period) * time.Second)
			validateFund()
		}
	}()

	go func() {
		for {
			time.Sleep(time.Duration(3) * time.Minute)
			adjustUpTrendPercentage()
		}
	}()

	for i := 0; i < numWorker; i++ {
		<-results
	}
	_, err := exchangeClient.CancelAllOrder(coin)
	if err != nil {
		log.Printf("Err CancelAllOrder: %s", err.Error())
	}
	resetOpenBuySellOrders()
	log.Println("==================")
	log.Println("Finish trading bot")
}
