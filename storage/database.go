package storage

import (
	"database/sql"
	"fmt"
	"net/url"
	"spendon/models"

	mssql "github.com/denisenkom/go-mssqldb"
)

const (
	insertTransaction = "INSERT INTO dbo.Transactions (Id, Amount, SpentAt, Note, CategoryId) VALUES (newid(), @AMOUNT, @SpentAt, @Note, @Category)"
	selectCategories  = "SELECT * FROM dbo.Categories"
)

var databaseConnection *sql.DB

func StartConnection(driverName, host, user, password string) {
	tst := &url.URL{
		Scheme: driverName,
		Host:   host,
		User:   url.UserPassword(user, password),
	}
	connString := tst.String()
	dbConnector, err := mssql.NewConnector(connString)
	if err != nil {
		fmt.Println("DB connection error", err)
	}
	dbConnector.SessionInitSQL = "USE SpendonDB"
	databaseConnection = sql.OpenDB(dbConnector)
}

func InsertTransaction(transaction *models.Transaction) {
	if databaseConnection == nil {
		fmt.Println("DB not connected!")
		return
	}
	rslt, err := databaseConnection.Exec(insertTransaction,
		sql.Named("AMOUNT", transaction.Amount),
		sql.Named("SpentAt", transaction.SpentAt),
		sql.Named("Note", transaction.Note),
		sql.Named("Category", transaction.CategoryId))
	if err != nil {
		fmt.Println(err)
	}
	rowsAffectedCount, err := rslt.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Rows affected:", rowsAffectedCount)
}

func GetCategories() (models.Categories, error) {
	if databaseConnection == nil {
		return nil, fmt.Errorf("DB not connected!")
	}
	categories := make(models.Categories, 0)
	rows, err := databaseConnection.Query(selectCategories)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		category := models.Category{}
		err := rows.Scan(&category.Id, &category.Name)
		if err != nil {
			fmt.Println(err)
			return categories, err
		} else {
			categories = append(categories, category)
		}
	}
	return categories, err
}
