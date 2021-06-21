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

func trade_ask(botID string, coin string, askF float64, askB float64) {
	baseVntQuantity, _ := strconv.Atoi(os.Getenv("BASE_VNT_QUANTITY"))
	perCancel := redisClient.GetFloat64("per_cancel")
	randdomVntQuantity, _ := strconv.Atoi(os.Getenv("RANDOM_VNT_QUANTITY"))
	randNumber := rand.Intn(randdomVntQuantity)

	vntQuantity := float64(baseVntQuantity + randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/askF, decimalsToRound)
	priceSell := askF
	pricesellRandom := askF
	orderType := "AskOrder"
	askOrder := fiahub.OrderParams{
		Coin:              coin,
		CoinAmount:        originalCoinAmount,
		PricePerUnitCents: pricesellRandom,
		Currency:          "VNT",
		Type:              orderType,
	}
	log.Printf("make ask order: %v", askOrder)
	fiahubOrder, code, err := fia.CreateAskOrder(askOrder)
	if err != nil {
		text := fmt.Sprintf("Error CreateAskOrder %s %s %s Coin Amount: %v Price: %v Code: %d Error: %s. Proceed cancel all order",
			coin, botID, orderType, originalCoinAmount, pricesellRandom, code, err)
		time.Sleep(60 * time.Second)
		log.Println(text)
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
		order, err := fia.GetOrderDetails(fiahubOrderID)
		if err != nil {
			text := fmt.Sprintf("Error %s IDTrade: %s type: %s GetOrderDetails %s fiahubOrderID: %d", coin, botID, orderType, err, fiahubOrderID)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(1 * time.Second)
			continue
		}
		state := order.State
		coinAmount := order.CoinAmount
		if coinAmount > 0 {
			executedQty = originalCoinAmount - coinAmount
		}
		if state == fiahub.ORDER_CANCELLED || state == fiahub.ORDER_FINISHED {
			matchingTX, _ := fia.GetSelfMatchingTransaction(order.UserID, order.ID)
			matching = matchingTX != nil
			break
		}

		// Trigger cancel process
		_, askPriceByQuantity := binance.GetPriceByQuantity(coin+"USDT", quantityToGetPrice)
		perChange := math.Abs((askPriceByQuantity - askB) / askB)
		if perChange > perCancel || executedQty > 0 {
			lastestCancelAllTime := fia.GetCancelTime()
			now := time.Now()
			miliTime := now.UnixNano() / int64(time.Millisecond)
			elapsedTime := miliTime - lastestCancelAllTime
			if elapsedTime < 10000 {
				text := fmt.Sprintf("%s IDTrade: %s, CancelTime < 10s continue ElapsedTime: %v Starttime: %v", coin, botID, elapsedTime, lastestCancelAllTime)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(3000 * time.Millisecond)
				continue
			}

			log.Printf("Bot: %s cancel fiahub ask order %d due to: perChange: %v, executedQty: %v", botID, fiahubOrderID, perChange, executedQty)
			orderDetails, _, err := fia.CancelOrder(fiahubOrderID)
			if err != nil {
				text := fmt.Sprintf("Error %s IDTrade: %s, type: %s, ERROR!!! CancelOrder: %d with error: %s", coin, botID, orderType, fiahubOrderID, err)
				log.Println(text)
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
		time.Sleep(5 * time.Second)
	}

	// If newSellVNTQuantity < 50.000 ignore
	// If newSellVNTQuantity > 250.000 mới tạo lệnh mua bù trên binance không thì tạo lệnh bán lại luôn giá + rand từ 1->3000
	newSellQuantity := executedQty - totalSell
	newSellVNTQuantity := newSellQuantity * priceSell
	if newSellVNTQuantity <= 50000 {
		return
	}

	if newSellVNTQuantity < 250000 {
		text := fmt.Sprintf("%s %s  Take profit < 10$ %s Quant: %v Price: %v ID: %d", coin, botID, orderType, newSellQuantity, pricesellRandom, fiahubOrderID)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(0.3 * 60 * 1000 * time.Millisecond)
		return
	}

	if matching {
		text := fmt.Sprintf("%s %s Self Matching", coin, botID)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5000 * time.Millisecond)
		return
	}

	newSellQuantity = utils.RoundTo(newSellQuantity, decimalsToRound)
	placeBinanceOrder(botID, newSellQuantity, askB, askF, "BUY")
}
