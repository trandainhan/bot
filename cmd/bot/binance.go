package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
)

func placeBinanceOrder(botID string, newSellQuantity, priceB, priceF float64, side string) {
	var orderDetails *binance.OrderDetailsResp
	var err error
	if side == "BUY" {
		orderDetails, err = bn.BuyLimit(coin+"USDT", priceB, newSellQuantity)
	} else if side == "SELL" {
		orderDetails, err = bn.SellLimit(coin+"USDT", priceB, newSellQuantity)
	}
	if err != nil {
		totalUSDT := newSellQuantity * priceB
		notifyBinanceFailOrder(botID, newSellQuantity, totalUSDT, side, err)
		return
	}

	binanceOrderID := orderDetails.OrderID
	origClientOrderID := orderDetails.ClientOrderID
	if binanceOrderID != 0 {
		text := fmt.Sprintf("%s %s Take profit Binance %s Quant: %v Price: %v ID: %d", coin, botID, side, newSellQuantity, priceB, binanceOrderID)
		isPlaceSellOrder := side == "SELL"
		go calculateProfit(coin, newSellQuantity, priceF, priceB, botID, binanceOrderID, origClientOrderID, isPlaceSellOrder)
		text = fmt.Sprintf("%s Sleep %d seconds", text, defaultSleepSeconds)
		log.Println(text)
		go teleClient.SendMessage(text, chatProfitID)
		time.Sleep(time.Duration(defaultSleepSeconds) * time.Second)
		return
	} else {
		text := fmt.Sprintf("%s %s Err Take profit Binance %s Quant: %v Price: %v ID: %d", coin, botID, side, newSellQuantity, priceB, binanceOrderID)
		log.Println(text)
		go teleClient.SendMessage(text, chatErrorID)
	}
}

func notifyBinanceFailOrder(botID string, newSellQuantity, totalUSDT float64, orderType string, err error) {
	text := fmt.Sprintf("Error %s: Can not make order on Binance %s %s", os.Getenv("TELEGRAM_HANDLER"), coin, botID)
	text = fmt.Sprintf("%s =========================: %sLIMIT: Quantity: %v TotalUSDT %v Error: %s", text, orderType, newSellQuantity, totalUSDT, err)
	log.Println(text)
	go teleClient.SendMessage(text, chatErrorID)
	time.Sleep(5 * time.Second)
}
