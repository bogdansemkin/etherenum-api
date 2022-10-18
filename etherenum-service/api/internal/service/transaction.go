package service

import (
	"etherenum-api/etherenum-service/api/internal/entities"
	"fmt"
	"strconv"
)

var _ TransactionService = (*transactionService)(nil)

type transactionService struct {
	repos Repos
}

func NewTransactionService(repos Repos) *transactionService {
	return &transactionService{repos: repos}
}

func (s *transactionService) GetAll(query string) (*[]entities.Transaction, error) {
	page, err := strconv.Atoi(query)
	if err != nil {
		fmt.Printf("error during converting page, %s\n", err)
		return nil, err
	}

	transactions, err := s.repos.Transactions.GetAll(int64(page))
	if err != nil {
		return nil, fmt.Errorf("error during get all transactions , %s", err)
	}

	return transactions, nil
}

func (s *transactionService) GetByFilter(body string) (*entities.Transactions, error) {
	transactions, err := s.repos.Transactions.GetByFilter(body)
	fmt.Println(transactions)
	if err != nil {
		return nil, fmt.Errorf("error during get all transactions , %s", err)
	}
	return transactions, nil
}
