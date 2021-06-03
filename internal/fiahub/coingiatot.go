package fiahub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	u "gitlab.com/fiahub/bot/internal/utils"
)

type CoinGiaTotParams struct {
	AutoMode          int     `json:"AutoMode"`
	ProfitMax         float64 `json:"ProfitMax"`
	ProfitPerThousand float64 `json:"ProfitPerThousand"`
	Spread            float64 `json:"Spread"`
	USDTMax           float64 `json:"USDTMax"`
	USDTMidPoint      float64 `json:"USDTMidPoint"`
	USDTOffset2       float64 `json:"USDTOffset2"`
}

func GetCoinGiaTotParams() *CoinGiaTotParams {
	BASE_URL := os.Getenv("COINGIATOT_URL")
	BOT_NAME := os.Getenv("BOT_NAME")
	url := fmt.Sprintf("%s/bot_vars?bot_name=%s", BASE_URL, BOT_NAME)
	body, _, err := u.HttpGet(url, nil)
	if err != nil {
		log.Printf("Err in RenewParam with Body: %s, Err: %s", body, err.Error())
		return nil
	}

	var params *CoinGiaTotParams
	if err := json.Unmarshal([]byte(body), params); err != nil {
		panic(err)
	}
	return params
}
