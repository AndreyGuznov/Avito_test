package entity

import "time"

type User struct {
	Id      int     `json:"id" db:"id"`
	Balance float64 `json:"balance" db:"balance"`
}

type Transanction struct {
	UserId    int       `json:"userId" db:"user_id"`
	Amount    float64   `json:"amount" db:"amount"`
	Operation string    `json:"operation" db:"operation"`
	Date      time.Time `json:"date" db:"date"`
}
