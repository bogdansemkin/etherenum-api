package entities

import "time"

type Transactions struct {
	Trans []Transaction
}

type Transaction struct {
	Hash         string `bson:"hash" json:"hash"`
	From         string `bson:"from" json:"from"`
	To           string `bson:"to" json:"to"`
	BlockNumber  string `bson:"blocknumber" json:"block_number"`
	Gas          string `bson:"gas" json:"gas"`
	GasPrice     string `bson:"gasprice" json:"gas_price"`
	Timestamp    string `bson:"timestamp" json:"timestamp"`
	CreateAt     time.Time `bson:"create_at" json:"-"`
}
