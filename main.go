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
		transaction := models.Transaction{}

		decoder := json.NewDecoder(r.Body)

		decoder.Decode(&transaction)
		storage.InsertTransaction(&transaction)
	})
}
