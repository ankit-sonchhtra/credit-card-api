package repository

//go:generate mockgen -source=user_repo.go -destination=mocks/mock_user_repo.go -package=mocks

import (
	"context"
	"errors"

	"github.com/credit-card-service/internal/repository/model"
	"github.com/credit-card-service/pkg/constants"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userDoc model.UserDocument) error
	GetUser(ctx context.Context, mobileNumber string) (*model.UserDocument, error)
}

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{collection: collection}
}

func (ar *userRepository) CreateUser(ctx context.Context, userDoc model.UserDocument) error {
	_, err := ar.collection.InsertOne(ctx, userDoc)
	if err != nil {
		logger.Errorf("error while create an user: %v", err)
		return err
	}
	logger.Info("user created successfully.")
	return nil
}

func (ar *userRepository) GetUser(ctx context.Context, mobileNumber string) (user *model.UserDocument, err error) {
	filter := bson.M{constants.MobileNumberFilter: mobileNumber}
	err = ar.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Infof("no user found with id: %s", mobileNumber)
			return nil, nil
		}

		logger.Errorf("error while fetching user: %v", err)
		return nil, err
	}
	logger.Info("user fetched successfully.")
	return user, nil
}
