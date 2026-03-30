package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/repository/mocks"
	"github.com/credit-card-api/pkg/constants"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testAccountId     int64
	testTransactionId int64
)

type TransactionServiceTestSuite struct {
	suite.Suite
	context                   context.Context
	mockController            *gomock.Controller
	mockTransactionRepository *mocks.MockTransactionRepository
	mockAccountRepository     *mocks.MockAccountRepository
	transactionService        TransactionService
}

func TestTransactionServiceTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.context = context.TODO()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockTransactionRepository = mocks.NewMockTransactionRepository(suite.mockController)
	suite.mockAccountRepository = mocks.NewMockAccountRepository(suite.mockController)
	suite.transactionService = NewTransactionService(suite.mockTransactionRepository, suite.mockAccountRepository)
	testAccountId = 1
	testTransactionId = 1
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_Success() {
	request := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 2,
		Amount:          2345.67,
	}

	account := &domain.Account{
		Id:             accountId,
		DocumentNumber: documentNumber,
		CreatedAt:      time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC),
	}
	transactionParam := domain.CreateTransactionParam{
		AccountId:       testAccountId,
		OperationTypeId: 2,
		Amount:          -2345.67,
	}
	expectedTransaction := &domain.Transaction{
		Id:              testTransactionId,
		AccountId:       accountId,
		OperationTypeId: 2,
		Amount:          -2345.67,
		CreatedAt:       time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC),
	}

	suite.mockAccountRepository.EXPECT().GetById(suite.context, testAccountId).Return(account, nil)
	suite.mockTransactionRepository.EXPECT().Create(suite.context, transactionParam).Return(expectedTransaction, nil).Times(1)

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.Nil(err)
	suite.Equal(expectedTransaction, response)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_When_OperationType_Is_Payment() {
	request := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 4,
		Amount:          4567.67,
	}

	account := &domain.Account{
		Id:             accountId,
		DocumentNumber: documentNumber,
		CreatedAt:      time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC),
	}
	transactionParam := domain.CreateTransactionParam{
		AccountId:       testAccountId,
		OperationTypeId: 4,
		Amount:          4567.67,
	}
	expectedTransaction := &domain.Transaction{
		Id:              testTransactionId,
		AccountId:       accountId,
		OperationTypeId: 4,
		Amount:          4567.67,
		CreatedAt:       time.Date(2026, time.January, 1, 10, 0, 0, 0, time.UTC),
	}

	suite.mockAccountRepository.EXPECT().GetById(suite.context, testAccountId).Return(account, nil)
	suite.mockTransactionRepository.EXPECT().Create(suite.context, transactionParam).Return(expectedTransaction, nil).Times(1)

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.Nil(err)
	suite.Equal(expectedTransaction, response)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_Return_Error_When_OperationTypeId_IsNotSupported() {
	request := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 5,
		Amount:          2345.67,
	}
	expectedErr := &domain.AppError{
		Code:    constants.InvalidOperationTypeErrCode,
		Message: "operation type ID provided is not supported by the system.",
	}

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.Nil(response)
	suite.Equal(expectedErr, err)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_ReturnErr_When_AccountRepo_Return_NotAccountFoundErr() {
	request := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 2,
		Amount:          2345.67,
	}

	accountErr := &domain.AppError{
		Code:    constants.AccountNotFoundErrCode,
		Message: "account does not exists with provided id.",
	}

	expectedErr := &domain.AppError{
		Code:    constants.TransactionAccountNotFoundErrCode,
		Message: "account does not exist with provided id.",
	}

	suite.mockAccountRepository.EXPECT().GetById(suite.context, testAccountId).Return(nil, accountErr)

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.Equal(expectedErr, err)
	suite.Nil(response)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_ReturnErr_When_AccountRepo_Returns_AnError() {
	request := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 2,
		Amount:          2345.67,
	}

	expectedErr := errors.New("failed to fetch an account")

	suite.mockAccountRepository.EXPECT().GetById(suite.context, testAccountId).Return(nil, expectedErr)

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.Equal(expectedErr, err)
	suite.Nil(response)
}

func (suite *TransactionServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
