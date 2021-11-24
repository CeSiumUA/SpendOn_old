package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"spendon/models"
	"spendon/settings"
	"spendon/storage"
)

func main() {
	settings := settings.LoadSettings()
	if settings.IsValid() {
		storage.StartConnection(settings.Driver, settings.Host, settings.User, settings.Password)
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
			rw.Write([]byte("Please, use UPDATE method to get categories!"))
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
			rw.Write([]byte("Please, use DELETE method to get categories!"))
			return
		}

		decoder := json.NewDecoder(r.Body)
		removeTransaction := models.TransactionRemove{}
		err := decoder.Decode(&removeTransaction)
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
}
