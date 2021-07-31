package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/utils"
)

func trade_ask(botID string, coin string, fiahubPrice float64, exchangePrice float64, cancelFactor int) {
	key := fmt.Sprintf("%s_vnt_quantity", botID)
	baseVntQuantity, _ := strconv.ParseFloat(os.Getenv(key), 64)
	perCancel := redisClient.GetFloat64("per_cancel") + float64(cancelFactor-1)*0.05/100
	randdomVntQuantity, _ := strconv.Atoi(os.Getenv("RANDOM_VNT_QUANTITY"))
	randNumber := rand.Intn(randdomVntQuantity)

	vntQuantity := baseVntQuantity + float64(randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/fiahubPrice, decimalsToRound)
	priceSell := fiahubPrice
	orderType := "AskOrder"
	askOrder := fiahub.OrderParams{
		Coin:              coin,
		CoinAmount:        originalCoinAmount,
		PricePerUnitCents: priceSell,
		Currency:          "VNT",
		Type:              orderType,
	}
	log.Printf("make ask order: %v", askOrder)
	fiahubOrder, code, err := fia.CreateAskOrder(askOrder)
	if err != nil {
		text := fmt.Sprintf("Error CreateAskOrder %s %s %s Coin Amount: %v Price: %v Code: %d Error: %s. Proceed cancel all order",
			coin, botID, orderType, originalCoinAmount, priceSell, code, err)
		time.Sleep(60 * time.Second)
		log.Println(text)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5 * time.Second)
		return
	}
	time.Sleep(3 * time.Second)

	// Loop to check order
	fiahubOrderID := fiahubOrder.ID
	executedQty, matching := checkFiahubOrder(botID, fiahubOrderID, originalCoinAmount, exchangePrice, perCancel, orderType)

	// If newSellVNTQuantity < 50.000 ignore
	// If newSellVNTQuantity > 250.000 mới tạo lệnh mua bù trên binance không thì tạo lệnh bán lại luôn giá + rand từ 1->3000
	newSellQuantity := executedQty
	newSellVNTQuantity := newSellQuantity * priceSell
	if newSellVNTQuantity <= 50000 {
		return
	}

	if newSellVNTQuantity < 250000 {
		text := fmt.Sprintf("%s %s  Take profit < 10$ %s Quant: %.6f Price: %.6f ID: %d", coin, botID, orderType, newSellQuantity, priceSell, fiahubOrderID)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(0.3 * 60 * time.Second)
		return
	}

	if matching {
		text := fmt.Sprintf("%s %s Self Matching", coin, botID)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5 * time.Second)
		return
	}

	newSellQuantity = utils.RoundTo(newSellQuantity, decimalsToRound)
	placeOrder(botID, newSellQuantity, exchangePrice, fiahubPrice, "buy")
}
