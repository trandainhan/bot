package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/fiahub"
	"gitlab.com/fiahub/bot/internal/utils"
)

func trade_bid(botID string, coin string, bidF float64, bidB float64) {
	baseVntQuantity, _ := strconv.Atoi(os.Getenv("BASE_VNT_QUANTITY"))
	perCancel := redisClient.GetFloat64("per_cancel")
	randdomVntQuantity, _ := strconv.Atoi(os.Getenv("RANDOM_VNT_QUANTITY"))
	randNumber := rand.Intn(randdomVntQuantity)

	vntQuantity := float64(baseVntQuantity + randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/bidF, decimalsToRound)
	priceBuy := bidF
	orderType := "BidOrder"
	bidOrder := fiahub.Order{
		Coin:              coin,
		CoinAmount:        originalCoinAmount,
		PricePerUnitCents: priceBuy,
		Currency:          "VNT",
		Type:              orderType,
	}
	log.Printf("make bid order: %v", bidOrder)
	fiahubOrder, code, err := fia.CreateBidOrder(bidOrder)
	if err != nil {
		text := fmt.Sprintf("Error CreateBidOrder %s %s %s Coin Amount: %v Price: %v, StatusCode: %d Err: %s. Proceed cancel all order",
			coin, botID, orderType, originalCoinAmount, priceBuy, code, err)
		time.Sleep(60 * time.Second)
		go teleClient.SendMessage(text, chatErrorID)
		fia.CancelAllOrder()
		time.Sleep(5 * time.Second)
		return
	}
	time.Sleep(3 * time.Second)

	// Loop to check order
	fiahubOrderID := fiahubOrder.ID
	executedQty := 0.0
	totalSell := 0.0
	matching := false
	for {
		orderDetails, code, err := fia.GetBidOrderDetails(fiahubOrderID)
		if err != nil {
			text := fmt.Sprintf("Error %s IDTrade: %s type: %s GetBidOrderDetails %s StatusCode: %d fiahubOrderID: %d", coin, botID, orderType, err, code, fiahubOrderID)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		state := orderDetails.State
		coinAmount := orderDetails.GetCoinAmount()
		matching = orderDetails.Matching
		if coinAmount > 0 {
			executedQty = originalCoinAmount - coinAmount
		}
		if state == fiahub.ORDER_CANCELLED || state == fiahub.ORDER_FINISHED {
			break
		}

		// Trigger cancel process
		bidPriceByQuantity, _ := binance.GetPriceByQuantity(coin+"USDT", quantityToGetPrice)
		perChange := math.Abs((bidPriceByQuantity - bidB) / bidB)
		if perChange > perCancel || executedQty > 0 {
			lastestCancelAllTime := fia.GetCancelTime()
			now := time.Now()
			miliTime := now.UnixNano() / int64(time.Millisecond)
			elapsedTime := miliTime - lastestCancelAllTime
			if elapsedTime < 10000 {
				text := fmt.Sprintf("%s IDTrade: %s, CancelTime < 10s continue ElapsedTime: %v Starttime: %v", coin, botID, elapsedTime, lastestCancelAllTime)
				go teleClient.SendMessage(text, chatID)
				time.Sleep(3000 * time.Millisecond)
				continue
			}

			log.Printf("Bot: %s cancel fiahub bid order %d due to: perChange: %v, executedQty: %v", botID, fiahubOrderID, perChange, executedQty)
			orderDetails, code, err = fia.CancelOrder(fiahubOrderID)
			if err != nil {
				text := fmt.Sprintf("Error! %s IDTrade: %s, type: %s, ERROR!!! CancelOrder: %d with error: %s", coin, botID, orderType, fiahubOrderID, err)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(3 * time.Second)
				continue
			}
			coinAmount = orderDetails.GetCoinAmount()
			matching = orderDetails.Matching
			if coinAmount > 0 {
				executedQty = originalCoinAmount - coinAmount
			}
			break
		}
		time.Sleep(5000 * time.Millisecond)
	}

	// If newSellVNTQuantity < 50.000 ignore
	// If newSellVNTQuantity > 250.000 mới tạo lệnh mua bù trên binance không thì tạo lệnh bán lại luôn giá + rand từ 1->3000
	newSellQuantity := executedQty - totalSell
	newSellVNTQuantity := newSellQuantity * bidF
	if newSellVNTQuantity <= 50000 {
		return
	}

	if newSellVNTQuantity < 250000 {
		text := fmt.Sprintf("%s %s  Chốt lời < 10$ %s Quant: %v Price: %v ID: %d", coin, botID, orderType, newSellQuantity, priceBuy, fiahubOrderID)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(0.3 * 60 * 1000 * time.Millisecond)
		return
	}

	if matching {
		text := fmt.Sprintf("%s %s self-matching  matching: %v", coin, botID, matching)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5000 * time.Millisecond)
		return
	}

	newSellQuantity = utils.RoundTo(newSellQuantity, decimalsToRound)
	go placeBinanceOrder(botID, newSellQuantity, bidB, bidF, "SELL")
}
