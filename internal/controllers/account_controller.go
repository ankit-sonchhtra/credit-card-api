package controllers

import (
	"net/http"
	"strings"

	"github.com/credit-card-service/internal/models"
	"github.com/credit-card-service/internal/services"
	"github.com/credit-card-service/pkg/constants"
	"github.com/credit-card-service/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AccountController interface {
	CreateAccount(ctx *gin.Context)
	GetAccount(ctx *gin.Context)
}

type accountController struct {
	accountService services.AccountService
}

func NewAccountController(accountService services.AccountService) AccountController {
	return &accountController{
		accountService: accountService,
	}
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
// @Failure      500  {object}  models.InternalServerError
// @Router       /api/credit-card-api/v1/accounts [post]
func (ac *accountController) CreateAccount(ctx *gin.Context) {
	var payload models.CreateAccountRequest
	err := ctx.BindJSON(&payload)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.InvalidRequestBodyErrMsg))
		return
	}

	response, createErr := ac.accountService.CreateAccount(ctx, payload)
	if createErr != nil {
		additionalData := createErr.AdditionalData.(models.AdditionalData)
		ctx.AbortWithStatusJSON(additionalData.StatusCode, createErr)
		return
	}

	ctx.JSON(http.StatusCreated, response)
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
func (ac *accountController) GetAccount(ctx *gin.Context) {
	accountId := ctx.Param(constants.AccountIdPathParam)

	if strings.TrimSpace(accountId) == constants.EmptyString {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.AccountIdMissingErrMsg))
		return
	}

	response, err := ac.accountService.GetAccount(ctx, accountId)
	if err != nil {
		additionalData := err.AdditionalData.(models.AdditionalData)
		ctx.AbortWithStatusJSON(additionalData.StatusCode, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
	return
}
