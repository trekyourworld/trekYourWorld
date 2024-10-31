package db

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"
	"trekyourworld/env"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"oss.nandlabs.io/golly/l3"
)

var (
	clientInstance    *mongo.Client
	clientInstanceErr error
	mongoOnce         sync.Once
	mongoDatabaseName string
)

var logger = l3.Get()

func Init() error {
	if err := env.LoadEnv(); err != nil {
		logger.Error("Error loading .env file")
	}

	mongoDatabaseName = os.Getenv("DB_NAME")
	mongoURI := os.Getenv("DB_URI")
	if mongoURI == "" || mongoDatabaseName == "" {
		return errors.New("DB URI and DB Name is not set in environment")
	}

	var initErr error
	mongoOnce.Do(func() {
		clientOptions := options.Client().ApplyURI(mongoURI)

		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			initErr = err
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Ping(ctx, nil); err != nil {
			initErr = err
			return
		}

		clientInstance = client
		logger.Info("DB connected successfully")
	})

	clientInstanceErr = initErr
	return initErr
}

func GetClient() (*mongo.Client, error) {
	if clientInstanceErr != nil {
		return nil, clientInstanceErr
	}
	return clientInstance, nil
}

func GetCollection(collectionName string) (*mongo.Collection, error) {
	client, err := GetClient()
	if err != nil {
		return nil, err
	}
	return client.Database(mongoDatabaseName).Collection(collectionName), nil
}
