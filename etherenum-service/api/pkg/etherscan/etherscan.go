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
	Config  *config.Config
	Logger  logger.Logger
	Service service.Service
}

func NewEtherscan(config *config.Config, logger logger.Logger, service service.Service) *etherscan {
	return &etherscan{
		Config:  config,
		Logger:  logger,
		Service: service,
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

func (e *etherscan) InputData() error {
	logger := e.Logger.Named("InputData")
	var transactions []entities.Transaction

	body, err := e.GetBlock()
	if body == nil {
		logger.Error("failed to get block: body is empty. ", "err", err)
		return fmt.Errorf("failed to get block: body is empty: %s", err)
	}

	if err != nil {
		logger.Error("failed to get block", "err", err)
		return fmt.Errorf("failed to get block, %s", err)
	}
	if len(log) == 0 {
		log = logger.GetLogs()
	}
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
		return fmt.Errorf("error during getting the transaction, %s\n", err)
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
	_, err = e.Service.Transaction.Insert(body.Result, transactions)
	if err != nil {
		logger.Error("error on transaction insert", "err", err)
		return fmt.Errorf("error on transaction insert, %s", err)
	}

	return nil
}
