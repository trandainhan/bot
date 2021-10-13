package main

import (
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
)

func monitorOrder(order *exchanges.OrderResp, orderChan chan<- *exchanges.OrderResp) {
	log.Printf("Start monitor order: %d", order.ID)
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
			log.Printf("%s %s Order %d is canceled at price %f", coin, side, orderDetails.ID, orderDetails.Price)
			orderChan <- orderDetails
			updateOrderCount(side)
			break
		}
		i++
		if i == 60 {
			log.Printf("%s %s Order %d at price %f is not filled after 1 hours, stop monitor", coin, side, orderDetails.ID, orderDetails.Price)
			updateOrderCount(side)
			break
		}
		time.Sleep(1 * time.Minute)
	}
}

func updateOrderCount(side string) {
	if side == "buy" {
		decreaseOpenBuyOrder()
	} else if side == "sell" {
		decreaseOpenSellOrder()
	}
}
