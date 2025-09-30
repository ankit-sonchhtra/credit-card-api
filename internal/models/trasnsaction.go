package models

type CreateTransactionRequest struct {
	AccountId     string  `json:"accountId" example:"92d68c0e-dafe-406a-a0f2-8faae2020947" binding:"required,uuid"`
	OperationType string  `json:"operationType" example:"cash purchase | installment purchase | withdrawal | payment"`
	Amount        float64 `json:"amount" example:"-123.45"`
}

type CreateTransactionResponse struct {
	TransactionId string `json:"transactionId" example:"82dfa288-28b5-430f-9b54-a4f99e546a40"`
}
