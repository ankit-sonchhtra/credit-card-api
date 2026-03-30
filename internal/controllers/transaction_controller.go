package controllers

import (
	"errors"
	"net/http"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/services"
	"github.com/credit-card-api/pkg/constants"
	"github.com/credit-card-api/pkg/utils"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

type TransactionController struct {
	transactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

// CreateTransaction godoc
// @Summary      Create transaction
// @Description  Create transaction by request payload
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param CreateTransactionRequest body models.TransactionRequest true "Request Body"
// @Success      201  {object} models.CreateTransactionResponse
// @Failure      400  {object}  models.BadRequestError
// @Failure      422  {object}  models.UnprocessableEntityError
// @Failure      500  {object}  models.InternalServerError
// @Router       /api/credit-card-api/v1/transactions [post]
func (tc *TransactionController) CreateTransaction(ctx *gin.Context) {
	var payload models.TransactionRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		logger.Error("failed to binding a request payload error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.InvalidRequestBodyErrMsg))
		return
	}

	validationErr := payload.Validate()
	if validationErr != nil {
		logger.Error("validation failure on request payload error:", validationErr)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(validationErr.Error()))
		return
	}

	transaction, txnErr := tc.transactionService.CreateTransaction(ctx, payload)
	if txnErr != nil {
		tc.respondWithError(ctx, txnErr)
		return
	}
	ctx.JSON(http.StatusCreated, mapToCreateTransactionResponse(*transaction))
}

func (tc *TransactionController) respondWithError(ctx *gin.Context, err error) {
	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		appErr = domain.ErrInternal
	}

	status := http.StatusInternalServerError
	switch appErr.Code {
	case constants.InvalidOperationTypeErrCode, constants.TransactionAccountNotFoundErrCode:
		status = http.StatusUnprocessableEntity
	}

	ctx.AbortWithStatusJSON(status, &models.CCError{
		ErrorCode:    appErr.Code,
		ErrorMessage: appErr.Message,
		StatusCode:   status,
	})
}

func mapToCreateTransactionResponse(transaction domain.Transaction) models.CreateTransactionResponse {
	return models.CreateTransactionResponse{
		TransactionId: transaction.Id,
	}
}
