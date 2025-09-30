package routes

import (
	"github.com/credit-card-api/bootstrap"
	"github.com/credit-card-api/internal/controllers"
	"github.com/credit-card-api/internal/repository"
	"github.com/credit-card-api/internal/services"
	"github.com/credit-card-api/pkg/constants"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes() *gin.Engine {
	router := gin.Default()

	userCollection := bootstrap.GetCollection(constants.UserCollection)
	userRepository := repository.NewUserRepository(userCollection)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	accountCollection := bootstrap.GetCollection(constants.AccountCollection)
	accountRepository := repository.NewAccountRepository(accountCollection)
	accountService := services.NewAccountService(accountRepository)
	accountController := controllers.NewAccountController(accountService)

	transactionCollection := bootstrap.GetCollection(constants.TransactionCollection)
	transactionRepository := repository.NewTransactionRepository(transactionCollection)
	transactionService := services.NewTransactionService(transactionRepository, accountRepository)
	transactionController := controllers.NewTransactionController(transactionService)

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routerGroup := router.Group("/api/credit-card-api/v1")
	routerGroup.POST("/users", userController.CreateUser)
	routerGroup.POST("/accounts", accountController.CreateAccount)
	routerGroup.GET("/accounts/:accountId", accountController.GetAccount)
	routerGroup.POST("/transactions", transactionController.CreateTransaction)

	return router
}
