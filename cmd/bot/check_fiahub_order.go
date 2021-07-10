package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
	"gitlab.com/fiahub/bot/internal/fiahub"
)

func checkFiahubOrder(botID string, fiahubOrderID int, originalCoinAmount float64,
	oldPriceOnExchange float64, perCancel float64, orderType string) (float64, bool) {

	executedQty := 0.0
	matching := false
	for {
		order, err := fia.GetOrderDetails(fiahubOrderID)
		if err != nil {
			text := fmt.Sprintf("Error %s IDTrade: %s GetOrderDetails %s fiahubOrderID: %d", coin, botID, err, fiahubOrderID)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(1 * time.Second)
			continue
		}
		state := order.State
		coinAmount := order.CoinAmount
		executedQty = originalCoinAmount - coinAmount
		if state == fiahub.ORDER_CANCELLED || state == fiahub.ORDER_FINISHED {
			matchingTX, _ := fia.GetSelfMatchingTransaction(order.UserID, order.ID)
			matching = matchingTX != nil
			break
		}

		// Trigger cancel process
		newPriceOnExchange := 0.0
		if orderType == "AskOrder" {
			newPriceOnExchange, err = exchanges.GetAskPriceByQuantity(coin, quantityToGetPrice)
		} else {
			newPriceOnExchange, err = exchanges.GetBidPriceByQuantity(coin, quantityToGetPrice)
		}

		if err != nil {
			text := fmt.Sprintf("%s Err GetPriceByQuantity inside the loop: %s", coin, err.Error())
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(1 * time.Second)
			continue
		}
		perChange := math.Abs((newPriceOnExchange - oldPriceOnExchange) / oldPriceOnExchange)
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
			_, _, err := fia.CancelOrder(fiahubOrderID)
			if err != nil {
				text := fmt.Sprintf("Error %s IDTrade: %s Err CancelOrder: %d with error: %s", coin, botID, fiahubOrderID, err)
				log.Println(text)
				go teleClient.SendMessage(text, chatErrorID)
				time.Sleep(3 * time.Second)
				continue
			}
			continue
		}
		time.Sleep(5 * time.Second)
	}
	return executedQty, matching
}
