package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/services/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	accountId      int64
	documentNumber string
)

type AccountControllerTestSuite struct {
	suite.Suite
	context            *gin.Context
	goContext          context.Context
	recorder           *httptest.ResponseRecorder
	mockController     *gomock.Controller
	mockAccountService *mocks.MockAccountService
	controller         *AccountController
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
	accountId = 1
	documentNumber = "0123456789"
}

func (suite *AccountControllerTestSuite) TestCreateAccount_Success() {
	payload := models.CreateAccountRequest{
		DocumentNumber: documentNumber,
	}
	accountResponse := &domain.Account{
		Id:             accountId,
		DocumentNumber: documentNumber,
	}

	expectedResponseBody := `{"account_id":1,"document_number":"0123456789"}`
	bodyBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req
	suite.mockAccountService.EXPECT().RegisterAccount(suite.context, payload).Return(accountResponse, nil)

	suite.controller.CreateAccount(suite.context)

	suite.Equal(http.StatusCreated, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_BindingFails_Returns_Error() {
	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"invalid request body","status_code":400}`
	bodyBytes, _ := json.Marshal("invalid payload")

	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.controller.CreateAccount(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_DocumentNumber_IsInvalid() {
	payload := models.CreateAccountRequest{
		DocumentNumber: "asdfasdfasd",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"The 'DocumentNumber' field will only accept numeric value.","status_code":400}`

	suite.context.Request = req
	suite.controller.CreateAccount(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_DocumentNumber_Length_IsUnexpected() {
	payload := models.CreateAccountRequest{
		DocumentNumber: "0998877665544",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"The 'DocumentNumber' field cannot exceed 12 digits.","status_code":400}`

	suite.context.Request = req
	suite.controller.CreateAccount(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_DocumentNumber_IsMissing() {
	payload := models.CreateAccountRequest{
		DocumentNumber: "",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"The 'DocumentNumber' field is mandatory.","status_code":400}`

	suite.context.Request = req
	suite.controller.CreateAccount(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_AccountService_Returns_Error() {
	payload := models.CreateAccountRequest{
		DocumentNumber: documentNumber,
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req

	errResponse := &domain.AppError{
		Code:    "ERR_CC_INTERNAL_SERVER_ERROR",
		Message: "internal server error",
	}

	suite.mockAccountService.EXPECT().RegisterAccount(suite.context, payload).Return(nil, errResponse)

	suite.controller.CreateAccount(suite.context)

	expectedResponseBody := `{"error_code":"ERR_CC_INTERNAL_SERVER_ERROR","error_message":"internal server error","status_code":500}`

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_AccountService_Returns_ConflictError() {
	payload := models.CreateAccountRequest{
		DocumentNumber: documentNumber,
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req

	errResponse := &domain.AppError{
		Code:    "ERR_CC_ACCOUNT_ALREADY_EXIST",
		Message: "account already exists.",
	}

	suite.mockAccountService.EXPECT().RegisterAccount(suite.context, payload).Return(nil, errResponse)

	suite.controller.CreateAccount(suite.context)

	expectedResponseBody := `{"error_code":"ERR_CC_ACCOUNT_ALREADY_EXIST","error_message":"account already exists.","status_code":409}`

	suite.Equal(http.StatusConflict, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestCreateAccount_When_AccountService_Returns_UnknownError() {
	payload := models.CreateAccountRequest{
		DocumentNumber: documentNumber,
	}

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/accounts", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req

	suite.mockAccountService.EXPECT().RegisterAccount(suite.context, payload).Return(nil, errors.New("failed to fetch"))

	suite.controller.CreateAccount(suite.context)

	expectedResponseBody := `{"error_code":"ERR_CC_INTERNAL_SERVER_ERROR","error_message":"an unexpected error occurred.","status_code":500}`

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestGetAccount_Success() {
	accountResponse := &domain.Account{
		Id:             accountId,
		DocumentNumber: documentNumber,
	}
	expectedResponseBody := `{"account_id":1,"document_number":"0123456789"}`

	req := httptest.NewRequest(http.MethodGet, "/api/credit-card-api/v1/accounts/1", nil)
	suite.context.Request = req
	suite.context.Params = gin.Params{gin.Param{
		Key:   "accountId",
		Value: "1",
	}}
	suite.mockAccountService.EXPECT().GetAccount(suite.context, accountId).Return(accountResponse, nil)

	suite.controller.GetAccount(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(string(expectedResponseBody), suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestGetAccount_When_AccountId_IsMissing() {
	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"accountId is missing in path params","status_code":400}`
	req := httptest.NewRequest(http.MethodGet, "/api/credit-card-api/v1/accounts/", nil)
	suite.context.Request = req

	suite.controller.GetAccount(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestGetAccount_When_AccountService_Returns_Error() {
	expectedErr := &domain.AppError{
		Code:    "ERR_CC_INTERNAL_SERVER_ERROR",
		Message: "internal server error",
	}
	expectedResponseBody := `{"error_code":"ERR_CC_INTERNAL_SERVER_ERROR","error_message":"internal server error","status_code":500}`

	req := httptest.NewRequest(http.MethodGet, "/api/credit-card-api/v1/accounts/1", nil)
	suite.context.Request = req
	suite.context.Params = gin.Params{gin.Param{
		Key:   "accountId",
		Value: "1",
	}}
	suite.mockAccountService.EXPECT().GetAccount(suite.context, accountId).Return(nil, expectedErr)

	suite.controller.GetAccount(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(expectedResponseBody), suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TestGetAccount_When_AccountService_Returns_AccountNotFoundError() {
	expectedErr := &domain.AppError{
		Code:    "ERR_CC_ACCOUNT_NOT_FOUND",
		Message: "account does not exists with provided id.",
	}
	expectedResponseBody := `{"error_code":"ERR_CC_ACCOUNT_NOT_FOUND","error_message":"account does not exists with provided id.","status_code":404}`

	req := httptest.NewRequest(http.MethodGet, "/api/credit-card-api/v1/accounts/1", nil)
	suite.context.Request = req
	suite.context.Params = gin.Params{gin.Param{
		Key:   "accountId",
		Value: "1",
	}}
	suite.mockAccountService.EXPECT().GetAccount(suite.context, accountId).Return(nil, expectedErr)

	suite.controller.GetAccount(suite.context)

	suite.Equal(http.StatusNotFound, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *AccountControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
