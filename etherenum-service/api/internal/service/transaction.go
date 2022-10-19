package service

import (
	"etherenum-api/etherenum-service/api/internal/entities"
	"fmt"
)

var _ TransactionService = (*transactionService)(nil)

type transactionService struct {
	repos Repos
}

func NewTransactionService(repos Repos) *transactionService {
	return &transactionService{repos: repos}
}

func (s *transactionService) GetAll(query int64) (*[]entities.Transaction, error) {
	if query < 1 {
		query = 1
	}
	transactions, err := s.repos.Transactions.GetAll(query)
	if err != nil {
		return nil, fmt.Errorf("error during get all transactions , %s", err)
	}

	return transactions, nil
}

func (s *transactionService) GetByFilter(body string) (*entities.Transactions, error) {
	transactions, err := s.repos.Transactions.GetByFilter(body)
	if err != nil {
		return nil, fmt.Errorf("error during get all transactions , %s", err)
	}
	if transactions.Trans == nil {
		return nil, NilPointerDataError{}
	}
	return transactions, nil
}
