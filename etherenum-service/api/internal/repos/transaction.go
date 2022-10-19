package repos

import (
	"context"
	"etherenum-api/etherenum-service/api/internal/entities"
	"etherenum-api/etherenum-service/api/internal/service"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var _ service.TransactionRepo = (*transactionRepo)(nil)

type transactionRepo struct {
	collection *mongo.Collection
}

func NewTransactionRepo(collection *mongo.Collection) *transactionRepo {
	return &transactionRepo{
		collection: collection,
	}
}

func (r *transactionRepo) GetAll(page int64) (*[]entities.Transaction, error) {
	var transactions []entities.Transaction

	cur, err := r.collection.Find(context.TODO(), bson.D{}, options.Find().SetSkip(5*page).SetLimit(5))
	if err != nil {
		return nil, fmt.Errorf("error during getting all transactions, %s", err)
	}

	for cur.Next(context.TODO()) {

		var transaction entities.Transaction
		err := cur.Decode(&transaction)
		if err != nil {
			return nil, fmt.Errorf("error during decoding bson, %s", err)
		}
		//TODO изменить реализацию
		//плохая реализация для большого объёма данных, потому что размер увеличивается на 2
		transactions = append(transactions, transaction)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return &transactions, nil
}

func (r *transactionRepo) GetByFilter(body string) (*entities.Transactions, error) {
	var transactions []entities.Transaction
	filter := bson.D{
		{"$or", bson.A{
			bson.D{{"hash", body}},
			bson.D{{"blocknumber", body}},
			bson.D{{"from", body}},
			bson.D{{"to", body}},
			bson.D{{"timestamp", body}},
		}},
	}

	cursor, err := r.collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, fmt.Errorf("error during finding data by filter, %s", err)
	}

	if err = cursor.All(context.TODO(), &transactions); err != nil {
		return nil, fmt.Errorf("error during finding data by filter, %s", err)
	}

	return &entities.Transactions{Trans: transactions}, nil
}

func (r *transactionRepo) Insert(data []interface{}) error {
	_, err := r.collection.InsertMany(context.TODO(), data)
	if err != nil {
		return fmt.Errorf("error during inserting many fields, %s", err)
	}
	return nil
}
