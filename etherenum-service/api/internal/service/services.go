package service

import (
	"context"
	"etherenum-api/etherenum-service/api/internal/entities"
)

type Service struct {
	Transaction TransactionService
}

type TransactionService interface {
	GetAll(ctx context.Context, query int64) (*[]entities.Transaction, error)
	GetByFilter(ctx context.Context, body string, page int64) (*entities.Transactions, error)
	Insert(result int64, transactions []entities.Transaction) (*entities.Transactions, error)
}

type (
	//LessPointerZeroError struct{ error }
	NilPointerDataError struct{ error }
)
