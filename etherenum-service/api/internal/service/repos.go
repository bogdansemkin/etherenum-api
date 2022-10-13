package service

import "etherenum-api/etherenum-service/api/internal/entities"

type Repos struct {
	Transactions TransactionRepo
}

type TransactionRepo interface {
	GetAll()(*[]entities.Transaction, error)
}
