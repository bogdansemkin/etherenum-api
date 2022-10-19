package etherscan

import (
	"etherenum-api/etherenum-service/api/internal/config"
	"etherenum-api/etherenum-service/api/internal/entities"
	"etherenum-api/etherenum-service/api/internal/service"
	"etherenum-api/etherenum-service/api/pkg/json"
	"etherenum-api/etherenum-service/api/pkg/logger"
	"fmt"
)

var _ Scanner = (*etherscan)(nil)

type etherscan struct {
	Config *config.Config
	Logger logger.Logger
	Repos  service.Service
}

func NewEtherscan(config *config.Config, logger logger.Logger, repos service.Service) *etherscan {
	return &etherscan{
		Config: config,
		Logger: logger,
		Repos: repos,
	}
}

type getBlockNumberBody struct {
	ID     int
	Result string
}

func (e *etherscan) GetBlock() (*getBlockNumberBody, error) {
	var body getBlockNumberBody

	err := json.GetJson(fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_blockNumber&apikey=%s", e.Config.Etherscan.Key), &body)
	if err != nil {
		return nil, fmt.Errorf("error during getting json from etherscan, %s", err)
	}

	return &body, nil
}

type GetTransactionsBody struct {
	ID     int
	Result Result
}

func (e *etherscan) GetTransactions(result string) ([]Transaction, error) {
	var body GetTransactionsBody

	err := json.GetJson(fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", result, e.Config.Etherscan.Key), &body)
	if err != nil {
		return nil, fmt.Errorf("error during getting json from etherscan, %s", err)
	}
	for i := range body.Result.Transactions {
		body.Result.Transactions[i].Timestamp = body.Result.Timestamp
	}
	return body.Result.Transactions, nil
}

func (e *etherscan) AcceptIncrement(newBlockTransactions, oldBlockTransactions []Transaction) ([]Transaction, error) {
	var i int
	for {
		switch true {
		case newBlockTransactions == nil || oldBlockTransactions == nil:
			fmt.Println("nil pointer exception")
			return nil, fmt.Errorf("nil pointer exception")
		case i == len(newBlockTransactions)-1 || i == len(oldBlockTransactions)-1:
			return newBlockTransactions, nil
		case e.CompareTransactions(&newBlockTransactions[i], &Transactions{Trans: oldBlockTransactions}):
			newBlockTransactions[i].AcceptNumber++
		}
		i++
	}
}

type Compare struct {
	From         string
	To           string
	Gas          string
	GasPrice     string
	AcceptNumber int
}

func (e *etherscan) CompareTransactions(block *Transaction, oldblock *Transactions) bool {
	compareNewBlock := Compare{
		From:     block.From,
		To:       block.To,
		Gas:      block.Gas,
		GasPrice: block.GasPrice,
	}
	for i := range oldblock.Trans {
		compareOldBlock := Compare{
			From:         oldblock.Trans[i].From,
			To:           oldblock.Trans[i].To,
			Gas:          oldblock.Trans[i].Gas,
			GasPrice:     oldblock.Trans[i].GasPrice,
			AcceptNumber: oldblock.Trans[i].AcceptNumber,
		}

		if compareNewBlock == compareOldBlock {
			block.AcceptNumber = compareOldBlock.AcceptNumber + 1
			return true
		}
	}
	return false
}

func (e *etherscan) InputData() ([]string, error) {
	logger := e.Logger.Named("InputData")
	var transactions []entities.Transaction

	body, err := e.GetBlock()
	if body == nil {
		logger.Error("failed to get block: body is empty. ", "err", err)
		return nil, fmt.Errorf("failed to get block: body is empty: %s", err)
	}
	if err != nil {
		logger.Error("failed to get block", "err", err)
		return nil, fmt.Errorf("failed to get block, %s", err)
	}
	getTransactions, err := e.GetTransactions(body.Result)
	if err != nil {
		logger.Error("error during getting the transaction", "err", err)
		return nil, fmt.Errorf("error during getting the transaction, %s\n", err)
	}

	for i := range getTransactions {
		transactions = append(transactions, entities.Transaction{
			Hash:         getTransactions[i].Hash,
			From:         getTransactions[i].From,
			To:           getTransactions[i].To,
			BlockNumber:  getTransactions[i].BlockNumber,
			Gas:          getTransactions[i].Gas,
			GasPrice:     getTransactions[i].GasPrice,
			Timestamp:    getTransactions[i].Timestamp,
			AcceptNumber: getTransactions[i].AcceptNumber,
		})
	}
	_, err = e.Repos.Transaction.Insert(body.Result, transactions)
	if err != nil {
		fmt.Errorf("error on inserting")
		return nil, err
	}
	return nil, nil
}
