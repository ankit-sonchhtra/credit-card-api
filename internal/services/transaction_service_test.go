package services

import (
	"context"
	"errors"

	"math/rand"
	"testing"

	"github.com/credit-card-service/internal/models"
	"github.com/credit-card-service/internal/repository/mocks"
	"github.com/credit-card-service/internal/repository/model"
	"github.com/credit-card-service/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testAccountId     string
	testTransactionId string
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
	testAccountId = "cb9841df-c22e-4897-abfb-2411fad3e03d"
	testTransactionId = "52fdfc07-2182-454f-963f-5f0f9a621d72"
	currentTime = mockNow
	uuid.SetRand(rand.New(rand.NewSource(1)))
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_Success() {
	request := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "cash purchase",
		Amount:        -2345.67,
	}
	expectedResponse := &models.CreateTransactionResponse{
		TransactionId: testTransactionId,
	}
	accDocument := &model.AccountDocument{
		AccountId:      accountId,
		UserId:         userId,
		DocumentNumber: documentNumber,
	}

	transactionDoc := model.TransactionDocument{
		Id:            testTransactionId,
		AccountId:     testAccountId,
		OperationType: "CASH PURCHASE",
		Amount:        -2345.67,
		CreatedAt:     1759170600000,
		UpdatedAt:     1759170600000,
	}

	suite.mockAccountRepository.EXPECT().GetAccount(suite.context, testAccountId).Return(accDocument, nil)
	suite.mockTransactionRepository.EXPECT().CreateTransaction(suite.context, transactionDoc).Return(nil)

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.Nil(err)
	suite.Equal(expectedResponse, response)

}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_When_AccountRepo_Returns_AnError() {
	request := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "cash purchase",
		Amount:        -2345.67,
	}

	suite.mockAccountRepository.EXPECT().GetAccount(suite.context, testAccountId).Return(nil, errors.New("failed to fetch account"))

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.NotNil(err)
	suite.Equal(utils.NewCCInternalServerError(), err)
	suite.Nil(response)
}

func (suite *TransactionServiceTestSuite) TestCreateTransaction_When_Account_IsNotExist() {
	request := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "cash purchase",
		Amount:        -2345.67,
	}
	expectedErr := &models.CCError{
		ErrorCode:      "ERR_CC_BAD_REQUEST",
		ErrorMessage:   "account does not exists with requested accountId",
		AdditionalData: models.AdditionalData{StatusCode: 400},
	}

	suite.mockAccountRepository.EXPECT().GetAccount(suite.context, testAccountId).Return(nil, nil)

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.NotNil(err)
	suite.Equal(expectedErr, err)
	suite.Nil(response)
}

func (suite *TransactionServiceTestSuite) TestCreateAccount_When_TransactionRepo_Returns_AnError() {
	request := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "cash purchase",
		Amount:        -2345.67,
	}
	accDocument := &model.AccountDocument{
		AccountId:      accountId,
		UserId:         userId,
		DocumentNumber: documentNumber,
	}

	transactionDoc := model.TransactionDocument{
		Id:            testTransactionId,
		AccountId:     testAccountId,
		OperationType: "CASH PURCHASE",
		Amount:        -2345.67,
		CreatedAt:     1759170600000,
		UpdatedAt:     1759170600000,
	}

	suite.mockAccountRepository.EXPECT().GetAccount(suite.context, testAccountId).Return(accDocument, nil)
	suite.mockTransactionRepository.EXPECT().CreateTransaction(suite.context, transactionDoc).Return(errors.New("failed to insert"))

	response, err := suite.transactionService.CreateTransaction(suite.context, request)

	suite.Equal(utils.NewCCInternalServerError(), err)
	suite.Nil(response)

}

func (suite *TransactionServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
