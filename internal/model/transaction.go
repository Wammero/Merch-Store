package model

type Transaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type CoinHistory struct {
	Received []Transaction `json:"received"`
	Sent     []Transaction `json:"sent"`
}
