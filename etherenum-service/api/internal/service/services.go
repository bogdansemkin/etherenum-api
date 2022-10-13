package service

import "etherenum-api/etherenum-service/api/internal/entities"

type Service struct {
	Transaction TransactionService
}

type TransactionService interface {
	GetAll() (*transactionsOutput, error)
}

type transactionsOutput struct {
	//ID         string
	//From       string
	//To         string
	//Block      string
	//Accepts    string
	//Date       string
	//Value      string
	//Commission string
	transactions *[]entities.Transaction
}
