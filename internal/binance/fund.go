package binance

// Binance_CheckFundAllGetMessage()
// {
// 	content := BinanceAPI_CheckFund()
// 	parsed := JSON.Load(content)
// 	Maxindex := parsed.balances.Maxindex()
// 	Message1 := "Binance Funds:  "
// 	Message2 := "%0A Inorder:  "
//
//
// 	Loop, %Maxindex%
// 	{
// 		asset := parsed.balances[a_index].asset
// 		freefund := parsed.balances[a_index].free
// 		lockedfund := parsed.balances[a_index].locked
//
// 		if((freefund > 0) or (lockedfund >0))
// 		{
// 			Message1 := Message1 . freefund . " " . asset . "  -  "
// 			Message2 := Message2 . lockedfund . " " . asset . "  -  "
// 		}
// 	}
// 	Message := Message1 . Message2
// 	return Message
// }

type Balance struct {
	Asset      string  `json:"asset"`
	Free       float64 `json:"fee"`
	LockedFund float64 `json:"locked"`
}

type Fund struct {
	Balances []Balance `json:"balances"`
}

func CheckFund() {
	// URL := "https://api.binance.com/api/v3/account?"
	//
	// TimeUNIX = %A_NowUTC%
	// TimeUNIX -= 19700101000000,seconds
	// FileReadLine, offset, %A_ScriptDir%\file\Offsetime.txt, 1
	// TimeUNIX := (TimeUNIX-5)*1000 + offset
	//
	// queryString := "&recvWindow=59000&timestamp=" . TimeUNIX
	//
	// content := GetContentWithAPI("GET",link,queryString)
	//
	// return content

}

func GetFundsMessages() string {
	return ""
}

type UserAsset struct {
	Name     string
	NetAsset float64
}

type MarginDetailsResposne struct {
	UserAssets []UserAsset `json:"userAssets"`
}

func GetMarginDetails() MarginDetailsResposne {
	userAssets := []UserAsset{}
	return MarginDetailsResposne{UserAssets: userAssets}
}
