package service

import "etherenum-api/etherenum-service/api/internal/entities"

type Service struct {
	Transaction TransactionService
}

type TransactionService interface {
	GetAll(page int64) (*[]entities.Transaction, error)
	GetByFilter(body string) (*entities.Transactions, error)
}

type (
	//LessPointerZeroError struct{ error }
	NilPointerDataError struct{ error }
)
