package model

type AccountDocument struct {
	AccountId      string  `bson:"account_id"`
	UserId         string  `bson:"user_id"`
	DocumentNumber string  `bson:"document_number"`
	CurrentBalance float64 `bson:"current_balance"`
	Status         string  `bson:"status"`
	CreatedAt      int64   `bson:"created_at"`
	UpdatedAt      int64   `bson:"updated_at"`
}
