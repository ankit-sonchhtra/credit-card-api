package repository

//go:generate mockgen -source=account_repository.go -destination=mocks/mock_account_repository.go -package=mocks

import (
	"context"
	"errors"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/repository/sqlc"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	logger "github.com/sirupsen/logrus"
)

type AccountRepository interface {
	Create(ctx context.Context, accountParam domain.CreateAccountParam) (domainAccount *domain.Account, err error)
	GetById(ctx context.Context, id int64) (domainAccount *domain.Account, err error)
}

type accountRepository struct {
	querier sqlc.Querier
}

func NewAccountRepository(querier sqlc.Querier) AccountRepository {
	return &accountRepository{querier: querier}
}

func (ar *accountRepository) Create(ctx context.Context, accountParam domain.CreateAccountParam) (domainAccount *domain.Account, err error) {
	account, err := ar.querier.CreateAccount(ctx, accountParam.DocumentNumber)
	if err != nil {
		logger.Error("error while create an account: ", err.Error())
		if isUniqueViolation(err) {
			return nil, domain.ErrAccountAlreadyExist
		}
		return nil, err
	}
	logger.Info("account created successfully in db.")
	return mapToDomainAccount(account), err
}

func (ar *accountRepository) GetById(ctx context.Context, id int64) (domainAccount *domain.Account, err error) {
	account, err := ar.querier.GetAccountByID(ctx, id)
	if err != nil {
		logger.Errorf("error while fetch account by id:%d, error: %s", id, err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}
		return nil, err
	}
	logger.Info("account fetched successfully from db.")
	return mapToDomainAccount(account), err
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == pgerrcode.UniqueViolation
	}
	return false
}

func mapToDomainAccount(account sqlc.Account) *domain.Account {
	return &domain.Account{
		Id:             account.AccountID,
		DocumentNumber: account.DocumentNumber,
		CreatedAt:      account.CreatedAt.Time,
	}
}
