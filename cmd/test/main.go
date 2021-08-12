package main

import (
	"context"
	// "log"
	"os"

	// "time"

	// "gitlab.com/fiahub/bot/internal/binance"

	// "gitlab.com/fiahub/bot/internal/exchanges/ftx"
	"gitlab.com/fiahub/bot/internal/rediswrapper"
)

type NewOrder struct {
	Future string `json:"future"`
	Market string `json:"market"`
}

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
	redisClient := rediswrapper.NewRedisClient(ctx, redisURL, 1)
	redisClient.Set("nhantran", "hello")
	// log.Println(redisClient.Get("nhantran"))

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

	// db := pg.Connect(&pg.Options{
	// 	Addr:     os.Getenv("DATABASE_ADDR"),
	// 	User:     os.Getenv("DATABASE_USERNAME"),
	// 	Password: os.Getenv("DATABASE_PASSWORD"),
	// 	Database: os.Getenv("DATABASE_NAME"),
	// })
	//
	// fia := fiahub.Fiahub{
	// 	RedisClient: redisClient,
	// 	DB:          db,
	// }
	// result, _ := fia.GetOrderDetails(100000000)
	// log.Println(result)
	// log.Println(result.ID)
	// log.Println(result.UserID)

	// tx, err := fia.GetSelfMatchingTransaction(result.UserID, result.ID)
	// if err == pg.ErrNoRows {
	// 	log.Println("Should be fine")
	// }
	// matching := tx != nil
	// log.Println(matching)

	// detail, _, err := fia.GetAskOrderDetails(103411475)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println(detail)

	// Test binance api
	// offset := binance.GetOffsetTimeUnix()
	// log.Println(offset)

	// ftx := ftx.FtxClient{}
	// log.Println(string(res))
	// usdtFund, err := ftx.CheckFund("USDT")
	// log.Println(err)
	// log.Println(usdtFund)
	// order, err := ftx.SellLimit("ETH/USDT", 3200, 0.01)
	// log.Println(err)
	// log.Println(order)

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

	// now := time.Now().UTC().AddDate(0, 0, -3)
	// log.Println(now.Format("2006-01-02 15:04:05"))
	// miliTime := now.UnixNano() / int64(time.Millisecond)
	// log.Println(now.UnixNano())
	// log.Println(miliTime)
	// time.Sleep(1 * time.Second)
	// now = time.Now()
	// miliTime2 := now.UnixNano() / int64(time.Millisecond)
	// log.Println(miliTime2 - miliTime)

}
