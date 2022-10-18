package etherscan

import (
	"etherenum-api/etherenum-service/api/internal/config"
	"etherenum-api/etherenum-service/api/pkg/json"
	"fmt"
)

var _ Scanner = (*Etherscan)(nil)

type Etherscan struct {
	Config *config.Config
}

func NewEtherscan(config *config.Config) *Etherscan {
	return &Etherscan{Config: config}
}

type getBlockNumberBody struct {
	ID     int
	Result string
}

func (e *Etherscan) GetBlock() (*getBlockNumberBody, error) {
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

func (e *Etherscan) GetTransactions(result string) ([]Transaction, error) {
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

func (e *Etherscan) AcceptIncrement(newBlockTransactions, oldBlockTransactions []Transaction) ([]Transaction, error) {
	var i int
	var testBlockTransactions []Transaction
	for {
		switch true {
		case newBlockTransactions == nil || oldBlockTransactions == nil:
			return nil, fmt.Errorf("nil pointer exception")
		case e.CompareTransactions(&newBlockTransactions[i], &Transactions{Trans: oldBlockTransactions}):
			newBlockTransactions[i].AcceptNumber++
			fmt.Printf("Transactions from: %s, to: %s, iterations: %d\n", newBlockTransactions[i].From, newBlockTransactions[i].To, newBlockTransactions[i].AcceptNumber)
		case i == len(newBlockTransactions)-1 || i == len(oldBlockTransactions)-1:
			testBlockTransactions = newBlockTransactions
			return testBlockTransactions, nil
		}
		i++
	}
}

type Compare struct {
	From     string
	To       string
	Gas      string
	GasPrice string
	AcceptNumber int
}

func (e *Etherscan) CompareTransactions(block *Transaction, oldblock *Transactions) bool {
	compareNewBlock := Compare{
		From:     block.From,
		To:       block.To,
		Gas:      block.Gas,
		GasPrice: block.GasPrice,
	}
	for i := range oldblock.Trans {
		compareOldBlock := Compare{
			From:     oldblock.Trans[i].From,
			To:       oldblock.Trans[i].To,
			Gas:      oldblock.Trans[i].Gas,
			GasPrice: oldblock.Trans[i].GasPrice,
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
