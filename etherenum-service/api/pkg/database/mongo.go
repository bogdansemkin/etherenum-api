package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	Name   string
	User   string
	Pass   string
	DBname string

}

type MongoDB struct {
	*mongo.Collection
}

func NewMongo(config MongoDBConfig) (*mongo.Collection, error) {
	//"mongodb+srv://bogdansemkin:<password>@cluster0.kcby05n.mongodb.net/test"
	opt := options.Client().ApplyURI(fmt.Sprintf("%s+srv://%s:%s@cluster0.kcby05n.mongodb.net/%s", config.Name, config.User, config.Pass, config.DBname))
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

	return transactionsCollection, nil
}
