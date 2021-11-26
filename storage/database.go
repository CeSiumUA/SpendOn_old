package storage

import (
	"crypto"
	"database/sql"
	"fmt"
	"net/url"
	"spendon/models"

	mssql "github.com/denisenkom/go-mssqldb"
)

const (
	insertTransaction = "INSERT INTO dbo.Transactions (Id, Amount, SpentAt, Note, CategoryId, UserId) VALUES (newid(), @AMOUNT, @SpentAt, @Note, @Category, @UserId)"
	selectCategories  = "SELECT * FROM dbo.Categories"
	updateTransaction = "UPDATE dbo.Transactions SET Amount=@AMOUNT, SpentAt=@SPENTAT, Note=@NOTE, CategoryId=@CATEGORYID where Id=@ID and UserId=@UserId"
	removeTransaction = "DELETE FROM dbo.Transactions WHERE Id=@ID and UserId=@UserId"
	getUserByPassword = "SELECT Id, Login from dbo.Users WHERE Login=@LOGIN and PasswordHash=@PWD"
	getUserByLogin    = "SELECT Id, Login from dbo.Users WHERE Login=@LOGIN"
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

func InsertTransaction(transaction *models.Transaction, userId string) error {
	if databaseConnection == nil {
		fmt.Println("DB not connected!")
		return fmt.Errorf("DB not connected")
	}
	rslt, err := databaseConnection.Exec(insertTransaction,
		sql.Named("AMOUNT", transaction.Amount),
		sql.Named("SpentAt", transaction.SpentAt),
		sql.Named("Note", transaction.Note),
		sql.Named("Category", transaction.CategoryId),
		sql.Named("UserId", userId))
	if err != nil {
		fmt.Println(err)
		return err
	}
	rowsAffectedCount, err := rslt.RowsAffected()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Rows affected:", rowsAffectedCount)
	return nil
}

func GetCategories() (models.Categories, error) {
	if databaseConnection == nil {
		return nil, fmt.Errorf("DB not connected")
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

func UpdateTransaction(transaction *models.Transaction, userId string) (*models.Transaction, error) {
	if databaseConnection == nil {
		return &models.Transaction{}, fmt.Errorf("DB not connected")
	}
	result, err := databaseConnection.Exec(updateTransaction,
		sql.Named("AMOUNT", transaction.Amount),
		sql.Named("SPENTAT", transaction.SpentAt),
		sql.Named("NOTE", transaction.Note),
		sql.Named("CATEGORYID", transaction.CategoryId),
		sql.Named("ID", transaction.Id),
		sql.Named("UserId", userId))
	if err != nil {
		fmt.Println(err)
		return &models.Transaction{}, err
	}
	fmt.Println("Update result:", result)
	return transaction, nil
}

func RemoveTransaction(id string, userId string) error {
	if databaseConnection == nil {
		return fmt.Errorf("DB not connected")
	}
	result, err := databaseConnection.Exec(removeTransaction,
		sql.Named("ID", id),
		sql.Named("UserId", userId))
	if err != nil {
		return err
	}
	fmt.Println("Delete result:", result)
	return nil
}

func GetUserByPassword(password, login string) (*models.DbLogin, error) {
	if databaseConnection == nil {
		return &models.DbLogin{}, fmt.Errorf("DB not connected")
	}
	dbLogin := models.DbLogin{}
	pwdHasher := crypto.SHA512.New()
	pwdHash := pwdHasher.Sum([]byte(password))
	row := databaseConnection.QueryRow(getUserByPassword,
		sql.Named("LOGIN", login),
		sql.Named("PWD", string(pwdHash)))
	err := row.Scan(&dbLogin.Id, &dbLogin.Login)
	if err != nil {
		return &dbLogin, err
	}
	return &dbLogin, nil
}

func GetUserByLogin(login string) (*models.DbLogin, error) {
	dbLogin := models.DbLogin{}
	if databaseConnection == nil {
		return &dbLogin, fmt.Errorf("DB not connected")
	}
	row := databaseConnection.QueryRow(getUserByLogin,
		sql.Named("LOGIN", login))
	err := row.Scan(&dbLogin.Id, &dbLogin.Login)
	if err != nil {
		return &dbLogin, err
	}
	return &dbLogin, nil
}
