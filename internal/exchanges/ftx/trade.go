package ftx

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gitlab.com/fiahub/bot/internal/utils"
)

func (ftx FtxClient) SellLimit(marketParam string, price, quantity float64) (*Order, error) {
	newOrder := NewOrder{
		Market:     marketParam,
		Side:       "sell",
		Type:       "limit",
		Price:      price,
		Size:       quantity,
		ReduceOnly: false,
		Ioc:        false,
		PostOnly:   false,
	}
	return ftx.makeTradeRequest(newOrder)
}

func (ftx FtxClient) BuyLimit(marketParam string, price float64, quantity float64) (*Order, error) {
	newOrder := NewOrder{
		Market:     marketParam,
		Side:       "buy",
		Type:       "limit",
		Price:      price,
		Size:       quantity,
		ReduceOnly: false,
		Ioc:        false,
		PostOnly:   false,
	}
	return ftx.makeTradeRequest(newOrder)
}

func (ftx FtxClient) makeTradeRequest(data NewOrder) (*Order, error) {
	requestBody, _ := json.Marshal(data)
	body, _, err := ftx.makeRequest("POST", "/orders", string(requestBody))
	if err != nil {
		log.Printf("Error ftx makeTradeRequest: %s", err.Error())
		return nil, err
	}
	var resp NewOrderResponse
	err = utils.ProcessResponse(body, &resp)
	if err != nil {
		return nil, err
	}
	log.Printf("Successfully make an order in ftx: %v", resp)
	return &resp.Result, nil
}

func (ftx FtxClient) makeRequest(method, path, requestBody string) (string, int, error) {
	BASE_URL := os.Getenv("FTX_URL")

	ts := strconv.FormatInt(time.Now().UTC().Unix()*1000, 10)
	signaturePayload := ts + method + "/api" + path + requestBody
	ftxAPIKey := os.Getenv("FTX_API_KEY")
	ftxAPISecret := os.Getenv("FTX_API_SECRET")

	hmac := utils.GenerateHmac(signaturePayload, ftxAPISecret)

	final_url := fmt.Sprintf("%s/api%s", BASE_URL, path)

	headers := map[string]string{
		"FTX-KEY":  ftxAPIKey,
		"FTX-TS":   ts,
		"FTX-SIGN": hmac,
	}
	if method == "POST" {
		return utils.HttpPost(final_url, requestBody, &headers)
	}
	return utils.HttpGet(final_url, &headers)
}
