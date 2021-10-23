package technical_analysis

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gitlab.com/fiahub/bot/internal/utils"
)

type RSI struct {
	Value     float64 `json:"value"`
	Backtrack int     `json:"Backtrack"`
}

/**
	Value is RSI value
**/
func GetRSI(exchange, coin, interval string, backtracks int) ([]RSI, error) {
	template := BASE_ENDPOINT + "/rsi?secret=%s&exchange=%s&symbol=%s/USDT&interval=%s&backtracks=%d"
	finalUrl := fmt.Sprintf(template, os.Getenv("TAAPI_KEY"), exchange, coin, interval, backtracks)
	body, code, err := utils.HttpGet(finalUrl, nil)
	if err != nil || code >= 400 {
		log.Printf("Err Get RSI value %s", body)
		return nil, err
	}
	var rsi []RSI
	err = json.Unmarshal([]byte(body), &rsi)
	if err != nil {
		log.Printf("Err getRSI Unmarshal %s", body)
		return nil, err
	}
	return rsi, nil
}
