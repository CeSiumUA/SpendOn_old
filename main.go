package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	registerHandlers()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Listener creation error:", err)
	}
}

func registerHandlers() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "Hello from SpendOn server!")
	})
	http.HandleFunc("/add", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use POST method to add new transactions!"))
			return
		}
		authTokenHeader := r.Header.Get("Token")
		err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}
		transaction := models.Transaction{}

		decoder := json.NewDecoder(r.Body)

		decoder.Decode(&transaction)
		storage.InsertTransaction(&transaction)
	})
	http.HandleFunc("/getcategories", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
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
	http.HandleFunc("/updatetransaction", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use UPDATE method to update transactions!"))
			return
		}
		authTokenHeader := r.Header.Get("Token")
		err := ValidateLoginToken(authTokenHeader)
		if err != nil {
			fmt.Println(err)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("Authorize failure!"))
			return
		}
		transaction := models.Transaction{}

		decoder := json.NewDecoder(r.Body)

		decoder.Decode(&transaction)

		resultTransaction, err := storage.UpdateTransaction(&transaction)

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
	http.HandleFunc("/removetransaction", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use DELETE method to remove transaction!"))
			return
		}

		authTokenHeader := r.Header.Get("Token")
		err := ValidateLoginToken(authTokenHeader)
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
		err = storage.RemoveTransaction(removeTransaction.TransactionId)
		if err != nil {
			fmt.Println("Remove transaction error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte("An error occured on the server! This message is already delivered to developer ;)"))
			return
		}
	})
	http.HandleFunc("/login", func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			rw.WriteHeader(http.StatusMethodNotAllowed)
			rw.Write([]byte("Please, use POST method to login!"))
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

		if loadedSettings.AllowedLogin == "" || loadedSettings.AllowedPassword == "" || loadedSettings.SigningSecret == "" {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("User, password or secret are not set on the server!"))
			return
		}
		if loadedSettings.AllowedLogin != loginModel.UserName || loadedSettings.AllowedPassword != loginModel.Password {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("User or password are incorrect!"))
			return
		}
		expireDate := time.Now().Add(24 * time.Hour)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user": loadedSettings.AllowedLogin,
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
}

func ValidateLoginToken(token string) error {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(loadedSettings.SigningSecret), nil
	})
	if err != nil {
		return err
	}
	claims, ok := jwtToken.Claims.(*jwt.MapClaims)
	if !ok {
		return fmt.Errorf("Validation error")
	}
	validationError := claims.Valid()
	if validationError != nil {
		return validationError
	}
	userName := (*claims)["user"]
	if userName != loadedSettings.AllowedLogin {
		return fmt.Errorf("User not found")
	}
	return nil
}
