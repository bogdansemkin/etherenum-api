package etherscan

import (
	"etherenum-api/etherenum-service/api/internal/entities"
)

type Scanner interface {
	//GetBlock - Returns the number of most recent block
	GetBlock() (*getBlockNumberBody, error)
	//GetTransactions - returns slice of block's transactions
	GetTransactions(result string) ([]entities.Transaction, error)
	//InputTransactions - sends transactions to service
	InputTransactions() error
	//InitBlocks - init blocks while collection is empty
	InitBlocks() error
}

type Result struct {
	Difficulty   string
	ExtraData    string
	GasLimit     string
	GasUsed      string
	Hash         string
	Timestamp    string
	Transactions []Transaction
}

type (
	NilSliceError struct{ error }
)

type Transactions struct {
	Trans []Transaction
}

type Transaction struct {
	Blockhash        string
	BlockNumber      string
	From             string
	Gas              string
	GasPrice         string
	Hash             string
	Input            string
	Nonce            string
	To               string
	Value        string
	TransactionIndex string
	ChainId          string
	Timestamp        string
}
