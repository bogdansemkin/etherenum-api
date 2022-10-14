package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	Name string
	Port string
	Host string
}

type MongoDB struct {
	*mongo.Collection
}

func NewMongo(config MongoDBConfig) (*mongo.Collection, error) {
	opt := options.Client().ApplyURI(fmt.Sprintf("%s://%s:%s", config.Name, config.Host, config.Port))
	client, err := mongo.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("error during connection to mongo db: %s", err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error during client connection to mongo db: %s", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, fmt.Errorf("error during ping : %s", err)
	}
	database := client.Database("etherenum-api")
	transactionsCollection := database.Collection("transactions")
	//delete while finished
	defer transactionsCollection.Drop(context.TODO())

	return transactionsCollection, nil
}
