package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
	"gitlab.com/fiahub/bot/internal/telegram"
)

var (
	coin                  string
	minPrice              float64
	maxPrice              float64
	decimalsToRound       int
	defaultSleepSeconds   int64
	quantityToGetPrice    float64
	numWorker             int
	chatID                int64
	chatErrorID           int64
	binanceTimeDifference int64
	lastestCancelAllTime  int64
	redisClient           *rediswrapper.MyRedis
	teleClient            *telegram.TeleBot
	fia                   *fiahub.Fiahub
	bn                    *binance.Binance
)

func main() {
	log.Println("=================")
	log.Println("Start trading bot")

	results := make(chan bool, numWorker)

	// Ask trading
	var perProfitStep float64

	perProfitStep = 1.0
	log.Println("Start ask worker riki1")
	go ask_worker("riki1", coin, perProfitStep, results)

	perProfitStep = 2.0
	go ask_worker("riki2", coin, perProfitStep, results)
	//
	// perProfitStep = 3.0
	// go ask_worker("riki3", coin, askPriceByQuantity, perProfitStep)
	//
	// perProfitStep = 4.0
	// go ask_worker("riki4", coin, askPriceByQuantity, perProfitStep)

	// go bid_worker

	perProfitStep = 1.0
	log.Println("Start bid worker rikiatb1")
	go bid_worker("rikiatb1", coin, perProfitStep, results)

	perProfitStep = 2.0
	go bid_worker("rikiatb2", coin, perProfitStep, results)
	//
	// perProfitStep = 3.0
	// go bid_worker("rikiatb3", coin, bidPriceByQuantity, perProfitStep)
	//
	// perProfitStep = 4.0
	// go bid_worker("rikiatb4", coin, bidPriceByQuantity, perProfitStep)

	// go renew params, env, token
	go func() {
		for {
			period, err := strconv.Atoi(os.Getenv("SET_COINGIATOT_PERIOD"))
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Duration(period) * time.Second)
			setCoinGiatotParams()
		}
	}()

	go func() {
		for {
			period, err := strconv.Atoi(os.Getenv("CALCULATE_PER_PROFIT_PERIOD"))
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Duration(period) * time.Second)
			calculatePerProfit()
		}
	}()

	go func() {
		for {
			period, err := strconv.Atoi(os.Getenv("RESET_ALL_ORDER_PERIOD"))
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Duration(period) * time.Second)
			log.Printf("Cancal all order after %d Second", period)
			fia.CancelAllOrder()
		}
	}()

	go func() {
		for {
			period, err := strconv.Atoi(os.Getenv("RESET_FIAHUB_TOKEN_PERIOD"))
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Duration(period) * time.Second)
			log.Printf("Reset token after %d seconds", period)
			token := login()
			fia.SetToken(token)
		}
	}()

	for i := 0; i < numWorker; i++ {
		<-results
	}
	log.Println("==================")
	log.Println("Finish trading bot")
}
