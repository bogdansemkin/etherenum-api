package service

import "fmt"

var _ TransactionService = (*transactionService) (nil)

type transactionService struct {
	repos Repos
}

func NewTransactionService(repos Repos) *transactionService {
	return &transactionService{repos: repos}
}

func (s *transactionService) GetAll()(*transactionsOutput, error) {
	fmt.Println("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	transactions, err := s.repos.Transactions.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error during get all transactions , %s", err)
	}

	return &transactionsOutput{transactions: transactions}, nil
}
