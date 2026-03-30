package services

import (
	"context"
	"testing"
	"time"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	accountId      int64
	documentNumber string
)

type AccountServiceTestSuite struct {
	suite.Suite
	context               context.Context
	mockController        *gomock.Controller
	mockAccountRepository *mocks.MockAccountRepository
	accountService        AccountService
}

func TestAccountServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AccountServiceTestSuite))
}

func (suite *AccountServiceTestSuite) SetupTest() {
	suite.context = context.TODO()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockAccountRepository = mocks.NewMockAccountRepository(suite.mockController)
	suite.accountService = NewAccountService(suite.mockAccountRepository)
	accountId = 1
	documentNumber = "0123456789"
}

func (suite *AccountServiceTestSuite) TestCreateAccount_Success() {
	requestPayload := models.CreateAccountRequest{DocumentNumber: documentNumber}

	accountParam := domain.CreateAccountParam{
		DocumentNumber: documentNumber,
	}

	expectedResponse := &domain.Account{
		Id:             accountId,
		DocumentNumber: documentNumber,
		CreatedAt:      time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC),
	}

	suite.mockAccountRepository.EXPECT().Create(suite.context, accountParam).Return(expectedResponse, nil).Times(1)

	response, err := suite.accountService.RegisterAccount(suite.context, requestPayload)

	suite.Nil(err)
	suite.Equal(expectedResponse, response)
}

func (suite *AccountServiceTestSuite) TestCreateAccount_When_AccountRepo_Returns_ConflictError() {
	request := models.CreateAccountRequest{
		DocumentNumber: documentNumber,
	}
	accountParam := domain.CreateAccountParam{
		DocumentNumber: documentNumber,
	}
	expectedErr := domain.ErrAccountAlreadyExist

	suite.mockAccountRepository.EXPECT().Create(suite.context, accountParam).Return(nil, expectedErr).Times(1)

	response, err := suite.accountService.RegisterAccount(suite.context, request)

	suite.Nil(response)
	suite.NotNil(err)
	suite.Equal(expectedErr, err)

}

func (suite *AccountServiceTestSuite) TestGetAccount_Success() {
	expectedResponse := &domain.Account{
		Id:             accountId,
		DocumentNumber: documentNumber,
		CreatedAt:      time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC),
	}

	suite.mockAccountRepository.EXPECT().GetById(suite.context, accountId).Return(expectedResponse, nil)

	response, err := suite.accountService.GetAccount(suite.context, accountId)

	suite.Nil(err)
	suite.Equal(expectedResponse, response)
}

func (suite *AccountServiceTestSuite) TestGetAccount_When_AccountRepo_Returns_Error() {
	expectedErr := domain.ErrAccountNotFound
	suite.mockAccountRepository.EXPECT().GetById(suite.context, accountId).Return(nil, expectedErr).Times(1)

	response, err := suite.accountService.GetAccount(suite.context, accountId)

	suite.Nil(response)
	suite.NotNil(err)
	suite.Equal(expectedErr, err)
}

func (suite *AccountServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
