package main

func incrementOpenBuyOrder() {
	key := coin + "_open_buy_order"
	total := redisClient.GetInt64(key)
	total = total + 1
	redisClient.Set(key, total, 0)
}

func decreaseOpenBuyOrder() {
	key := coin + "_open_buy_order"
	total := redisClient.GetInt64(key)
	total = total - 1
	redisClient.Set(key, total, 0)
}

func incrementOpenSellOrder() {
	key := coin + "_open_sell_order"
	total := redisClient.GetInt64(key)
	total = total + 1
	redisClient.Set(key, total, 0)
}

func decreaseOpenSellOrder() {
	key := coin + "_open_sell_order"
	total := redisClient.GetInt64(key)
	total = total - 1
	redisClient.Set(key, total, 0)
}

func resetOpenBuySellOrders() {
	redisClient.Set(coin+"_open_sell_order", 0, 0)

	redisClient.Set(coin+"_open_buy_order", 0, 0)
}
