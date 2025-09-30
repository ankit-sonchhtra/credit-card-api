package main

// @title           Credit Card API
// @version         1.0
// @description     credit card api which support to make payment.
// @host      localhost:8080
// @BasePath  /api/v1
// @schemes http
import (
	"net/http"
	"os"

	"github.com/credit-card-api/bootstrap"
	_ "github.com/credit-card-api/docs"
	"github.com/credit-card-api/internal/routes"
	"github.com/credit-card-api/pkg/constants"
	logger "github.com/sirupsen/logrus"
)

func main() {

	uri := os.Getenv(constants.MongoUri)

	bootstrap.ConnectToDB(uri, constants.MongoDatabaseName)
	router := routes.RegisterRoutes()

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatal("failed to start server")
	}
}
