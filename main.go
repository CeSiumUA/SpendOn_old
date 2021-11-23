package main

import (
	"fmt"
	"net/http"
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

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "Hello from Go server!")
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Listener creation error:", err)
	}
}
