package models

type Category struct {
	Id   int32
	Name string
}

type Transaction struct {
	Id         string
	Amount     float32
	SpentAt    string
	Note       string
	CategoryId int32
}
