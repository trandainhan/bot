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

func trade_ask(id string, coin string, askF float64, askB float64, perProfitStep float64) {
	baseVntQuantity, _ := strconv.Atoi(os.Getenv("BASE_VNT_QUANTITY"))
	perCancel := redisClient.GetFloat64("per_cancel")
	perProfit := redisClient.GetFloat64("per_profit_ask")
	fiahubToken := redisClient.Get("fiahub_token")
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	chatErrorID, _ := strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)

	perProfit = perProfit + perProfitStep*0.6/100
	randNumber := rand.Intn(4000000)

	vntQuantity := float64(baseVntQuantity + randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/askF, decimalsToRound)
	coinAmount := originalCoinAmount
	priceSell := askF
	pricesellRandom := askF
	orderType := "AskOrder"
	askOrder := fiahub.Order{
		Coin:               coin,
		OriginalCoinAmount: originalCoinAmount,
		PricePerUnitCents:  pricesellRandom,
		Currency:           "VNT",
		Type:               orderType,
	}
	fiahubOrder, statusCode, err := fiahub.CreateAskOrder(fiahubToken, askOrder)
	if err != nil {
		text := fmt.Sprintf("fiahubAPI_AskOrder Error! %s %s %s Coin Amount: %v Price: %v, StatusCode: %d %s", coin, id, orderType, coinAmount, pricesellRandom, statusCode, err)
		time.Sleep(60000 * time.Millisecond)
		log.Println(text)
		go teleClient.SendMessage(text, chatErrorID)
		fia.CancelAllOrder(fiahubToken)
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
		orderDetails, code, err := fiahub.GetOrderDetails(fiahubToken, fiahubOrderID)
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
			lastestCancelAllTime := redisClient.GetTime("lastest_cancel_all_time")
			tnow := time.Now()
			elapsedTime := tnow.Sub(lastestCancelAllTime)
			if elapsedTime < 10000*time.Millisecond {
				text := fmt.Sprintf("%s IDTrade: %s, CancelTime < 10s continue ElapsedTime: %v Starttime: %v", coin, id, elapsedTime, lastestCancelAllTime)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(3000 * time.Millisecond)
				continue
			}
		}
		orderDetails, code, err = fiahub.CancelOrder(fiahubToken, fiahubOrderID)
		if err != nil {
			text := fmt.Sprintf("Error! %s IDTrade: %s, type: %s, ERROR!!! Cancelorder: %d with error: %s", coin, id, orderType, fiahubOrderID, err)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(3000 * time.Millisecond)
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
		go teleClient.SendMessage(text, chatErrorID)
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
		text := fmt.Sprintf("Error! %s %s %s %s Không Thực hiện được lệnh", os.Getenv("TELEGRAM_HANDLER"), coin, id, orderType)
		log.Println(text)
		btcQuantity := newSellQuantity * askB
		text = fmt.Sprintf("%s  ===   =====   ========   ======   ===   BuyLimit: %v TotalUSDT %v Error: %s", text, newSellQuantity, btcQuantity, err)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5000 * time.Millisecond)
	}
	if binanceOrderID != nil {
		text := fmt.Sprintf("%s %s Chot loi Binance BuyLimit Quant: %v Price: %v ID: %s", coin, id, newSellQuantity, askB, *binanceOrderID)
		log.Println(text)
		isLiquidBaseBinanceTradeBid := false
		calculateProfit(coin, newSellQuantity, askF, askB, id, binanceOrderID, origClientOrderID, isLiquidBaseBinanceTradeBid)

		text = fmt.Sprintf("%s Sleep %d seconds", text, defaultSleepSeconds)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(time.Duration(defaultSleepSeconds) * time.Second)
		return
	}
}
