package bootstrap

import (
	"context"

	logger "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"sync"
	"time"
)

var (
	dbInstance *mongo.Database
	once       sync.Once
)

func ConnectToDB(uri string, dbName string) *mongo.Database {
	once.Do(func() {
		clientOptions := options.Client().ApplyURI(uri)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			logger.Fatalf("mongoDB connection error: %v", err)
		}

		if err := client.Ping(ctx, nil); err != nil {
			logger.Fatalf("mongoDB ping failed: %v", err)
		}

		logger.Println("connected to mongoDB successfully.")
		dbInstance = client.Database(dbName)
	})

	return dbInstance
}

func GetCollection(collectionName string) *mongo.Collection {
	if dbInstance == nil {
		logger.Fatal("database not connected.")
	}
	return dbInstance.Collection(collectionName)
}
