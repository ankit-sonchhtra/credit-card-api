package repository

import (
	"context"
	"testing"

	"github.com/credit-card-service/internal/repository/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	testTransactionId = "52fdfc07-2182-454f-963f-5f0f9a621d72"
	transactionDoc    = model.TransactionDocument{
		Id:            testTransactionId,
		AccountId:     accountId,
		OperationType: "CASH PURCHASE",
		Amount:        -2345.67,
		CreatedAt:     1759170600000,
		UpdatedAt:     1759170600000,
	}
)

func TestCreateTransaction(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("it should insert successfully", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		repo := transactionRepository{collection: mt.Coll}
		err := repo.CreateTransaction(context.Background(), transactionDoc)

		require.Nil(t, err)
	})

	mt.Run("it should return an error on failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index: 0, Code: 123, Message: "insert error",
		}))
		repo := transactionRepository{collection: mt.Coll}
		err := repo.CreateTransaction(context.Background(), transactionDoc)

		require.NotNil(t, err)
		require.Equal(t, "write exception: write errors: [insert error]", err.Error())

	})
}
