package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/utils"
)

func trade_bid(botID string, coin string, fiahubPrice float64, exchangePrice float64, cancelFactor int) {
	key := fmt.Sprintf("%s_VNT_QUANTITY", strings.ToUpper(botID))
	baseVntQuantity, _ := strconv.ParseFloat(os.Getenv(key), 64)
	perCancel := redisClient.GetFloat64("per_cancel") + float64(cancelFactor-1)*0.05/100
	randdomVntQuantity, _ := strconv.Atoi(os.Getenv("RANDOM_VNT_QUANTITY"))
	randNumber := rand.Intn(randdomVntQuantity)

	vntQuantity := baseVntQuantity + float64(randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/fiahubPrice, decimalsToRound)
	priceBuy := fiahubPrice
	orderType := "BidOrder"
	bidOrder := fiahub.OrderParams{
		Coin:              coin,
		CoinAmount:        originalCoinAmount,
		PricePerUnitCents: priceBuy,
		Currency:          "VNT",
		Type:              orderType,
	}
	log.Printf("Make bid order: %v", bidOrder)
	fiahubOrder, code, err := fia.CreateBidOrder(bidOrder)
	if err != nil {
		text := fmt.Sprintf("Error CreateBidOrder %s %s %s Coin Amount: %v Price: %v, StatusCode: %d Err: %s",
			coin, botID, orderType, originalCoinAmount, priceBuy, code, err)
		time.Sleep(60 * time.Second)
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
	newSellVNTQuantity := newSellQuantity * fiahubPrice
	if newSellVNTQuantity <= 50000 {
		return
	}

	if newSellVNTQuantity < 250000 {
		text := fmt.Sprintf("%s %s  Take profit < 10$ %s Quant: %.6f Price: %.6f ID: %d", coin, botID, orderType, newSellQuantity, priceBuy, fiahubOrderID)
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
	placeOrder(botID, newSellQuantity, exchangePrice, fiahubPrice, "sell")
}
