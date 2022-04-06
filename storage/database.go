package storage

import (
	"crypto/sha256"
	"fmt"
	"github.com/jackc/pgx"
	"spendon/models"
)

const (
	insertTransaction           = "INSERT INTO transactions (amount, spentat, note, categoryid, userid) VALUES ($1, $2, $3, $4, $5)"
	insertUser                  = "IF NOT EXISTS (select * from dbo.Users where login=$1) THEN INSERT INTO users (login, passwordhash, currency) VALUES($1, $2, $3) END IF"
	selectCategories            = "SELECT * FROM categories"
	updateTransaction           = "UPDATE transactions SET amount=$1, spentat=$2, note=$3, categoryid=$4 where id=$5 and userid=$6"
	removeTransaction           = "DELETE FROM transactions WHERE id=$1 and userid=$2"
	getPaginatedTransactions    = "SELECT id, amount::numeric, spentat::text, note, categoryid FROM transactions WHERE %s userId=$%d ORDER BY spentat DESC OFFSET $%d ROWS FETCH NEXT $%d ROWS ONLY"
	getUserByPassword           = "SELECT id, login from users WHERE login=$1 and passwordhash=$2"
	getUserByLogin              = "SELECT id, login from users WHERE login=$1"
	getStatistics               = "SELECT categoryid , SUM(amount)::numeric from transactions where %s userid=$%d GROUP BY categoryid"
	getTransactionsCountForUser = "SELECT COUNT(*) as cnt FROM transactions WHERE %s userid=$1"
)

var databaseConnection *pgx.Conn

func StartConnection(connectionUrl string) *pgx.Conn {
	connStr, err := pgx.ParseConnectionString(connectionUrl)
	if err != nil {
		fmt.Println("error parsing db url", err)
	}
	conn, err := pgx.Connect(connStr)
	if err != nil {
		fmt.Println("DB connection error", err)
	}
	databaseConnection = conn
	return conn
}

