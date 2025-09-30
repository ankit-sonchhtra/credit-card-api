package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	accountId      string
	userId         string
	documentNumber string
)

type AccountControllerTestSuite struct {
	suite.Suite
	context            *gin.Context
	goContext          context.Context
	recorder           *httptest.ResponseRecorder
	mockController     *gomock.Controller
	mockAccountService *mocks.MockAccountService
	controller         AccountController
}

func TestAccountControllerTestSuite(t *testing.T) {
	suite.Run(t, new(AccountControllerTestSuite))
}

func (suite *AccountControllerTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.goContext = context.TODO()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockAccountService = mocks.NewMockAccountService(suite.mockController)
	suite.controller = NewAccountController(suite.mockAccountService)
	accountId = "92d68c0e-dafe-406a-a0f2-8faae2020947"
	userId = "cb9841df-c22e-4897-abfb-2411fad3e03d"
	documentNumber = "0123456789"
}

func (suite *AccountControllerTestSuite) TestCreateAccount_Success() {
	payload := models.CreateAccountRequest{
		UserId:         userId,
		DocumentNumber: documentNumber,
	}
	expectedResponse := &models.CreateAccountResponse{
		AccountId:      accountId,
		DocumentNumber: documentNumber,
	}
	expectedResponseBody := `{"accountId":"92d68c0e-dafe-406a-a0f2-8faae2020947","documentNumber":"0123456789"}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockAccountService.EXPECT().CreateAccount(suite.context, payload).Return(expectedResponse, nil)

	suite.controller.CreateAccount(suite.context)

	suite.Equal(http.StatusCreated, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_DocumentNumber_IsMissing() {
	payload := models.CreateAccountRequest{
		UserId: userId,
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req
	suite.controller.CreateAccount(suite.context)
	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"invalid request body","additionalData":{"statusCode":400}}`
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_UserId_IsMissing() {
	payload := models.CreateAccountRequest{
		DocumentNumber: documentNumber,
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req
	suite.controller.CreateAccount(suite.context)
	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"invalid request body","additionalData":{"statusCode":400}}`
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_AccountService_Returns_Error() {
	payload := models.CreateAccountRequest{
		UserId:         userId,
		DocumentNumber: documentNumber,
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req

	expectedErr := &models.CCError{
		ErrorCode:      "ERR_CC_INTERNAL_SERVER_ERROR",
		ErrorMessage:   "internal server error",
		AdditionalData: models.AdditionalData{StatusCode: 500},
	}

	suite.mockAccountService.EXPECT().CreateAccount(suite.context, payload).Return(nil, expectedErr)

	suite.controller.CreateAccount(suite.context)

	expectedResponseBody := `{"errorCode":"ERR_CC_INTERNAL_SERVER_ERROR","errorMessage":"internal server error","additionalData":{"statusCode":500}}`

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestGetAccount_Success() {
	expectedResponse := &models.GetAccountResponse{
		AccountId:      accountId,
		DocumentNumber: documentNumber,
		UserId:         userId,
	}
	expectedResponseBody, _ := json.Marshal(expectedResponse)

	req := httptest.NewRequest(http.MethodGet, "/accounts/92d68c0e-dafe-406a-a0f2-8faae2020947", nil)
	suite.context.Request = req
	suite.context.Params = gin.Params{gin.Param{
		Key:   "accountId",
		Value: accountId,
	}}
	suite.mockAccountService.EXPECT().GetAccount(suite.context, accountId).Return(expectedResponse, nil)

	suite.controller.GetAccount(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(string(expectedResponseBody), suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestGetAccount_When_AccountId_IsMissing() {
	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"accountId is missing in path params","additionalData":{"statusCode":400}}`
	req := httptest.NewRequest(http.MethodGet, "/accounts/92d68c0e-dafe-406a-a0f2-8faae2020947", nil)
	suite.context.Request = req

	suite.controller.GetAccount(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestGetAccount_When_AccountService_Returns_Error() {
	expectedErr := &models.CCError{
		ErrorCode:      "ERR_CC_INTERNAL_SERVER_ERROR",
		ErrorMessage:   "internal server error",
		AdditionalData: models.AdditionalData{StatusCode: 500},
	}
	expectedResponseBody, _ := json.Marshal(expectedErr)

	req := httptest.NewRequest(http.MethodGet, "/accounts/92d68c0e-dafe-406a-a0f2-8faae2020947", nil)
	suite.context.Request = req
	suite.context.Params = gin.Params{gin.Param{
		Key:   "accountId",
		Value: accountId,
	}}
	suite.mockAccountService.EXPECT().GetAccount(suite.context, accountId).Return(nil, expectedErr)

	suite.controller.GetAccount(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(expectedResponseBody), suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
