package models

type TransactionRemove struct {
	TransactionId string
}
type Transaction struct {
	Id         string
	Amount     float32
	SpentAt    string
	Note       string
	CategoryId int32
}
