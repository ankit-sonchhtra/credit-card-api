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
	testAccountId int64
	testTxnId     int64
)

type TransactionControllerTestSuite struct {
	suite.Suite
	context                *gin.Context
	goContext              context.Context
	recorder               *httptest.ResponseRecorder
	mockController         *gomock.Controller
	mockTransactionService *mocks.MockTransactionService
	transactionController  *TransactionController
}

func TestTransactionControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}

func (suite *TransactionControllerTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.goContext = context.TODO()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockTransactionService = mocks.NewMockTransactionService(suite.mockController)
	suite.transactionController = NewTransactionController(suite.mockTransactionService)
	testAccountId = 1
	testTxnId = 1
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_Success() {
	payload := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 1,
		Amount:          12345.67,
	}
	transaction := &domain.Transaction{
		Id:              testTxnId,
		AccountId:       testAccountId,
		OperationTypeId: 1,
		Amount:          12345.67,
	}
	expectedResponseBody := `{"transaction_id":1}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockTransactionService.EXPECT().CreateTransaction(suite.context, payload).Return(transaction, nil)

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusCreated, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_When_Payload_Binding_Fails() {
	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"invalid request body","status_code":400}`

	bodyBytes, _ := json.Marshal("invalid payload")
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_When_AccountID_IsMissing() {
	payload := models.TransactionRequest{
		OperationTypeId: 1,
		Amount:          12345.67,
	}

	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"The 'AccountId' field is mandatory.","status_code":400}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_When_OperationTypeId_IsMissing() {
	payload := models.TransactionRequest{
		AccountId: testAccountId,
		Amount:    12345.67,
	}

	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"The 'OperationTypeId' field is mandatory.","status_code":400}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_When_Amount_IsMissing() {
	payload := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 2,
	}

	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"The 'Amount' field is mandatory.","status_code":400}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_When_Amount_Is_Invalid() {
	payload := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 2,
		Amount:          -123.05,
	}

	expectedResponseBody := `{"error_code":"ERR_CC_BAD_REQUEST","error_message":"The 'Amount' field value must be greater than 0.","status_code":400}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_When_Service_Return_AccountNotFoundError() {
	payload := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 2,
		Amount:          123.05,
	}

	appError := &domain.AppError{
		Code:    "ERR_CC_TRANSACTION_ACCOUNT_NOT_FOUND",
		Message: "account does not exist with provided id.",
	}

	expectedResponseBody := `{"error_code":"ERR_CC_TRANSACTION_ACCOUNT_NOT_FOUND","error_message":"account does not exist with provided id.","status_code":422}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockTransactionService.EXPECT().CreateTransaction(suite.context, payload).Return(nil, appError).Times(1)
	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusUnprocessableEntity, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_When_Service_Return_UnknownError() {
	payload := models.TransactionRequest{
		AccountId:       testAccountId,
		OperationTypeId: 2,
		Amount:          123.05,
	}

	err := errors.New("unknown error")
	expectedResponseBody := `{"error_code":"ERR_CC_INTERNAL_SERVER_ERROR","error_message":"an unexpected error occurred.","status_code":500}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/credit-card-api/v1/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockTransactionService.EXPECT().CreateTransaction(suite.context, payload).Return(nil, err).Times(1)
	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
