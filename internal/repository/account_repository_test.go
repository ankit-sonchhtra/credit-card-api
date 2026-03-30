package repository

import (
	"context"
	"testing"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/repository/mocks"
	"github.com/credit-card-api/internal/repository/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	accountId      int64
	documentNumber string
)

type AccountRepositoryTestSuite struct {
	suite.Suite
	context           context.Context
	mockController    *gomock.Controller
	mockQuerier       *mocks.MockQuerier
	accountRepository AccountRepository
}

func TestAccountRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(AccountRepositoryTestSuite))
}

func (suite *AccountRepositoryTestSuite) SetupTest() {
	suite.context = context.TODO()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockQuerier = mocks.NewMockQuerier(suite.mockController)
	suite.accountRepository = NewAccountRepository(suite.mockQuerier)
	accountId = 1
	documentNumber = "0123456789"
}

func (suite *AccountRepositoryTestSuite) TestAccountRepository_Create_Success() {
	expectedRow := sqlc.Account{
		AccountID:      1,
		DocumentNumber: documentNumber,
	}

	suite.mockQuerier.EXPECT().CreateAccount(suite.context, documentNumber).Return(expectedRow, nil)

	response, err := suite.accountRepository.Create(suite.context, domain.CreateAccountParam{DocumentNumber: documentNumber})

	suite.NoError(err)
	suite.Equal(int64(1), response.Id)
}

func (suite *AccountRepositoryTestSuite) TestAccountRepository_Create_Returns_ConflictError() {
	pgErr := &pgconn.PgError{Code: "23505"}
	suite.mockQuerier.EXPECT().CreateAccount(suite.context, documentNumber).Return(sqlc.Account{}, pgErr)

	res, err := suite.accountRepository.Create(suite.context, domain.CreateAccountParam{DocumentNumber: documentNumber})

	suite.Nil(res)
	suite.ErrorIs(err, domain.ErrAccountAlreadyExist)
}

func (suite *AccountRepositoryTestSuite) TestAccountRepository_GetById_Success() {
	suite.mockQuerier.EXPECT().GetAccountByID(suite.context, int64(1)).Return(sqlc.Account{AccountID: 1}, nil)

	res, err := suite.accountRepository.GetById(suite.context, accountId)

	suite.NoError(err)
	suite.Equal(int64(1), res.Id)
}

func (suite *AccountRepositoryTestSuite) TestAccountRepository_GetById_Account_Not_Found() {
	suite.mockQuerier.EXPECT().GetAccountByID(suite.context, int64(404)).Return(sqlc.Account{}, pgx.ErrNoRows)

	res, err := suite.accountRepository.GetById(suite.context, 404)

	suite.Nil(res)
	suite.ErrorIs(err, domain.ErrAccountNotFound)
}
