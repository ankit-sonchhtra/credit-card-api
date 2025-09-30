package repository

import (
	"context"
	"testing"

	"github.com/credit-card-service/internal/repository/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	accountId      = "52fdfc07-2182-454f-963f-5f0f9a621d72"
	userId         = "cb9841df-c22e-4897-abfb-2411fad3e03d"
	documentNumber = "0123456789"
	accDocument    = model.AccountDocument{
		AccountId:      accountId,
		UserId:         userId,
		DocumentNumber: documentNumber,
		CreatedAt:      1759170600000,
		UpdatedAt:      1759170600000,
	}
)

func TestCreateAccount(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("it should successfully insert", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		repo := accountRepository{collection: mt.Coll}
		err := repo.CreateAccount(context.Background(), accDocument)

		require.Nil(t, err)
	})

	mt.Run("it should return an error on failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index: 0, Code: 123, Message: "insert error",
		}))
		repo := accountRepository{collection: mt.Coll}
		err := repo.CreateAccount(context.Background(), accDocument)

		require.NotNil(t, err)
		require.Equal(t, "write exception: write errors: [insert error]", err.Error())

	})
}

func TestGetAccount(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("it should return an account successfully", func(mt *mtest.T) {
		accountDoc := bson.D{
			{Key: "account_id", Value: accountId},
			{Key: "user_id", Value: userId},
			{Key: "document_number", Value: documentNumber},
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "credit-card-api.accounts", mtest.FirstBatch, accountDoc))

		repo := accountRepository{collection: mt.Coll}
		account, err := repo.GetAccount(context.Background(), accountId)

		require.NoError(t, err)
		require.NotNil(t, account)
		require.Equal(t, accountId, account.AccountId)
	})

	mt.Run("it should return an error when account not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "credit-card-api.accounts", mtest.FirstBatch))

		repo := accountRepository{collection: mt.Coll}
		account, err := repo.GetAccount(context.Background(), accountId)

		require.NoError(mt, err)
		require.Nil(mt, account)
	})

	mt.Run("it should return an error when it mongo gives an error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
			mtest.WriteError{Code: 11000, Message: "some mongo error"},
		))

		repo := accountRepository{collection: mt.Coll}
		account, err := repo.GetAccount(context.Background(), accountId)

		require.Error(mt, err)
		require.Nil(mt, account)
	})
}
