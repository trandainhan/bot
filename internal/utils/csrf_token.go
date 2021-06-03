package utils

import (
	"encoding/base64"
	"strconv"
	"time"
)

func GenerateFiahubCSRFToken(email string) string {
	seed := "t.a=function(e,t){return i()(r()(e,{raw:{value:i()(t)}}))};"
	loc, _ := time.LoadLocation("Asia/Singapore")
	dayOfMonth := time.Now().In(loc).Day()
	pre_key := seed + "_" + strconv.Itoa(dayOfMonth)
	key := base64.StdEncoding.Strict().EncodeToString([]byte(pre_key))
	return GenerateHmac(email, key)
}
