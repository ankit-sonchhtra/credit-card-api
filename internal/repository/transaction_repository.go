package repository

//go:generate mockgen -source=transaction_repository.go -destination=mocks/mock_transaction_repository.go -package=mocks

import (
	"context"
	"fmt"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/repository/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	logger "github.com/sirupsen/logrus"
)

type TransactionRepository interface {
	Create(ctx context.Context, transactionParam domain.CreateTransactionParam) (*domain.Transaction, error)
}

type transactionRepository struct {
	querier sqlc.Querier
}

func NewTransactionRepository(querier sqlc.Querier) TransactionRepository {
	return &transactionRepository{querier: querier}
}

func (tr *transactionRepository) Create(ctx context.Context, transactionParam domain.CreateTransactionParam) (*domain.Transaction, error) {
	transaction, err := tr.getQuerier(ctx).CreateTransaction(ctx, sqlc.CreateTransactionParams{
		AccountID:       transactionParam.AccountId,
		OperationTypeID: transactionParam.OperationTypeId,
		Amount:          float64ToNumeric(transactionParam.Amount),
	})

	if err != nil {
		logger.Errorf("error while create transaction: %s", err.Error())
		return nil, err
	}
	logger.Info("transaction created successfully in db.")
	return mapToDomainTransaction(transaction), nil
}

func float64ToNumeric(val float64) pgtype.Numeric {
	var n pgtype.Numeric
	err := n.Scan(fmt.Sprintf("%.2f", val))
	if err != nil {
		n.Valid = false
	}
	return n
}

func numericToFloat64(n pgtype.Numeric) float64 {
	if !n.Valid {
		return 0.0
	}
	f, _ := n.Float64Value()
	return f.Float64
}

func mapToDomainTransaction(transaction sqlc.Transaction) *domain.Transaction {
	return &domain.Transaction{
		Id:              transaction.TransactionID,
		AccountId:       transaction.AccountID,
		OperationTypeId: transaction.OperationTypeID,
		Amount:          numericToFloat64(transaction.Amount),
		CreatedAt:       transaction.CreatedAt.Time,
	}
}

func (tr *transactionRepository) getQuerier(ctx context.Context) sqlc.Querier {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return sqlc.New(tx)
	}
	return tr.querier
}
