package main

import (
	"log"

	"gitlab.com/fiahub/bot/internal/binance"
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
)

func main() {
	// perFeeBinance := redisClient.Get("per_fee_binance").(float64) // 0.075 / 100
	// perProfitAsk := redisClient.Get("per_profit_ask").(float64)
	// perProfitBid := redisClient.Get("per_profit_bid").(float64)

	marketParam := coin + "USDT"
	bidPriceByQuantity, askPriceByQuantity := binance.GetPriceByQuantity(marketParam, quantityToGetPrice)

	log.Println(bidPriceByQuantity)
	// log.Println(perProfitBid)

	// go renew params

	// Ask trading
	var perProfitStep float64
	perProfitStep = 1.0
	go ask_worker("riki1", coin, askPriceByQuantity, perProfitStep)

	perProfitStep = 2.0
	go ask_worker("riki1", coin, askPriceByQuantity, perProfitStep)

	perProfitStep = 3.0
	go ask_worker("riki1", coin, askPriceByQuantity, perProfitStep)

	perProfitStep = 4.0
	go ask_worker("riki1", coin, askPriceByQuantity, perProfitStep)

	// go bid_worker
	perProfitStep = 1.0
	go bid_worker("riki1", coin, bidPriceByQuantity, perProfitStep)

	perProfitStep = 2.0
	go bid_worker("riki1", coin, bidPriceByQuantity, perProfitStep)

	perProfitStep = 3.0
	go bid_worker("riki1", coin, bidPriceByQuantity, perProfitStep)

	perProfitStep = 4.0
	go bid_worker("riki1", coin, bidPriceByQuantity, perProfitStep)
}
