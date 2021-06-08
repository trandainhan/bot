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

func trade_bid(id string, coin string, bidF float64, bidB float64, perProfitStep float64) {
	baseVntQuantity, _ := strconv.Atoi(os.Getenv("BASE_VNT_QUANTITY"))
	perCancel := redisClient.GetFloat64("per_cancel")
	perProfit := redisClient.GetFloat64("per_profit_bid")
	fiahubToken := redisClient.Get("fiahub_token")
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	chatErrorID, _ := strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)

	perProfit = perProfit + perProfitStep*0.6/100
	randdomVntQuantity, _ := strconv.Atoi(os.Getenv("RANDOM_VNT_QUANTITY"))
	randNumber := rand.Intn(randdomVntQuantity)

	vntQuantity := float64(baseVntQuantity + randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/bidF, decimalsToRound)
	coinAmount := originalCoinAmount
	priceBuy := bidF
	orderType := "BidOrder"
	bidOrder := fiahub.Order{
		Coin:               coin,
		Type:               orderType,
		Currency:           "VNT",
		PricePerUnitCents:  priceBuy,
		OriginalCoinAmount: originalCoinAmount,
	}
	log.Printf("make bid order: %v", bidOrder)
	fiahubOrder, err := fiahub.CreateBidOrder(fiahubToken, bidOrder)
	if err != nil {
		text := fmt.Sprintf("fiahubAPI_AskOrder Error! %s %s %s Coin Amount: %v Price: %v, %s", coin, id, orderType, coinAmount, priceBuy, err)
		time.Sleep(60000 * time.Millisecond)
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
		perchange := math.Abs((bidPriceByQuantity - bidB) / bidB)
		if perchange > perCancel || executedQty > 0 {
			lastestCancelAllTime := redisClient.GetTime("lastest_cancel_all_time")
			tnow := time.Now()
			elapsedTime := tnow.Sub(lastestCancelAllTime)
			if elapsedTime < 10000*time.Millisecond {
				text := fmt.Sprintf("%s IDTrade: %s, CancelTime < 10s continue ElapsedTime: %v Starttime: %v", coin, id, elapsedTime, lastestCancelAllTime)
				go teleClient.SendMessage(text, chatID)
				time.Sleep(3000 * time.Millisecond)
				continue
			}
		}
		orderDetails, code, err = fiahub.CancelOrder(fiahubToken, fiahubOrderID)
		if err != nil {
			text := fmt.Sprintf("Error! %s IDTrade: %s, type: %s, ERROR!!! Cancelorder: %d with error: %s", coin, id, orderType, fiahubOrderID, err)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(3000 * time.Millisecond)
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
		text := fmt.Sprintf("%s %s  Chốt lời < 10$ %s Quant: %v Price: %v ID: %d", coin, id, orderType, newSellQuantity, priceBuy, fiahubOrderID)
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
	orderDetails, err := bn.SellLimit(coin+"USDT", bidB, newSellQuantity)
	binanceOrderID := orderDetails.OrderID
	origClientOrderID := orderDetails.ClientOrderID
	if err != nil {
		text := fmt.Sprintf("Error! %s %s %s %s Không Thực hiện được lệnh", os.Getenv("TELEGRAM_HANDLER"), coin, id, orderType)
		btcQuantity := newSellQuantity * bidB
		text = fmt.Sprintf("%s  ===   =====   ========   ======   ===   BuyLimit: %v TotalUSDT %v Error: %s", text, newSellQuantity, btcQuantity, err)
		go teleClient.SendMessage(text, chatErrorID)
		time.Sleep(5000 * time.Millisecond)
	}
	if binanceOrderID != nil {
		text := fmt.Sprintf("%s %s Chot loi Binance BuyLimit Quant: %v Price: %v ID: %s", coin, id, newSellQuantity, bidB, *binanceOrderID)
		isLiquidBaseBinanceTradeBid := true
		calculateProfit(coin, newSellQuantity, bidF, bidB, id, binanceOrderID, origClientOrderID, isLiquidBaseBinanceTradeBid)

		text = fmt.Sprintf("%s Sleep %d seconds", text, defaultSleepSeconds)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(time.Duration(defaultSleepSeconds) * time.Second)
		return
	}
}
