package domain

import "time"

type CreateTransactionParam struct {
	AccountId       int64
	OperationTypeId int64
	Amount          float64
}

type Transaction struct {
	Id              int64
	AccountId       int64
	OperationTypeId int64
	Amount          float64
	CreatedAt       time.Time
}

type TransactionType struct {
	Id         int64
	IsNegative bool
}

var ValidOperations = map[int64]TransactionType{
	1: {Id: 1, IsNegative: true},  // Normal Purchase
	2: {Id: 2, IsNegative: true},  // Purchase with Installments
	3: {Id: 3, IsNegative: true},  // Withdrawal
	4: {Id: 4, IsNegative: false}, // Credit Voucher
}
