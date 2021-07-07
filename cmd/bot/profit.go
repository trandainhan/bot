package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/utils"
)

func calculateProfit(coin string, newSellQuantity, askF, askB float64, id string, binanceOrderID int, origClientOrderID string, isLiquidBaseBinanceTradeBid bool) {
	orderDetails := getBinanceOrderDetail(id, coin, binanceOrderID, origClientOrderID)
	if orderDetails == nil {
		return
	}

	bidB := orderDetails.GetPrice()
	askB = orderDetails.GetPrice()
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
		text := fmt.Sprintf("%s %s \n %s: %s %s - %v USDT(-) - %v VNT(+) \n Status: %s Price %v \n Profit: %v - Perprofit %v %% \n", coin, id, side,
			origQty, coin, totalUSDTGive, totalVNTRecieve, status, askB, profit, perProfit)

		allFundMessage := bn.GetFundsMessages()
		text = fmt.Sprintf("%s %s", text, allFundMessage)

		// Calculate USDT Margin
		name := "USDT"
		marginDetails, _ := bn.GetMarginDetails()
		netAsset := calculateUSDTMargin(marginDetails, name)

		text = fmt.Sprintf("%s \n USDT(Margin): %.6f", text, netAsset)
		go teleClient.SendMessage(text, chatProfitID)
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
		text := fmt.Sprintf("%s %s \n %s: %s %s - %.6f USDT(+) - %.6f VNT(-) \n Status: %s Price %v \n Profit: %v - Perprofit %v %% \n", coin, id, side,
			origQty, coin, totalUSDTRecieve, totalVNTGive, status, bidB, profit, perProfit)

		allFundMessage := bn.GetFundsMessages()
		text = fmt.Sprintf("%s %s", text, allFundMessage)

		// Calculate USDT Margin
		name := "USDT"
		marginDetails, _ := bn.GetMarginDetails()
		netAsset := calculateUSDTMargin(marginDetails, name)

		text = fmt.Sprintf("%s \n USDT(Margin): %v", text, netAsset)
		go teleClient.SendMessage(text, chatProfitID)
		time.Sleep(2000 * time.Millisecond)

		notifyWhenAssetIsLow(netAsset, text)
	}
}

func calculateUSDTMargin(marginDetails *binance.MarginDetails, name string) float64 {
	netAsset := 0.0
	userAssets := marginDetails.UserAssets
	for _, userAsset := range userAssets {
		if userAsset.Asset == name {
			netAsset = userAsset.GetNetAsset()
			break
		}
	}
	return netAsset
}

func notifyWhenAssetIsLow(netAsset float64, baseText string) {
	baseNetAsset, _ := strconv.ParseFloat(os.Getenv("BASE_NET_ASSET"), 64)
	if netAsset < baseNetAsset {
		text := fmt.Sprintf("%s %s", os.Getenv("TELEGRAM_HANDLER"), baseText)
		go teleClient.SendMessage(text, chatID)
		time.Sleep(1000 * time.Millisecond)
	}
	text := fmt.Sprintf("USDT(Margin) %v", netAsset)
	go teleClient.SendMessage(text, chatID)
}

func getBinanceOrderDetail(id string, coin string, binanceOrderID int, origClientOrderID string) *binance.OrderDetailsResp {
	var orderDetails *binance.OrderDetailsResp
	var err error
	for j := 0; j <= 2; j++ {
		orderDetails, err = bn.GetOrder(coin+"USDT", binanceOrderID, origClientOrderID)
		if err != nil {
			text := fmt.Sprintf("%s %s ERROR getBinanceOrderDetail: %s", coin, id, err)
			go teleClient.SendMessage(text, chatErrorID)
		} else {
			status := orderDetails.Status
			if status == binance.ORDER_FILLED {
				break
			}
		}

		if j == 2 {
			log.Println("Unsucessfully sell in binance after 1 minute")
			break
		}
		time.Sleep(30 * time.Second)
	}
	return orderDetails
}
