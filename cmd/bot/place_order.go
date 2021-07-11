package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
)

func placeOrder(botID string, quantity, exchancePrice, fiahubPrice float64, side string) {
	var orderDetails *exchanges.OrderResp
	var err error
	if side == "buy" {
		orderDetails, err = exchangeClient.BuyLimit(coin, exchancePrice, quantity)
	} else if side == "sell" {
		orderDetails, err = exchangeClient.SellLimit(coin, exchancePrice, quantity)
	}
	if err != nil {
		totalUSDT := quantity * exchancePrice
		notifyFailedOrder(botID, quantity, totalUSDT, side, err)
		return
	}

	ID := orderDetails.ID
	clientID := orderDetails.ClientID

	text := fmt.Sprintf("%s %s Take profit From Exchange %s Quantity: %v Price: %v ID: %d", coin, botID, side, quantity, exchancePrice, ID)
	isPlaceSellOrder := side == "sell"
	go calculateProfit(coin, quantity, fiahubPrice, botID, ID, clientID, isPlaceSellOrder)
	text = fmt.Sprintf("%s Sleep %d seconds", text, defaultSleepSeconds)
	log.Println(text)
	go teleClient.SendMessage(text, chatID)
	time.Sleep(time.Duration(defaultSleepSeconds) * time.Second)
	return
}

func notifyFailedOrder(botID string, sellQuantity, totalUSDT float64, orderType string, err error) {
	text := fmt.Sprintf("Err %s: Can not make order on Exchange %s %s", os.Getenv("TELEGRAM_HANDLER"), coin, botID)
	text = fmt.Sprintf("%s =========================: %s LIMIT: Quantity: %v TotalUSDT %v Error: %s", text, orderType, sellQuantity, totalUSDT, err)
	log.Println(text)
	go teleClient.SendMessage(text, chatErrorID)
	time.Sleep(5 * time.Second)
}
