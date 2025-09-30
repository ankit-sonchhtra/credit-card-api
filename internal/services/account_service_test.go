package services

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/credit-card-service/internal/models"
	"github.com/credit-card-service/internal/repository/mocks"
	"github.com/credit-card-service/internal/repository/model"
	"github.com/credit-card-service/pkg/constants"
	"github.com/credit-card-service/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	accountId      string
	userId         string
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
	accountId = "52fdfc07-2182-454f-963f-5f0f9a621d72"
	userId = "cb9841df-c22e-4897-abfb-2411fad3e03d"
	documentNumber = "0123456789"
	currentTime = mockNow
	uuid.SetRand(rand.New(rand.NewSource(1)))
}

func (suite *AccountServiceTestSuite) TestCreateAccount_Success() {
	request := models.CreateAccountRequest{
		UserId:         userId,
		DocumentNumber: documentNumber,
	}
	expectedResponse := &models.CreateAccountResponse{
		AccountId:      accountId,
		DocumentNumber: documentNumber,
	}
	accDocument := model.AccountDocument{
		AccountId:      accountId,
		UserId:         userId,
		DocumentNumber: documentNumber,
		CreatedAt:      1759170600000,
		UpdatedAt:      1759170600000,
	}

	suite.mockAccountRepository.EXPECT().CreateAccount(suite.context, accDocument).Return(nil)

	response, err := suite.accountService.CreateAccount(suite.context, request)

	suite.Nil(err)
	suite.Equal(expectedResponse, response)

}

func (suite *AccountServiceTestSuite) TestCreateAccount_When_AccountRepo_Returns_AnError() {
	request := models.CreateAccountRequest{
		UserId:         userId,
		DocumentNumber: documentNumber,
	}
	expectedErr := utils.NewCCInternalServerError()
	accDocument := model.AccountDocument{
		AccountId:      accountId,
		UserId:         userId,
		DocumentNumber: documentNumber,
		CreatedAt:      1759170600000,
		UpdatedAt:      1759170600000,
	}

	suite.mockAccountRepository.EXPECT().CreateAccount(suite.context, accDocument).Return(errors.New("failed to create"))

	response, err := suite.accountService.CreateAccount(suite.context, request)

	suite.Nil(response)
	suite.NotNil(err)
	suite.Equal(expectedErr, err)

}

func (suite *AccountServiceTestSuite) TestGetAccount_Success() {
	expectedResponse := &models.GetAccountResponse{
		AccountId:      accountId,
		DocumentNumber: documentNumber,
		UserId:         userId,
	}
	accDocument := &model.AccountDocument{
		AccountId:      accountId,
		UserId:         userId,
		DocumentNumber: documentNumber,
		CreatedAt:      1759170600000,
		UpdatedAt:      1759170600000,
	}

	suite.mockAccountRepository.EXPECT().GetAccount(suite.context, accountId).Return(accDocument, nil)

	response, err := suite.accountService.GetAccount(suite.context, accountId)

	suite.Nil(err)
	suite.Equal(expectedResponse, response)
}

func (suite *AccountServiceTestSuite) TestGetAccount_When_AccountRepo_Returns_Error() {
	expectedErr := utils.NewCCInternalServerError()

	suite.mockAccountRepository.EXPECT().GetAccount(suite.context, accountId).Return(nil, errors.New("failed to fetch"))

	response, err := suite.accountService.GetAccount(suite.context, accountId)

	suite.Nil(response)
	suite.NotNil(err)
	suite.Equal(expectedErr, err)
}

func (suite *AccountServiceTestSuite) TestGetAccount_When_Account_IsNotExist() {
	expectedErr := &models.CCError{
		ErrorCode:      constants.AccountNotPresentErrCode,
		ErrorMessage:   constants.AccountNotPresentErrMsg,
		AdditionalData: models.AdditionalData{StatusCode: 404},
	}

	suite.mockAccountRepository.EXPECT().GetAccount(suite.context, accountId).Return(nil, nil)

	response, err := suite.accountService.GetAccount(suite.context, accountId)

	suite.Nil(response)
	suite.NotNil(err)
	suite.Equal(expectedErr, err)
}

func (suite *AccountServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func mockNow() time.Time {
	IST, _ := time.LoadLocation(constants.AsiaKolkataTimeZone)
	return time.Date(2025, 9, 30, 0, 0, 0, 0, IST)
}
