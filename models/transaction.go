package models

type Transaction struct {
	Id         string
	Amount     float32
	SpentAt    string
	Note       string
	CategoryId int32
}
