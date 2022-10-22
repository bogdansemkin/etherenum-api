package service

import (
	"context"
	"etherenum-api/etherenum-service/api/internal/entities"
	"etherenum-api/etherenum-service/api/pkg/commission"
	"etherenum-api/etherenum-service/api/pkg/logger"
	"fmt"
)

var _ TransactionService = (*transactionService)(nil)

type transactionService struct {
	repos      Repos
	logger     logger.Logger
	calculator *commission.Calculator
}

func NewTransactionService(repos Repos, logger logger.Logger, calculator *commission.Calculator) *transactionService {
	return &transactionService{
		repos:      repos,
		logger:     logger,
		calculator: calculator,
	}
}

func (s *transactionService) GetAll(ctx context.Context, query int64) (*[]entities.Transaction, error) {
	logger := s.logger.
		Named("GetAll").
		WithContext(ctx).
		With("query", query)

	if query < 1 {
		logger.Info("query is less then 1")
		query = 1
	}
	transactions, err := s.repos.Transactions.GetAll(query)
	if err != nil {
		logger.Error("error during get all transactions", "err", err)
		return nil, fmt.Errorf("error during get all transactions , %s", err)
	}

	return transactions, nil
}

func (s *transactionService) GetByFilter(ctx context.Context, body string, page int64) (*entities.Transactions, error) {
	logger := s.logger.
		Named("GetByFilter").
		WithContext(ctx).
		With("body", body)

	if page < 1 {
		logger.Info("query is less then 1")
		page = 1
	}
	transactions, err := s.repos.Transactions.GetByFilter(body, page)
	if err != nil {
		logger.Error("error during get transactions by filter", "err", err)
		return nil, fmt.Errorf("error during get transactions by filter, %s", err)
	}
	if transactions.Trans == nil {
		logger.Error("transactions by filter are empty")
		return nil, NilPointerDataError{}
	}

	return transactions, nil
}

func (s *transactionService) Insert(result int64, transactions []entities.Transaction) (*entities.Transactions, error) {
	logger := s.logger.
		Named("GetByFilter").
		With("result", result)

	var trainers []interface{}
	ok := s.repos.Transactions.CheckOnDuplicate(result)
	if !ok {
		logger.Error("duplicate found")
		return nil, fmt.Errorf("duplicate found")
	}

	for i := range transactions {
		transactions[i].Commission = s.calculator.GetCommission(transactions[i].Gas, transactions[i].GasPrice)
		trainers = append(trainers, transactions[i])
	}
	err := s.repos.Transactions.Insert(trainers)
	if err != nil {
		logger.Error("error during insert data", "err", err)
		return nil, fmt.Errorf("error during insert data, %s", err)
	}

	return &entities.Transactions{Trans: transactions}, nil
}
