package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spendon/models"
	"spendon/settings"
	"spendon/storage"
	"time"

	"github.com/golang-jwt/jwt"
)

var loadedSettings *settings.Settings

func main() {
	loadedSettings = settings.LoadSettings()
	if loadedSettings.IsValid() {
		storage.StartConnection(loadedSettings.Driver, loadedSettings.Host, loadedSettings.User, loadedSettings.Password)
	} else {
		fmt.Println("Settings were not loaded")
	}
	//serveSPA()
	registerHandlers()
	port := loadedSettings.Port
	if port == "" {
		port = "8080"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Listener creation error:", err)
	}
}

func registerHandlers() {
	http.HandleFunc("/api/add", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use POST method to add new transactions!"))
			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		authTokenHeader := r.Header.Get("Token")
		dbLogin, err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}
		transaction := models.Transaction{}

		decoder := json.NewDecoder(r.Body)

		err = decoder.Decode(&transaction)
		if err != nil {
			fmt.Println("Encoding response error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
		}
		storage.InsertTransaction(&transaction, dbLogin.Id)
	})

	http.HandleFunc("/api/bulkadd", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use POST method to add new transactions!"))
			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		authTokenHeader := r.Header.Get("Token")
		dbLogin, err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}
		transactions := make(models.BulkTransactions, 0)

		decoder := json.NewDecoder(r.Body)

		err = decoder.Decode(&transactions)
		if err != nil {
			fmt.Println("Encoding response error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
		}
		for _, transaction := range transactions {
			storage.InsertTransaction(&transaction, dbLogin.Id)
		}
	})
	http.HandleFunc("/api/getcategories", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method != http.MethodGet && r.Method != http.MethodOptions {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use GET method to get categories!"))
			return
		}
		categories, err := storage.GetCategories()
		if err != nil {
			fmt.Println("Category getching error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
		encoder := json.NewEncoder(rw)
		err = encoder.Encode(categories)
		if err != nil {
			fmt.Println("Encoding response error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
	})
	http.HandleFunc("/api/updatetransaction", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method != http.MethodPut && r.Method != http.MethodOptions {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use UPDATE method to update transactions!"))
			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		authTokenHeader := r.Header.Get("Token")
		dbLogin, err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}
		transaction := models.Transaction{}

		decoder := json.NewDecoder(r.Body)

		decoder.Decode(&transaction)

		resultTransaction, err := storage.UpdateTransaction(&transaction, dbLogin.Id)

		if err != nil {
			fmt.Println("Update transaction error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
		} else {
			encoder := json.NewEncoder(rw)
			err := encoder.Encode(*resultTransaction)
			if err != nil {
				fmt.Println("Encoding response error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			}
		}
	})
	http.HandleFunc("/api/removetransaction", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method != http.MethodDelete && r.Method != http.MethodOptions {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use DELETE method to remove transaction!"))
			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		authTokenHeader := r.Header.Get("Token")
		dbLogin, err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}

		decoder := json.NewDecoder(r.Body)
		removeTransaction := models.TransactionRemove{}
		err = decoder.Decode(&removeTransaction)
		if err != nil {
			fmt.Println("Decode body error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
		err = storage.RemoveTransaction(removeTransaction.TransactionId, dbLogin.Id)
		if err != nil {
			fmt.Println("Remove transaction error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
	})
	http.HandleFunc("/api/login", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method != http.MethodPost && r.Method != http.MethodOptions {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use POST method to login!"))
			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		decoder := json.NewDecoder(r.Body)
		loginModel := models.Login{}
		err := decoder.Decode(&loginModel)
		if err != nil {
			fmt.Println("Decode body error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}

		if loadedSettings.SigningSecret == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("Secret is not set on the server!"))
			return
		}

		dbLogin, err := storage.GetUserByPassword(loginModel.Password, loginModel.UserName)

		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if dbLogin.Login != loginModel.UserName {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("User or password are incorrect!"))
			return
		}
		expireDate := time.Now().Add(7 * 24 * time.Hour)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": dbLogin.Login,
			"exp":  expireDate.Unix(),
		})
		secretKey := []byte(loadedSettings.SigningSecret)
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			fmt.Println("Token create error!:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
		loginResult := models.LoginResult{
			Token:      tokenString,
			ExpireDate: expireDate,
		}
		encoder := json.NewEncoder(rw)
		err = encoder.Encode(loginResult)
		if err != nil {
			fmt.Println("Encoding response error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
	})
	http.HandleFunc("/api/checkauth", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method == http.MethodOptions {
			return
		}
		authTokenHeader := r.Header.Get("Token")
		dbLogin, err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}
		rw.WriteHeader(http.StatusOK)
		fmt.Println("User found: " + dbLogin.Login)
	})
	http.HandleFunc("/api/fetchtransactions", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use POST method to fetch transactions!"))
			return
		}
		authTokenHeader := r.Header.Get("Token")
		dbLogin, err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}
		transactions, err := storage.GetAllTransactions(dbLogin.Id)
		if err != nil {
			fmt.Println("Fetching transactions error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
		encoder := json.NewEncoder(rw)
		err = encoder.Encode(transactions)
		if err != nil {
			fmt.Println("Encoding response error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
	})
	http.HandleFunc("/api/getcategoriesstats", func(rw http.ResponseWriter, r *http.Request) {
		SetCORS(&rw)
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use POST method to get stats!"))
			return
		}

		authTokenHeader := r.Header.Get("Token")
		dbLogin, err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}

		categorySummaries, err := storage.GetTransactionsSummary(dbLogin.Id)
		if err != nil {
			fmt.Println("Fetching transactions error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
		encoder := json.NewEncoder(rw)
		err = encoder.Encode(categorySummaries)
		if err != nil {
			fmt.Println("Encoding response error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
	})
}

func ValidateLoginToken(token string) (*models.DbLogin, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(loadedSettings.SigningSecret), nil
	})
	if err != nil {
		return &models.DbLogin{}, err
	}
	claims, ok := jwtToken.Claims.(*jwt.MapClaims)
	if !ok {
		return &models.DbLogin{}, fmt.Errorf("validation error")
	}
	validationError := claims.Valid()
	if validationError != nil {
		return &models.DbLogin{}, validationError
	}
	userName := (*claims)["user"]
	userNameStr := fmt.Sprintf("%v", userName)
	dbLogin, err := storage.GetUserByLogin(userNameStr)
	if err != nil {
		return dbLogin, err
	}
	if userNameStr != dbLogin.Login {
		return dbLogin, fmt.Errorf("user not found")
	}
	return dbLogin, nil
}

func SetCORS(rw *http.ResponseWriter) {
	(*rw).Header().Add("Access-Control-Allow-Origin", "*")
	(*rw).Header().Add("Access-Control-Allow-Methods", "*")
	(*rw).Header().Add("Access-Control-Allow-Headers", "*")
}
