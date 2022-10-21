package etherscan

import (
	"context"
	"etherenum-api/etherenum-service/api/internal/config"
	"etherenum-api/etherenum-service/api/internal/entities"
	"etherenum-api/etherenum-service/api/internal/service"
	"etherenum-api/etherenum-service/api/pkg/hex"
	"etherenum-api/etherenum-service/api/pkg/json"
	"etherenum-api/etherenum-service/api/pkg/logger"
	"fmt"
	"time"
)

var _ Scanner = (*etherscan)(nil)

type etherscan struct {
	Config    *config.Config
	Logger    logger.Logger
	Service   service.Service
	Converter *hex.Converter
}

func NewEtherscan(config *config.Config, logger logger.Logger, service service.Service, converter *hex.Converter) *etherscan {
	return &etherscan{
		Config:    config,
		Logger:    logger,
		Service:   service,
		Converter: converter,
	}
}

type getBlockNumberBody struct {
	ID     int
	Result string
}

func (e *etherscan) GetBlock() (*getBlockNumberBody, error) {
	logger := e.Logger.Named("GetBlock")
	var body getBlockNumberBody

	err := json.GetJson(fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=%s", e.Config.Etherscan.Key), &body)
	if err != nil {
		logger.Error("error during getting json from etherscan", "err", err)
		return nil, fmt.Errorf("error during getting json from etherscan, %s", err)
	}

	return &body, nil
}

type GetTransactionsBody struct {
	ID     int
	Result Result
}

func (e *etherscan) GetTransactions(result string) ([]entities.Transaction, error) {
	logger := e.Logger.
		Named("GetTransactions").
		With("result", result)

	var body GetTransactionsBody
	var transactions []entities.Transaction

	err := json.GetJson(fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", result, e.Config.Etherscan.Key), &body)
	if err != nil {
		logger.Error("error during getting json from etherscan", "err", err)
		return nil, fmt.Errorf("error during getting json from etherscan, %s", err)
	}

	for i := range body.Result.Transactions {
		body.Result.Transactions[i].Timestamp = body.Result.Timestamp

		transactions = append(transactions, entities.Transaction{
			Hash:        body.Result.Transactions[i].Hash,
			From:        body.Result.Transactions[i].From,
			To:          body.Result.Transactions[i].To,
			BlockNumber: e.Converter.HexaNumberToInteger(body.Result.Transactions[i].BlockNumber),
			Gas:         body.Result.Transactions[i].Gas,
			GasPrice:    e.Converter.HexaNumberToInteger(body.Result.Transactions[i].GasPrice),
			Timestamp:   body.Result.Transactions[i].Timestamp,
			CreateAt:    time.Now(),
		})
	}

	return transactions, nil
}

func (e *etherscan) InputTransactions() error {
	logger := e.Logger.Named("InputTransactions")

	body, err := e.GetBlock()
	if body == nil {
		logger.Error("failed to get block: body is empty. ", "err", err)
		return fmt.Errorf("failed to get block: body is empty: %s", err)
	}
	if err != nil {
		logger.Error("failed to get block", "err", err)
		return fmt.Errorf("failed to get block, %s", err)
	}

	getTransactions, err := e.GetTransactions(body.Result)
	if err != nil {
		logger.Error("error during getting the transaction", "err", err)
		return fmt.Errorf("error during getting the transaction, %s\n", err)
	}

	_, err = e.Service.Transaction.Insert(e.Converter.HexaNumberToInteger(body.Result), getTransactions)
	if err != nil {
		logger.Error("error on transaction insert", "err", err)
		return fmt.Errorf("error on transaction insert, %s", err)
	}

	return nil
}

func (e *etherscan) InitBlocks() error {
	logger := e.Logger.Named("InitBlocks")

	body, err := e.GetBlock()
	if body == nil {
		logger.Error("failed to get block: body is empty. ", "err", err)
		return fmt.Errorf("failed to get block: body is empty: %s", err)
	}
	if err != nil {
		logger.Error("failed to get block", "err", err)
		return fmt.Errorf("failed to get block, %s", err)
	}

	//todo create custom struct for getAll
	transactions, err := e.Service.Transaction.GetAll(context.TODO(), 1)
	if err != nil {
		logger.Error("error on get all transactions", "err", err)
		return fmt.Errorf("error on get all transactions, %s", err)
	}

	if *transactions == nil {
		logger.Info("There are any transactions onto database. Init blocks")

		//we can easily change point from 10 block to 1000, like in technical task
		for i := 0; i <= 10; i++ {
			*transactions, err = e.GetTransactions(body.Result)
			if err != nil {
				logger.Error("error during getting the transaction", "err", err)
				return fmt.Errorf("error during getting the transaction, %s\n", err)
			}

			_, err = e.Service.Transaction.Insert(e.Converter.HexaNumberToInteger(body.Result)-int64(i), *transactions)
			if err != nil {
				logger.Error("error on transaction insert", "err", err)
				return fmt.Errorf("error on transaction insert, %s", err)
			}
		}
		return nil
	}

	return nil
}
