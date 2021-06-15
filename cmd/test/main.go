package main

import (
	"context"
	// "gitlab.com/fiahub/bot/internal/fiahub"
	"log"
	"os"
	"time"

	// "gitlab.com/fiahub/bot/internal/binance"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
)

func main() {
	// Test login
	// email := "trdainhan@gmail.com"
	// password := ""
	// token := fiahub.Login(email, password)
	// log.Println(token)
	// rate, _ := fiahub.GetUSDVNDRate()
	// log.Println(rate)

	// params := fiahub.GetCoinGiaTotParams()
	// log.Println(params)

	// Test redis
	ctx := context.Background()
	redisURL := os.Getenv("REDIS_URL")
	redisClient := rediswrapper.NewRedisClient(ctx, redisURL)
	redisClient.Set("nhantran", "hello")
	log.Println(redisClient.Get("nhantran"))

	// bid_order: {coin_amount: 3, price_per_unit_cents: 23731, type: "BidOrder"
	// coin: "USDT", currency: "VNT"}
	// order := fiahub.Order{
	// 	CoinAmount:        3,
	// 	PricePerUnitCents: 7848,
	// 	Type:              "AskOrder",
	// 	Coin:              "DOGE",
	// 	Currency:          "VNT",
	// }
	// resp, _ := fiahub.CreateAskOrder(token, order)
	// log.Println(resp)
	// resp, _, _ := fiahub.CancelOrder(token, 103184704)
	// log.Println(resp)

	// fia := fiahub.Fiahub{
	// 	RedisClient: redisClient,
	// 	Token:       token,
	// }
	//
	// detail, _, err := fia.GetAskOrderDetails(103411475)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(detail)

	// Test binance api
	// offset := binance.GetOffsetTimeUnix()
	// redisClient.Set("local_binance_time_difference", offset)

	// bn := binance.Binance{
	// 	RedisClient: redisClient,
	// }
	//
	// usdtFund := bn.CheckFund("USDT")
	// log.Println(usdtFund)
	//
	// msg := bn.GetFundsMessages()
	// log.Println(msg)

	// detail, _ := bn.GetMarginDetails()
	// log.Println(detail)

	// bidPriceByQuantity, askPriceByQuantity := binance.GetPriceByQuantity("DOGEUSDT", 8.0)
	// log.Printf("bidPriceByQuantity: %v", bidPriceByQuantity)
	// log.Printf("askPriceByQuantity: %v", askPriceByQuantity)
	// resp, err := bn.BuyLimit("DOGEUSDT", 0.3, 100)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(resp.OrderID)
	// log.Println(resp)

	// resp, _ := bn.GetOrder("DOGEUSDT", 1239188099, "SLYJI2yBT99GaIo4qc35iM")
	// log.Println(resp)

	now := time.Now()
	miliTime := now.UnixNano() / int64(time.Millisecond)
	log.Println(now.UnixNano())
	log.Println(miliTime)
	time.Sleep(1 * time.Second)
	now = time.Now()
	miliTime2 := now.UnixNano() / int64(time.Millisecond)
	log.Println(miliTime2 - miliTime)

}
