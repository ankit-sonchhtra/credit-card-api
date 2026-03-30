package main

import (
	"context"
	"log"
	"net/http"
	"os"

	_ "github.com/credit-card-api/docs"
	"github.com/credit-card-api/internal/repository/sqlc"
	"github.com/credit-card-api/internal/routes"
	"github.com/credit-card-api/pkg/constants"
	"github.com/jackc/pgx/v5/pgxpool"
	logger "github.com/sirupsen/logrus"
)

func main() {
	dbURL := os.Getenv(constants.DBUrl)
	if dbURL == constants.EmptyString {
		log.Fatal("please set DB_URL env as mentioned in README.md")
	}

	dbPool, dbErr := pgxpool.New(context.Background(), dbURL)
	if dbErr != nil {
		log.Fatal("unable to connect with database:", dbErr)
	}

	logger.Info("database connected successfully.")
	defer dbPool.Close()

	queries := sqlc.New(dbPool)
	router := routes.RegisterRoutes(queries)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		logger.Fatal("failed to start server")
	}
}
