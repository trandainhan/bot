package fiahub

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	u "gitlab.com/fiahub/bot/internal/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

func Login(email, password string) string {
	url := os.Getenv("FIAHUB_URL")
	data := LoginRequest{
		Email:    email,
		Password: password,
	}
	body, _, err := u.HttpPost(fmt.Sprintf("%s/sessions", url), data, nil)
	if err != nil {
		log.Println(body)
		panic(err)
	}
	var result LoginResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		panic(err)
	}
	return result.Token
}
