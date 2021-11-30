package models

type TransactionRemove struct {
	TransactionId int64
}

type BulkTransactions []Transaction

type PagedTransactions struct {
	Transactions []Transaction
	Count        int64
}
type Transaction struct {
	Id         int64
	Amount     float32
	SpentAt    string
	Note       string
	CategoryId int32
}
