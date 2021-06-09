package main

import (
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
	"gitlab.com/fiahub/bot/internal/telegram"
)

var (
	coin                string
	minPrice            float64
	maxPrice            float64
	redisClient         *rediswrapper.MyRedis
	teleClient          *telegram.TeleBot
	decimalsToRound     int
	defaultSleepSeconds int64
	quantityToGetPrice  float64
	fiahubToken         string
	fia                 fiahub.Fiahub
	bn                  binance.Binance
)

func main() {
	log.Println("Start trading bot")
	marketParam := coin + "USDT"
	bidPriceByQuantity, askPriceByQuantity := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)
	log.Printf("bidPriceByQuantity: %v", bidPriceByQuantity)
	log.Printf("askPriceByQuantity: %v", askPriceByQuantity)

	const numWorker = 2
	results := make(chan bool, numWorker)

	// Ask trading
	var perProfitStep float64

	perProfitStep = 1.0
	log.Println("Start ask worker riki1")
	go ask_worker("riki1", coin, askPriceByQuantity, perProfitStep, results)

	// perProfitStep = 2.0
	// go ask_worker("riki2", coin, askPriceByQuantity, perProfitStep)
	//
	// perProfitStep = 3.0
	// go ask_worker("riki3", coin, askPriceByQuantity, perProfitStep)
	//
	// perProfitStep = 4.0
	// go ask_worker("riki4", coin, askPriceByQuantity, perProfitStep)

	// go bid_worker

	perProfitStep = 1.0
	log.Println("Start bid worker rikiatb1")
	go bid_worker("rikiatb1", coin, bidPriceByQuantity, perProfitStep, results)

	// perProfitStep = 2.0
	// go bid_worker("rikiatb2", coin, bidPriceByQuantity, perProfitStep)
	//
	// perProfitStep = 3.0
	// go bid_worker("rikiatb3", coin, bidPriceByQuantity, perProfitStep)
	//
	// perProfitStep = 4.0
	// go bid_worker("rikiatb4", coin, bidPriceByQuantity, perProfitStep)

	// go renew params, env, token
	go func() {
		for {
			time.Sleep(30 * time.Second)
			setCoinGiatotParams()
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			calculatePerProfit()
		}
	}()

	go func() {
		for {
			time.Sleep(176 * time.Second)
			fia.CancelAllOrder(fiahubToken)
		}
	}()

	go func() {
		for {
			time.Sleep(3600 * time.Second)
			login()
		}
	}()

	for i := 0; i < numWorker; i++ {
		<-results
	}
	log.Println("Finish trading bot")
}
