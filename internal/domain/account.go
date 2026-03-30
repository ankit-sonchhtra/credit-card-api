package domain

import "time"

type Account struct {
	Id             int64
	DocumentNumber string
	CreatedAt      time.Time
}

type CreateAccountParam struct {
	DocumentNumber string
}
