package service

import "etherenum-api/etherenum-service/api/internal/entities"

type Service struct {
	Transaction TransactionService
}

type TransactionService interface {
	GetAll() (*[]entities.Transaction, error)
}

type transactionsOutput struct {
	transactions *[]entities.Transaction
}
