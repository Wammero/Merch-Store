package model

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserInfoResponse struct {
	Coins       int64           `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}
