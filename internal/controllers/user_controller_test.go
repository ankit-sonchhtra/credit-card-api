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
	testUserId       string
	testEmail        string
	testMobileNumber string
)

type UserControllerTestSuite struct {
	suite.Suite
	context         *gin.Context
	goContext       context.Context
	recorder        *httptest.ResponseRecorder
	mockController  *gomock.Controller
	mockUserService *mocks.MockUserService
	controller      UserController
}

func TestUserControllerTestSuite(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.recorder = httptest.NewRecorder()
	suite.goContext = context.TODO()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockController = gomock.NewController(suite.T())
	suite.mockUserService = mocks.NewMockUserService(suite.mockController)
	suite.controller = NewUserController(suite.mockUserService)

	testUserId = "52fdfc07-2182-454f-963f-5f0f9a621d72"
	testMobileNumber = "+918908011223"
	testEmail = "abc@xyz.com"
}

func (suite *UserControllerTestSuite) TestCreateUser_Success() {
	payload := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
	}
	expectedResponse := &models.CreateUserResponse{
		UserId: testUserId,
	}
	expectedResponseBody := `{"userId":"52fdfc07-2182-454f-963f-5f0f9a621d72"}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockUserService.EXPECT().CreateUser(suite.context, payload).Return(expectedResponse, nil)

	suite.controller.CreateUser(suite.context)

	suite.Equal(http.StatusCreated, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *UserControllerTestSuite) TestCreateAccount_Failed_When_MobileNumber_IsMissing() {
	payload := models.CreateUserRequest{
		Name:  "John Doe",
		Email: testEmail,
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req
	suite.controller.CreateUser(suite.context)
	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"invalid request body","additionalData":{"statusCode":400}}`
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *UserControllerTestSuite) TestCreateAccount_Failed_When_Email_Is_Invalid() {
	payload := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        "###$$$$@@asdfasdf",
		MobileNumber: testMobileNumber,
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req
	suite.controller.CreateUser(suite.context)
	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"invalid request body","additionalData":{"statusCode":400}}`
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *UserControllerTestSuite) TestCreateAccount_When_MobileNumber_Is_Invalid() {
	payload := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: "0987654321",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	suite.context.Request = req
	suite.controller.CreateUser(suite.context)
	expectedResponseBody := `{"errorCode":"ERR_CC_BAD_REQUEST","errorMessage":"invalid mobile number","additionalData":{"statusCode":400}}`
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *UserControllerTestSuite) TestCreateUser_Failed_When_UserService_Returns_Error() {
	payload := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
	}

	expectedResponseBody := `{"errorCode":"ERR_CC_INTERNAL_SERVER_ERROR","errorMessage":"internal server error","additionalData":{"statusCode":500}}`

	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	suite.context.Request = req

	suite.mockUserService.EXPECT().CreateUser(suite.context, payload).Return(nil, utils.NewCCInternalServerError())

	suite.controller.CreateUser(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(expectedResponseBody, suite.recorder.Body.String())
}

func (suite *UserControllerTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
