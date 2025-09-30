package model

type TransactionDocument struct {
	Id            string  `bson:"transaction_id"`
	AccountId     string  `bson:"account_id"`
	OperationType string  `bson:"operation_type"`
	Amount        float64 `bson:"amount"`
	CreatedAt     int64   `bson:"created_at"`
	UpdatedAt     int64   `bson:"updated_at"`
}
