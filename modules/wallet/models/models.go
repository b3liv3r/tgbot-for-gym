package models

import "time"

type Transaction struct {
	Id          int
	UserId      int
	Amount      float64
	Type        string
	Description string
	Date        time.Time
}

type Wallet struct {
	UserID  int
	Balance float64
}
