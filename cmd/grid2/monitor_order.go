package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
)

func monitorOrder(order *exchanges.OrderResp, orderChan chan<- *exchanges.OrderResp) {
	i := 0
	side := order.GetSide()
	for {
		orderDetails, err := exchangeClient.GetOrder(coin, order.ID, order.ClientID)
		if err != nil {
			text := fmt.Sprintf("%s Order %d Err getOrderDetails: %s", coin, order.ID, err)
			log.Println(text)
			go teleClient.SendMessage(text, chatErrorID)
			time.Sleep(60 * time.Second)
			continue
		}
		if orderDetails.IsFilled() {
			calculateProfit(orderDetails.ID, orderDetails.ExecutedQty, orderDetails.Price, side)
			orderChan <- orderDetails
			break
		} else if orderDetails.IsCanceled() {
			log.Printf("%s %s Order %d is canceled at price %f", coin, orderDetails.GetSide(), orderDetails.ID, orderDetails.Price)
			if side == "buy" {
				decreaseOpenBuyOrder()
			} else if side == "sell" {
				decreaseOpenSellOrder()
			}
			orderChan <- orderDetails
			break
		}
		i++
		if i%60 == 0 {
			log.Printf("%s %s Order %d at price %f is not filled after %d hours", coin, side, orderDetails.ID, orderDetails.Price, i)
		}
		time.Sleep(1 * time.Minute)
	}
}
