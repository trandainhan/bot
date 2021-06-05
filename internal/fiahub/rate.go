package fiahub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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
	var BASE_URL = os.Getenv("FIAHUB_URL")
	url := fmt.Sprintf("%s/vars/currency_rates", BASE_URL)

	body, _, err := u.HttpGet(url, nil)
	if err != nil {
		log.Printf("Err GetUSDVNDRate: %s", err.Error())
		return 0.0, err
	}
	var rates Rates
	err = json.Unmarshal([]byte(body), &rates)
	if err != nil {
		return 0.0, err
	}

	teleClient := telegram.NewTeleBot(os.Getenv("TELE_BOT_TOKEN"))
	vndRate := rates.USD.Rates["VND"]
	var text string
	if vndRate > 22000 && vndRate < 25000 {
		text = fmt.Sprintf("USD->VND: %v", vndRate)
		chatID, _ := strconv.ParseInt(os.Getenv("CHAT_ID"), 10, 64)
		go teleClient.SendMessage(text, chatID)
	} else {
		text = fmt.Sprintf("BAO DONG USD_VNDRATE khong nam trong range 22000-25000 %s", os.Getenv("TELEGRAM_HANDLER"))
		chatErrorID, _ := strconv.ParseInt(os.Getenv("CHAT_ERROR_ID"), 10, 64)
		go teleClient.SendMessage(text, chatErrorID)
	}
	return vndRate, nil
}
