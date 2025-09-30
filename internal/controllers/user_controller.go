package controllers

import (
	"net/http"
	"regexp"

	"github.com/credit-card-service/internal/models"
	"github.com/credit-card-service/internal/services"
	"github.com/credit-card-service/pkg/constants"
	"github.com/credit-card-service/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
	//GetUser(ctx *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

// CreateUser godoc
// @Summary      Create an user
// @Description  Create an user by request payload
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param CreateUserRequest body models.CreateUserRequest true "Request Body"
// @Success      201  {object}  models.CreateUserResponse
// @Failure      400  {object}  models.BadRequestError
// @Failure      409  {object}  models.ConflictError
// @Failure      500  {object}  models.InternalServerError
// @Router       /api/credit-card-api/v1/users [post]
func (ac *userController) CreateUser(ctx *gin.Context) {
	var payload models.CreateUserRequest
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.InvalidRequestBodyErrMsg))
		return
	}
	if !isValidMobileNumber(payload.MobileNumber) {
		ctx.JSON(http.StatusBadRequest, utils.NewCCBadRequestError(constants.InvalidMobileNumberErrMsg))
		return
	}

	response, createErr := ac.userService.CreateUser(ctx, payload)
	if createErr != nil {
		additionalData := createErr.AdditionalData.(models.AdditionalData)
		ctx.JSON(additionalData.StatusCode, createErr)
		return
	}

	ctx.JSON(http.StatusCreated, response)
	return
}

func isValidMobileNumber(mobile string) bool {
	pattern := `^\+91[6-9][0-9]{9}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(mobile)
}
