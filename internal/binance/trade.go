package binance

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"gitlab.com/fiahub/bot/internal/utils"
)

func (binance Binance) SellLimit(marketParam string, price, quantity float64) (*OrderDetailsResp, error) {
	params := map[string]string{
		"symbol":      marketParam,
		"side":        "SELL",
		"type":        "LIMIT",
		"timeInForce": "GTC",
		"quantity":    fmt.Sprintf("%f", quantity),
		"price":       fmt.Sprintf("%f", price),
	}
	return binance.makeTradeRequest(params)
}

func (binance Binance) BuyLimit(marketParam string, price float64, quantity float64) (*OrderDetailsResp, error) {
	params := map[string]string{
		"symbol":      marketParam,
		"side":        "BUY",
		"type":        "LIMIT",
		"timeInForce": "GTC",
		"quantity":    fmt.Sprintf("%f", quantity),
		"price":       fmt.Sprintf("%f", price),
	}
	return binance.makeTradeRequest(params)
}

func (binance Binance) makeTradeRequest(params map[string]string) (*OrderDetailsResp, error) {
	body, _, err := binance.makeRequest("POST", params, "/api/v3/order")
	if err != nil {
		log.Printf("Err makeTradeRequest: %s", err.Error())
	}
	var order *OrderDetailsResp
	err = json.Unmarshal([]byte(body), order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (binance Binance) makeRequest(httpType string, params map[string]string, postURL string) (string, int, error) {
	var BASE_URL = os.Getenv("BINANCE_URL")
	now := time.Now()
	sec := now.UnixNano()
	timeDifferences := binance.RedisClient.GetInt64("local_binance_time_difference")
	mili := (sec-5)/int64(time.Millisecond) + timeDifferences
	params["recvWindow"] = "59000"
	params["timestamp"] = fmt.Sprintf("%v", mili)

	queryString := utils.BuildQueryStringFromMap(params)

	binanceAPIKey := os.Getenv("BINANCE_API_KEY")
	binanceAPISecret := os.Getenv("BINANCE_API_SECRET")

	hmac := utils.GenerateHmac(queryString, binanceAPISecret)
	params["signature"] = hmac

	final_url, _ := utils.BuildUrlWithParams(fmt.Sprintf("%s/order", BASE_URL), params)

	headers := map[string]string{
		"X-MBX-APIKEY": binanceAPIKey,
	}
	if httpType == "POST" {
		return utils.HttpPost(final_url, nil, &headers)
	}
	return utils.HttpGet(final_url, &headers)
}
