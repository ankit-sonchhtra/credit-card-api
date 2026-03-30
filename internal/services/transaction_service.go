package services

//go:generate mockgen -source=transaction_service.go -destination=mocks/mock_transaction_service.go -package=mocks

import (
	"context"
	"errors"
	"math"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, request models.TransactionRequest) (*domain.Transaction, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	accountRepo     repository.AccountRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository, accountRepo repository.AccountRepository) TransactionService {
	return &transactionService{transactionRepo: transactionRepo, accountRepo: accountRepo}
}

func (ts *transactionService) CreateTransaction(ctx context.Context, request models.TransactionRequest) (*domain.Transaction, error) {
	logger.Infof("Started to create transaction with accountId: %d and operationTypeId :%d", request.AccountId, request.OperationTypeId)
	operationType, exists := domain.ValidOperations[request.OperationTypeId]
	if !exists {
		logger.Errorf("error: operation type id %d is not supported by the system.", request.OperationTypeId)
		return nil, domain.ErrInvalidOperationType
	}

	_, err := ts.accountRepo.GetById(ctx, request.AccountId)
	if err != nil {
		if isAccountNotFoundError(err) {
			logger.Errorf("error: account is not exist with provided id: %d", request.AccountId)
			return nil, domain.ErrTransactionAccountNotFound
		}
		return nil, err
	}

	transaction := domain.CreateTransactionParam{
		AccountId:       request.AccountId,
		OperationTypeId: request.OperationTypeId,
		Amount:          normalizeAmountByOperation(request.Amount, operationType),
	}
	return ts.transactionRepo.Create(ctx, transaction)
}

func normalizeAmountByOperation(amount float64, opType domain.TransactionType) float64 {
	abs := math.Abs(amount)
	if opType.IsNegative {
		return -abs
	}
	return abs
}

func isAccountNotFoundError(err error) bool {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		return appErr.Code == domain.ErrAccountNotFound.Code
	}
	return false
}
