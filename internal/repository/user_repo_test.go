package repository

import (
	"context"
	"testing"

	"github.com/credit-card-api/internal/repository/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var (
	testUserId   = "cb9841df-c22e-4897-abfb-2411fad3e03d"
	mobileNumber = "+919898161616"
	userDocument = model.UserDocument{
		UserId:    testUserId,
		CreatedAt: 1759170600000,
		UpdatedAt: 1759170600000,
	}
)

func TestCreateUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("it should insert successfully", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		repo := userRepository{collection: mt.Coll}
		err := repo.CreateUser(context.Background(), userDocument)

		require.Nil(t, err)
	})

	mt.Run("it should return an error on failure", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index: 0, Code: 123, Message: "insert error",
		}))
		repo := userRepository{collection: mt.Coll}
		err := repo.CreateUser(context.Background(), userDocument)

		require.NotNil(t, err)
		require.Equal(t, "write exception: write errors: [insert error]", err.Error())

	})
}

func TestGetUser(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("it should return an user successfully", func(mt *mtest.T) {
		userDoc := bson.D{
			{Key: "user_id", Value: testUserId},
			{Key: "name", Value: "John Doe"},
		}
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "credit-card-api.users", mtest.FirstBatch, userDoc))

		repo := userRepository{collection: mt.Coll}
		user, err := repo.GetUser(context.Background(), mobileNumber)

		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, testUserId, user.UserId)
	})

	mt.Run("it should return an error when user not found", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "credit-card-api.users", mtest.FirstBatch))

		repo := userRepository{collection: mt.Coll}
		user, err := repo.GetUser(context.Background(), mobileNumber)

		require.NoError(mt, err)
		require.Nil(mt, user)
	})

	mt.Run("it should return an error when it mongo gives an error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(
			mtest.WriteError{Code: 11000, Message: "some mongo error"},
		))

		repo := userRepository{collection: mt.Coll}
		user, err := repo.GetUser(context.Background(), mobileNumber)

		require.Error(mt, err)
		require.Nil(mt, user)
	})
}
