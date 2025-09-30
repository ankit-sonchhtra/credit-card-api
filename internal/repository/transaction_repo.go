package repository

//go:generate mockgen -source=transaction_repo.go -destination=mocks/mock_transaction_repo.go -package=mocks

import (
	"context"

	"github.com/credit-card-service/internal/repository/model"
	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, txnDocument model.TransactionDocument) error
}

type transactionRepository struct {
	collection *mongo.Collection
}

func NewTransactionRepository(collection *mongo.Collection) TransactionRepository {
	return &transactionRepository{collection: collection}
}

func (tr *transactionRepository) CreateTransaction(ctx context.Context, txnDocument model.TransactionDocument) error {
	_, err := tr.collection.InsertOne(ctx, txnDocument)
	if err != nil {
		logger.Errorf("error while inserting a transaction: %v", err)
		return err
	}
	return nil
}
