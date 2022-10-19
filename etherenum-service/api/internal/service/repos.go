package service

import "etherenum-api/etherenum-service/api/internal/entities"

type Repos struct {
	Transactions TransactionRepo
}

type TransactionRepo interface {
	GetAll(page int64) (*[]entities.Transaction, error)
	GetByFilter(body string) (*entities.Transactions, error)
	Insert(data []interface{}) error
	CheckOnDuplicate(body string) bool
}
