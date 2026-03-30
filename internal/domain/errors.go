package domain

import "github.com/credit-card-api/pkg/constants"

type AppError struct {
	Code    string
	Message string
}

var (
	ErrAccountAlreadyExist        = &AppError{Code: constants.AccountAlreadyExistErrCode, Message: "account already exists."}
	ErrAccountNotFound            = &AppError{Code: constants.AccountNotFoundErrCode, Message: "account does not exists with provided id."}
	ErrInvalidOperationType       = &AppError{Code: constants.InvalidOperationTypeErrCode, Message: "operation type ID provided is not supported by the system."}
	ErrTransactionAccountNotFound = &AppError{Code: constants.TransactionAccountNotFoundErrCode, Message: "account does not exist with provided id."}
	ErrInternal                   = &AppError{Code: constants.InternalServerErrCode, Message: "an unexpected error occurred."}
)

func (e *AppError) Error() string {
	return e.Message
}
