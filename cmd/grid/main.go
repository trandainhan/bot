package main

import (
	"fmt"
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
	quantityToGetPrice    float64
	numWorker             int
	defaultQuantity       float64
	jumpPricePercentage   float64
	upTrendPercentage     float64
	chatID                int64
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
		exchangeClient.CancelAllOrder(coin)
		log.Println("==================")
		log.Println("Finish trading bot")
		os.Exit(0)
	}()

	log.Println("=================")
	log.Println("Start trading bot")

	results := make(chan bool, numWorker)
	var id string

	for i := 1; i <= numWorker/2; i++ {
		// buy
		id = fmt.Sprintf("buy%d", i)
		go buy_worker(id, coin, i, results)

		// sell
		id = fmt.Sprintf("sell%d", i)
		go sell_worker(id, coin, i, results)
	}

	go func() {
		for {
			time.Sleep(time.Duration(60) * time.Second)
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
			time.Sleep(time.Duration(15) * time.Minute)
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
	log.Println("==================")
	log.Println("Finish trading bot")
}
