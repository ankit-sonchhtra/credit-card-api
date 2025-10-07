package services

//go:generate mockgen -source=account_service.go -destination=mocks/mock_account_service.go -package=mocks

import (
	"context"
	"time"

	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/repository"
	"github.com/credit-card-api/internal/repository/model"
	"github.com/credit-card-api/pkg/constants"
	"github.com/credit-card-api/pkg/utils"
	"github.com/google/uuid"
)

var (
	currentTime = time.Now
)

type AccountService interface {
	CreateAccount(ctx context.Context, request models.CreateAccountRequest) (*models.CreateAccountResponse, *models.CCError)
	GetAccount(ctx context.Context, id string) (*models.GetAccountResponse, *models.CCError)
}

type accountService struct {
	accountsRepo repository.AccountRepository
	userRepo     repository.UserRepository
}

func NewAccountService(accountRepo repository.AccountRepository, userRepo repository.UserRepository) AccountService {
	return &accountService{accountsRepo: accountRepo, userRepo: userRepo}
}

func (as *accountService) CreateAccount(ctx context.Context, request models.CreateAccountRequest) (*models.CreateAccountResponse, *models.CCError) {
	err := as.validateUserExist(ctx, request)
	if err != nil {
		return nil, err
	}

	accountDocument := buildAccountDocument(request)
	createErr := as.accountsRepo.CreateAccount(ctx, accountDocument)
	if createErr != nil {
		return nil, utils.NewCCInternalServerError()
	}

	return &models.CreateAccountResponse{
		AccountId:      accountDocument.AccountId,
		DocumentNumber: accountDocument.DocumentNumber,
	}, nil
}

func (as *accountService) GetAccount(ctx context.Context, id string) (*models.GetAccountResponse, *models.CCError) {
	accountDocument, err := as.accountsRepo.GetAccount(ctx, id)
	if err != nil {
		return nil, utils.NewCCInternalServerError()
	}

	if accountDocument == nil {
		return nil, &models.CCError{
			ErrorCode:      constants.AccountNotPresentErrCode,
			ErrorMessage:   constants.AccountNotPresentErrMsg,
			AdditionalData: models.AdditionalData{StatusCode: 404},
		}
	}

	return &models.GetAccountResponse{
		AccountId:      accountDocument.AccountId,
		DocumentNumber: accountDocument.DocumentNumber,
		UserId:         accountDocument.UserId,
	}, nil
}

func (as *accountService) validateUserExist(ctx context.Context, request models.CreateAccountRequest) *models.CCError {
	filters := make(map[string]interface{})
	filters[constants.UserIdFilter] = request.UserId

	user, err := as.userRepo.GetUserByFilters(ctx, filters)
	if err != nil {
		return utils.NewCCInternalServerError()
	}

	if user == nil {
		return utils.NewCCBadRequestError(constants.InvalidUserIdErrMsg)
	}

	return nil
}

func buildAccountDocument(request models.CreateAccountRequest) model.AccountDocument {
	accountId := uuid.New().String()
	documentNumber := request.DocumentNumber
	return model.AccountDocument{
		AccountId:      accountId,
		UserId:         request.UserId,
		DocumentNumber: documentNumber,
		CreatedAt:      currentTime().UnixMilli(),
		UpdatedAt:      currentTime().UnixMilli(),
	}
}
