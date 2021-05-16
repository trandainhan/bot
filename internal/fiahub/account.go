package fiahub

import (
	"encoding/json"
	u "gitlab.com/fiahub/bot/internal/utils"
	"os"
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
	url := os.Getenv("fiahubURL")
	data := LoginRequest{
		Email:    email,
		Password: password,
	}
	body, _, _ := u.HttpPost(url, data)

	var result LoginResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		panic(err)
	}
	return result.Token
}
