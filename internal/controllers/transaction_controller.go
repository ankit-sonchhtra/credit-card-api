package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/services"
	"github.com/credit-card-api/pkg/constants"
	"github.com/credit-card-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TransactionController interface {
	CreateTransaction(ctx *gin.Context)
}

type transactionController struct {
	transactionService services.TransactionService
}

func NewTransactionController(transactionService services.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

// CreateTransaction godoc
// @Summary      Create transaction
// @Description  Create transaction by request payload
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Param CreateTransactionRequest body models.CreateTransactionRequest true "Request Body"
// @Success      201  {object} models.CreateTransactionResponse
// @Failure      400  {object}  models.BadRequestError
// @Failure      500  {object}  models.InternalServerError
// @Router       /api/credit-card-api/v1/transactions [post]
func (tc *transactionController) CreateTransaction(ctx *gin.Context) {
	var payload models.CreateTransactionRequest
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.InvalidRequestBodyErrMsg))
		return
	}

	validationErr := validateAmount(payload)
	if validationErr != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(validationErr.Error()))
		return
	}

	txnResponse, txnErr := tc.transactionService.CreateTransaction(ctx, payload)
	if txnErr != nil {
		additionalData := txnErr.AdditionalData.(models.AdditionalData)
		ctx.AbortWithStatusJSON(additionalData.StatusCode, txnErr)
		return
	}
	ctx.JSON(http.StatusCreated, txnResponse)
}

func validateAmount(payload models.CreateTransactionRequest) error {
	operationType := strings.ToLower(payload.OperationType)

	switch operationType {
	case constants.OpTypeCashPurchase, constants.OpTypeInstallmentPurchase, constants.OpTypeWithdrawal:
		if payload.Amount >= 0 {
			return errors.New(constants.AmountMustBeNegativeErrMsg)
		}
	case constants.OpTypePayment:
		if payload.Amount <= 0 {
			return errors.New(constants.AmountMustBePositiveErrMsg)
		}
	default:
		return errors.New(constants.InvalidOperationTypeErrMsg)
	}
	return nil
}
