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
	GetAllTransactions(ctx context.Context, accountId int64) ([]domain.Transaction, error)
	UpdateTransactionById(ctx context.Context, transactionId int64, balance float64) error
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
		Balance:         float64ToNumeric(transactionParam.Balance),
	})

	if err != nil {
		logger.Errorf("error while create transaction: %s", err.Error())
		return nil, err
	}
	logger.Info("transaction created successfully in db.")
	return mapToDomainTransaction(transaction), nil
}

func (tr *transactionRepository) GetAllTransactions(ctx context.Context, accountId int64) ([]domain.Transaction, error) {
	var transactionList []domain.Transaction
	transactions, err := tr.getQuerier(ctx).GetAllTransactionById(ctx, accountId)
	if err != nil {
		logger.Errorf("error while fetch all transactions: %s", err.Error())
		return nil, err
	}
	for _, tx := range transactions {
		transactionList = append(transactionList, *mapToDomainTransaction(tx))
	}
	return transactionList, nil
}

func (tr *transactionRepository) UpdateTransactionById(ctx context.Context, transactionId int64, balance float64) error {
	_, err := tr.getQuerier(ctx).UpdateTransaction(ctx, sqlc.UpdateTransactionParams{
		TransactionID: transactionId,
		Balance:       float64ToNumeric(balance),
	})
	if err != nil {
		logger.Errorf("error while fetch all transactions: %s", err.Error())
		return err
	}
	return nil
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
		Balance:         numericToFloat64(transaction.Balance),
		CreatedAt:       transaction.CreatedAt.Time,
	}
}

func (tr *transactionRepository) getQuerier(ctx context.Context) sqlc.Querier {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return sqlc.New(tx)
	}
	return tr.querier
}
