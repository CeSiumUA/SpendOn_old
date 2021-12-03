package models

import (
	"database/sql"
	"fmt"
	"strings"
)

type FilterBatch []FilterModel

type FilterSettings struct {
	Fields []map[string]interface{}
	Signs  []map[string]interface{}
}

type FilteredRequest struct {
	PageNumber int64
	Pagination int64
	Filters    FilterBatch
}

type FilterModel struct {
	Property int
	Value    string
	Operator int
}

func (filterModel *FilterModel) Build(nameForParameter string) (string, sql.NamedArg, error) {
	sign, ok := signsMap[filterModel.Operator]
	if !ok {
		return "", sql.NamedArg{}, fmt.Errorf("sign was not found")
	}
	namedParam := sql.Named(nameForParameter, filterModel.Value)
	paramName, ok := fieldsMap[filterModel.Property]
	if !ok {
		return "", sql.NamedArg{}, fmt.Errorf("field was not found")
	}
	return paramName + sign + fmt.Sprintf("@%s", nameForParameter), namedParam, nil
}

func (filterBatch *FilterBatch) Build() (string, []sql.NamedArg, error) {
	if len(*filterBatch) == 0 {
		return "", nil, nil
	}
	namedArgs := make([]sql.NamedArg, 0)
	params := make([]string, 0)
	for idx, el := range *filterBatch {
		paramName, namedArg, err := el.Build(fmt.Sprintf("P%d", idx+1))
		if err != nil {
			return "", nil, err
		}
		namedArgs = append(namedArgs, namedArg)
		params = append(params, paramName)
	}
	joinedParams := strings.Join(params, " and ")
	return fmt.Sprintf("%s and", joinedParams), namedArgs, nil
}

const (
	Equal = iota
	Less
	Greater
	NotEqual
	LessOrEqual
	GreaterOrEqual
)

const (
	Amount = iota
	SpentAt
	Note
	CategoryId
)

var signsMap map[int]string = map[int]string{
	Equal:          "=",
	Less:           "<",
	Greater:        ">",
	NotEqual:       "<>",
	LessOrEqual:    "<=",
	GreaterOrEqual: ">=",
}

var fieldsMap map[int]string = map[int]string{
	Amount:     "Amount",
	SpentAt:    "SpentAt",
	Note:       "Note",
	CategoryId: "CategoryId",
}

func GetFilterSettings() FilterSettings {
	fieldsData := make([]map[string]interface{}, 0, 0)

	for counter, field := range fieldsMap {
		singleMap := make(map[string]interface{})
		singleMap["index"] = counter
		singleMap["value"] = field
		fieldsData = append(fieldsData, singleMap)
	}

	signsData := make([]map[string]interface{}, 0, 0)

	for counter, field := range signsMap {
		singleMap := make(map[string]interface{})
		singleMap["index"] = counter
		singleMap["value"] = field
		signsData = append(signsData, singleMap)
	}

	return FilterSettings{
		Signs:  signsData,
		Fields: fieldsData,
	}
}
