package routes

import (
	"github.com/credit-card-api/internal/controllers"
	"github.com/credit-card-api/internal/repository"
	"github.com/credit-card-api/internal/repository/sqlc"
	"github.com/credit-card-api/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(queries *sqlc.Queries) *gin.Engine {
	router := gin.Default()

	accountRepository := repository.NewAccountRepository(queries)
	accountService := services.NewAccountService(accountRepository)
	accountController := controllers.NewAccountController(accountService)

	transactionRepository := repository.NewTransactionRepository(queries)
	transactionService := services.NewTransactionService(transactionRepository, accountRepository)
	transactionController := controllers.NewTransactionController(transactionService)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routerGroup := router.Group("/api/credit-card-api/v1")
	routerGroup.POST("/accounts", accountController.CreateAccount)
	routerGroup.GET("/accounts/:accountId", accountController.GetAccount)
	routerGroup.POST("/transactions", transactionController.CreateTransaction)

	return router
}