func InsertTransaction(transaction *models.Transaction, userId int64) error {
	if databaseConnection == nil {
		fmt.Println("DB not connected!")
		return fmt.Errorf("DB not connected")
	}
	rslt, err := databaseConnection.Exec(insertTransaction,
		fmt.Sprintf("$%f", transaction.Amount),
		transaction.SpentAt,
		transaction.Note,
		transaction.CategoryId,
		userId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	rowsAffectedCount := rslt.RowsAffected()
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

func UpdateTransaction(transaction *models.Transaction, userId int64) (*models.Transaction, error) {
	if databaseConnection == nil {
		return &models.Transaction{}, fmt.Errorf("DB not connected")
	}
	result, err := databaseConnection.Exec(updateTransaction,
		transaction.Amount,
		transaction.SpentAt,
		transaction.Note,
		transaction.CategoryId,
		transaction.Id,
		userId)
	if err != nil {
		fmt.Println(err)
		return &models.Transaction{}, err
	}
	fmt.Println("Update result:", result.RowsAffected())
	return transaction, nil
}

func RemoveTransaction(id, userId int64) error {
	if databaseConnection == nil {
		return fmt.Errorf("DB not connected")
	}
	result, err := databaseConnection.Exec(removeTransaction,
		id,
		userId)
	if err != nil {
		return err
	}
	fmt.Println("Delete result:", result.RowsAffected())
	return nil
}

func GetUserByPassword(password, login string) (*models.DbLogin, error) {
	if databaseConnection == nil {
		return &models.DbLogin{}, fmt.Errorf("DB not connected")
	}
	dbLogin := models.DbLogin{}
	pwdHash := sha256.Sum256([]byte(password))
	pwdHashString := fmt.Sprintf("%x", pwdHash)
	row := databaseConnection.QueryRow(getUserByPassword,
		login,
		pwdHashString)
	err := row.Scan(&dbLogin.Id, &dbLogin.Login)
	if err != nil {
		return &dbLogin, err
	}

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
	row := databaseConnection.QueryRow(getUserByLogin, login)

	err := row.Scan(&dbLogin.Id, &dbLogin.Login)
	if err != nil {
		return &dbLogin, err
	}

	if err != nil {
		return &dbLogin, err
	}
	return &dbLogin, nil
}

func GetFilteredTransactions(userId, pageNumber, pagination int64, filterBatch *models.FilterBatch) (models.PagedTransactions, error) {
	if databaseConnection == nil {
		return models.PagedTransactions{}, fmt.Errorf("DB not connected")
	}

	filterString, namedArgs, err := filterBatch.Build()

	if err != nil {
		return models.PagedTransactions{}, err
	}

	offset := pageNumber * pagination

	parameterIndex := len(namedArgs)

	namedArgs = append(namedArgs, userId)

	countArgs := namedArgs

	namedArgs = append(namedArgs, offset)
	namedArgs = append(namedArgs, pagination)

	interfaceArgs := make([]interface{}, 0)

	for _, arg := range namedArgs {
		interfaceArgs = append(interfaceArgs, arg)
	}

	formattedTransaction := fmt.Sprintf(getPaginatedTransactions, filterString, parameterIndex+1, parameterIndex+2, parameterIndex+3)
	rows, err := databaseConnection.Query(formattedTransaction,
		interfaceArgs...)
	if err != nil {
		fmt.Println(err)
		return models.PagedTransactions{}, err
	}
	bulkTransactions := models.PagedTransactions{}
	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		transaction := models.Transaction{}
		err := rows.Scan(&transaction.Id, &transaction.Amount, &transaction.SpentAt, &transaction.Note, &transaction.CategoryId)
		if err != nil {
			fmt.Println(err)
			return models.PagedTransactions{}, err
		} else {
			transactions = append(transactions, transaction)
		}
	}
	bulkTransactions.Transactions = transactions

	countInterfaceArgs := make([]interface{}, 0)

	for _, arg := range countArgs {
		countInterfaceArgs = append(countInterfaceArgs, arg)
	}

	row := databaseConnection.QueryRow(fmt.Sprintf(getTransactionsCountForUser, filterString),
		countInterfaceArgs...)
	err = row.Scan(&bulkTransactions.Count)
	if err != nil {
		return models.PagedTransactions{}, nil
	}
	return bulkTransactions, nil
}

func GetTransactionsSummary(userId int64, filterBatch models.FilterBatch) (models.CategoriesSummary, error) {
	if databaseConnection == nil {
		return nil, fmt.Errorf("DB not connected")
	}

	filterString, namedArgs, err := filterBatch.Build()

	if err != nil {
		return nil, err
	}

	parameterIndex := len(namedArgs)

	namedArgs = append(namedArgs, userId)

	interfaceArgs := make([]interface{}, 0)

	for _, arg := range namedArgs {
		interfaceArgs = append(interfaceArgs, arg)
	}

	formattedRequest := fmt.Sprintf(getStatistics, filterString, parameterIndex+1)

	rows, err := databaseConnection.Query(formattedRequest,
		interfaceArgs...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	categoriesSummary := make(models.CategoriesSummary, 0)
	for rows.Next() {
		categorySummary := models.CategorySummary{}
		err := rows.Scan(&categorySummary.CategoryId, &categorySummary.Sum)
		if err != nil {
			fmt.Println(err)
			return categoriesSummary, err
		} else {
			categoriesSummary = append(categoriesSummary, categorySummary)
		}
	}
	return categoriesSummary, nil
}

func AddUser(registerModel *models.RegisterModel) (bool, error) {
	if databaseConnection == nil {
		return false, fmt.Errorf("DB not connected")
	}

	pwdHash := sha256.Sum256([]byte(registerModel.Password))
	pwdHashString := fmt.Sprintf("%x", pwdHash)

	sqlResult, err := databaseConnection.Exec(insertUser,
		registerModel.Login,
		pwdHashString,
		"UAH")
	if err != nil {
		return false, err
	}
	rowResult := sqlResult.RowsAffected()
	if err != nil {
		return false, err
	}
	fmt.Println("Rows affected", rowResult)
	return rowResult == 1, nil
}
