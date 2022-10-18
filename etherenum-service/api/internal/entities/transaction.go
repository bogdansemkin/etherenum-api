package entities

type Transactions struct {
	Trans []Transaction
}

type Transaction struct {
	Blockhash        string `bson:"blockhash" json:"blockhash"`
	BlockNumber      string `bson:"block_number" json:"block_number"`
	From             string `bson:"from" json:"from"`
	Gas              string `bson:"gas" json:"gas"`
	GasPrice         string `bson:"gas_price" json:"gas_price"`
	Hash             string `bson:"hash" json:"hash"`
	Input            string `bson:"input" json:"input"`
	Nonce            string `bson:"nonce" json:"nonce"`
	To               string `bson:"to" json:"to"`
	TransactionIndex string `bson:"transaction_index" json:"transaction_index"`
	ChainId          string `bson:"chain_id" json:"chain_id"`
	AcceptNumber     int    `bson:"accept_number" json:"accept_number"`
}
