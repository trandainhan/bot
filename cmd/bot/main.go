package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
	"gitlab.com/fiahub/bot/internal/exchanges"
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
	chatProfitID          int64
	chatErrorID           int64
	binanceTimeDifference int64
	lastestCancelAllTime  int64
	redisClient           *rediswrapper.MyRedis
	teleClient            *telegram.TeleBot
	fia                   *fiahub.Fiahub
	db                    *pg.DB
	exchangeClient        *exchanges.ExchangeClient
)

func main() {
	// run init.go

	log.Println("=================")
	log.Println("Start trading bot")

	results := make(chan bool, numWorker)
	var id string
	var perProfitStep float64

	for i := 1; i <= numWorker/2; i++ {
		// Ask trading
		perProfitStep = float64(i)

		id = fmt.Sprintf("ask%d", i)
		go ask_worker(id, coin, perProfitStep, i, results)

		// Bid trading
		id = fmt.Sprintf("bid%d", i)
		go bid_worker(id, coin, perProfitStep, i, results)
	}

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
			validatePerProfit()
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

	go func() {
		for {
			period, err := strconv.Atoi(os.Getenv("RESET_RATES_PERIOD"))
			if err != nil {
				log.Panic(err)
			}
			time.Sleep(time.Duration(period) * time.Second)
			getRates()
		}
	}()

	for i := 0; i < numWorker; i++ {
		<-results
	}
	defer db.Close()
	log.Println("==================")
	log.Println("Finish trading bot")
}
