package service

import "etherenum-api/etherenum-service/api/internal/entities"

type Service struct {
	Transaction TransactionService
}

type TransactionService interface {
	GetAll(page string) (*[]entities.Transaction, error)
	GetByFilter(body string) (*entities.Transactions, error)
}

type transactionsOutput struct {
	transactions *[]entities.Transaction
}
