package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/repository/mocks"
	"github.com/credit-card-api/internal/repository/sqlc"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	context               context.Context
	mockController        *gomock.Controller
	mockQuerier           *mocks.MockQuerier
	transactionRepository TransactionRepository
}

func TestTransactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	suite.context = context.TODO()
	suite.mockController = gomock.NewController(suite.T())
	suite.mockQuerier = mocks.NewMockQuerier(suite.mockController)
	suite.transactionRepository = NewTransactionRepository(suite.mockQuerier)
	accountId = 1
	documentNumber = "0123456789"
}

func (suite *TransactionRepositoryTestSuite) TestTransactionRepository_Create() {
	params := domain.CreateTransactionParam{
		AccountId:       1,
		OperationTypeId: 4,
		Amount:          -123.45,
	}

	dbResult := sqlc.Transaction{
		TransactionID:   100,
		AccountID:       1,
		OperationTypeID: 4,
		Amount:          float64ToNumeric(-123.45), // Using your helper
	}

	suite.mockQuerier.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(dbResult, nil).Times(1)

	res, err := suite.transactionRepository.Create(suite.context, params)

	suite.NoError(err)
	suite.NotNil(res)
	suite.Equal(int64(100), res.Id)
	suite.Equal(-123.45, res.Amount)
}

func (suite *TransactionRepositoryTestSuite) TestTransactionRepository_Create_Returns_Database_Error() {

	params := domain.CreateTransactionParam{AccountId: 1}
	err := errors.New("failed to create transaction")
	suite.mockQuerier.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(sqlc.Transaction{}, err)

	res, err := suite.transactionRepository.Create(suite.context, params)

	suite.Error(err)
	suite.Nil(res)

}
