package repository

//go:generate mockgen -source=account_repo.go -destination=mocks/mock_account_repo.go -package=mocks

import (
	"context"
	"errors"

	"github.com/credit-card-api/internal/repository/model"
	"github.com/credit-card-api/pkg/constants"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, accountDoc model.AccountDocument) error
	GetAccount(ctx context.Context, id string) (*model.AccountDocument, error)
}

type accountRepository struct {
	collection *mongo.Collection
}

func NewAccountRepository(collection *mongo.Collection) AccountRepository {
	return &accountRepository{collection: collection}
}

func (ar *accountRepository) CreateAccount(ctx context.Context, accountDoc model.AccountDocument) error {
	_, err := ar.collection.InsertOne(ctx, accountDoc)
	if err != nil {
		logger.Errorf("error while create an account: %v", err)
		return err
	}
	logger.Info("account created successfully.")
	return nil
}

func (ar *accountRepository) GetAccount(ctx context.Context, id string) (account *model.AccountDocument, err error) {
	filter := bson.M{constants.AccountIdFilter: id}
	err = ar.collection.FindOne(ctx, filter).Decode(&account)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Infof("no account found with id: %s", id)
			return nil, nil
		}

		logger.Errorf("error while fetching account: %v", err)
		return nil, err
	}
	logger.Info("account fetched successfully.")
	return account, nil
}
