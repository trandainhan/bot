package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/exchanges/binance"
	"gitlab.com/fiahub/bot/internal/exchanges/ftx"
	"gitlab.com/fiahub/bot/internal/utils"
)

func calculateProfit(coin string, newSellQuantity, fiahubPrice float64, id string, exchangeOrderID int64, origClientOrderID string, isLiquidBaseBinanceTradeBid bool) {
	orderDetails := getOrderDetails(id, coin, exchangeOrderID, origClientOrderID)
	if orderDetails == nil {
		return
	}

	exchangePrice := orderDetails.Price
	rate, _ := redisClient.GetFloat64("usdtvnd_rate")
	origQty := orderDetails.OriginQty
	status := orderDetails.Status
	side := orderDetails.Side
	text := ""
	if isLiquidBaseBinanceTradeBid == false {
		totalVNTRecieve := fiahubPrice * newSellQuantity
		feeUSDT := exchangePrice * newSellQuantity * 0.075 / 100
		totalUSDTGive := exchangePrice*newSellQuantity + feeUSDT

		profit := utils.RoundTo((totalVNTRecieve - totalUSDTGive*rate), 0)
		perProfit := utils.RoundTo((profit/totalVNTRecieve)*100, 2)
		text = fmt.Sprintf("%s %s %s \n %s: %f %s - %v USDT(-) - %v VNT(+) \n Status: %s Price %v \n Profit: %v - Perprofit %v %% \n", currentExchange, coin, id, side,
			origQty, coin, totalUSDTGive, totalVNTRecieve, status, exchangePrice, profit, perProfit)

		allFundMessage := exchangeClient.GetFundsMessages()
		text = fmt.Sprintf("%s %s", text, allFundMessage)
	}

	if isLiquidBaseBinanceTradeBid == true {
		totalVNTGive := fiahubPrice * newSellQuantity
		feeUSDT := exchangePrice * newSellQuantity * 0.075 / 100
		totalUSDTRecieve := exchangePrice*newSellQuantity - feeUSDT

		profit := utils.RoundTo((totalUSDTRecieve*rate - totalVNTGive), 0)
		perProfit := utils.RoundTo((profit/(totalUSDTRecieve*rate))*100, 2)
		text = fmt.Sprintf("%s %s %s \n %s: %f %s - %.6f USDT(+) - %.6f VNT(-) \n Status: %s Price %v \n Profit: %v - Perprofit %v %% \n", currentExchange, coin, id, side,
			origQty, coin, totalUSDTRecieve, totalVNTGive, status, exchangePrice, profit, perProfit)

		allFundMessage := exchangeClient.GetFundsMessages()
		text = fmt.Sprintf("%s %s", text, allFundMessage)
	}

	// Calculate USDT Margin
	if currentExchange == "BINANCE" {
		name := "USDT"
		marginDetails, _ := exchangeClient.Bn.GetMarginDetails()
		netAsset := calculateUSDTMargin(marginDetails, name)
		text = fmt.Sprintf("%s \n USDT(Margin): %v", text, netAsset)
		notifyWhenAssetIsLow(netAsset, text)
	}
	go teleClient.SendMessage(text, chatProfitID)
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

func getOrderDetails(id string, coin string, exchangeOrderID int64, origClientOrderID string) *exchanges.OrderResp {
	var orderDetails *exchanges.OrderResp
	var err error
	for j := 0; j <= 2; j++ {
		orderDetails, err = exchangeClient.GetOrder(coin, exchangeOrderID, origClientOrderID)
		if err != nil {
			text := fmt.Sprintf("%s %s Err getOrderDetails: %s", coin, id, err)
			go teleClient.SendMessage(text, chatErrorID)
		} else {
			if isFilledStatus(orderDetails.Status) {
				break
			}
		}

		if j == 2 {
			log.Printf("Unsucessfully sell in %s after 1 minute", currentExchange)
			break
		}
		time.Sleep(30 * time.Second)
	}
	return orderDetails
}

func isFilledStatus(status string) bool {
	if currentExchange == "FTX" {
		return status == ftx.ORDER_CLOSED
	}
	return status == binance.ORDER_FILLED
}
