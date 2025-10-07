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
	"github.com/credit-card-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

var (
	testAccountId string
	testTxnId     string
)

type TransactionControllerTestSuite struct {
	suite.Suite
	context                *gin.Context
	goContext              context.Context
	recorder               *httptest.ResponseRecorder
	mockController         *gomock.Controller
	mockTransactionService *mocks.MockTransactionService
	transactionController  TransactionController
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
	testAccountId = "92d68c0e-dafe-406a-a0f2-8faae2020947"
	testTxnId = "cb9841df-c22e-4897-abfb-2411fad3e03d"
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_Success() {
	payload := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "payment",
		Amount:        12345.67,
	}
	expectedResponse := &models.CreateTransactionResponse{
		TransactionId: testTxnId,
	}
	expectedResponseBody := `{"transactionId":"cb9841df-c22e-4897-abfb-2411fad3e03d"}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockTransactionService.EXPECT().CreateTransaction(suite.context, payload).Return(expectedResponse, nil)

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusCreated, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_Invalid_RequestPayload() {
	payload := models.CreateTransactionRequest{
		OperationType: "payment",
		Amount:        12345.67,
	}

	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"invalid request body","additionalData":{"statusCode":400}}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_Invalid_OperationType() {
	payload := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "",
		Amount:        12345.67,
	}

	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"invalid request body","additionalData":{"statusCode":400}}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_Invalid_Amount() {
	payload := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "cash purchase",
		Amount:        0,
	}

	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"amount value can not be zero","additionalData":{"statusCode":400}}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TestCreateTransaction_Failed_When_TransactionService_Returns_Error() {
	payload := models.CreateTransactionRequest{
		AccountId:     testAccountId,
		OperationType: "payment",
		Amount:        12345.67,
	}

	expectedResponseBody := `{"errorCode":"ERR_CC_INTERNAL_SERVER_ERROR","errorMessage":"internal server error","additionalData":{"statusCode":500}}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockTransactionService.EXPECT().CreateTransaction(suite.context, payload).Return(nil, utils.NewCCInternalServerError())

	suite.transactionController.CreateTransaction(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *TransactionControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
