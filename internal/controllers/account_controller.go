package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/credit-card-api/internal/domain"
	"github.com/credit-card-api/internal/models"
	"github.com/credit-card-api/internal/services"
	"github.com/credit-card-api/pkg/constants"
	"github.com/credit-card-api/pkg/utils"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

type AccountController struct {
	accountService services.AccountService
}

func NewAccountController(accountService services.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}

// CreateAccount godoc
// @Summary      Create an account
// @Description  Create an account by request payload
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param CreateAccountRequest body models.CreateAccountRequest true "Request Body"
// @Success      201  {object}  models.CreateAccountResponse
// @Failure      400  {object}  models.BadRequestError
// @Failure      409  {object}  models.ConflictError
// @Failure      500  {object}  models.InternalServerError
// @Router       /api/credit-card-api/v1/accounts [post]
func (ac *AccountController) CreateAccount(ctx *gin.Context) {
	var payload models.CreateAccountRequest
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		logger.Error("failed to binding a request payload error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.InvalidRequestBodyErrMsg))
		return
	}

	validationErr := payload.Validate()
	if validationErr != nil {
		logger.Error("validation failure on request payload error: ", validationErr)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(validationErr.Error()))
		return
	}

	account, accountErr := ac.accountService.RegisterAccount(ctx, payload)
	if accountErr != nil {
		ac.respondWithError(ctx, accountErr)
		return
	}
	ctx.JSON(http.StatusCreated, mapToCreateAccountResponse(*account))
	return
}

// GetAccount godoc
// @Summary      Get an account
// @Description  Get an account by accountId
// @Tags         Accounts
// @Produce      json
// @Param accountId path string true "accountId"
// @Success      200  {object}  models.GetAccountResponse
// @Failure      404  {object}  models.NotFoundError
// @Failure      500  {object}  models.InternalServerError
// @Router       /api/credit-card-api/v1/accounts/{accountId} [Get]
func (ac *AccountController) GetAccount(ctx *gin.Context) {
	accountIdStr := ctx.Param(constants.AccountIdPathParam)
	id, err := strconv.ParseInt(accountIdStr, 10, 64)
	if err != nil {
		logger.Error("failed to read path param error: ", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.AccountIdMissingErrMsg))
		return
	}

	account, accountErr := ac.accountService.GetAccount(ctx, id)
	if accountErr != nil {
		ac.respondWithError(ctx, accountErr)
		return
	}
	ctx.JSON(http.StatusOK, mapToGetAccountResponse(*account))
	return
}

func (ac *AccountController) respondWithError(ctx *gin.Context, err error) {
	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		appErr = domain.ErrInternal
	}

	status := http.StatusInternalServerError
	switch appErr.Code {
	case constants.AccountNotFoundErrCode:
		status = http.StatusNotFound
	case constants.AccountAlreadyExistErrCode:
		status = http.StatusConflict
	}

	ctx.AbortWithStatusJSON(status, &models.CCError{
		ErrorCode:    appErr.Code,
		ErrorMessage: appErr.Message,
		StatusCode:   status,
	})
}

func mapToCreateAccountResponse(account domain.Account) models.CreateAccountResponse {
	return models.CreateAccountResponse{
		AccountId:      account.Id,
		DocumentNumber: account.DocumentNumber,
	}
}

// This is duplicate code currently, but account would have more field to return in response.
func mapToGetAccountResponse(account domain.Account) models.GetAccountResponse {
	return models.GetAccountResponse{
		AccountId:      account.Id,
		DocumentNumber: account.DocumentNumber,
	}
}
