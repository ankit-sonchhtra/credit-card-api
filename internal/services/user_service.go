package services

//go:generate mockgen -source=user_service.go -destination=mocks/mock_user_service.go -package=mocks

import (
	"context"

	"github.com/credit-card-service/internal/models"
	"github.com/credit-card-service/internal/repository"
	"github.com/credit-card-service/internal/repository/model"
	"github.com/credit-card-service/pkg/constants"
	"github.com/credit-card-service/pkg/utils"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, request models.CreateUserRequest) (*models.CreateUserResponse, *models.CCError)
}

type userService struct {
	usersRepo repository.UserRepository
}

func NewUserService(usersRepo repository.UserRepository) UserService {
	return &userService{usersRepo: usersRepo}
}

func (as *userService) CreateUser(ctx context.Context, request models.CreateUserRequest) (*models.CreateUserResponse, *models.CCError) {
	user, err := as.usersRepo.GetUser(ctx, request.MobileNumber)
	if err != nil {
		return nil, utils.NewCCInternalServerError()
	}

	if user != nil {
		return nil, &models.CCError{
			ErrorCode:      constants.UserAlreadyExistErrCode,
			ErrorMessage:   constants.UserAlreadyExistErrMsg,
			AdditionalData: models.AdditionalData{StatusCode: 409},
		}
	}

	userDocument := buildUserDocument(request)
	err = as.usersRepo.CreateUser(ctx, userDocument)
	if err != nil {
		return nil, utils.NewCCInternalServerError()
	}

	return &models.CreateUserResponse{UserId: userDocument.UserId}, nil
}

func buildUserDocument(request models.CreateUserRequest) model.UserDocument {
	userId := uuid.New().String()
	return model.UserDocument{
		UserId:       userId,
		Name:         request.Name,
		Email:        request.Email,
		MobileNumber: request.MobileNumber,
		CreatedAt:    currentTime().UnixMilli(),
		UpdatedAt:    currentTime().UnixMilli(),
	}
}
