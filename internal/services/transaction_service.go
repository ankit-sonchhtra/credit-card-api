package services

//go:generate mockgen -source=transaction_service.go -destination=mocks/mock_transaction_service.go -package=mocks

import (
	"context"
	"net/http"
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
		return &models.CCError{
			ErrorCode:      constants.BadRequestErrCode,
			ErrorMessage:   constants.InvalidAccountIdErrMsg,
			AdditionalData: models.AdditionalData{StatusCode: http.StatusBadRequest},
		}
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
