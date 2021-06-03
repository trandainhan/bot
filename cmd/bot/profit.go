package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/utils"
)

var bn = binance.Binance{
	RedisClient: redisClient,
}

func calculateProfit(coin string, newSellQuantity, askF, askB float64, id string, binanceOrderID *string, origClientOrderID string, isLiquidBaseBinanceTradeBid bool) {
	orderDetails := getBinanceOrderDetail(id, coin, binanceOrderID, origClientOrderID)

	bidB := orderDetails.Price
	askB = orderDetails.Price
	rate := redisClient.GetFloat64("usdtvnd_rate")
	origQty := orderDetails.OriginQty
	status := orderDetails.Status
	side := orderDetails.Side
	if isLiquidBaseBinanceTradeBid == false {
		totalVNTRecieve := askF * newSellQuantity
		feeUSDT := askB * newSellQuantity * 0.075 / 100
		totalUSDTGive := askB*newSellQuantity + feeUSDT

		profit := utils.RoundTo((totalVNTRecieve - totalUSDTGive*rate), 0)
		perProfit := utils.RoundTo((profit/totalVNTRecieve)*100, 2)
		text := fmt.Sprintf("%s %s \n %s: %v %s - %v USDT(-) - %v VNT(+) \n Status: %s Price %v \n Profit: %v - Perprofit %v %% \n", coin, id, side,
			origQty, coin, totalUSDTGive, totalVNTRecieve, status, askB, profit, perProfit)

		allFundMessage := bn.GetFundsMessages()
		text = fmt.Sprintf("%s %s", text, allFundMessage)

		// ;Tinh USDT Margin
		name := "USDT"
		marginDetails, _ := bn.GetMarginDetails()
		netAsset := calculateUSDTMargin(marginDetails, name)

		text = fmt.Sprintf("%s \n USDT(Margin): %v", text, netAsset)
		teleClient.SendMessage(text, -465055332)
		time.Sleep(2000 * time.Millisecond)

		notifyWhenAssetIsLow(netAsset, text)
	}

	if isLiquidBaseBinanceTradeBid == true {
		bidF := askF
		totalVNTGive := bidF * newSellQuantity
		feeUSDT := bidB * newSellQuantity * 0.075 / 100
		totalUSDTRecieve := bidB*newSellQuantity - feeUSDT

		profit := utils.RoundTo((totalUSDTRecieve*rate - totalVNTGive), 0)
		perProfit := utils.RoundTo((profit/(totalUSDTRecieve*rate))*100, 2)
		text := fmt.Sprintf("%s %s \n %s: %v %s - %v USDT(+) - %v VNT(-) \n Status: %s Price %v \n Profit: %v - Perprofit %v %% \n", coin, id, side,
			origQty, coin, totalUSDTRecieve, totalVNTGive, status, askB, profit, perProfit)

		allFundMessage := bn.GetFundsMessages()
		text = fmt.Sprintf("%s %s", text, allFundMessage)

		// ;Tinh USDT Margin
		name := "USDT"
		marginDetails, _ := bn.GetMarginDetails()
		netAsset := calculateUSDTMargin(marginDetails, name)

		text = fmt.Sprintf("%s \n USDT(Margin): %v", text, netAsset)
		teleClient.SendMessage(text, -465055332)
		time.Sleep(2000 * time.Millisecond)

		notifyWhenAssetIsLow(netAsset, text)
	}
}

func calculateUSDTMargin(marginDetails *binance.MarginDetails, name string) float64 {
	netAsset := 0.0
	userAssets := marginDetails.UserAssets
	for _, userAsset := range userAssets {
		if userAsset.Name == name {
			netAsset = userAsset.NetAsset
			break
		}
	}
	return netAsset
}

func notifyWhenAssetIsLow(netAsset float64, baseText string) {
	if netAsset < 10000 {
		text := fmt.Sprintf("%s @ndtan", baseText)
		teleClient.SendMessage(text, -357553425)
		time.Sleep(1000 * time.Millisecond)
	}
	text := fmt.Sprintf("USDT(Margin) %v", netAsset)
	teleClient.SendMessage(text, -357553425)
}

func getBinanceOrderDetail(id string, coin string, binanceOrderID *string, origClientOrderID string) *binance.OrderDetailsResp {
	chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
	var orderDetails *binance.OrderDetailsResp
	var err error
	for j := 0; j <= 2; j++ {
		orderDetails, err = bn.GetOrder(coin+"USDT", *binanceOrderID, origClientOrderID)
		if err != nil {
			text := fmt.Sprintf("%s %s ERROR!!! Queryorder %s", coin, id, err)
			teleClient.SendMessage(text, chatID)
		}

		status := orderDetails.Status
		if status == "FILLED" {
			break
		}

		// Sell uncessfully in 1 minutes
		if j == 2 {
			break
		}
		time.Sleep(30000 * time.Millisecond)
	}
	return orderDetails
}
