package main

import (
	"fmt"
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
	baseVntQuantity, _ := strconv.Atoi(os.Getenv("base_vnt_quantity")) // 18000000
	perCancel := redisClient.Get("per_cancel").(float64)
	perProfit := redisClient.Get("per_profit_ask").(float64) // this is ask worker
	fiahubToken := redisClient.Get("fiahub_token").(string)
	chatID, _ := strconv.ParseInt(os.Getenv("chat_id"), 10, 64)
	chatErrorID, _ := strconv.ParseInt(os.Getenv("chat_error_id"), 10, 64)

	perProfit = perProfit + perProfitStep*0.6/100
	randNumber := rand.Intn(4000000)

	vntQuantity := float64(baseVntQuantity + randNumber)

	originalCoinAmount := utils.RoundTo(vntQuantity/bidF, decimalsToRound)
	coinAmount := originalCoinAmount
	priceBuy := bidF
	// pricesellRandom := bidF
	orderType := "BidOrder"
	bidOrder := fiahub.Order{
		Coin:               coin,
		OriginalCoinAmount: originalCoinAmount,
		Currency:           "VNT",
		Type:               orderType,
	}
	fiahubOrder, statusCode, err := fiahub.CreateBidOrder(fiahubToken, bidOrder)
	fiahubOrderID := fiahubOrder.ID
	if err != nil {
		text := fmt.Sprintf("fiahubAPI_AskOrder Error! %s %s %s Coin Amount: %v Price: %v, StatusCode: %d %s", coin, id, orderType, coinAmount, priceBuy, statusCode, err)
		time.Sleep(60000 * time.Millisecond)
		go teleClient.SendMessage(text, chatErrorID)
		fiahub.CancelAllOrder(fiahubToken)
		time.Sleep(5000 * time.Millisecond)
		return
	}
	time.Sleep(3000 * time.Millisecond)

	// Loop to check order
	executedQty := 0.0
	totalSell := 0.0
	matching := false
	for {
		orderDetails, code, err := fiahub.GetOrderDetails(fiahubToken, fiahubOrderID)
		if err != nil {
			text := fmt.Sprintf("Error! %s IDTrade: %s, type: %s ERROR!!! Queryorder %s StatusCode: %d fiahubOrderID: %s", coin, id, orderType, err, code, fiahubOrderID)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		state := orderDetails.State
		coinAmount := orderDetails.CoinAmount
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
			lastestCancelAllTime := redisClient.Get("lastest_cancel_all_time").(time.Time)
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
			text := fmt.Sprintf("Error! %s IDTrade: %s, type: %s, ERROR!!! Cancelorder: %s with error: %s", coin, id, orderType, fiahubOrderID, err)
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
		text := fmt.Sprintf("%s %s  Chốt lời < 10$ %s Quant: %v Price: %v ID: %s", coin, id, orderType, newSellQuantity, priceBuy, fiahubOrderID)
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
	orderDetails, err := binance.SellLimit(coin+"USDT", bidB, newSellQuantity)
	binanceOrderID := orderDetails.ID
	origClientOrderID := orderDetails.ClientOrderID
	if err != nil {
		text := fmt.Sprintf("Error! @ndtan %s %s %s Không Thực hiện được lệnh", coin, id, orderType)
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
