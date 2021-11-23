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
	}

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
