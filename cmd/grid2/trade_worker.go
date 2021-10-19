package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/utils"
)

func trade_worker(id string, coin string, results chan<- bool) {
	for {
		autoMode := redisClient.GetBool(currentExchange + "_auto_mode")
		coinRunable := redisClient.GetBool(currentExchange + coin + "_worker_runable")
		if !autoMode || !coinRunable {
			time.Sleep(30 * time.Second)
			continue
		}

		_, err := redisClient.GetTime(coin + "_latest_place_order_time")
		if err == nil { // mean the key is existed, Only start new worker after 5 minutes
			continue
		}

		totalOpenBuyOrders := redisClient.GetInt(coin + "_open_buy_order")
		totalOpenSellOrders := redisClient.GetInt(coin + "_open_sell_order")

		if totalOpenBuyOrders >= numWorker || totalOpenSellOrders >= numWorker {
			log.Printf("%s Ignore trade worker due to %d couple of orders is running", coin, numWorker)
			time.Sleep(1 * time.Minute)
			continue
		}

		exchangeBidPrice, err := exchanges.GetBidPriceByQuantity(coin, quantityToGetPrice)
		if err != nil {
			time.Sleep(15 * time.Second)
			continue
		}
		exchangeAskPrice, err := exchanges.GetAskPriceByQuantity(coin, quantityToGetPrice)
		if err != nil {
			time.Sleep(15 * time.Second)
			continue
		}

		jumpPrice := exchangeBidPrice * jumpPricePercentage / 100

		key := fmt.Sprintf("%s_up_trend_percentage", coin)
		upTrendPercentage, _ := redisClient.GetFloat64(key)
		upTrendPriceAdjust := jumpPrice * upTrendPercentage / 100

		// When market is up trend, upTrendPercentage > 0 => upTrendPriceAdjust > 0, Buy order price should be closed to the current market price
		// When market is down trend, upTrendPercentage < 0 => upTrendPriceAdjust < 0, Buy order price should be distance from the current market price
		finalBuyPrice := utils.RoundTo(exchangeBidPrice-jumpPrice+upTrendPriceAdjust, decimalsToRound)

		// When market is up trend, upTrendPercentage > 0 => upTrendPriceAdjust > 0, Sell order price should be distanced from the current market price
		// When market is down trend, upTrendPercentage < 0 => upTrendPriceAdjust < 0, Buy order price should be closed to the current market price
		finalSellPrice := utils.RoundTo(exchangeAskPrice+jumpPrice+upTrendPriceAdjust, decimalsToRound)

		// Contraint order maximum quanity
		maxOrderQuantity := maximumOrderUsdt / currentBidPrice
		if orderQuantity > maxOrderQuantity {
			orderQuantity = utils.RoundTo(maxOrderQuantity, orderQuantityToRound)
		}

		buyOrder, err := placeOrder(id, orderQuantity, finalBuyPrice, "buy")
		if err != nil {
			log.Printf("%s Err Can not place buy order %s", coin, err.Error())
			time.Sleep(30 * time.Second)
			continue
		}
		sellOrder, err := placeOrder(id, orderQuantity, finalSellPrice, "sell")

		if err != nil {
			log.Printf("%s Err Can not place sell order, will cancel buy order %s", coin, err.Error())
			exchangeClient.CancelOrder(coin, buyOrder.ID, buyOrder.ClientID)
			time.Sleep(30 * time.Second)
			continue
		}

		redisClient.Set(coin+"_latest_place_order_time", time.Now(), time.Duration(5)*time.Minute)
		time.Sleep(15 * time.Second)

		orderChan := make(chan *exchanges.OrderResp)
		go monitorOrder(buyOrder, orderChan)
		go monitorOrder(sellOrder, orderChan)

		// only wait for either buy or sell order is filled, then can start another loop
		<-orderChan
	}
	results <- true
}
