package etherscan

type Scanner interface {
	GetBlock() (*getBlockNumberBody, error)
	GetTransactions(result string) ([]Transaction, error)
	InputData()  error
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
	NilSliceError struct { error }
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
	TransactionIndex string
	ChainId          string
	Timestamp        string
	AcceptNumber     int
}
