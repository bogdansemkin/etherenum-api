package etherscan

import (
	"etherenum-api/etherenum-service/api/internal/config"
	"etherenum-api/etherenum-service/api/internal/service"
	"etherenum-api/etherenum-service/api/pkg/json"
	"etherenum-api/etherenum-service/api/pkg/logger"
	"fmt"
)

var _ Scanner = (*etherscan)(nil)

type etherscan struct {
	Config *config.Config
	Logger logger.Logger
	Repos  service.Repos
}

func NewEtherscan(config *config.Config, logger logger.Logger, repos service.Repos) *etherscan {
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
	var testBlockTransactions []Transaction
	for {
		switch true {
		case newBlockTransactions == nil || oldBlockTransactions == nil:
			return nil, fmt.Errorf("nil pointer exception")
		case i == len(newBlockTransactions)-1 || i == len(oldBlockTransactions)-1:
			testBlockTransactions = newBlockTransactions
			return testBlockTransactions, nil

		case e.CompareTransactions(&newBlockTransactions[i], &Transactions{Trans: oldBlockTransactions}):
			newBlockTransactions[i].AcceptNumber++
			fmt.Printf("Transactions from: %s, to: %s, iterations: %d\n", newBlockTransactions[i].From, newBlockTransactions[i].To, newBlockTransactions[i].AcceptNumber)
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
			fmt.Printf("hash new block: %s\n", block.Hash)
			fmt.Printf("hash old block: %s\n", oldblock.Trans[i].Hash)
			block.AcceptNumber = compareOldBlock.AcceptNumber + 1
			return true
		}
	}
	return false
}

func (e *etherscan) HandlingTransactions([]string)  ([]string, error) {
	logger := e.Logger.Named("handlingTransactions")

	var transactions []Transaction
	var trainers []interface{}

	body, err := e.GetBlock()
	if body == nil {
		logger.Error("failed to get block: body is empty. ", "err", err)
		return nil, fmt.Errorf("failed to get block: body is empty: %s", err)
	}


	if err != nil {
		logger.Error("failed to get block", "err", err)
		return nil, fmt.Errorf("failed to get block, %s", err)
	}
	log := logger.GetLogs()
	if log[len(log)-1] == body.Result {
		logger.Error("repetitive object on result")
		return nil, fmt.Errorf("repetitive object on result")
	}
	log = logger.CreateLog(body.Result)
	fmt.Println(logger.GetLogs())
	logger.Info("logs", log)

	transactionsOutNewBlock, err := e.GetTransactions(body.Result)
	if err != nil {
		logger.Error("error during getting the transaction", "err", err)
		return nil, fmt.Errorf("error during getting the transaction, %s\n", err)
	}

	transactionsOutPreviousBlock, err := e.Repos.Transactions.GetByFilter(log[len(log)-1])
	if err != nil {
		logger.Error("error during getting the transaction", "err", err)
		return nil, fmt.Errorf("error during getting the transaction, %s\n", err)
	}

	for i := range transactionsOutPreviousBlock.Trans {
		transactions = append(transactions, Transaction{
			BlockNumber:  transactionsOutPreviousBlock.Trans[i].BlockNumber,
			From:         transactionsOutPreviousBlock.Trans[i].From,
			Gas:          transactionsOutPreviousBlock.Trans[i].Gas,
			GasPrice:     transactionsOutPreviousBlock.Trans[i].GasPrice,
			Hash:         transactionsOutPreviousBlock.Trans[i].Hash,
			To:           transactionsOutPreviousBlock.Trans[i].To,
			Timestamp:    transactionsOutPreviousBlock.Trans[i].Timestamp,
			AcceptNumber: transactionsOutPreviousBlock.Trans[i].AcceptNumber,
		})
	}
	incrementedTransactions, err := e.AcceptIncrement(transactionsOutNewBlock, transactions)
	if incrementedTransactions == nil {
		incrementedTransactions = transactionsOutNewBlock
	}
	for i := range incrementedTransactions {
		trainers = append(trainers, incrementedTransactions[i])
	}

	err = e.Repos.Transactions.Insert(trainers)
	if err != nil {
		logger.Error("error during inserting", "err", err)
		return nil, fmt.Errorf("error during inserting, %s\n", err)
	}
	return log, nil
}
