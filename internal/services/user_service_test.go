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
	testUserId       string
	testEmail        string
	testMobileNumber string
)

type UserServiceTestSuite struct {
	suite.Suite
	context            context.Context
	mockController     *gomock.Controller
	mockUserRepository *mocks.MockUserRepository
	userService        UserService
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.context = context.TODO()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockUserRepository = mocks.NewMockUserRepository(suite.mockController)
	suite.userService = NewUserService(suite.mockUserRepository)
	testUserId = "52fdfc07-2182-454f-963f-5f0f9a621d72"
	testMobileNumber = "+918908011223"
	testEmail = "abc@xyz.com"
	currentTime = mockNow
	uuid.SetRand(rand.New(rand.NewSource(1)))
}

func (suite *UserServiceTestSuite) TestCreateUser_Success() {
	request := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
	}
	expectedResponse := &models.CreateUserResponse{
		UserId: testUserId,
	}
	userDocument := model.UserDocument{
		UserId:       testUserId,
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
		CreatedAt:    1759170600000,
		UpdatedAt:    1759170600000,
	}
	suite.mockUserRepository.EXPECT().GetUser(suite.context, testMobileNumber).Return(nil, nil)
	suite.mockUserRepository.EXPECT().CreateUser(suite.context, userDocument).Return(nil)
	response, err := suite.userService.CreateUser(suite.context, request)

	suite.Nil(err)
	suite.Equal(expectedResponse, response)
}

func (suite *UserServiceTestSuite) TestCreateUser_Failed_When_GetUser_Fails() {
	request := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
	}

	suite.mockUserRepository.EXPECT().GetUser(suite.context, testMobileNumber).Return(nil, errors.New("failed to fetch"))
	_, err := suite.userService.CreateUser(suite.context, request)

	suite.NotNil(err)
	suite.Equal(utils.NewCCInternalServerError(), err)
}

func (suite *UserServiceTestSuite) TestCreateUser_Failed_When_User_IsAlready_Exist() {
	request := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
	}
	expectedErr := &models.CCError{
		ErrorCode:      "ERR_CC_USER_ALREADY_EXIST",
		ErrorMessage:   "user already exist with requested mobile number",
		AdditionalData: models.AdditionalData{StatusCode: 409},
	}
	userDocument := &model.UserDocument{
		UserId:       testUserId,
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
		CreatedAt:    1759170600000,
		UpdatedAt:    1759170600000,
	}

	suite.mockUserRepository.EXPECT().GetUser(suite.context, testMobileNumber).Return(userDocument, nil)

	_, err := suite.userService.CreateUser(suite.context, request)

	suite.NotNil(err)
	suite.Equal(expectedErr, err)
}

func (suite *UserServiceTestSuite) TestCreateUser_When_UserRepo_Returns_AnError() {
	request := models.CreateUserRequest{
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
	}
	userDocument := model.UserDocument{
		UserId:       testUserId,
		Name:         "John Doe",
		Email:        testEmail,
		MobileNumber: testMobileNumber,
		CreatedAt:    1759170600000,
		UpdatedAt:    1759170600000,
	}
	suite.mockUserRepository.EXPECT().GetUser(suite.context, testMobileNumber).Return(nil, nil)
	suite.mockUserRepository.EXPECT().CreateUser(suite.context, userDocument).Return(errors.New("failed to create"))
	response, err := suite.userService.CreateUser(suite.context, request)

	suite.NotNil(err)
	suite.Nil(response)
	suite.Equal(utils.NewCCInternalServerError(), err)
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.mockController.Finish()
}
