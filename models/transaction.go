package models

type TransactionRemove struct {
	TransactionId string
}

type BulkTransactions []Transaction
type Transaction struct {
	Id         int64
	Amount     float32
	SpentAt    string
	Note       string
	CategoryId int32
}
