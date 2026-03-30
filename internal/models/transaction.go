package models

import (
	"github.com/go-playground/validator/v10"
)

type TransactionRequest struct {
	AccountId       int64   `json:"account_id" example:"1" validate:"required"`
	OperationTypeId int64   `json:"operation_type_id" example:"1" validate:"required"`
	Amount          float64 `json:"amount" example:"123.45" validate:"required,gt=0"`
}

type CreateTransactionResponse struct {
	TransactionId int64 `json:"transaction_id" example:"1"`
}

func (request TransactionRequest) Validate() error {
	err := validator.New().Struct(&request)
	return translateError(err)
}
