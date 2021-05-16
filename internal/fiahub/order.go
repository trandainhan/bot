package fiahub

import (
	// u "gitlab.com/fiahub/bot/internal/utils"
	"os"
)

var BASE_URL = os.Getenv("fiahub_url")

const (
	ORDER_CANCELLED = "cancelled"
	ORDER_FINISHED  = "finished"
)

type Order struct {
	Coin               string
	OriginalCoinAmount float64
	CoinAmount         float64
	PriceSellRandom    float64
	Currency           string
	Type               string
}

type OrderDetails struct {
	ID         string  `json:"id"`
	State      string  `json:"state"`
	CoinAmount float64 `json:"coin_amount"`
	Matching   bool    `json:"matching"`
}

func CancelAllOrder(token string) {
	// write cancel time here
}

func CancelOrder(token string, orderID string) (OrderDetails, int, error) {
	result := OrderDetails{}
	return result, 200, nil
}

func CreateAskOrder(token string, askOrder Order) (string, int, error) {
	return "", 200, nil
}

func GetOrderDetails(token string, orderID string) (OrderDetails, int, error) {
	// u.HttpGet(BASE_URL)
	result := OrderDetails{}
	return result, 200, nil
}
