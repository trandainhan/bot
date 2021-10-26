package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gitlab.com/fiahub/bot/internal/exchanges"
)

func reviveMonitorOrder() {
	redisValue, err := redisClient.Get(coin + "_open_orders")
	if err != nil { // Does not exist
		return
	}
	var orders []exchanges.OrderResp
	err = json.Unmarshal([]byte(redisValue), &orders)
	if err != nil {
		log.Printf("Err reviveMonitorOrder: %s", err.Error())
		return
	}
	num := len(orders)
	if num == 0 {
		return
	}
	for i, _ := range orders {
		order := &orders[i] // Notice here
		go monitorOrder(order)
	}
}

func monitorOrder(order *exchanges.OrderResp) {
	log.Printf("Start monitor order: %d", order.ID)
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
			break
		} else if orderDetails.IsCanceled() {
			log.Printf("%s %s Order %d is canceled at price %f", coin, side, orderDetails.ID, orderDetails.Price)
			break
		}
		time.Sleep(1 * time.Minute)
	}
}
