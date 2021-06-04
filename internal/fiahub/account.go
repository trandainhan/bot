package fiahub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gitlab.com/fiahub/bot/internal/utils"
)

type LoginRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	CSRFToken string `json:"crfs_token"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func Login(email, password string) string {
	url := os.Getenv("FIAHUB_URL")
	data := LoginRequest{
		Email:     email,
		Password:  password,
		CSRFToken: utils.GenerateFiahubCSRFToken(email),
	}
	body, _, err := utils.HttpPost(fmt.Sprintf("%s/sessions", url), data, nil)
	if err != nil {
		log.Printf("Err Login Body: %s", body)
		panic(err)
	}
	var result LoginResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		panic(err)
	}
	return result.Token
}
