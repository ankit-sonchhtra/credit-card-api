package services

//go:generate mockgen -source=transaction_service.go -destination=mocks/mock_transaction_service.go -package=mocks

import (
	"context"
	"errors"
	"strings"

	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/repository"
	"github.com/credit-card-api/internal/repository/model"
	"github.com/credit-card-api/pkg/constants"
	"github.com/credit-card-api/pkg/utils"
	"github.com/google/uuid"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, request models.CreateTransactionRequest) (*models.CreateTransactionResponse, *models.CCError)
}

type transactionService struct {
	transactionsRepo repository.TransactionRepository
	accountsRepo     repository.AccountRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository,
	accountsRepo repository.AccountRepository) TransactionService {
	return &transactionService{transactionsRepo: transactionRepo, accountsRepo: accountsRepo}
}

func (as *transactionService) CreateTransaction(ctx context.Context, request models.CreateTransactionRequest) (*models.CreateTransactionResponse, *models.CCError) {
	txnValidationErr := validateTransactionRules(request)
	if txnValidationErr != nil {
		return nil, utils.NewCCBadRequestError(txnValidationErr.Error())
	}

	err := as.validateAccountExist(ctx, request)
	if err != nil {
		return nil, err
	}

	transactionDocument := buildTransactionDocument(request)
	txnErr := as.transactionsRepo.CreateTransaction(ctx, transactionDocument)
	if txnErr != nil {
		return nil, utils.NewCCInternalServerError()
	}

	return &models.CreateTransactionResponse{TransactionId: transactionDocument.Id}, nil
}

func (as *transactionService) validateAccountExist(ctx context.Context, request models.CreateTransactionRequest) *models.CCError {
	account, err := as.accountsRepo.GetAccount(ctx, request.AccountId)
	if err != nil {
		return utils.NewCCInternalServerError()
	}

	if account == nil {
		return utils.NewCCBadRequestError(constants.InvalidAccountIdErrMsg)
	}

	return nil
}

func buildTransactionDocument(request models.CreateTransactionRequest) model.TransactionDocument {
	transactionId := uuid.New().String()
	return model.TransactionDocument{
		Id:            transactionId,
		AccountId:     request.AccountId,
		OperationType: strings.ToUpper(request.OperationType),
		Amount:        request.Amount,
		CreatedAt:     currentTime().UnixMilli(),
		UpdatedAt:     currentTime().UnixMilli(),
	}
}

func validateTransactionRules(request models.CreateTransactionRequest) error {
	operationType := strings.ToLower(request.OperationType)

	switch operationType {
	case constants.OpTypeCashPurchase, constants.OpTypeInstallmentPurchase, constants.OpTypeWithdrawal:
		if request.Amount > 0 {
			return errors.New(constants.AmountMustBeNegativeErrMsg)
		}
	case constants.OpTypePayment:
		if request.Amount < 0 {
			return errors.New(constants.AmountMustBePositiveErrMsg)
		}
	default:
		return errors.New(constants.InvalidOperationTypeErrMsg)
	}
	return nil
}
