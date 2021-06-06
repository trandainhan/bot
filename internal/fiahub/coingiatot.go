package fiahub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	u "gitlab.com/fiahub/bot/internal/utils"
)

type CoinGiaTotParams struct {
	AutoMode          string `json:"AutoMode"`
	ProfitMax         string `json:"ProfitMax"`
	ProfitPerThousand string `json:"ProfitPerThousand"`
	Spread            string `json:"Spread"`
	USDTMax           string `json:"USDTMax"`
	USDTMidPoint      string `json:"USDTMidPoint"`
	USDTOffset2       string `json:"USDTOffset2"`
}

func (cgt CoinGiaTotParams) GetAutoMode() int {
	autoMode, _ := strconv.Atoi(cgt.AutoMode)
	return autoMode
}

func (cgt CoinGiaTotParams) GetProfitMax() float64 {
	profitMax, _ := strconv.ParseFloat(cgt.ProfitMax, 64)
	return profitMax
}

func (cgt CoinGiaTotParams) GetProfitPerThousand() float64 {
	profitPerThousand, _ := strconv.ParseFloat(cgt.ProfitPerThousand, 64)
	return profitPerThousand
}

func (cgt CoinGiaTotParams) GetSpread() float64 {
	spread, _ := strconv.ParseFloat(cgt.Spread, 64)
	return spread
}

func (cgt CoinGiaTotParams) GetUSDTMax() float64 {
	usdtMax, _ := strconv.ParseFloat(cgt.USDTMax, 64)
	return usdtMax
}

func (cgt CoinGiaTotParams) GetUSDTMidPoint() float64 {
	uSDTMidPoint, _ := strconv.ParseFloat(cgt.USDTMidPoint, 64)
	return uSDTMidPoint
}

func (cgt CoinGiaTotParams) GetUSDTOffset2() float64 {
	uSDTOffset2, _ := strconv.ParseFloat(cgt.USDTOffset2, 64)
	return uSDTOffset2
}

type Coingiatotresp struct {
	Params CoinGiaTotParams `json:"bot_vars"`
}

func GetCoinGiaTotParams() *CoinGiaTotParams {
	BASE_URL := os.Getenv("COINGIATOT_URL")
	BOT_NAME := os.Getenv("BOT_NAME")
	url := fmt.Sprintf("%s/bot_vars?bot_name=%s", BASE_URL, BOT_NAME)
	body, _, err := u.HttpGet(url, nil)
	log.Println(body)
	if err != nil {
		log.Printf("Err in RenewParam with Body: %s, Err: %s", body, err.Error())
		return nil
	}

	var resp Coingiatotresp
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		panic(err)
	}
	return &resp.Params
}
