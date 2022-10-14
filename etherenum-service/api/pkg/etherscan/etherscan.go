package etherscan

import (
	"etherenum-api/etherenum-service/api/internal/config"
	"etherenum-api/etherenum-service/api/pkg/json"
	"fmt"
)

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

func (e *Etherscan) getBlock() (*getBlockNumberBody, error) {
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

type Result struct {
	Difficulty   string
	ExtraData    string
	GasLimit     string
	GasUsed      string
	Hash         string
	Transactions []transaction
}
type transaction struct {
	Blockhash        string
	BlockNumber      string
	From             string
	Gas              string
	GasPrice         string
	Hash             string
	Input            string
	Nonce            string
	To               string
	TransactionIndex string
	ChainId          string
}

func (e *Etherscan) GetTransactions() ([]transaction, error) {
	var body GetTransactionsBody
	blockResult, err := e.getBlock()
	if err != nil {
		return nil, fmt.Errorf("error during getting json from etherscan, %s", err)
	}

	err = json.GetJson(fmt.Sprintf("https://api.etherscan.io/api?module=proxy&action=eth_getBlockByNumber&tag=%s&boolean=true&apikey=%s", blockResult.Result, e.Config.Etherscan.Key), &body)
	if err != nil {
		return nil, fmt.Errorf("error during getting json from etherscan, %s", err)
	}

	return body.Result.Transactions, nil
}
