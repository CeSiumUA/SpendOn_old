package models

type Categories []Category

type CategoriesSummary []CategorySummary

type Category struct {
	Id   int32
	Name string
}

type CategorySummary struct {
	CategoryId int64
	Sum        float32
}
