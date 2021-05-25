package fiahub

import (
	"encoding/json"
	"fmt"
	"os"

	"gitlab.com/fiahub/bot/internal/telegram"
	u "gitlab.com/fiahub/bot/internal/utils"
)

type RateData struct {
	Rates map[string]float64 `json:"rates"`
}

type Rates struct {
	USD RateData `json:"usd"`
	VND RateData `json:"vnd"`
}

func GetUSDVNDRate() (float64, error) {
	var BASE_URL = os.Getenv("fiahub_url")
	url := fmt.Sprintf("%s/vars/currency_rates", BASE_URL)

	body, _, err := u.HttpGet(url)
	if err != nil {
		return 0.0, err
	}
	var rates *Rates
	err = json.Unmarshal([]byte(body), rates)
	if err != nil {
		return 0.0, err
	}

	teleClient := telegram.NewTeleBot(os.Getenv("tele_bot_token"))
	vndRate := rates.USD.Rates["vnd"]
	var text string
	if vndRate > 22000 && vndRate < 25000 {
		text = fmt.Sprintf("USD->VND: %v", vndRate)
		go teleClient.SendMessage(text, -357553425)
	} else {
		text = fmt.Sprintf("BAO DONG USD_VNDRATE khong nam trong range 22000-25000 @ndtan")
		go teleClient.SendMessage(text, -357553425)
	}
	return vndRate, nil
}
