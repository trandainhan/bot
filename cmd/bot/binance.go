package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func notifyBinanceFailOrder(botID string, newSellQuantity, totalUSDT float64, orderType string, err error) {
	text := fmt.Sprintf("Error %s: Can not make order on Binance %s %s", os.Getenv("TELEGRAM_HANDLER"), coin, botID)
	text = fmt.Sprintf("%s =========================: %sLIMIT: Quantity: %v TotalUSDT %v Error: %s", text, orderType, newSellQuantity, totalUSDT, err)
	log.Println(text)
	go teleClient.SendMessage(text, chatErrorID)
	time.Sleep(5000 * time.Millisecond)
}
