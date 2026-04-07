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

	if request.OperationTypeId == 4 {
		//Need to fetch transactions for same accountID based on eventDate in ascending order.
		transactions, fetchTransactionErr := ts.transactionRepo.GetAllTransactions(ctx, request.AccountId)

		if fetchTransactionErr != nil {
			if !isAccountNotFoundError(err) {
				logger.Errorf("error: account is not exist with provided id: %d", request.AccountId)
				return nil, domain.ErrTransactionAccountNotFound
			}
		}
		// 60,

		// -50. 23.5, -18.7
		remainingBalance := request.Amount
		var updatedBalance float64

		for _, tx := range transactions {

			if remainingBalance > tx.Balance {
				remainingBalance = remainingBalance + tx.Balance
				if remainingBalance > 0 {
					transactionUpdateErr := ts.transactionRepo.UpdateTransactionById(ctx, tx.Id, 0.0)
					if transactionUpdateErr != nil {
						return nil, domain.ErrInternal
					}
				} else if remainingBalance != 0 {
					transactionUpdateErr := ts.transactionRepo.UpdateTransactionById(ctx, tx.Id, remainingBalance)
					if transactionUpdateErr != nil {
						return nil, domain.ErrInternal
					}
				}
				continue
			}
			finalAmount := normalizeAmountByOperation(request.Amount, operationType)
			if remainingBalance > 0 {
				updatedBalance = request.Amount + remainingBalance
			} else {
				updatedBalance = 0.0
			}

			transaction := domain.CreateTransactionParam{
				AccountId:       request.AccountId,
				OperationTypeId: request.OperationTypeId,
				Amount:          finalAmount,
				Balance:         updatedBalance,
			}
			return ts.transactionRepo.Create(ctx, transaction)
			//else if tx.Balance < 0 { // 60 > -50 && -50 < 0 ,  // 10 > -23.5 &&
			//	remainingBalance = remainingBalance + tx.Balance // 60 - 50 = 10, // 10 - 23.5 == 13.5
			//	// Update query to update transaction balance with the correct value
			//	if remainingBalance > 0 {
			//		updatedBalance = 0
			//	} else {
			//		updatedBalance = remainingBalance
			//	}
			//
			//	transactionUpdateErr := ts.transactionRepo.UpdateTransactionById(ctx, tx.Id, updatedBalance)
			//
			//	if transactionUpdateErr != nil {
			//		return nil, domain.ErrInternal
			//	}
			//}
		}
	}

	finalAmount := normalizeAmountByOperation(request.Amount, operationType)
	transaction := domain.CreateTransactionParam{
		AccountId:       request.AccountId,
		OperationTypeId: request.OperationTypeId,
		Amount:          finalAmount,
		Balance:         finalAmount,
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
