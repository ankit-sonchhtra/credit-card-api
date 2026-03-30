package services

//go:generate mockgen -source=account_service.go -destination=mocks/mock_account_service.go -package=mocks

import (
	"context"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/repository"
	logger "github.com/sirupsen/logrus"
)

type AccountService interface {
	RegisterAccount(ctx context.Context, request models.CreateAccountRequest) (*domain.Account, error)
	GetAccount(ctx context.Context, id int64) (*domain.Account, error)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return &accountService{accountRepository: accountRepository}
}

func (as *accountService) RegisterAccount(ctx context.Context, request models.CreateAccountRequest) (*domain.Account, error) {
	logger.Infof("Started to create account with documentNumber: %s", request.DocumentNumber)
	accountParam := domain.CreateAccountParam{DocumentNumber: request.DocumentNumber}
	return as.accountRepository.Create(ctx, accountParam)

}

func (as *accountService) GetAccount(ctx context.Context, id int64) (*domain.Account, error) {
	logger.Infof("Started to get account by id: %d", id)
	return as.accountRepository.GetById(ctx, id)
}
