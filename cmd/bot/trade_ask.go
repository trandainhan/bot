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

func trade_ask(id string, coin string, askF float64, askB float64) {
	baseVntQuantity, _ := strconv.Atoi(os.Getenv("BASE_VNT_QUANTITY"))
	perCancel := redisClient.GetFloat64("per_cancel")
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	chatErrorID, _ := strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)

	randdomVntQuantity, _ := strconv.Atoi(os.Getenv("RANDOM_VNT_QUANTITY"))
	randNumber := rand.Intn(randdomVntQuantity)

	vntQuantity := float64(baseVntQuantity + randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/askF, decimalsToRound)
	priceSell := askF
	pricesellRandom := askF
	orderType := "AskOrder"
	askOrder := fiahub.Order{
		Coin:              coin,
		CoinAmount:        originalCoinAmount,
		PricePerUnitCents: pricesellRandom,
		Currency:          "VNT",
		Type:              orderType,
	}
	log.Printf("make ask order: %v", askOrder)
	fiahubOrder, code, err := fia.CreateAskOrder(askOrder)
	if err != nil {
		text := fmt.Sprintf("Error CreateAskOrder %s %s %s Coin Amount: %v Price: %v Code: %d Error: %s", coin, id, orderType, originalCoinAmount, pricesellRandom, code, err)
		time.Sleep(60000 * time.Millisecond)
		log.Println(text)
		go teleClient.SendMessage(text, chatErrorID)
		fia.CancelAllOrder()
		time.Sleep(5000 * time.Millisecond)
		return
	}
	time.Sleep(3000 * time.Millisecond)

	// Loop to check order
	fiahubOrderID := fiahubOrder.ID
	executedQty := 0.0
	totalSell := 0.0
	matching := false
	for {
		orderDetails, code, err := fia.GetAskOrderDetails(fiahubOrderID)
		if err != nil {
			text := fmt.Sprintf("Error! %s IDTrade: %s, type: %s ERROR!!! Queryorder %s StatusCode: %d fiahubOrderID: %d", coin, id, orderType, err, code, fiahubOrderID)
			log.Println(text)
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
		_, askPriceByQuantity := binance.GetPriceByQuantity(coin+"USDT", quantityToGetPrice)
		perchange := math.Abs((askPriceByQuantity - askB) / askB)
		if perchange > perCancel || executedQty > 0 {
			lastestCancelAllTime := redisClient.GetInt64("lastest_cancel_all_time")
			now := time.Now()
			miliTime := now.UnixNano() / int64(time.Millisecond)
			elapsedTime := miliTime - lastestCancelAllTime
			if elapsedTime < 10000 {
				text := fmt.Sprintf("%s IDTrade: %s, CancelTime < 10s continue ElapsedTime: %v Starttime: %v", coin, id, elapsedTime, lastestCancelAllTime)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(3000 * time.Millisecond)
				continue
			}

			orderDetails, code, err = fia.CancelOrder(fiahubOrderID)
			if err != nil {
				text := fmt.Sprintf("Error %s IDTrade: %s, type: %s, ERROR!!! Cancelorder: %d with error: %s", coin, id, orderType, fiahubOrderID, err)
				log.Println(text)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(3000 * time.Millisecond)
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
	newSellVNTQuantity := newSellQuantity * priceSell
	if newSellVNTQuantity <= 50000 {
		return
	}

	if newSellVNTQuantity < 250000 {
		text := fmt.Sprintf("%s %s  Chốt lời < 10$ %s Quant: %v Price: %v ID: %d", coin, id, orderType, newSellQuantity, pricesellRandom, fiahubOrderID)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(0.3 * 60 * 1000 * time.Millisecond)
		return
	}

	if matching {
		text := fmt.Sprintf("%s %s self-matching  matching: %v", coin, id, matching)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5000 * time.Millisecond)
		return
	}

	newSellQuantity = utils.RoundTo(newSellQuantity, decimalsToRound)
	bn := binance.Binance{
		RedisClient: redisClient,
	}
	orderDetails, err := bn.BuyLimit(coin+"USDT", askB, newSellQuantity)
	binanceOrderID := orderDetails.OrderID
	origClientOrderID := orderDetails.ClientOrderID
	if err != nil {
		text := fmt.Sprintf("Error BuyLimit: %s %s %s %s", os.Getenv("TELEGRAM_HANDLER"), coin, id, orderType)
		log.Println(text)
		btcQuantity := newSellQuantity * askB
		text = fmt.Sprintf("%s  ===   =====   ========   ======   ===   BuyLimit: %v TotalUSDT %v Error: %s", text, newSellQuantity, btcQuantity, err)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5000 * time.Millisecond)
	}
	if binanceOrderID != 0 {
		text := fmt.Sprintf("%s %s Take profit Binance BuyLimit Quant: %v Price: %v ID: %d", coin, id, newSellQuantity, askB, binanceOrderID)
		log.Println(text)
		isLiquidBaseBinanceTradeBid := false
		go calculateProfit(coin, newSellQuantity, askF, askB, id, binanceOrderID, origClientOrderID, isLiquidBaseBinanceTradeBid)

		text = fmt.Sprintf("%s Sleep %d seconds", text, defaultSleepSeconds)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(time.Duration(defaultSleepSeconds) * time.Second)
		return
	}
}
