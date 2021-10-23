package technical_analysis

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gitlab.com/fiahub/bot/internal/utils"
)

type MACD struct {
	Value     float64 `json:"valueMACD"`
	Signal    float64 `json:"valueMACDSignal"`
	Hint      float64 `json:"valueMACDHist"`
	Backtrack int     `json:"backtrack"`
}

/**
	Value is MACD value: EMA (12) – EMA (26)
	Signal: EMA (9)
	Hint = MACD - Signal
  Hint from negative to positive, up trend signal => Buy
  Hint from negative to positive, up trend signal => Buy
**/
func GetMACD(exchange, coin, interval string, backtracks int) ([]MACD, error) {
	template := BASE_ENDPOINT + "/macd?secret=%s&exchange=%s&symbol=%s/USDT&interval=%s&backtracks=%d"
	finalUrl := fmt.Sprintf(template, os.Getenv("TAAPI_KEY"), exchange, coin, interval, backtracks)
	body, code, err := utils.HttpGet(finalUrl, nil)
	if err != nil || code >= 400 {
		log.Printf("Err Get MACD value %s", body)
		return nil, err
	}
	var macd []MACD
	err = json.Unmarshal([]byte(body), &macd)
	if err != nil {
		log.Printf("Err getMACD Unmarshal %s", body)
		return nil, err
	}
	return macd, nil
}
