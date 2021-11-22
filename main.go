package main

import (
	"fmt"
	"net/http"
	"spendon/settings"
	"spendon/storage"
)

func main() {
	settings := settings.LoadSettings()
	storage.StartConnection(settings.Driver, settings.Host, settings.User, settings.Password)

	err := http.ListenAndServe(":1462", nil)
	if err != nil {
		fmt.Println(err)
	}
}
